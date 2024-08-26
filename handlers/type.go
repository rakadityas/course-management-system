package handlers

// HandlerStatus represents the default response structure for error on handler.
type HandlerStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
