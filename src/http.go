package src

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func GetBody(c *gin.Context) (map[string]interface{}, error) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	defer c.Request.Body.Close()

	var data []map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, err
	}
	return data[0], nil
}
