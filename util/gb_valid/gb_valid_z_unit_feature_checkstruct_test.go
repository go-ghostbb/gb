package gbvalid_test

import (
	"context"
	gbvar "ghostbb.io/gb/container/gb_var"
	"ghostbb.io/gb/frame/g"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbconv "ghostbb.io/gb/util/gb_conv"
	"testing"
)

func Test_CheckStruct(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Object struct {
			Name string
			Age  int
		}
		rules := []string{
			"@required|length:6,16",
			"@between:18,30",
		}
		msgs := map[string]interface{}{
			"Name": map[string]string{
				"required": "名稱不能為空",
				"length":   "名稱長度為{min}到{max}個字符",
			},
			"Age": "年齡為18到30周歲",
		}
		obj := &Object{"john", 16}
		err := g.Validator().Data(obj).Rules(rules).Messages(msgs).Run(context.TODO())
		t.AssertNil(err)
	})

	gbtest.C(t, func(t *gbtest.T) {
		type Object struct {
			Name string
			Age  int
		}
		rules := []string{
			"Name@required|length:6,16#名稱不能為空",
			"Age@between:18,30",
		}
		msgs := map[string]interface{}{
			"Name": map[string]string{
				"required": "名稱不能為空",
				"length":   "名稱長度為{min}到{max}個字符",
			},
			"Age": "年齡為18到30周歲",
		}
		obj := &Object{"john", 16}
		err := g.Validator().Data(obj).Rules(rules).Messages(msgs).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(len(err.Maps()), 2)
		t.Assert(err.Maps()["Name"]["required"], "")
		t.Assert(err.Maps()["Name"]["length"], "名稱長度為6到16個字符")
		t.Assert(err.Maps()["Age"]["between"], "年齡為18到30周歲")
	})

	gbtest.C(t, func(t *gbtest.T) {
		type Object struct {
			Name string
			Age  int
		}
		rules := []string{
			"Name@required|length:6,16#名稱不能為空|",
			"Age@between:18,30",
		}
		msgs := map[string]interface{}{
			"Name": map[string]string{
				"required": "名稱不能為空",
				"length":   "名稱長度為{min}到{max}個字符",
			},
			"Age": "年齡為18到30周歲",
		}
		obj := &Object{"john", 16}
		err := g.Validator().Data(obj).Rules(rules).Messages(msgs).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(len(err.Maps()), 2)
		t.Assert(err.Maps()["Name"]["required"], "")
		t.Assert(err.Maps()["Name"]["length"], "名稱長度為6到16個字符")
		t.Assert(err.Maps()["Age"]["between"], "年齡為18到30周歲")
	})

	gbtest.C(t, func(t *gbtest.T) {
		type Object struct {
			Name string
			Age  int
		}
		rules := map[string]string{
			"Name": "required|length:6,16",
			"Age":  "between:18,30",
		}
		msgs := map[string]interface{}{
			"Name": map[string]string{
				"required": "名稱不能為空",
				"length":   "名稱長度為{min}到{max}個字符",
			},
			"Age": "年齡為18到30周歲",
		}
		obj := &Object{"john", 16}
		err := g.Validator().Data(obj).Rules(rules).Messages(msgs).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(len(err.Maps()), 2)
		t.Assert(err.Maps()["Name"]["required"], "")
		t.Assert(err.Maps()["Name"]["length"], "名稱長度為6到16個字符")
		t.Assert(err.Maps()["Age"]["between"], "年齡為18到30周歲")
	})

	gbtest.C(t, func(t *gbtest.T) {
		type LoginRequest struct {
			Username string `json:"username" valid:"username@required#使用者名稱不能為空"`
			Password string `json:"password" valid:"password@required#登入密碼不能為空"`
		}
		var login LoginRequest
		err := g.Validator().Data(login).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(len(err.Maps()), 2)
		t.Assert(err.Maps()["username"]["required"], "使用者名稱不能為空")
		t.Assert(err.Maps()["password"]["required"], "登入密碼不能為空")
	})

	gbtest.C(t, func(t *gbtest.T) {
		type LoginRequest struct {
			Username string `json:"username" valid:"@required#使用者名稱不能為空"`
			Password string `json:"password" valid:"@required#登入密碼不能為空"`
		}
		var login LoginRequest
		err := g.Validator().Data(login).Run(context.TODO())
		t.AssertNil(err)
	})

	gbtest.C(t, func(t *gbtest.T) {
		type LoginRequest struct {
			username string `json:"username" valid:"username@required#使用者名稱不能為空"`
			Password string `json:"password" valid:"password@required#登入密碼不能為空"`
		}
		var login LoginRequest
		err := g.Validator().Data(login).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(err.Maps()["password"]["required"], "登入密碼不能為空")
	})

	// gvalid tag
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Id       int    `valid:"uid@required|min:10#|ID不能為空"`
			Age      int    `valid:"age@required#年齡不能為空"`
			Username string `json:"username" valid:"username@required#使用者名稱不能為空"`
			Password string `json:"password" valid:"password@required#登入密碼不能為空"`
		}
		user := &User{
			Id:       1,
			Username: "john",
			Password: "123456",
		}
		err := g.Validator().Data(user).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(len(err.Maps()), 1)
		t.Assert(err.Maps()["uid"]["min"], "ID不能為空")
	})

	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Id       int    `valid:"uid@required|min:10#|ID不能為空"`
			Age      int    `valid:"age@required#年齡不能為空"`
			Username string `json:"username" valid:"username@required#使用者名稱不能為空"`
			Password string `json:"password" valid:"password@required#登入密碼不能為空"`
		}
		user := &User{
			Id:       1,
			Username: "john",
			Password: "123456",
		}

		rules := []string{
			"username@required#使用者名稱不能為空",
		}

		err := g.Validator().Data(user).Rules(rules).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(len(err.Maps()), 1)
		t.Assert(err.Maps()["uid"]["min"], "ID不能為空")
	})

	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Id       int    `valid:"uid@required|min:10#ID不能為空"`
			Age      int    `valid:"age@required#年齡不能為空"`
			Username string `json:"username" valid:"username@required#使用者名稱不能為空"`
			Password string `json:"password" valid:"password@required#登入密碼不能為空"`
		}
		user := &User{
			Id:       1,
			Username: "john",
			Password: "123456",
		}
		err := g.Validator().Data(user).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(len(err.Maps()), 1)
	})

	// valid tag
	gbtest.C(t, func(t *gbtest.T) {
		type User struct {
			Id       int    `valid:"uid@required|min:10#|ID不能為空"`
			Age      int    `valid:"age@required#年齡不能為空"`
			Username string `json:"username" valid:"username@required#使用者名稱不能為空"`
			Password string `json:"password" valid:"password@required#登入密碼不能為空"`
		}
		user := &User{
			Id:       1,
			Username: "john",
			Password: "123456",
		}
		err := g.Validator().Data(user).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(len(err.Maps()), 1)
		t.Assert(err.Maps()["uid"]["min"], "ID不能為空")
	})
}

