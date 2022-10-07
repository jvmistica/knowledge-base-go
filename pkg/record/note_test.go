package record

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	mocket.Catcher.Register()
	db, _ := gorm.Open(postgres.New(postgres.Config{
		DriverName: mocket.DriverName,
		DSN:        "user:test@tcp(127.0.0.1:3306)",
	}), &gorm.Config{})
	return db
}

func TestListNotes(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	t.Run("successful: no records", func(t *testing.T) {
		mocket.Catcher.Reset().NewMock().WithReply(nil)
		rw := httptest.NewRecorder()
		r.ListNotes(rw, &http.Request{})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)
		assert.Equal(t, "", string(res))
	})

	t.Run("successful: one record", func(t *testing.T) {
		records := []map[string]interface{}{{"name": "sample note #123", "description": "noice"}}
		mocket.Catcher.Reset().NewMock().WithReply(records)
		rw := httptest.NewRecorder()
		r.ListNotes(rw, &http.Request{})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)
		assert.Equal(t, "sample note #123\n", string(res))
	})

	t.Run("successful: multiple records", func(t *testing.T) {
		records := []map[string]interface{}{{"name": "sample note #123", "description": "noice"},
			{"name": "sample note #234", "description": "wow"}}
		mocket.Catcher.Reset().NewMock().WithReply(records)
		rw := httptest.NewRecorder()
		r.ListNotes(rw, &http.Request{})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)
		assert.Equal(t, "sample note #123\nsample note #234\n", string(res))
	})
}
