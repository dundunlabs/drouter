package validation_test

import (
	"fmt"
	"testing"

	"github.com/dundunlabs/drouter/validate/validation"
)

func init() {
	validate.RegisterValidation("enum", validation.ValidateEnum)
}

type VehicleType string

func (t VehicleType) Values() []VehicleType {
	return []VehicleType{"car", "truck"}
}

type Vehicle struct {
	Type VehicleType `validate:"enum"`
}

func TestValidateEnum(t *testing.T) {
	tests := []struct {
		vehicle Vehicle
		result  bool
	}{
		{Vehicle{"car"}, true},
		{Vehicle{"truck"}, true},
		{Vehicle{"bike"}, false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprint(test.vehicle.Type), func(t *testing.T) {
			if result := validate.Struct(test.vehicle) == nil; result != test.result {
				t.Errorf("want %v, got %v", test.result, result)
			}
		})
	}
}
