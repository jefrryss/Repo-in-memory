package parser

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParserQuery(t *testing.T) {
	parser := NewLineParser()

	testCases := []struct {
		name        string
		query       string
		expected    *Query
		expectedErr error
	}{
		//БАЗОВЫЕ ПРОВЕРКИ
		{
			name:        "empty query",
			query:       "   ",
			expected:    nil,
			expectedErr: ErrEmptyQuery,
		},
		{
			name:        "unknown command",
			query:       "PING",
			expected:    nil,
			expectedErr: ErrInvalidCommand,
		},

		//КОМАНДА SET
		{
			name:        "valid SET command",
			query:       "SET user:1 Alex",
			expected:    &Query{Cmd: CmdSet, Key: "user:1", Value: "Alex"},
			expectedErr: nil,
		},
		{
			name:        "SET with invalid case",
			query:       "Set user:1 Alex",
			expected:    nil,
			expectedErr: ErrInvalidCommand, 
		},
		{
			name:        "SET with not enough args (1 arg)",
			query:       "SET user:1",
			expected:    nil,
			expectedErr: ErrNotEnoughArgs,
		},
		{
			name:        "SET with too many args",
			query:       "SET user:1 Alex age 25",
			expected:    nil,
			expectedErr: ErrNotEnoughArgs, 
		},

		//КОМАНДА GET
	
		{
			name:        "valid GET command",
			query:       "GET user:1",
			expected:    &Query{Cmd: CmdGet, Key: "user:1"},
			expectedErr: nil,
		},
		{
			name:        "GET with not enough args (0 args)",
			query:       "GET",
			expected:    nil,
			expectedErr: ErrNotEnoughArgs,
		},
		{
			name:        "GET with too many args",
			query:       "GET user:1 extra_arg",
			expected:    nil,
			expectedErr: ErrNotEnoughArgs,
		},

		//КОМАНДА DEL
	
		{
			name:        "valid DEL command",
			query:       "DEL user:1",
			expected:    &Query{Cmd: CmdDel, Key: "user:1"},
			expectedErr: nil,
		},
		{
			name:        "DEL with not enough args (0 args)",
			query:       "DEL",
			expected:    nil,
			expectedErr: ErrNotEnoughArgs,
		},
		{
			name:        "DEL with too many args",
			query:       "DEL user:1 user:2",
			expected:    nil,
			expectedErr: ErrNotEnoughArgs, 
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := parser.Parse(context.Background(), tt.query)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, res, "при ошибке результат должен быть nil")
				return 
			}

			require.NoError(t, err, "ошибка не ожидалась")
			assert.Equal(t, tt.expected, res, "распарсенный запрос не совпадает с ожидаемым")
		})
	}
}


func TestParserWithContext(t *testing.T){
	parser := NewLineParser()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := parser.Parse(ctx, "SET 123 123")
	require.ErrorIs(t, err, context.Canceled)
}