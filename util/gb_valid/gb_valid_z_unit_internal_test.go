package gbvalid

import (
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func Test_parseSequenceTag(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := "name@required|length:2,20|password3|same:password1#||密碼強度不足|兩次密碼不一致"
		field, rule, msg := ParseTagValue(s)
		t.Assert(field, "name")
		t.Assert(rule, "required|length:2,20|password3|same:password1")
		t.Assert(msg, "||密碼強度不足|兩次密碼不一致")
	})
	gbtest.C(t, func(t *gbtest.T) {
		s := "required|length:2,20|password3|same:password1#||密碼強度不足|兩次密碼不一致"
		field, rule, msg := ParseTagValue(s)
		t.Assert(field, "")
		t.Assert(rule, "required|length:2,20|password3|same:password1")
		t.Assert(msg, "||密碼強度不足|兩次密碼不一致")
	})
	gbtest.C(t, func(t *gbtest.T) {
		s := "required|length:2,20|password3|same:password1"
		field, rule, msg := ParseTagValue(s)
		t.Assert(field, "")
		t.Assert(rule, "required|length:2,20|password3|same:password1")
		t.Assert(msg, "")
	})
	gbtest.C(t, func(t *gbtest.T) {
		s := "required"
		field, rule, msg := ParseTagValue(s)
		t.Assert(field, "")
		t.Assert(rule, "required")
		t.Assert(msg, "")
	})
}

func Test_GetTags(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(structTagPriority, GetTags())
	})
}
