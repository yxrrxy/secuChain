package error0

import (
	"github.com/cloudwego/hertz/pkg/protocol"
)

func ErrorResponse(message string, err error) protocol.H {
	return protocol.H{
		"success": false,
		"message": message,
		"error":   err.Error(),
	}
}

func SuccessResponse(message string, data interface{}) protocol.H {
	return protocol.H{
		"success": true,
		"message": message,
		"data":    data,
	}
}

func PageResponse(message string, data interface{}, total int64) protocol.H {
	return protocol.H{
		"success": true,
		"message": message,
		"data":    data,
		"total":   total,
	}
}
