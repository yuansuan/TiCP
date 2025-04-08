package validation

import (
	"errors"
	"fmt"
	"reflect"

	"gopkg.in/go-playground/validator.v9"
)

// Error 错误
type Error struct {
	errorMessage string
	validator.FieldError
}

func (e Error) Error() string {
	return e.errorMessage
}

// FieldErrorHandler 字段错误处理函数
type FieldErrorHandler func(fe Error) bool

// HandleFieldError 处理字段错误
func (e Error) HandleFieldError(handlers ...FieldErrorHandler) bool {
	for _, handler := range handlers {
		if handler != nil && handler(e) {
			return true
		}
	}
	return false
}

// HandleError 处理错误
func HandleError(err error, request interface{}) error {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err
	}

	reqType := reflect.TypeOf(request)
	if reqType.Kind() == reflect.Ptr {
		reqType = reqType.Elem()
	}
	for _, fieldError := range validationErrors {
		filed, found := reqType.FieldByName(fieldError.Field())
		if !found {
			return Error{fmt.Sprintf("field %s not found", fieldError.Field()), fieldError}
		}

		// 首先判断error是否有配置，有配置直接返回
		err := handleCustomError(filed)
		if err != nil {
			return Error{err.Error(), fieldError}
		}

		// 判断是否是通用错误，是通用错误返回
		err = handleCommonError(fieldError)
		if err != nil {
			return Error{err.Error(), fieldError}
		}
	}

	return err
}

func handleCustomError(filed reflect.StructField) error {
	errMessage := filed.Tag.Get("error")
	if errMessage != "" {
		return fmt.Errorf(errMessage)
	}

	return nil
}

func handleCommonError(fe validator.FieldError) error {
	fieldName := fe.Field()
	switch fe.Tag() {
	case "required":
		return fmt.Errorf("field %s is required", fieldName)
	case "max":
		return fmt.Errorf("field %s maximum character num is %s", fieldName, fe.Param())
	case "min":
		return fmt.Errorf("field %s minimum character num is %s", fieldName, fe.Param())
	}

	return nil
}
