package model

type SystemError struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

func (s SystemError) Error() string {
	return s.Message
}
