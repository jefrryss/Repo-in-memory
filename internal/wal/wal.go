package wal

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"context"

	"in-memory/config"
	"in-memory/internal/compute/parser"

	"go.uber.org/zap"
)


type LogEntry struct {
	Query   *parser.Query
	ErrChan chan error
}

type WAL interface {
	Recovery(outChan chan<- *parser.Query) error
	Write(ctx context.Context, query *parser.Query) error
}

type FileWAL struct {
	countCommands  int
	maxSegmentSize int
	batchTimeOut   time.Duration
	directory      string

	currentSegmentID int
	file             *os.File
	currentSize      int

	queueChan chan LogEntry
	logger    *zap.Logger
}

func NewWal(cnf *config.Config, l *zap.Logger) (WAL, error) {
	if !cnf.WAL.TurnOn {
		return nil, nil
	}
	maxSegmentSize, err := config.PasreSize(cnf.WAL.MaxSegmentSize)
	if err != nil {
		return nil, err
	}
	w := &FileWAL{
		countCommands:  cnf.WAL.FlushingBatchSize,
		maxSegmentSize: maxSegmentSize,
		batchTimeOut:   cnf.WAL.FlushingBatchTimeout,
		directory:      cnf.WAL.Directory,
		queueChan:      make(chan LogEntry, 100),

		logger: l.With(zap.String("component", "wal")), 
	}

	if err := w.createIfNotExistsPath(); err != nil {
		w.logger.Error("Failed to create WAL directory", zap.Error(err), zap.String("dir", w.directory))
		return nil, err
	}
	if err := w.getCurrentIndexSegment(); err != nil {
		w.logger.Error("Failed to get current segment index", zap.Error(err))
		return nil, err
	}

	if err := w.openOrCreateCurrentSegment(); err != nil {
		w.logger.Error("Failed to open current segment", zap.Error(err))
		return nil, err
	}

	w.logger.Info("WAL successfully initialized",
		zap.String("directory", w.directory),
		zap.Int("max_segment_size", w.maxSegmentSize),
		zap.Int("batch_size", w.countCommands),
	)

	go w.workerLoop()

	return w, nil
}

func (w *FileWAL) Recovery(outChan chan<- *parser.Query) error {
	w.logger.Info("Started data recovery from WAL", zap.String("directory", w.directory))
	defer close(outChan)
	
	files, err := os.ReadDir(w.directory)
	if err != nil {
		w.logger.Error("Error reading directory", zap.Error(err), zap.String("dir", w.directory))
		return err
	}
	
	var walFiles []string
	for _, file := range files {
		name := file.Name()
		if !file.IsDir() && strings.HasPrefix(name, "wal_") && strings.HasSuffix(name, ".log") {
			walFiles = append(walFiles, name)
		}
	}
	sort.Strings(walFiles)

	w.logger.Info("Found files for recovery", zap.Int("files_count", len(walFiles)))

	totalQueries := 0
	for _, file := range walFiles {
		fullPath := filepath.Join(w.directory, file)
		openFile, err := os.Open(fullPath)

		if err != nil {
			w.logger.Error("Error opening file during recovery", zap.Error(err), zap.String("file", fullPath))
			return err
		}
		
		scanner := bufio.NewScanner(openFile)
		fileQueries := 0
		for scanner.Scan() {
			query := parser.UnMarshal(scanner.Text())
			if query != nil {
				outChan <- query
				fileQueries++
				totalQueries++
			}
		}
		openFile.Close()
		w.logger.Debug("File read successfully", zap.String("file", file), zap.Int("queries_loaded", fileQueries))
	}
	
	w.logger.Info("Recovery successfully completed", zap.Int("total_queries_loaded", totalQueries))
	return nil
}

func (w *FileWAL) Write(ctx context.Context, query *parser.Query) error {
	logEntry := LogEntry{Query: query, ErrChan: make(chan error, 1)}
	select {
	case <- ctx.Done():
		return ctx.Err()
	case w.queueChan <- logEntry:

	}
	select{
	case <-ctx.Done():
		return ctx.Err()
	case res := <-logEntry.ErrChan:
	return res
	}
}

