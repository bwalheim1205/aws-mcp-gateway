package server

import (
	"net/http"

	"github.com/bwalheim1205/aws-mcp-gateway/internal/handlers"

	"github.com/gin-gonic/gin"
)

type MCPRequest struct {
	Tool string                 `json:"tool"`
	Args map[string]interface{} `json:"args"`
}

// StartServer initializes the Gin HTTP server and routes
func StartServer(port string, lambdaInvoker *handlers.LambdaInvoker) error {
	r := gin.Default()

	r.POST("/mcp", func(c *gin.Context) {
		var req MCPRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := lambdaInvoker.InvokeLambda(c.Request.Context(), req.Tool, req.Args)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	return r.Run(":" + port)
}
