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
)

var (
	testScript = []map[string]interface{}{
		{"name": "Sample script #123", "description": "A bash script that does something"},
		{"name": "Sample script #234", "description": "Handy SQL scripts"},
	}
)

func TestListScripts(t *testing.T) {
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
			dbResult:           []map[string]interface{}{testScript[0]},
			wantErr:            false,
			expectedCount:      1,
			expectedStatusCode: http.StatusOK,
		},
		successMultRecords: {
			method:             http.MethodGet,
			dbResult:           testScript,
			wantErr:            false,
			expectedCount:      2,
			expectedStatusCode: http.StatusOK,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			mocket.Catcher.Reset().NewMock().WithReply(test.dbResult)
			rw := httptest.NewRecorder()

			r.ListScripts(rw, &http.Request{Method: test.method})
			assert.Equal(t, test.expectedStatusCode, rw.Code)

			if !test.wantErr {
				result, err := io.ReadAll(rw.Body)
				assert.Nil(t, err)

				var scripts []Script
				err = json.Unmarshal(result, &scripts)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedCount, len(scripts))

				for i, script := range scripts {
					assert.Equal(t, test.dbResult[i]["name"], script.Name)
					assert.Equal(t, test.dbResult[i]["description"], script.Description)
				}
			}
		})
	}
}

func TestCreateScript(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	t.Run(errInvalidMethod, func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.CreateScript(rw, &http.Request{Method: http.MethodGet})

		assert.Equal(t, http.StatusMethodNotAllowed, rw.Code)
	})

	t.Run(successOneRecord, func(t *testing.T) {
		req := io.NopCloser(strings.NewReader(`{"name": "Sample script #345", "description": "Automation script"}`))
		rw := httptest.NewRecorder()
		r.CreateScript(rw, &http.Request{
			Method: http.MethodPost,
			Body:   req,
		})

		assert.Equal(t, http.StatusCreated, rw.Code)
	})
}

func TestDeleteScript(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	t.Run(errInvalidMethod, func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.DeleteScript(rw, &http.Request{Method: http.MethodPost})

		assert.Equal(t, http.StatusMethodNotAllowed, rw.Code)
	})

	t.Run(errMissingParam, func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.DeleteScript(rw, &http.Request{
			Method: http.MethodDelete,
			URL:    &url.URL{},
		})

		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	t.Run(errRecordNotFound, func(t *testing.T) {
		rw := httptest.NewRecorder()
		mocket.Catcher.Reset().NewMock().WithRowsNum(0)
		r.DeleteScript(rw, &http.Request{
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
		r.DeleteScript(rw, &http.Request{
			Method: http.MethodDelete,
			URL: &url.URL{
				RawQuery: "id=23",
			},
		})

		assert.Equal(t, http.StatusOK, rw.Code)
	})
}

func TestGetScript(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	t.Run(errInvalidMethod, func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.GetScript(rw, &http.Request{Method: http.MethodPost})

		assert.Equal(t, http.StatusMethodNotAllowed, rw.Code)
	})

	t.Run(errMissingParam, func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.GetScript(rw, &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{},
		})

		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	t.Run(errRecordNotFound, func(t *testing.T) {
		rw := httptest.NewRecorder()
		mocket.Catcher.Reset().NewMock().WithRowsNum(0)
		r.GetScript(rw, &http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				RawQuery: "id=99",
			},
		})

		assert.Equal(t, http.StatusNotFound, rw.Code)
	})

	t.Run(successRecordFound, func(t *testing.T) {
		rw := httptest.NewRecorder()
		records := []map[string]interface{}{testScript[0]}
		mocket.Catcher.Reset().NewMock().WithReply(records)
		r.GetScript(rw, &http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				RawQuery: "id=123",
			},
		})
		assert.Equal(t, http.StatusOK, rw.Code)

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)

		var script Script
		err = json.Unmarshal(res, &script)
		assert.Nil(t, err)

		assert.Equal(t, "Sample script #123", script.Name)
		assert.Equal(t, "A bash script that does something", script.Description)
	})
}
