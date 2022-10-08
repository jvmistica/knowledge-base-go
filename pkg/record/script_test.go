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

func TestListScripts(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	t.Run("error: invalid method", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.ListScripts(rw, &http.Request{Method: http.MethodPost})

		assert.Equal(t, http.StatusMethodNotAllowed, rw.Code)
	})

	t.Run("successful: no records", func(t *testing.T) {
		mocket.Catcher.Reset().NewMock().WithReply(nil)
		rw := httptest.NewRecorder()
		r.ListScripts(rw, &http.Request{Method: http.MethodGet})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)

		var scripts []Script
		err = json.Unmarshal(res, &scripts)
		assert.Nil(t, err)

		assert.Equal(t, 0, len(scripts))
	})

	t.Run("successful: one record", func(t *testing.T) {
		records := []map[string]interface{}{{"name": "Sample script #123", "description": "A bash script that does something"}}
		mocket.Catcher.Reset().NewMock().WithReply(records)
		rw := httptest.NewRecorder()
		r.ListScripts(rw, &http.Request{Method: http.MethodGet})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)

		var scripts []Script
		err = json.Unmarshal(res, &scripts)
		assert.Nil(t, err)

		assert.Equal(t, 1, len(scripts))
		assert.Equal(t, "Sample script #123", scripts[0].Name)
		assert.Equal(t, "A bash script that does something", scripts[0].Description)
	})

	t.Run("successful: multiple records", func(t *testing.T) {
		records := []map[string]interface{}{{"name": "Sample script #123", "description": "A bash script that does something"},
			{"name": "Sample script #234", "description": "Handy SQL scripts"}}
		mocket.Catcher.Reset().NewMock().WithReply(records)
		rw := httptest.NewRecorder()
		r.ListScripts(rw, &http.Request{Method: http.MethodGet})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)

		var scripts []Script
		err = json.Unmarshal(res, &scripts)
		assert.Nil(t, err)

		assert.Equal(t, 2, len(scripts))
		assert.Equal(t, "Sample script #123", scripts[0].Name)
		assert.Equal(t, "A bash script that does something", scripts[0].Description)
		assert.Equal(t, "Sample script #234", scripts[1].Name)
		assert.Equal(t, "Handy SQL scripts", scripts[1].Description)
	})
}

func TestCreateScript(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	t.Run("error: invalid method", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.CreateScript(rw, &http.Request{Method: http.MethodGet})

		assert.Equal(t, http.StatusMethodNotAllowed, rw.Code)
	})

	t.Run("successful: one record", func(t *testing.T) {
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

	t.Run("error: invalid method", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.DeleteScript(rw, &http.Request{Method: http.MethodPost})

		assert.Equal(t, http.StatusMethodNotAllowed, rw.Code)
	})

	t.Run("error: missing parameter", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.DeleteScript(rw, &http.Request{
			Method: http.MethodDelete,
			URL:    &url.URL{},
		})

		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	t.Run("error: record not found", func(t *testing.T) {
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

	t.Run("successful: script deleted", func(t *testing.T) {
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

	t.Run("error: invalid method", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.GetScript(rw, &http.Request{Method: http.MethodPost})

		assert.Equal(t, http.StatusMethodNotAllowed, rw.Code)
	})

	t.Run("error: missing parameter", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.GetScript(rw, &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{},
		})

		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	t.Run("error: record not found", func(t *testing.T) {
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

	t.Run("successful: record found", func(t *testing.T) {
		rw := httptest.NewRecorder()
		records := []map[string]interface{}{{"name": "Sample script #123", "description": "A bash script that does something"}}
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
