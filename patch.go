package patchstruct

import (
	"errors"
	"reflect"
)

func Patch(src interface{}, dst interface{}, lockTag string) error {
	if reflect.TypeOf(src).Kind() != reflect.Ptr {
		return errors.New("src_not_struct_pointer")
	}

	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return errors.New("dst_not_struct_pointer")
	}

	// not necessary because we search by field name
	// if reflect.TypeOf(src).Kind() != reflect.TypeOf(dst).Kind() {
	// 	return errors.New("src_dst_different_types")
	// }

	srcStructValue := reflect.ValueOf(src).Elem()
	// srcStructType := reflect.TypeOf(src).Elem()

	dstStructValue := reflect.ValueOf(dst).Elem()
	dstStructType := reflect.TypeOf(dst).Elem()

	if srcStructValue.Kind() != reflect.Struct {
		return errors.New("src underlying value must be a struct (*T)")
	}

	if dstStructValue.Kind() != reflect.Struct {
		return errors.New("dst underlying value must be a struct (*T)")
	}


	for i := 0; i < srcStructValue.NumField(); i++ {
		fieldName := srcStructValue.Type().Field(i).Name

		dstFieldValue := dstStructValue.FieldByName(fieldName)

		if dstFieldValue.IsZero() {
			// dst field not found
			continue
		}

		if !dstFieldValue.CanAddr() {
			// dst field is not addressable return error
			return errors.New(`cant address field on dst`)
		}

		if _, locked := dstStructType.Field(i).Tag.Lookup(lockTag); locked {
			// field is locked by tag
			continue
		}

		// not needed anymore. if field is pointer, skip only if nil (later)
		// if srcStructValue.Field(i).Kind() == reflect.Ptr && !srcStructValue.Field(i).IsNil() && srcStructValue.Field(i).Elem().IsZero() {
		// 	// field is ptr and value is zero
		// 	continue
		// }


		// same name, different type // maybe if destination is ptr, set value
		if dstStructValue.FieldByName(fieldName).Kind() != srcStructValue.Field(i).Kind() {
			// fields have different types
			continue
		}

		// src field is nil ptr
		if srcStructValue.Field(i).Kind() == reflect.Ptr && srcStructValue.Field(i).IsNil() {
			// skipping src nil pointer // nil means it was not set, example from json
			continue
		}

		// src field value (or ptr underlying value) is valid and not zero
		if srcStructValue.Field(i).IsValid() && !srcStructValue.Field(i).IsZero() {
			dstStructValue.FieldByName(fieldName).Set(srcStructValue.Field(i))
		}
	}

	return nil
}
