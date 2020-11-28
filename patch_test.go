package patchstruct

import (
	"fmt"
	"testing"
)

type T1 struct {
	A	string
	B	*string
	C	string `lock:""`
	D	*string `lock:""`
}

func (t T1) String() string {
	return fmt.Sprintf("A:%+v, B:%+v, C:%+v, D:%+v", t.A, t.B, t.C, t.D)
}

var s1 = `a`
var s2 = `z`
var sE = ``


var dst = &T1 {A:s1, B:&s1, C:s1, D:&s1}
var src1 = &T1 {A:s2, B:&s2, C:s2, D:&s2}
var res1 = &T1 {A:s2, B:&s2, C:s1, D:&s1}

var srcE = &T1 {A:sE, B:&sE, C:sE, D:&sE}
var resE = &T1 {A:s1, B:&sE, C:s1, D:&sE}

var srcP = &T1 {B:nil, D:nil}
var resP = &T1 {A:s1, B:&s1, C:s1, D:&s1}


func TestPatchStruct(t *testing.T) {
	type args struct {
		src     interface{}
		dst     interface{}
		res 	interface{}
		lockTag string
	}

	dst1 := new(T1)
	*dst1 = *dst
	dst2 := new(T1)
	*dst2 = *dst
	dst3 := new(T1)
	*dst3 = *dst

	tests := []struct {
		name string
		args args
	}{
		{`strings z`, args{src1, dst1, res1, `lock`}},
		{`strings e`, args{srcE, dst2, resE, ``}},
		{`strings p`, args{srcP, dst3, resP, ``}},

	}

	for _, tt := range tests {
		// log.Printf("%+v %+v %+v", tt, tt.name, tt.args)
		t.Run(tt.name, func(t *testing.T) {
			if err := Patch(tt.args.src, tt.args.dst, tt.args.lockTag); err != nil {
				t.Error(err)
			}
			if tt.args.dst.(*T1).String() != tt.args.res.(*T1).String() {
				t.Errorf("%+v != %+v", tt.args.dst.(*T1).String(), tt.args.res.(*T1).String())
			}
		})
	}

}