func (w *FileWAL) workerLoop() {
	batch := make([]LogEntry, 0, w.countCommands)
	ticker := time.NewTicker(w.batchTimeOut)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if len(batch) > 0 {
				w.logger.Debug("Writing batch due to timeout", zap.Int("batch_size", len(batch)))
				w.flushBatch(batch)
				batch = batch[:0]
			}

		case res := <-w.queueChan:
			batch = append(batch, res)
			if len(batch) >= w.countCommands {
				w.logger.Debug("Writing batch due to operation limit", zap.Int("batch_size", len(batch)))
				w.flushBatch(batch)
				batch = batch[:0]
				ticker.Reset(w.batchTimeOut)
			}
		}
	}
}

func (w *FileWAL) flushBatch(batch []LogEntry) {
	var buffer bytes.Buffer
	for _, log := range batch {
		buffer.WriteString(log.Query.Marshal())
	}

	data := buffer.Bytes()
	batchLen := int(len(data))

	if w.currentSize+batchLen > w.maxSegmentSize {
		w.logger.Info("Segment limit reached, rotating",
			zap.Int("current_size", w.currentSize),
			zap.Int("batch_len", batchLen),
			zap.Int("max_limit", w.maxSegmentSize),
		)
		
		if err := w.rotateSegment(); err != nil {
			w.logger.Error("Error rotating segment", zap.Error(err))
			w.sendError(batch, err)
			return
		}
	}

	n, err := w.file.Write(data)
	if err != nil {
		w.logger.Error("Error writing data to file", zap.Error(err))
		w.sendError(batch, err)
		return
	}

	if err := w.file.Sync(); err != nil {
		w.logger.Error("Error syncing file (fsync)", zap.Error(err))
		w.sendError(batch, err)
		return
	}

	w.currentSize += int(n)
	w.sendError(batch, nil)
}

func (w *FileWAL) rotateSegment() error {
	w.file.Sync()
	w.file.Close()
	w.currentSegmentID++
	return w.openOrCreateCurrentSegment()
}

func (w *FileWAL) openOrCreateCurrentSegment() error {
	if w.currentSegmentID == 0 {
		w.currentSegmentID = 1
	}
	fileName := fmt.Sprintf("wal_%04d.log", w.currentSegmentID)
	fullpath := filepath.Join(w.directory, fileName)
	
	file, err := os.OpenFile(fullpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, _ := file.Stat()
	w.file = file
	w.currentSize = int(info.Size())

	w.logger.Info("Opened WAL segment",
		zap.String("file", fileName),
		zap.Int("current_size", w.currentSize),
	)
	
	return nil
}

func (w *FileWAL) sendError(batch []LogEntry, err error) {
	for _, log := range batch {
		log.ErrChan <- err
	}
}

func (w *FileWAL) getCurrentIndexSegment() error {
	files, err := os.ReadDir(w.directory)
	if err != nil {
		return err
	}

	var maxID int = 0
	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, "wal_") && strings.HasSuffix(name, ".log") {

			correctName := strings.TrimSuffix(strings.TrimPrefix(file.Name(), "wal_"), ".log")

			number, err := strconv.Atoi(correctName)
			if err != nil {
				continue
			}

			if maxID < number {
				maxID = number
			}
		}
	}
	w.currentSegmentID = maxID
	return nil
}

func (b *FileWAL) createIfNotExistsPath() error {
	info, err := os.Stat(b.directory)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(b.directory, 0755)
			if err == nil {
				b.logger.Info("Created directory for WAL", zap.String("dir", b.directory))
			}
			return err
		} else {
			return err
		}
	}

	if !info.IsDir() {
		return fmt.Errorf("path %s is a file, not a directory", b.directory)
	}

	return nil
}