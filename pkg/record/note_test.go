package record

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

	t.Run("error: no invalid method", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.ListNotes(rw, &http.Request{Method: http.MethodPost})

		assert.Equal(t, http.StatusMethodNotAllowed, rw.Code)
	})

	t.Run("successful: no records", func(t *testing.T) {
		mocket.Catcher.Reset().NewMock().WithReply(nil)
		rw := httptest.NewRecorder()
		r.ListNotes(rw, &http.Request{Method: http.MethodGet})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)

		var notes []Note
		err = json.Unmarshal(res, &notes)
		assert.Nil(t, err)

		assert.Equal(t, 0, len(notes))
	})

	t.Run("successful: one record", func(t *testing.T) {
		records := []map[string]interface{}{{"title": "Sample note #123", "content": "A reminder to buy a list of grocery items"}}
		mocket.Catcher.Reset().NewMock().WithReply(records)
		rw := httptest.NewRecorder()
		r.ListNotes(rw, &http.Request{Method: http.MethodGet})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)

		var notes []Note
		err = json.Unmarshal(res, &notes)
		assert.Nil(t, err)

		assert.Equal(t, 1, len(notes))
		assert.Equal(t, "Sample note #123", notes[0].Title)
		assert.Equal(t, "A reminder to buy a list of grocery items", notes[0].Content)
	})

	t.Run("successful: multiple records", func(t *testing.T) {
		records := []map[string]interface{}{{"title": "Sample note #123", "content": "A reminder to buy a list of grocery items"},
			{"title": "Sample note #234", "content": "Notes on how to do something"}}
		mocket.Catcher.Reset().NewMock().WithReply(records)
		rw := httptest.NewRecorder()
		r.ListNotes(rw, &http.Request{Method: http.MethodGet})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)

		var notes []Note
		err = json.Unmarshal(res, &notes)
		assert.Nil(t, err)

		assert.Equal(t, 2, len(notes))
		assert.Equal(t, "Sample note #123", notes[0].Title)
		assert.Equal(t, "A reminder to buy a list of grocery items", notes[0].Content)
		assert.Equal(t, "Sample note #234", notes[1].Title)
		assert.Equal(t, "Notes on how to do something", notes[1].Content)
	})
}

func TestCreateNote(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	t.Run("error: no invalid method", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.CreateNote(rw, &http.Request{Method: http.MethodGet})

		assert.Equal(t, http.StatusMethodNotAllowed, rw.Code)
	})

	t.Run("successful: one record", func(t *testing.T) {
		req := io.NopCloser(strings.NewReader(`{"title": "Sample note #345", "content": "Grocery list"}`))
		rw := httptest.NewRecorder()
		r.CreateNote(rw, &http.Request{
			Method: http.MethodPost,
			Body:   req,
		})

		assert.Equal(t, http.StatusCreated, rw.Code)
	})
}
