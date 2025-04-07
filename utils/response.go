package utils

type Response struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// NewSuccessResponse başarılı bir response döndürür
func NewSuccessResponse(status int, data interface{}, message string) Response {
	return Response{
		Status:  status,
		Data:    data,
		Message: message,
	}
}

// NewErrorResponse hata response'u döndürür
func NewErrorResponse(status int, error string, message string) Response {
	return Response{
		Status:  status,
		Error:   error,
		Message: message,
	}
}
