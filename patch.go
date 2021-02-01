package patchstruct

import (
	"errors"
	"reflect"
)

func Patch(src interface{}, dst interface{}, lockTag string) error {
	if reflect.TypeOf(src).Kind() != reflect.Ptr {
		return errors.New("src_not_pointer: " + reflect.TypeOf(src).Kind().String())
	}

	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return errors.New("dst_not_pointer: " + reflect.TypeOf(dst).Kind().String())
	}

	A := reflect.ValueOf(src).Elem()
	AT := reflect.TypeOf(src).Elem()

	if A.Kind() != reflect.Struct {
		return errors.New("src underlying value must be a struct (*Struct)")
	}

	B := reflect.ValueOf(dst).Elem()
	BT := reflect.TypeOf(dst).Elem()

	if B.Kind() != reflect.Struct {
		return errors.New("dst underlying value must be a struct (*Struct)")
	}


	for i := 0; i < A.NumField(); i++ {
		/*
			Will accept: Field with name F (public) and type or pointer to type.
			Rules for copy

			 a zero			 skip		ok
			*a nil			 skip		ok
			 a != (*)b		 skip		ok
			 a 				 b =>  a	ok
			 a		 		*b => *a	ok
		  (*)a != b       	 skip		ok
			*a      		 b =>  a	ok
			*a           	*b => *a	ok
		  (*)a != (*)b		 skip		ok
		*/

		fieldName := AT.Field(i).Name

		Af, _ := AT.FieldByName(fieldName)
		Bf, ok := BT.FieldByName(fieldName)
		if !ok {
			// B does not contain field
			continue
		}

		// field pointer kind
		var Afpk = reflect.Invalid
		var Bfpk = reflect.Invalid

		// field final kind (real value kind, not pointer)
		var Affk = reflect.Invalid
		var Bffk = reflect.Invalid

		Afk := Af.Type.Kind()
		Affk = Afk
		if Afk == reflect.Ptr {
			Afpk = Af.Type.Elem().Kind()
			Affk = Afpk
		}

		Bfk := Bf.Type.Kind()
		Bffk = Bfk
		if Bfk == reflect.Ptr {
			Bfpk = Bf.Type.Elem().Kind()
			Bffk = Bfpk
		}

		// underlying field kinds are invalid or different
		//    a != (*)b		skip
		// (*)a != b		skip
		// (*)a != (*)b		skip
		if Affk == reflect.Invalid || Bffk == reflect.Invalid || Affk != Bffk {
			continue
		}

		if _, locked := Bf.Tag.Lookup(lockTag); locked {
			continue
		}

		Av := A.FieldByName(fieldName)
		// not found
		// a zero			 skip
		if Av.IsZero() {
			continue
		}

		Bv := B.FieldByName(fieldName)

		if Av.Kind() != reflect.Ptr {
			// a
			// a 				 b =>  a
			if Bv.Kind() != reflect.Ptr {
				Bv.Set(Av)
				continue
			}
			// a		 		*b => *a
			if Bv.IsNil() {
				Bv = reflect.New(Bv.Type())
			}
			Bv.Elem().Set(Av)
			continue
		} else {
			// *a
			// *a nil			 skip
			if Av.IsNil() {
				continue
			}
			// *a      		 b =>  a
			if Bv.Kind() != reflect.Ptr {
				Bv.Set(Av.Elem())
				continue
			}
			// *a           *b => *a
			Bv.Set(Av)
			continue
		}
	}

	return nil
}
