package validation

import (
	"github.com/go-playground/validator/v10"
	"reflect"
)

type Validator interface {
	Validate(obj any) error
}

type DefaultValidator struct {
	validate *validator.Validate
}

func NewDefaultValidator() *DefaultValidator {
	return &DefaultValidator{
		validate: validator.New(),
	}
}

func (v *DefaultValidator) Validate(obj any) error {
	objectsSlice := make([]any, 0)

	switch reflect.TypeOf(obj).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(obj)

		for i := 0; i < s.Len(); i++ {
			objectsSlice = append(objectsSlice, s.Index(i).Interface())
		}
	default:
		objectsSlice = append(objectsSlice, obj)
	}

	for _, object := range objectsSlice {
		if err := v.validate.Struct(object); err != nil {
			return err
		}
	}

	return nil
}
