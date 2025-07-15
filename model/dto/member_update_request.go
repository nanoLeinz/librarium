package dto

import (
	"reflect"

	"github.com/google/uuid"
)

type MemberUpdateRequest struct {
	ID            *uuid.UUID
	Email         *string
	Password      *string
	FullName      *string
	AccountStatus *string
}

func StructToMap(data any) (result map[string]interface{}) {

	result = make(map[string]interface{})

	val := reflect.ValueOf(data).Elem()

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		types := typ.Field(i)

		if !field.IsNil() {
			result[types.Tag.Get("json")] = field.Elem().Interface()
		}

	}

	return

}
