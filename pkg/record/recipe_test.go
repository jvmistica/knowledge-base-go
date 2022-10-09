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

func TestListRecipes(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	tests := map[string]struct {
		method             string
		dbResult           []map[string]interface{}
		wantErr            bool
		expectedCount      int
		expectedStatusCode int
	}{
		"error: invalid method": {
			method:             http.MethodPost,
			dbResult:           nil,
			wantErr:            true,
			expectedCount:      0,
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		"successful: no records": {
			method:             http.MethodGet,
			dbResult:           nil,
			wantErr:            false,
			expectedCount:      0,
			expectedStatusCode: http.StatusOK,
		},
		"successful: one record": {
			method:             http.MethodGet,
			dbResult:           []map[string]interface{}{{"name": "Sample recipe #123", "description": "A very delicious dish"}},
			wantErr:            false,
			expectedCount:      1,
			expectedStatusCode: http.StatusOK,
		},
		"successful: multiple records": {
			method: http.MethodGet,
			dbResult: []map[string]interface{}{{"name": "Sample recipe #123", "description": "A very delicious dish"},
				{"name": "Sample recipe #234", "description": "An exotic dish"}},
			wantErr:            false,
			expectedCount:      2,
			expectedStatusCode: http.StatusOK,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			mocket.Catcher.Reset().NewMock().WithReply(test.dbResult)
			rw := httptest.NewRecorder()

			r.ListRecipes(rw, &http.Request{Method: test.method})
			assert.Equal(t, test.expectedStatusCode, rw.Code)

			if !test.wantErr {
				result, err := io.ReadAll(rw.Body)
				assert.Nil(t, err)

				var recipes []Recipe
				err = json.Unmarshal(result, &recipes)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedCount, len(recipes))

				for i, recipe := range recipes {
					assert.Equal(t, test.dbResult[i]["name"], recipe.Name)
					assert.Equal(t, test.dbResult[i]["description"], recipe.Description)
				}
			}
		})
	}

}

func TestCreateRecipe(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	t.Run("error: invalid method", func(t *testing.T) {
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

func TestDeleteRecipe(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	t.Run("error: invalid method", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.DeleteRecipe(rw, &http.Request{Method: http.MethodPost})

		assert.Equal(t, http.StatusMethodNotAllowed, rw.Code)
	})

	t.Run("error: missing parameter", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.DeleteRecipe(rw, &http.Request{
			Method: http.MethodDelete,
			URL:    &url.URL{},
		})

		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	t.Run("error: record not found", func(t *testing.T) {
		rw := httptest.NewRecorder()
		mocket.Catcher.Reset().NewMock().WithRowsNum(0)
		r.DeleteRecipe(rw, &http.Request{
			Method: http.MethodDelete,
			URL: &url.URL{
				RawQuery: "id=99",
			},
		})

		assert.Equal(t, http.StatusNotFound, rw.Code)
	})

	t.Run("successful: recipe deleted", func(t *testing.T) {
		rw := httptest.NewRecorder()
		mocket.Catcher.Reset().NewMock().WithRowsNum(1)
		r.DeleteRecipe(rw, &http.Request{
			Method: http.MethodDelete,
			URL: &url.URL{
				RawQuery: "id=23",
			},
		})

		assert.Equal(t, http.StatusOK, rw.Code)
	})
}

func TestGetRecipe(t *testing.T) {
	db := setupTestDB()
	r := &Record{DB: db}

	t.Run("error: invalid method", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.GetRecipe(rw, &http.Request{Method: http.MethodPost})

		assert.Equal(t, http.StatusMethodNotAllowed, rw.Code)
	})

	t.Run("error: missing parameter", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r.GetRecipe(rw, &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{},
		})

		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	t.Run("error: record not found", func(t *testing.T) {
		rw := httptest.NewRecorder()
		mocket.Catcher.Reset().NewMock().WithRowsNum(0)
		r.GetRecipe(rw, &http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				RawQuery: "id=99",
			},
		})

		assert.Equal(t, http.StatusNotFound, rw.Code)
	})

	t.Run("successful: record found", func(t *testing.T) {
		rw := httptest.NewRecorder()
		records := []map[string]interface{}{{"name": "Sample recipe #123", "description": "A very delicious dish"}}
		mocket.Catcher.Reset().NewMock().WithReply(records)
		r.GetRecipe(rw, &http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				RawQuery: "id=123",
			},
		})
		assert.Equal(t, http.StatusOK, rw.Code)

		res, err := io.ReadAll(rw.Body)
		assert.Nil(t, err)

		var recipe Recipe
		err = json.Unmarshal(res, &recipe)
		assert.Nil(t, err)

		assert.Equal(t, "Sample recipe #123", recipe.Name)
		assert.Equal(t, "A very delicious dish", recipe.Description)
	})
}
