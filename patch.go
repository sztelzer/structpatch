package patchstruct

import (
	"errors"
	"reflect"
)

func PatchStruct(src interface{}, dst interface{}, lockTag string) error {
	var srcStructValue reflect.Value
	var dstStructValue reflect.Value
	var dstStructType reflect.Type

	if reflect.TypeOf(src).Kind() != reflect.Ptr {
		return errors.New("source must be a pointer to struct (*T)")
	}

	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return errors.New("destination must be a pointer to struct (*T)")
	}

	if reflect.TypeOf(src).Kind() != reflect.TypeOf(dst).Kind() {
		return errors.New("source type (*T) should be same of destination type (*T)")
	}

	srcStructValue = reflect.ValueOf(src).Elem()
	dstStructValue = reflect.ValueOf(dst).Elem()
	dstStructType = reflect.TypeOf(dst).Elem()

	if srcStructValue.Kind() == reflect.Ptr {
		return errors.New("source cant be a pointer to a pointer (**T) ")
	}

	if dstStructValue.Kind() == reflect.Ptr {
		return errors.New("destination cant be a pointer to a pointer (**T)")
	}


	for i := 0; i < srcStructValue.NumField(); i++ {
		// if dst field is not addressable return error
		if !dstStructValue.Field(i).CanAddr() {
			return errors.New(`cant address field on dst`)
		}

		// field is locked
		if _, locked := dstStructType.Field(i).Tag.Lookup(lockTag); locked {
			continue
		}

		// field is ptr and value is zero
		if srcStructValue.Field(i).Kind() == reflect.Ptr && srcStructValue.Field(i).Elem().IsZero() {
			continue
		}

		// field value (and ptr underlying value) is valid and not zero
		if srcStructValue.Field(i).IsValid() && !srcStructValue.Field(i).IsZero() {
			dstStructValue.Field(i).Set(srcStructValue.Field(i))
		}
	}

	return nil
}
