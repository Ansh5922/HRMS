package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func OK(c *gin.Context, data interface{}, msg ...string) {
	message := "success"
	if len(msg) > 0 {
		message = msg[0]
	}
	c.JSON(http.StatusOK, Response{Success: true, Message: message, Data: data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{Success: true, Message: "created", Data: data})
}

func BadRequest(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, Response{Success: false, Error: err})
}

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{Success: false, Error: "unauthorized"})
}

func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{Success: false, Error: "forbidden"})
}

func NotFound(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, Response{Success: false, Error: resource + " not found"})
}

func InternalError(c *gin.Context, err string) {
	c.JSON(http.StatusInternalServerError, Response{Success: false, Error: err})
}
