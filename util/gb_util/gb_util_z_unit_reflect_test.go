package gbutil_test

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	gbutil "ghostbb.io/gb/util/gb_util"
	"reflect"
	"testing"
)

func Test_OriginValueAndKind(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var s = "s"
		out := gbutil.OriginValueAndKind(s)
		t.Assert(out.InputKind, reflect.String)
		t.Assert(out.OriginKind, reflect.String)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var s = "s"
		out := gbutil.OriginValueAndKind(&s)
		t.Assert(out.InputKind, reflect.Ptr)
		t.Assert(out.OriginKind, reflect.String)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var s []int
		out := gbutil.OriginValueAndKind(s)
		t.Assert(out.InputKind, reflect.Slice)
		t.Assert(out.OriginKind, reflect.Slice)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var s []int
		out := gbutil.OriginValueAndKind(&s)
		t.Assert(out.InputKind, reflect.Ptr)
		t.Assert(out.OriginKind, reflect.Slice)
	})
}

func Test_OriginTypeAndKind(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var s = "s"
		out := gbutil.OriginTypeAndKind(s)
		t.Assert(out.InputKind, reflect.String)
		t.Assert(out.OriginKind, reflect.String)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var s = "s"
		out := gbutil.OriginTypeAndKind(&s)
		t.Assert(out.InputKind, reflect.Ptr)
		t.Assert(out.OriginKind, reflect.String)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var s []int
		out := gbutil.OriginTypeAndKind(s)
		t.Assert(out.InputKind, reflect.Slice)
		t.Assert(out.OriginKind, reflect.Slice)
	})
	gbtest.C(t, func(t *gbtest.T) {
		var s []int
		out := gbutil.OriginTypeAndKind(&s)
		t.Assert(out.InputKind, reflect.Ptr)
		t.Assert(out.OriginKind, reflect.Slice)
	})
}
