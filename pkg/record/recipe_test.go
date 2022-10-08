package record

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
)

func TestListRecipes(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	t.Run("successful: no records", func(t *testing.T) {
		mocket.Catcher.Reset().NewMock().WithReply(nil)
		rw := httptest.NewRecorder()
		r.ListRecipes(rw, &http.Request{Method: http.MethodGet})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)
		assert.Equal(t, "", string(res))
	})

	t.Run("successful: one record", func(t *testing.T) {
		records := []map[string]interface{}{{"name": "Sample recipe #123", "description": "A very delicious dish"}}
		mocket.Catcher.Reset().NewMock().WithReply(records)
		rw := httptest.NewRecorder()
		r.ListRecipes(rw, &http.Request{Method: http.MethodGet})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)
		assert.Equal(t, "<b>Sample recipe #123</b></br>A very delicious dish</br></br>", string(res))
	})

	t.Run("successful: multiple records", func(t *testing.T) {
		records := []map[string]interface{}{{"name": "Sample recipe #123", "description": "A very delicious dish"},
			{"name": "Sample recipe #234", "description": "An exotic dish"}}
		mocket.Catcher.Reset().NewMock().WithReply(records)
		rw := httptest.NewRecorder()
		r.ListRecipes(rw, &http.Request{Method: http.MethodGet})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)
		assert.Equal(t, "<b>Sample recipe #123</b></br>A very delicious dish</br></br><b>Sample recipe #234</b></br>An exotic dish</br></br>", string(res))
	})
}
