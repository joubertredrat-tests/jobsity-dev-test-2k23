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

type UserLoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type MessageResponse struct {
	ID          string  `json:"id"`
	UserName    string  `json:"userName"`
	UserEmail   string  `json:"userEmail"`
	MessageText string  `json:"messageText"`
	Datetime    *string `json:"datetime"`
}
