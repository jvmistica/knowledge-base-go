package record

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
)

func TestListScripts(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	t.Run("successful: no records", func(t *testing.T) {
		mocket.Catcher.Reset().NewMock().WithReply(nil)
		rw := httptest.NewRecorder()
		r.ListScripts(rw, &http.Request{})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)
		assert.Equal(t, "", string(res))
	})

	t.Run("successful: one record", func(t *testing.T) {
		records := []map[string]interface{}{{"name": "Sample script #123", "description": "A bash script that does something"}}
		mocket.Catcher.Reset().NewMock().WithReply(records)
		rw := httptest.NewRecorder()
		r.ListScripts(rw, &http.Request{})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)
		assert.Equal(t, "<b>Sample script #123</b></br>A bash script that does something</br></br>", string(res))
	})

	t.Run("successful: multiple records", func(t *testing.T) {
		records := []map[string]interface{}{{"name": "Sample script #123", "description": "A bash script that does something"},
			{"name": "Sample script #234", "description": "Handy SQL scripts"}}
		mocket.Catcher.Reset().NewMock().WithReply(records)
		rw := httptest.NewRecorder()
		r.ListScripts(rw, &http.Request{})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)
		assert.Equal(t, "<b>Sample script #123</b></br>A bash script that does something</br></br><b>Sample script #234</b></br>Handy SQL scripts</br></br>", string(res))
	})
}
