package storage 


import(
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"context"
	"in-memory/internal/storage/engine"
)

func TestStorageLogic(t *testing.T){
	ht := engine.NewHashTable()
	st := NewStorage(ht)
	ctx := context.Background()

	_, err := st.Get(ctx, "unknown")
	require.ErrorIs(t, err, ErrNotExists)

	err = st.Set(ctx, "key1", "value")
	require.NoError(t, err, "write must be without errors")

	val, err := st.Get(ctx, "key1")
	require.NoError(t, err, "get must be without error")
	assert.Equal(t, "value", val, "the value must be equal")

	err  = st.Del(ctx, "key1")
	require.NoError(t, err, "del must be without error")

	_, err = st.Get(ctx, "key1")
	require.ErrorIs(t, err, ErrNotExists)
}