package server 


import(
	"errors"
	"net"
	"go.uber.org/zap"
	"bufio"
	"unicode"
	"strconv"
	"strings"
	"context"
	"in-memory/internal/compute"
	"in-memory/config"
	"time"
)

type ServerTCP struct {
	logger *zap.Logger 
	com *compute.Compute

	maxConnections int
	maxMessageSize int
	addres string
	timeout time.Duration
}



func NewServerTSP(cnf *config.ServerConfig, c *compute.Compute, l *zap.Logger) *ServerTCP{
	maxMessage, err := parseMaxSize(cnf.MaxMessageSize)
	if err != nil {
		l.Fatal("Invalid max message size in config", zap.Error(err))
	}
	return &ServerTCP{
		maxMessageSize: maxMessage,
		maxConnections: cnf.MaxConnections,
		addres: cnf.Address,
		com: c,
		logger: l,
		timeout: cnf.IdleTimeout,
	}
}

func (s *ServerTCP) StartServer(){
	listener, err := net.Listen("tcp", s.addres)
	if err != nil {
		s.logger.Fatal("Error with starting server", zap.Error(err))
	}
	defer listener.Close()
	s.logger.Info("TCP server started", zap.String("address", s.addres))

	semaphore := make(chan struct{}, s.maxConnections)
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		semaphore <- struct{}{}

		go func(c net.Conn) {
			defer func (){
				<- semaphore
			}()
			s.handleConnection(context.Background(), c)

		}(conn)
	}
}


func (s *ServerTCP) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	initialBuf := make([]byte, 1024)
	scanner.Buffer(initialBuf, s.maxMessageSize)

	ctxVal := context.WithValue(ctx, compute.ClientIpKey, conn.RemoteAddr().String())

	for {
		conn.SetReadDeadline(time.Now().Add(s.timeout))
		if !scanner.Scan(){
			break
		}

		ctxTime, cancel := context.WithTimeout(ctxVal, time.Second * 2)

		line := scanner.Text()
		ans, err := s.com.HandleQuery(ctxTime, line)

		cancel()

		if err != nil {
			str := "Error: " + err.Error() + "\n"
			conn.Write([]byte(str))
			continue
		}
		conn.Write([]byte(ans + "\n"))
	}

	if err := scanner.Err(); err != nil {
		if errors.Is(err, bufio.ErrTooLong) {
			s.logger.Warn("Client exceeded the limit of Message", 
				zap.String("client", conn.RemoteAddr().String()),
				zap.Int("limit_bytes", s.maxMessageSize))
			conn.Write([]byte("ERROR: Message too large\n"))
		} else {
			s.logger.Error("Error reading conecction", zap.Error(err))
		}
	}
}



func parseMaxSize(sizeConf string) (int, error) {
	var size string
	var bytes string
	sizeConf = strings.TrimSpace(sizeConf)
	for _, val := range sizeConf {
		if unicode.IsDigit(val) {
			size += string(val)
		} else {
			bytes += string(val)
		}
	}

	if size == "" {
		return 0, errors.New("size empty")
	}

	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		return 0, err
	}
	bytes = strings.TrimSpace(strings.ToUpper(bytes))
	switch bytes {
	case "", "Б", "B", "BYTE", "BYTES":
		return sizeInt, nil
	case "КБ", "KB", "K":
		return sizeInt * 1024, nil
	case "МБ", "MB", "M":
		return sizeInt * 1024 * 1024, nil
	case "ГБ", "GB", "G":
		return sizeInt * 1024 * 1024 * 1024, nil
	default:
		return 0, errors.New("Unknown format: " + bytes)
	}
}
