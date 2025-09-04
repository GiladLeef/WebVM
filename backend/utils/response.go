package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type paginatedData struct {
	Items []interface{} `json:"items"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
	Pages int           `json:"pages"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	resp := response{
		Success: true,
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, resp)
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	resp := response{
		Success: false,
		Error:   message,
	}
	c.JSON(statusCode, resp)
}

func ValidationErrorResponse(c *gin.Context, err error) {
	var message string
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		message = "Validation failed: " + validationErrors[0].Tag()
	} else {
		message = "Invalid request data"
	}
	ErrorResponse(c, http.StatusBadRequest, message)
}

func UnauthorizedResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	ErrorResponse(c, http.StatusUnauthorized, message)
}

func NotFoundResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	ErrorResponse(c, http.StatusNotFound, message)
}

func ServerErrorResponse(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
}

func PaginatedResponse(c *gin.Context, items interface{}, total int64, page, limit int) {
	pages := int((total + int64(limit) - 1) / int64(limit))

	var itemsSlice []interface{}
	if slice, ok := items.([]interface{}); ok {
		itemsSlice = slice
	} else {
		itemsSlice = []interface{}{items}
	}

	data := paginatedData{
		Items: itemsSlice,
		Total: total,
		Page:  page,
		Limit: limit,
		Pages: pages,
	}

	SuccessResponse(c, http.StatusOK, "Retrieved successfully", data)
}


