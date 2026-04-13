package storage 

import(
    "context"
    "testing"

    "github.com/stretchr/testify/mock" 
    "github.com/stretchr/testify/require"
)

type MockEngine struct{
    mock.Mock
}

func (m *MockEngine) Set(key, value string) {
    m.Called(key, value)
}

func (m *MockEngine) Get(key string) (string, bool) {
    args := m.Called(key)
    return args.String(0), args.Bool(1)
}

func (m *MockEngine) Del(key string) {
    m.Called(key)
}

func TestStorageLogic(t *testing.T){
    mockEn := new(MockEngine)
    storage := NewStorage(mockEn)
    ctx := context.Background()

    mockEn.On("Get", "unknown").Return("", false).Once()
    _, err := storage.Get(ctx, "unknown")
    require.ErrorIs(t, err, ErrNotExists, "Ожидалась ошибка ErrNotExists при запросе несуществующего ключа")

    mockEn.On("Set", "key1", "value").Return().Once()
    err = storage.Set(ctx, "key1", "value") 
    require.NoError(t, err, "Метод Set не должен возвращать ошибку при первичном добавлении валидных данных")

    mockEn.On("Get", "key1").Return("value", true).Once()
    value, err := storage.Get(ctx, "key1")
    require.NoError(t, err, "Метод Get не должен возвращать ошибку при запросе существующего ключа")
    require.Equal(t, "value", value, "Полученное значение должно точно совпадать с тем, которое было записано ранее")

    mockEn.On("Del", "key1").Return().Once()
    err = storage.Del(ctx, "key1")
    require.NoError(t, err, "Метод Del не должен возвращать ошибку при удалении существующего ключа")

    mockEn.AssertExpectations(t)
}

func TestStorageSetLogic(t *testing.T) {
    mockDb := new(MockEngine) 
    storage := NewStorage(mockDb)
    ctx := context.Background() 

    mockDb.On("Set", "key1", "value1").Return().Once()
    err := storage.Set(ctx, "key1", "value1")
    require.NoError(t, err, "При первичной записи по ключу не должно возникать ошибок")

    mockDb.On("Get", "key1").Return("value1", true).Once()
    value, err := storage.Get(ctx, "key1")
    require.NoError(t, err, "При получении только что записанных данных не должно быть ошибок")
    require.Equal(t, "value1", value, "Прочитанные данные должны быть полностью эквивалентны записанным")

    mockDb.On("Set", "key1", "value2").Return().Once()
    err = storage.Set(ctx, "key1", "value2")
    require.NoError(t, err, "При перезаписи существующего ключа новым значением не должно возникать ошибок")

    mockDb.On("Get", "key1").Return("value2", true).Once()
    value, err = storage.Get(ctx, "key1")
    require.NoError(t, err, "При получении обновленных данных не должно быть ошибок")
    require.Equal(t, "value2", value, "Значение по ключу должно быть успешно перезаписано на новое")

    mockDb.AssertExpectations(t)
}

func TestStorageLogicWithContext(t *testing.T) {
    mockDb := new(MockEngine)
    storage := NewStorage(mockDb)

    ctx, cancel := context.WithCancel(context.Background())
    cancel()
    err := storage.Set(ctx, "key1", "valu1")
    require.ErrorIs(t, err, context.Canceled, "Метод Set должен немедленно возвращать ошибку context.Canceled, если контекст был отменен до вызова")

    ctx, cancel = context.WithCancel(context.Background())
    cancel()
    _, err = storage.Get(ctx, "key1")
    require.ErrorIs(t, err, context.Canceled, "Метод Get должен немедленно возвращать ошибку context.Canceled, если контекст был отменен до вызова")

    ctx, cancel = context.WithCancel(context.Background())
    mockDb.On("Set", "key2", "value2").Return().Once()
    err = storage.Set(ctx, "key2", "value2")
    require.NoError(t, err, "При записи с активным контекстом метод Set должен отрабатывать без ошибок")
    
    cancel()
    _, err = storage.Get(ctx, "key2")
    require.ErrorIs(t, err, context.Canceled, "Метод Get должен прервать работу и вернуть ошибку context.Canceled, так как контекст был отменен перед чтением")

    ctx, cancel = context.WithCancel(context.Background())
    cancel()
    err = storage.Del(ctx, "unknown")
    require.ErrorIs(t, err, context.Canceled, "Метод Del должен немедленно возвращать ошибку context.Canceled, если контекст был отменен до вызова")

    mockDb.AssertExpectations(t)
}