func Test_CheckStruct_EmbeddedObject_Attribute(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Base struct {
			Time *gbtime.Time
		}
		type Object struct {
			Base
			Name string
			Type int
		}
		rules := map[string]string{
			"Name": "required",
			"Type": "required",
		}
		ruleMsg := map[string]interface{}{
			"Name": "名稱必填",
			"Type": "類型必填",
		}
		obj := &Object{}
		obj.Type = 1
		obj.Name = "john"
		obj.Time = gbtime.Now()
		err := g.Validator().Data(obj).Rules(rules).Messages(ruleMsg).Run(context.TODO())
		t.AssertNil(err)
	})
	gbtest.C(t, func(t *gbtest.T) {
		type Base struct {
			Name string
			Type int
		}
		type Object struct {
			Base Base
			Name string
			Type int
		}
		rules := map[string]string{
			"Name": "required",
			"Type": "required",
		}
		ruleMsg := map[string]interface{}{
			"Name": "名稱必填",
			"Type": "類型必填",
		}
		obj := &Object{}
		obj.Type = 1
		obj.Name = "john"
		err := g.Validator().Data(obj).Rules(rules).Messages(ruleMsg).Run(context.TODO())
		t.AssertNil(err)
	})
}

func Test_CheckStruct_With_EmbeddedObject(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Pass struct {
			Pass1 string `valid:"password1@required|same:password2#請輸入您的密碼|您兩次輸入的密碼不一致"`
			Pass2 string `valid:"password2@required|same:password1#請再次輸入您的密碼|您兩次輸入的密碼不一致"`
		}
		type User struct {
			Id   int
			Name string `valid:"name@required#請輸入您的姓名"`
			Pass
		}
		user := &User{
			Name: "",
			Pass: Pass{
				Pass1: "1",
				Pass2: "2",
			},
		}
		err := g.Validator().Data(user).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(err.Maps()["name"], g.Map{"required": "請輸入您的姓名"})
		t.Assert(err.Maps()["password1"], g.Map{"same": "您兩次輸入的密碼不一致"})
		t.Assert(err.Maps()["password2"], g.Map{"same": "您兩次輸入的密碼不一致"})
	})
}

