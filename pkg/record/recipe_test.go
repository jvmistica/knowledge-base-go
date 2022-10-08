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
)

func TestListRecipes(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	t.Run("error: no invalid method", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.ListRecipes(rw, &http.Request{Method: http.MethodPost})

		assert.Equal(t, http.StatusMethodNotAllowed, rw.Code)
	})

	t.Run("successful: no records", func(t *testing.T) {
		mocket.Catcher.Reset().NewMock().WithReply(nil)
		rw := httptest.NewRecorder()
		r.ListRecipes(rw, &http.Request{Method: http.MethodGet})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)

		var recipes []Recipe
		err = json.Unmarshal(res, &recipes)
		assert.Nil(t, err)

		assert.Equal(t, 0, len(recipes))
	})

	t.Run("successful: one record", func(t *testing.T) {
		records := []map[string]interface{}{{"name": "Sample recipe #123", "description": "A very delicious dish"}}
		mocket.Catcher.Reset().NewMock().WithReply(records)
		rw := httptest.NewRecorder()
		r.ListRecipes(rw, &http.Request{Method: http.MethodGet})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)

		var recipes []Recipe
		err = json.Unmarshal(res, &recipes)
		assert.Nil(t, err)

		assert.Equal(t, 1, len(recipes))
		assert.Equal(t, "Sample recipe #123", recipes[0].Name)
		assert.Equal(t, "A very delicious dish", recipes[0].Description)
	})

	t.Run("successful: multiple records", func(t *testing.T) {
		records := []map[string]interface{}{{"name": "Sample recipe #123", "description": "A very delicious dish"},
			{"name": "Sample recipe #234", "description": "An exotic dish"}}
		mocket.Catcher.Reset().NewMock().WithReply(records)
		rw := httptest.NewRecorder()
		r.ListRecipes(rw, &http.Request{Method: http.MethodGet})

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)

		var recipes []Recipe
		err = json.Unmarshal(res, &recipes)
		assert.Nil(t, err)

		assert.Equal(t, 2, len(recipes))
		assert.Equal(t, "Sample recipe #123", recipes[0].Name)
		assert.Equal(t, "A very delicious dish", recipes[0].Description)
		assert.Equal(t, "Sample recipe #234", recipes[1].Name)
		assert.Equal(t, "An exotic dish", recipes[1].Description)
	})
}

func TestCreateRecipe(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	t.Run("error: no invalid method", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.CreateRecipe(rw, &http.Request{Method: http.MethodGet})

		assert.Equal(t, http.StatusMethodNotAllowed, rw.Code)
	})

	t.Run("successful: one record", func(t *testing.T) {
		req := io.NopCloser(strings.NewReader(`{"name": "Sample recipe #345", "description": "A soup dish"}`))
		rw := httptest.NewRecorder()
		r.CreateRecipe(rw, &http.Request{
			Method: http.MethodPost,
			Body:   req,
		})

		assert.Equal(t, http.StatusCreated, rw.Code)
	})
}
