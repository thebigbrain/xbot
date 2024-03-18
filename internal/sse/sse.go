package sse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendSse(c *gin.Context, msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	c.Writer.WriteString(fmt.Sprintf("data: %s\n\n", data))
	// c.Writer.WriteString(fmt.Sprintf("%s", data))
	c.Writer.Flush()
}
