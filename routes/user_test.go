package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mend/db"
)

var (
	router *gin.Engine
	uri    = "https://127.0.0.1:8010/user"
	data   = []map[string]any{
		{"id": 0, "first_name": "matt", "last_name": "Pszczółkowski", "age": 100},
		{"id": 0, "first_name": "Albert", "last_name": "Einstein", "age": 100},
	}
)

func setup() {
	if router == nil {
		db.Use(db.InMemory)

		// sets gin to test mode
		gin.SetMode(gin.TestMode)

		router = gin.Default()
		group := router.Group("/user")
		{
			group.POST("", add())
			group.GET("count", count())
			group.GET(":id", get())
			group.GET("", getAll())
			group.DELETE(":id", remove())
		}
	}
	db.Db().Reset()
}

// addRequest send request to the server
// to add rows from 'data' to database and
// after that check number of rows in database.
// Return: number of added rows
func addRequest(t *testing.T) int64 {
	// save to database all rows of 'data' (global variable).
	for _, user := range data {
		w := httptest.NewRecorder()

		user["id"] = int64(0)
		jsonData, err := json.Marshal(user)
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "https://127.0.0.1:8010/user", bytes.NewBuffer(jsonData))
		assert.NoError(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		if w.Code == http.StatusOK {
			var store map[string]any
			err = json.Unmarshal(w.Body.Bytes(), &store)
			assert.NoError(t, err)

			if count, ok := store["ID of new user"]; ok {
				if value, ok := count.(float64); ok {
					user["id"] = int64(value)
				}
			}
		}
	}

	// ask server for current rows number
	// and check if number of rows in database
	// is equal to number of add operations
	currentCount := countRequest(t)
	assert.Equal(t, int64(len(data)), currentCount)

	return currentCount
}

// countRequest ask server about number of rows in database.
// Return: number of rows or -1 when something was wrong.
func countRequest(t *testing.T) int64 {
	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "https://127.0.0.1:8010/user/count", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	if w.Code == http.StatusOK {
		var store map[string]any
		err = json.Unmarshal(w.Body.Bytes(), &store)
		assert.NoError(t, err)

		if countValue, ok := store["count"]; ok {
			if value, ok := countValue.(float64); ok {
				return int64(value)
			}
		}
	}
	return -1
}

// getWithIDRequest ask server about user with passed ID.
// Return: row data or nil when something was wrong.
func getWithIDRequest(t *testing.T, id int64) map[string]any {
	w := httptest.NewRecorder()

	uri := fmt.Sprintf("https://127.0.0.1:8010/user/%v", id)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var store map[string]any
		err = json.Unmarshal(w.Body.Bytes(), &store)
		assert.NoError(t, err)
		typesUpdate(store)
		return store
	}
	return nil
}

func getAllRequest(t *testing.T) []map[string]any {
	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var store []map[string]any
		err = json.Unmarshal(w.Body.Bytes(), &store)
		assert.NoError(t, err)
		for _, user := range store {
			typesUpdate(user)
		}
		return store
	}
	return nil
}

// delRequest send to server request to remove row with passed ID.
func delRequest(t *testing.T, id int64) {
	w := httptest.NewRecorder()

	uri := fmt.Sprintf("https://127.0.0.1:8010/user/%v", id)
	req, err := http.NewRequest(http.MethodDelete, uri, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// Test_add tests add requests/operations.
func Test_add(t *testing.T) {
	setup()
	addRequest(t)
}

// Test_delete tests del requests/operations.
func Test_delete(t *testing.T) {
	setup()
	counter := addRequest(t)

	for _, user := range data {
		delRequest(t, user["id"].(int64))
		currentCount := countRequest(t)
		assert.Equal(t, counter-1, currentCount)
		counter--
	}
}

// Test_getWithID tests rows reading requests/operations.
func Test_getWithID(t *testing.T) {
	setup()
	addRequest(t)

	for _, user := range data {
		result := getWithIDRequest(t, user["id"].(int64))
		assert.Equal(t, user, result)
	}
}

// Test_getAll tests all rows reading requests/operations.
func Test_getAll(t *testing.T) {
	setup()
	addRequest(t)

	result := getAllRequest(t)
	assert.Equal(t, data, result)
}

// ******************************************************************
// *                                                                *
// *                          H E L P E R S                         *
// *                                                                *
// ******************************************************************

func typesUpdate(store map[string]any) {
	store["id"] = int64(store["id"].(float64))
	store["age"] = int(store["age"].(float64))
}
