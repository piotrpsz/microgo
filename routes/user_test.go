package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mend/db"
)

var (
	data = []map[string]any{
		{"id": 0, "first_name": "matt", "last_name": "Pszczółkowski", "age": 100},
		{"id": 0, "first_name": "Albert", "last_name": "Einstein", "age": 100},
	}
)

var router *gin.Engine

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
			group.DELETE(":id", remove())
		}
	}
}

// addTest adds rows from 'data' and checks number of rows
// in database.
func addRequest(t *testing.T, router *gin.Engine) int64 {
	w := httptest.NewRecorder()

	for _, user := range data {
		jsonData, err := json.Marshal(user)
		assert.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "https://127.0.0.1:8010/user", bytes.NewBuffer(jsonData))
		assert.NoError(t, err)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	shouldCount := int64(len(data))
	currentCount := countRequest(t, router)
	assert.Equal(t, shouldCount, currentCount)

	return shouldCount
}

func delRequest(t *testing.T, router *gin.Engine) {
	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodDelete, "https://127.0.0.1:8010/user/1", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func countRequest(t *testing.T, router *gin.Engine) int64 {
	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "https://127.0.0.1:8010/user/count", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	if w.Code == http.StatusOK {
		var store map[string]any
		err = json.Unmarshal(w.Body.Bytes(), &store)
		assert.NoError(t, err)

		if count, ok := store["count"]; ok {
			if value, ok := count.(float64); ok {
				return int64(value)
			}
		}
	}
	return -1
}

func Test_add(t *testing.T) {
	setup()
	addRequest(t, router)
}

func Test_delelete(t *testing.T) {
	setup()
	addedCount := addRequest(t, router)
	delRequest(t, router)
	currentCount := countRequest(t, router)
	assert.Equal(t, addedCount-1, currentCount)
}
