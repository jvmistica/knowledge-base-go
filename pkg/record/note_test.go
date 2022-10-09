package record

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	errInvalidMethod  = "error: invalid method"
	errMissingParam   = "error: missing parameter"
	errRecordNotFound = "error: record not found"

	successNoRecord      = "successful: no records"
	successOneRecord     = "successful: one record"
	successMultRecords   = "successful: multiple records"
	successRecordFound   = "successful: record found"
	successRecordDeleted = "successful: record deleted"
)

var (
	testNote = []map[string]interface{}{
		{"title": "Sample note #123", "content": "A reminder to buy a list of grocery items"},
		{"title": "Sample note #234", "content": "Note on how to do something"},
	}
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

	tests := map[string]struct {
		method             string
		dbResult           []map[string]interface{}
		wantErr            bool
		expectedCount      int
		expectedStatusCode int
	}{
		errInvalidMethod: {
			method:             http.MethodPost,
			dbResult:           nil,
			wantErr:            true,
			expectedCount:      0,
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		successNoRecord: {
			method:             http.MethodGet,
			dbResult:           nil,
			wantErr:            false,
			expectedCount:      0,
			expectedStatusCode: http.StatusOK,
		},
		successOneRecord: {
			method:             http.MethodGet,
			dbResult:           []map[string]interface{}{testNote[0]},
			wantErr:            false,
			expectedCount:      1,
			expectedStatusCode: http.StatusOK,
		},
		successMultRecords: {
			method:             http.MethodGet,
			dbResult:           testNote,
			wantErr:            false,
			expectedCount:      2,
			expectedStatusCode: http.StatusOK,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			mocket.Catcher.Reset().NewMock().WithReply(test.dbResult)
			rw := httptest.NewRecorder()

			r.ListNotes(rw, &http.Request{Method: test.method})
			assert.Equal(t, test.expectedStatusCode, rw.Code)

			if !test.wantErr {
				result, err := io.ReadAll(rw.Body)
				assert.Nil(t, err)

				var notes []Note
				err = json.Unmarshal(result, &notes)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedCount, len(notes))

				for i, note := range notes {
					assert.Equal(t, test.dbResult[i]["title"], note.Title)
					assert.Equal(t, test.dbResult[i]["content"], note.Content)
				}
			}
		})
	}
}

func TestCreateNote(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	t.Run(errInvalidMethod, func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.CreateNote(rw, &http.Request{Method: http.MethodGet})

		assert.Equal(t, http.StatusMethodNotAllowed, rw.Code)
	})

	t.Run(successOneRecord, func(t *testing.T) {
		req := io.NopCloser(strings.NewReader(`{"title": "Sample note #345", "content": "Grocery list"}`))
		rw := httptest.NewRecorder()
		r.CreateNote(rw, &http.Request{
			Method: http.MethodPost,
			Body:   req,
		})

		assert.Equal(t, http.StatusCreated, rw.Code)
	})
}

func TestDeleteNote(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	t.Run(errInvalidMethod, func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.DeleteNote(rw, &http.Request{Method: http.MethodPost})

		assert.Equal(t, http.StatusMethodNotAllowed, rw.Code)
	})

	t.Run(errMissingParam, func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.DeleteNote(rw, &http.Request{
			Method: http.MethodDelete,
			URL:    &url.URL{},
		})

		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	t.Run(errRecordNotFound, func(t *testing.T) {
		rw := httptest.NewRecorder()
		mocket.Catcher.Reset().NewMock().WithRowsNum(0)
		r.DeleteNote(rw, &http.Request{
			Method: http.MethodDelete,
			URL: &url.URL{
				RawQuery: "id=99",
			},
		})

		assert.Equal(t, http.StatusNotFound, rw.Code)
	})

	t.Run(successRecordDeleted, func(t *testing.T) {
		rw := httptest.NewRecorder()
		mocket.Catcher.Reset().NewMock().WithRowsNum(1)
		r.DeleteNote(rw, &http.Request{
			Method: http.MethodDelete,
			URL: &url.URL{
				RawQuery: "id=23",
			},
		})

		assert.Equal(t, http.StatusOK, rw.Code)
	})
}

func TestGetNote(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	t.Run(errInvalidMethod, func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.GetNote(rw, &http.Request{Method: http.MethodPost})

		assert.Equal(t, http.StatusMethodNotAllowed, rw.Code)
	})

	t.Run(errMissingParam, func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.GetNote(rw, &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{},
		})

		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	t.Run(errRecordNotFound, func(t *testing.T) {
		rw := httptest.NewRecorder()
		mocket.Catcher.Reset().NewMock().WithRowsNum(0)
		r.GetNote(rw, &http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				RawQuery: "id=99",
			},
		})

		assert.Equal(t, http.StatusNotFound, rw.Code)
	})

	t.Run(successRecordFound, func(t *testing.T) {
		rw := httptest.NewRecorder()
		records := []map[string]interface{}{testNote[0]}
		mocket.Catcher.Reset().NewMock().WithReply(records)
		r.GetNote(rw, &http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				RawQuery: "id=123",
			},
		})
		assert.Equal(t, http.StatusOK, rw.Code)

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)

		var note Note
		err = json.Unmarshal(res, &note)
		assert.Nil(t, err)

		assert.Equal(t, "Sample note #123", note.Title)
		assert.Equal(t, "A reminder to buy a list of grocery items", note.Content)
	})
}
