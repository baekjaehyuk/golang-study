package api

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *apiError   `json:"error,omitempty"`
}

type apiError struct {
	Message string `json:"message"`
}

func successResponse(data interface{}) response {
	return response{Success: true, Data: data}
}

func errorResponse(message string) response {
	return response{
		Success: false,
		Error: &apiError{
			Message: message,
		},
	}
}
