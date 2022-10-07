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
		r.ListRecipes(rw, &http.Request{})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)
		assert.Equal(t, "", string(res))
	})

	t.Run("successful: one record", func(t *testing.T) {
		records := []map[string]interface{}{{"name": "sample recipe #123", "description": "noice"}}
		mocket.Catcher.Reset().NewMock().WithReply(records)
		rw := httptest.NewRecorder()
		r.ListNotes(rw, &http.Request{})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)
		assert.Equal(t, "sample recipe #123\n", string(res))
	})

	t.Run("successful: multiple records", func(t *testing.T) {
		records := []map[string]interface{}{{"name": "sample recipe #123", "description": "noice"},
			{"name": "sample recipe #234", "description": "wow"}}
		mocket.Catcher.Reset().NewMock().WithReply(records)
		rw := httptest.NewRecorder()
		r.ListNotes(rw, &http.Request{})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)
		assert.Equal(t, "sample recipe #123\nsample recipe #234\n", string(res))
	})
}