func Test_CheckStruct_With_StructAttribute(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Pass struct {
			Pass1 string `valid:"password1@required|same:password2#請輸入您的密碼|您兩次輸入的密碼不一致"`
			Pass2 string `valid:"password2@required|same:password1#請再次輸入您的密碼|您兩次輸入的密碼不一致"`
		}
		type User struct {
			Pass
			Id   int
			Name string `valid:"name@required#請輸入您的姓名"`
		}
		user := &User{
			Name: "",
			Pass: Pass{
				Pass1: "1",
				Pass2: "2",
			},
		}
		err := g.Validator().Data(user).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(err.Maps()["name"], g.Map{"required": "請輸入您的姓名"})
		t.Assert(err.Maps()["password1"], g.Map{"same": "您兩次輸入的密碼不一致"})
		t.Assert(err.Maps()["password2"], g.Map{"same": "您兩次輸入的密碼不一致"})
	})
}

func Test_CheckStruct_Optional(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Params struct {
			Page      int    `v:"required|min:1         # page is required"`
			Size      int    `v:"required|between:1,100 # size is required"`
			ProjectId string `v:"between:1,10000        # project id must between {min}, {max}"`
		}
		obj := &Params{
			Page: 1,
			Size: 10,
		}
		err := g.Validator().Data(obj).Run(context.TODO())
		t.AssertNil(err)
	})
	gbtest.C(t, func(t *gbtest.T) {
		type Params struct {
			Page      int        `v:"required|min:1         # page is required"`
			Size      int        `v:"required|between:1,100 # size is required"`
			ProjectId *gbvar.Var `v:"between:1,10000        # project id must between {min}, {max}"`
		}
		obj := &Params{
			Page: 1,
			Size: 10,
		}
		err := g.Validator().Data(obj).Run(context.TODO())
		t.AssertNil(err)
	})
	gbtest.C(t, func(t *gbtest.T) {
		type Params struct {
			Page      int `v:"required|min:1         # page is required"`
			Size      int `v:"required|between:1,100 # size is required"`
			ProjectId int `v:"between:1,10000        # project id must between {min}, {max}"`
		}
		obj := &Params{
			Page: 1,
			Size: 10,
		}
		err := g.Validator().Data(obj).Run(context.TODO())
		t.Assert(err.String(), "project id must between 1, 10000")
	})
}

func Test_CheckStruct_NoTag(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Params struct {
			Page      int
			Size      int
			ProjectId string
		}
		obj := &Params{
			Page: 1,
			Size: 10,
		}
		err := g.Validator().Data(obj).Run(context.TODO())
		t.AssertNil(err)
	})
}

