package handlers

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse 错误响应
func ErrorResponse(message string, err error) Response {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return Response{
		Code:    400,
		Message: message,
		Error:   errMsg,
	}
}

// PageResponse 分页响应
func PageResponse(message string, items interface{}, total int64) Response {
	return Response{
		Code:    200,
		Message: message,
		Data: map[string]interface{}{
			"items": items,
			"total": total,
		},
	}
}
