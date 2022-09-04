package url

import (
	"testing"

	"github.com/fakefloordiv/indigo/http"

	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	// here we test laziness of query

	// just test that passed buffer's content will not be used
	query := NewQuery(func() map[string][]byte {
		return make(map[string][]byte)
	})
	query.Set([]byte("hello=world"))
	require.Equal(t, "hello=world", string(query.rawQuery))
	require.Nil(t, query.parsedQuery)

	t.Run("GetExistingKey", func(t *testing.T) {
		value, err := query.Get("hello")
		require.NoError(t, err)
		require.Equal(t, "world", string(value))
	})

	t.Run("GetNonExistingKey", func(t *testing.T) {
		_, err := query.Get("lorem")
		require.ErrorIs(t, err, http.ErrNoSuchKey)
	})
}
