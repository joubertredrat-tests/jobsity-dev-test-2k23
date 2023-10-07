package application

type UsecaseUserRegisterInput struct {
	Name     string
	Email    string
	Password string
}

type UsecaseUserLoginInput struct {
	Email    string
	Password string
}

type UsecaseMessageCreateInput struct {
	UserName    string
	UserEmail   string
	MessageText string
}

type UsecaseMessagesListInput struct {
	Page         uint
	ItemsPerPage uint
}
