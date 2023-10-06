package infra

type RequestValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

type UserRegisterResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
