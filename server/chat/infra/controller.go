package infra

import (
	"errors"
	"fmt"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/application"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"
	"joubertredrat-tests/jobsity-dev-test-2k23/pkg"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

const (
	QUERY_STRING_PAGE                 = "page"
	QUERY_STRING_ITEMS_PER_PAGE       = "itemsPerPage"
	DEFAULT_PAGINATION_PAGE           = 1
	DEFAULT_PAGINATION_ITEMS_PER_PAGE = 50
)

type ApiBaseController struct {
}

func NewApiBaseController() ApiBaseController {
	return ApiBaseController{}
}

func (c ApiBaseController) HandleStatus(ctx *gin.Context) {
	t := time.Now()
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   pkg.DatetimeCanonical(&t),
	})
}

func (c ApiBaseController) HandleNotFound(ctx *gin.Context) {
	t := time.Now()
	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "page not found",
		"time":  pkg.DatetimeCanonical(&t),
	})
}

type UserController struct {
}

func NewUserController() UserController {
	return UserController{}
}

func (c UserController) HandleCreate(usecase application.UsecaseUserRegister) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request UserRegisterRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			responseWithError(ctx, err)
			return
		}

		user, err := usecase.Execute(ctx, application.UsecaseUserRegisterInput{
			Name:     request.Name,
			Email:    request.Email,
			Password: request.Password,
		})

		if err != nil {
			switch err.(type) {
			case domain.ErrInvalidEmail,
				domain.ErrInvalidPasswordLength:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
			case application.ErrUserAlreadyRegistered:
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": err.Error(),
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
			return
		}

		ctx.JSON(http.StatusCreated, UserRegisterResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Password: "it's a secret :)",
		})
	}
}

func (c UserController) HandleLogin(usecase application.UsecaseUserLogin) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request UserLoginRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			responseWithError(ctx, err)
			return
		}

		userToken, err := usecase.Execute(ctx, application.UsecaseUserLoginInput{
			Email:    request.Email,
			Password: request.Password,
		})

		if err != nil {
			switch err.(type) {
			case domain.ErrInvalidEmail,
				domain.ErrInvalidPasswordLength:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
			case domain.ErrUserNotFoundByEmail:
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
			case domain.ErrUserNotAuthenticated:
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": err.Error(),
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
			return
		}

		ctx.JSON(http.StatusOK, UserLoginResponse{
			AccessToken: userToken.AccessToken,
		})
	}
}

type MessagesController struct {
}

func NewMessagesController() MessagesController {
	return MessagesController{}
}

func (c MessagesController) HandleCreate(usecase application.UsecaseMessageCreate) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request MessageCreateRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			responseWithError(ctx, err)
			return
		}

		var user, ok = ctx.Get("userAuth")
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
		userAuth := user.(domain.User)

		message, err := usecase.Execute(ctx, application.UsecaseMessageCreateInput{
			UserName:    userAuth.Name,
			UserEmail:   userAuth.Email,
			MessageText: request.MessageText,
		})

		if err != nil {
			switch err.(type) {
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
			return
		}

		ctx.JSON(http.StatusCreated, MessageResponse{
			ID:          message.ID,
			UserName:    message.UserName,
			UserEmail:   message.UserEmail,
			MessageText: message.Text,
			Datetime:    pkg.DatetimeCanonical(&message.Datetime),
		})
	}
}

func (c MessagesController) HandleList(usecase application.UsecaseMessagesList) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page := stringToUint(
			ctx.DefaultQuery(QUERY_STRING_PAGE, fmt.Sprintf("%d", DEFAULT_PAGINATION_PAGE)),
		)
		itemsPerPage := stringToUint(
			ctx.DefaultQuery(QUERY_STRING_ITEMS_PER_PAGE, fmt.Sprintf("%d", DEFAULT_PAGINATION_ITEMS_PER_PAGE)),
		)

		messages, err := usecase.Execute(ctx, application.UsecaseMessagesListInput{
			Page:         page,
			ItemsPerPage: itemsPerPage,
		})

		if err != nil {
			switch err.(type) {
			case domain.ErrPaginationPage,
				domain.ErrPaginationItemsPerPage:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
			return
		}

		response := []MessageResponse{}

		for _, message := range messages {
			response = append(response, MessageResponse{
				ID:          message.ID,
				UserName:    message.UserName,
				UserEmail:   message.UserEmail,
				MessageText: message.Text,
				Datetime:    pkg.DatetimeCanonical(&message.Datetime),
			})
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func stringToUint(s string) uint {
	i, _ := strconv.Atoi(s)
	return uint(i)
}

func RegisterCustomValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

func responseWithError(c *gin.Context, err error) {
	var verr validator.ValidationErrors
	if errors.As(err, &verr) {
		c.JSON(http.StatusBadRequest, gin.H{"errors": getValidatorErrors(verr)})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func getValidatorErrors(verr validator.ValidationErrors) []RequestValidationError {
	var errs []RequestValidationError

	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}

		errs = append(errs, RequestValidationError{Field: f.Field(), Reason: err})
	}

	return errs
}
