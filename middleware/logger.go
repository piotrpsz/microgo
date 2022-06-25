package middleware

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// LogToFile middleware to save infos to a file
// Tip: go to proper directory and: sudo cat mend.cfg | jq
func LogToFile(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()

		data := make(map[string]any)
		data["status"] = c.Writer.Status()
		data["duration"] = end.Sub(start)
		data["ip"] = c.ClientIP()
		data["method"] = c.Request.Method
		data["uri"] = c.Request.RequestURI

		bytes, _ := json.Marshal(data)
		logger.Info(string(bytes))
	}
}
