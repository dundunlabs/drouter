package validation

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func ValidateEnum(fl validator.FieldLevel) bool {
	f := fl.Field()
	values := f.MethodByName("Values").Call([]reflect.Value{})[0]
	for i := 0; i < values.Len(); i++ {
		if values.Index(i).Equal(f) {
			return true
		}
	}
	return false
}
