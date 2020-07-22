package http

type HTTPError struct {
	//ErrorCode int
	Error ErrorMessage `json:"error"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}