func Test_CheckStruct_InvalidRule(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Params struct {
			Name  string
			Age   uint
			Phone string `v:"mobile"`
		}
		obj := &Params{
			Name:  "john",
			Age:   18,
			Phone: "123",
		}
		err := g.Validator().Data(obj).Run(context.TODO())
		t.AssertNE(err, nil)
	})
}

func TestValidator_CheckStructWithData(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type UserApiSearch struct {
			Uid      int64  `v:"required"`
			Nickname string `v:"required-with:uid"`
		}
		data := UserApiSearch{
			Uid:      1,
			Nickname: "john",
		}
		t.Assert(
			g.Validator().Data(data).Assoc(
				g.Map{"uid": 1, "nickname": "john"},
			).Run(context.TODO()),
			nil,
		)
	})
	gbtest.C(t, func(t *gbtest.T) {
		type UserApiSearch struct {
			Uid      int64  `v:"required"`
			Nickname string `v:"required-with:uid"`
		}
		data := UserApiSearch{}
		t.AssertNE(g.Validator().Data(data).Assoc(g.Map{}).Run(context.TODO()), nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		type UserApiSearch struct {
			Uid      int64  `json:"uid" v:"required"`
			Nickname string `json:"nickname" v:"required-with:Uid"`
		}
		data := UserApiSearch{
			Uid: 1,
		}
		t.AssertNE(g.Validator().Data(data).Assoc(g.Map{}).Run(context.TODO()), nil)
	})

	gbtest.C(t, func(t *gbtest.T) {
		type UserApiSearch struct {
			Uid       int64        `json:"uid"`
			Nickname  string       `json:"nickname" v:"required-with:Uid"`
			StartTime *gbtime.Time `json:"start_time" v:"required-with:EndTime"`
			EndTime   *gbtime.Time `json:"end_time" v:"required-with:StartTime"`
		}
		data := UserApiSearch{
			StartTime: nil,
			EndTime:   nil,
		}
		t.Assert(g.Validator().Data(data).Assoc(g.Map{}).Run(context.TODO()), nil)
	})
	gbtest.C(t, func(t *gbtest.T) {
		type UserApiSearch struct {
			Uid       int64        `json:"uid"`
			Nickname  string       `json:"nickname" v:"required-with:Uid"`
			StartTime *gbtime.Time `json:"start_time" v:"required-with:EndTime"`
			EndTime   *gbtime.Time `json:"end_time" v:"required-with:StartTime"`
		}
		data := UserApiSearch{
			StartTime: gbtime.Now(),
			EndTime:   nil,
		}
		t.AssertNE(g.Validator().Data(data).Assoc(g.Map{"start_time": gbtime.Now()}).Run(context.TODO()), nil)
	})
}

func Test_CheckStruct_PointerAttribute(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		type Req struct {
			Name string
			Age  *uint `v:"min:18"`
		}
		req := &Req{
			Name: "john",
			Age:  gbconv.PtrUint(0),
		}
		err := g.Validator().Data(req).Run(context.TODO())
		t.Assert(err.String(), "The Age value `0` must be equal or greater than 18")
	})
	gbtest.C(t, func(t *gbtest.T) {
		type Req struct {
			Name string `v:"min-length:3"`
			Age  *uint  `v:"min:18"`
		}
		req := &Req{
			Name: "j",
			Age:  gbconv.PtrUint(19),
		}
		err := g.Validator().Data(req).Run(context.TODO())
		t.Assert(err.String(), "The Name value `j` length must be equal or greater than 3")
	})
	gbtest.C(t, func(t *gbtest.T) {
		type Params struct {
			Age *uint `v:"min:18"`
		}
		type Req struct {
			Name   string
			Params *Params
		}
		req := &Req{
			Name: "john",
			Params: &Params{
				Age: gbconv.PtrUint(0),
			},
		}
		err := g.Validator().Data(req).Run(context.TODO())
		t.Assert(err.String(), "The Age value `0` must be equal or greater than 18")
	})
}
