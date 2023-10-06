package infra

import (
	"errors"
	"fmt"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/application"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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
		"time":   DatetimeCanonical(&t),
	})
}

func (c ApiBaseController) HandleNotFound(ctx *gin.Context) {
	t := time.Now()
	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "page not found",
		"time":  DatetimeCanonical(&t),
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
