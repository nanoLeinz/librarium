package dto

import (
	"reflect"

	"github.com/google/uuid"
)

type MemberUpdateRequest struct {
	ID            uuid.UUID `json:"-"`
	Email         string    `json:"email" validate:"omitempty,email"`
	Password      string    `json:"password" validate:"omitempty,min=4"`
	FullName      string    `json:"full_name" validate:"omitempty,max=50"`
	AccountStatus string    `json:"account_status" validate:"omitempty"`
}

func StructToMap(data any) (result map[string]interface{}) {

	values := reflect.ValueOf(data)

	result = make(map[string]interface{}, values.NumField())

	for i := 0; i < values.NumField(); i++ {
		if values.Field(i).CanInterface() && !values.Field(i).IsZero() {
			result[values.Type().Field(i).Name] = values.Field(i).Interface()
		}
	}

	return

}
