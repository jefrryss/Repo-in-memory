package compute

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"in-memory/internal/compute/parser"
)

type MockParser struct {
	mock.Mock
}

func (m *MockParser) Parse(ctx context.Context, val string) (*parser.Query, error) {
	args := m.Called(ctx, val)
	if q := args.Get(0); q != nil {
		return q.(*parser.Query), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Set(ctx context.Context, key, value string) error {
	return m.Called(ctx, key, value).Error(0)
}

func (m *MockStorage) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) Del(ctx context.Context, key string) error {
	return m.Called(ctx, key).Error(0)
}

func TestComputeHandleQuery(t *testing.T) {
	logger := zap.NewNop()
	errParse := errors.New("parse error mock")
	errStorage := errors.New("storage error mock")

	testCases := []struct {
		name        string
		queryStr    string
		setupMocks  func(mParser *MockParser, mStorage *MockStorage)
		expectedRes string
		expectedErr error
	}{
		{
			name:     "Parser Error",
			queryStr: "INVALID QUERY",
			setupMocks: func(mp *MockParser, ms *MockStorage) {
				mp.On("Parse", mock.Anything, "INVALID QUERY").Return(nil, errParse).Once()
			},
			expectedRes: "",
			expectedErr: errParse,
		},
		{
			name:     "SET Success",
			queryStr: "SET key1 val1",
			setupMocks: func(mp *MockParser, ms *MockStorage) {
				query := &parser.Query{Cmd: parser.CmdSet, Key: "key1", Value: "val1"}
				mp.On("Parse", mock.Anything, "SET key1 val1").Return(query, nil).Once()
				ms.On("Set", mock.Anything, "key1", "val1").Return(nil).Once()
			},
			expectedRes: "success",
			expectedErr: nil,
		},
		{
			name:     "SET Storage Error",
			queryStr: "SET key1 val1",
			setupMocks: func(mp *MockParser, ms *MockStorage) {
				query := &parser.Query{Cmd: parser.CmdSet, Key: "key1", Value: "val1"}
				mp.On("Parse", mock.Anything, "SET key1 val1").Return(query, nil).Once()
				ms.On("Set", mock.Anything, "key1", "val1").Return(errStorage).Once()
			},
			expectedRes: "",
			expectedErr: errStorage,
		},
		{
			name:     "GET Success",
			queryStr: "GET key1",
			setupMocks: func(mp *MockParser, ms *MockStorage) {
				query := &parser.Query{Cmd: parser.CmdGet, Key: "key1"}
				mp.On("Parse", mock.Anything, "GET key1").Return(query, nil).Once()
				ms.On("Get", mock.Anything, "key1").Return("my_data", nil).Once()
			},
			expectedRes: "my_data",
			expectedErr: nil,
		},
		{
			name:     "GET Storage Error",
			queryStr: "GET unknown",
			setupMocks: func(mp *MockParser, ms *MockStorage) {
				query := &parser.Query{Cmd: parser.CmdGet, Key: "unknown"}
				mp.On("Parse", mock.Anything, "GET unknown").Return(query, nil).Once()
				ms.On("Get", mock.Anything, "unknown").Return("", errStorage).Once()
			},
			expectedRes: "",
			expectedErr: errStorage,
		},
		{
			name:     "DEL Success",
			queryStr: "DEL key1",
			setupMocks: func(mp *MockParser, ms *MockStorage) {
				query := &parser.Query{Cmd: parser.CmdDel, Key: "key1"}
				mp.On("Parse", mock.Anything, "DEL key1").Return(query, nil).Once()
				ms.On("Del", mock.Anything, "key1").Return(nil).Once()
			},
			expectedRes: "success",
			expectedErr: nil,
		},
		{
			name:     "DEL Storage Error",
			queryStr: "DEL key1",
			setupMocks: func(mp *MockParser, ms *MockStorage) {
				query := &parser.Query{Cmd: parser.CmdDel, Key: "key1"}
				mp.On("Parse", mock.Anything, "DEL key1").Return(query, nil).Once()
				ms.On("Del", mock.Anything, "key1").Return(errStorage).Once()
			},
			expectedRes: "",
			expectedErr: errStorage,
		},
		{
			name:     "Unknown Command from Parser",
			queryStr: "MAGIC key1",
			setupMocks: func(mp *MockParser, ms *MockStorage) {
				query := &parser.Query{Cmd: "999", Key: "key1"}
				mp.On("Parse", mock.Anything, "MAGIC key1").Return(query, nil).Once()
			},
			expectedRes: "",
			expectedErr: errors.New("internal error: unknown command"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockParser := new(MockParser)
			mockStorage := new(MockStorage)

			tt.setupMocks(mockParser, mockStorage)

			comp := NewCompute(mockParser, mockStorage, logger)

			ctx := context.WithValue(context.Background(), ClientIpKey, "192.168.1.1")

			res, err := comp.HandleQuery(ctx, tt.queryStr)

			if tt.expectedErr != nil {
				require.Error(t, err, "При ожидании ошибки метод HandleQuery обязан вернуть ошибку")
				assert.Equal(t, tt.expectedErr.Error(), err.Error(), "Текст полученной ошибки должен точно совпадать с ожидаемым")
				assert.Empty(t, res, "При возникновении ошибки возвращаемая строка результата должна быть пустой")
			} else {
				require.NoError(t, err, "При успешном сценарии метод HandleQuery не должен возвращать ошибку")
				assert.Equal(t, tt.expectedRes, res, "Возвращаемый результат должен точно совпадать с данными из хранилища или сообщением 'success'")
			}

			mockParser.AssertExpectations(t)
			mockStorage.AssertExpectations(t)
		})
	}
}