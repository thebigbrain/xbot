package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func sendSse(c *gin.Context, msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteString(fmt.Sprintf("data: %s\n\n", data))
	// c.Writer.WriteString(fmt.Sprintf("%s", data))
	c.Writer.Flush()
}
