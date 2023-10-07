package application

import (
	"context"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"
)

type UsecaseUserRegister struct {
	userRepository domain.UserRepository
}

func NewUsecaseUserRegister(r domain.UserRepository) UsecaseUserRegister {
	return UsecaseUserRegister{
		userRepository: r,
	}
}

func (u UsecaseUserRegister) Execute(ctx context.Context, input UsecaseUserRegisterInput) (domain.User, error) {
	user, err := domain.NewUser("", input.Name, input.Email, input.Password)
	if err != nil {
		return domain.User{}, err
	}

	userGot, err := u.userRepository.GetByEmail(ctx, user.Email)
	if err != nil {
		return domain.User{}, err
	}
	if userGot.Email == user.Email {
		return domain.User{}, NewErrUserAlreadyRegistered(user.Email)
	}

	return u.userRepository.Create(ctx, user)
}

type UsecaseUserLogin struct {
	userRepository domain.UserRepository
	tokenService   domain.TokenService
}

func NewUsecaseUserLogin(r domain.UserRepository, t domain.TokenService) UsecaseUserLogin {
	return UsecaseUserLogin{
		userRepository: r,
		tokenService:   t,
	}
}

func (u UsecaseUserLogin) Execute(ctx context.Context, input UsecaseUserLoginInput) (domain.UserToken, error) {
	user, err := domain.NewUser("", "", input.Email, input.Password)
	if err != nil {
		return domain.UserToken{}, err
	}
	userGot, err := u.userRepository.GetAuthenticated(ctx, user)
	if err != nil {
		return domain.UserToken{}, err
	}

	return u.tokenService.Generate(ctx, userGot)
}

type UsecaseMessageCreate struct {
	messageRepository domain.MessageRepository
	messageEvent      domain.MessageEvent
}

func NewUsecaseMessageCreate(r domain.MessageRepository, e domain.MessageEvent) UsecaseMessageCreate {
	return UsecaseMessageCreate{
		messageRepository: r,
		messageEvent:      e,
	}
}

func (u UsecaseMessageCreate) Execute(ctx context.Context, input UsecaseMessageCreateInput) (domain.Message, error) {
	message, err := domain.NewMessage("", input.UserName, input.UserEmail, input.MessageText)
	if err != nil {
		return domain.Message{}, err
	}

	if message.StockCommand() {
		err := u.messageEvent.StockCommandReceived(ctx, message)
		if err != nil {
			return domain.Message{}, err
		}

		return message, nil
	}

	return u.messageRepository.Create(ctx, message)
}

type UsecaseMessagesList struct {
	messageRepository domain.MessageRepository
}

func NewUsecaseMessageList(r domain.MessageRepository) UsecaseMessagesList {
	return UsecaseMessagesList{
		messageRepository: r,
	}
}

func (u UsecaseMessagesList) Execute(ctx context.Context, input UsecaseMessagesListInput) ([]domain.Message, error) {
	pagination, err := domain.NewPagination(input.Page, input.ItemsPerPage)
	if err != nil {
		return []domain.Message{}, err
	}

	return u.messageRepository.List(ctx, pagination)
}
