package server

import (
	"bufio"
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"in-memory/internal/compute"
	"in-memory/internal/compute/parser"
	"in-memory/config"
)

type StubParser struct{}
func (s *StubParser) Parse(ctx context.Context, val string) (*parser.Query, error) {
	return &parser.Query{Cmd: parser.CmdGet, Key: "test"}, nil
}

type StubStorage struct{}
func (s *StubStorage) Set(ctx context.Context, key, value string) error { return nil }
func (s *StubStorage) Get(ctx context.Context, key string) (string, error) { return "PONG", nil }
func (s *StubStorage) Del(ctx context.Context, key string) error { return nil }

func SetupTestServer(port string, maxConn int) {
	logger := zap.NewNop()
	
	cnf := &config.ServerConfig{
		MaxMessageSize: "1024B", 
		Address:        port,
		MaxConnections: maxConn,
		IdleTimeout:    2 * time.Second,
	}
	stubParser := &StubParser{}
	stubStorage := &StubStorage{}
	realComputeWithStubs := compute.NewCompute(stubParser, stubStorage, logger, nil)
	server := NewServerTSP(cnf, realComputeWithStubs, logger)
	
	go server.StartServer() 

	time.Sleep(100 * time.Millisecond)
}

func TestServerTCP_HandleQuery(t *testing.T) {
	

	SetupTestServer("127.0.0.1:3224", 10)

	conn, err := net.Dial("tcp", "127.0.0.1:3224")
	require.NoError(t, err)
	defer conn.Close()

	_, err = conn.Write([]byte("PING\n"))
	require.NoError(t, err)

	line, err := bufio.NewReader(conn).ReadString('\n')
	require.NoError(t, err)
	require.Equal(t, "PONG\n", line) 

}

func TestServerTCP_MaxConnections(t *testing.T) {

	maxConn := 3
	
	SetupTestServer("127.0.0.1:3225", maxConn) 

	conns := []net.Conn{}
	
	for i := 0; i < maxConn; i++ {
		conn, err := net.Dial("tcp", "127.0.0.1:3225")
		require.NoError(t, err)
		conns = append(conns, conn)
	}

	extraConn, err := net.Dial("tcp", "127.0.0.1:3225")
	require.NoError(t, err)
	
	answer := make(chan struct{}) 
	
	go func() {
		extraConn.Write([]byte("PING\n"))
		bufio.NewReader(extraConn).ReadString('\n') 
		close(answer)
	}()

	select {
	case <-time.After(500 * time.Millisecond):
	case <-answer:
		t.Fatal("Extra connection was served, but it should have been blocked by maxConnections limit")
	}
	
	conns[0].Close()
	
	select {
	case <-time.After(1 * time.Second):
		t.Fatal("Extra connection was not served even after freeing up a slot")
	case <-answer:
	}
	
	for _, val := range conns[1:] {
		val.Close()
	}
	extraConn.Close()
}