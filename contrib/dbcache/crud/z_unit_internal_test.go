package crud

import (
	"ghostbb.io/gb/contrib/dbcache/cache"
	gbtest "ghostbb.io/gb/test/gb_test"
	"gorm.io/gorm"
	"testing"
	"time"
)

func Test_doFormatTag(t *testing.T) {
	type Model struct {
		ID        uint `gorm:"primarykey"`
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	type Test struct {
		T1 string `json:"t1" gorm:"not null;size:200"`
		T2 string `json:"t2"`
	}

	type User struct {
		Model      `dbcache:"true"`
		Username   string `json:"username" gorm:"index"`
		Password   string `json:"password"`
		TestStruct Test   `json:"testStruct"  gorm:"embedded"`
	}

	gbtest.C(t, func(t *gbtest.T) {
		handler := New(cache.New())
		result := handler.doFormat(&User{
			Model:      Model{ID: 1},
			Username:   "ghostbb",
			Password:   "123456",
			TestStruct: Test{T1: "test1", T2: "test2"},
		})

		t.Assert(result, `{"ID":1,"username":"ghostbb","password":"123456","testStruct":{"t1":"test1","t2":"test2"}}`)
	})
}

func Test_doFormatNoTag(t *testing.T) {
	type Model struct {
		ID        uint `gorm:"primarykey"`
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	type Test struct {
		T1 string `json:"t1" gorm:"not null;size:200"`
		T2 string `json:"t2"`
	}

	type User struct {
		Model
		Username   string `json:"username" gorm:"index"`
		Password   string `json:"password"`
		TestStruct Test   `json:"testStruct"  gorm:"embedded"`
	}

	gbtest.C(t, func(t *gbtest.T) {
		handler := New(cache.New())
		result := handler.doFormat(&User{
			Model:      Model{ID: 1},
			Username:   "ghostbb",
			Password:   "123456",
			TestStruct: Test{T1: "test1", T2: "test2"},
		})

		t.Assert(result, `{"username":"ghostbb","password":"123456","testStruct":{"t1":"test1","t2":"test2"}}`)
	})
}

func Test_doFormatEmbedded(t *testing.T) {
	type Test struct {
		T1 string `json:"t1" gorm:"not null;size:200"`
		T2 string `json:"t2"`
	}

	type User struct {
		gorm.Model
		Username   string `json:"username" gorm:"index"`
		Password   string `json:"password"`
		TestStruct Test   `json:"testStruct"  gorm:"embedded"`
	}

	gbtest.C(t, func(t *gbtest.T) {
		handler := New(cache.New())
		result := handler.doFormat(&User{
			Model:      gorm.Model{ID: 1},
			Username:   "ghostbb",
			Password:   "123456",
			TestStruct: Test{T1: "test1", T2: "test2"},
		})

		t.Assert(result, User{
			Model:      gorm.Model{ID: 1},
			Username:   "ghostbb",
			Password:   "123456",
			TestStruct: Test{T1: "test1", T2: "test2"},
		})
	})
}

func Test_doFormatNoEmbedded(t *testing.T) {
	type Test struct {
		T1 string `json:"t1" gorm:"not null;size:200"`
		T2 string `json:"t2"`
	}

	type User struct {
		gorm.Model
		Username   string `json:"username" gorm:"index"`
		Password   string `json:"password"`
		TestStruct Test   `json:"testStruct"`
	}

	gbtest.C(t, func(t *gbtest.T) {
		handler := New(cache.New())
		result := handler.doFormat(&User{
			Model:      gorm.Model{ID: 1},
			Username:   "ghostbb",
			Password:   "123456",
			TestStruct: Test{T1: "test1", T2: "test2"},
		})

		t.Assert(result, `{"ID":1,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"username":"ghostbb","password":"123456"}`)
	})
}

func Test_doFormatNoEmbeddedNoTag(t *testing.T) {
	type Model struct {
		ID        uint `gorm:"primarykey"`
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	type Test struct {
		T1 string `json:"t1" gorm:"not null;size:200"`
		T2 string `json:"t2"`
	}

	type User struct {
		Model
		Username   string `json:"username" gorm:"index"`
		Password   string `json:"password"`
		TestStruct Test   `json:"testStruct"`
	}

	gbtest.C(t, func(t *gbtest.T) {
		handler := New(cache.New())
		result := handler.doFormat(&User{
			Model:      Model{ID: 1},
			Username:   "ghostbb",
			Password:   "123456",
			TestStruct: Test{T1: "test1", T2: "test2"},
		})

		t.Assert(result, `{"username":"ghostbb","password":"123456"}`)
	})
}

func Test_doFormatMany(t *testing.T) {
	type Test struct {
		T1 string `json:"t1" gorm:"not null;size:200"`
		T2 string `json:"t2"`
	}

	type User struct {
		gorm.Model `dbcache:"true"`
		Username   string `json:"username" gorm:"index"`
		Password   string `json:"password"`
		TestStruct Test   `json:"testStruct"  gorm:"embedded"`
	}

	type TestRes struct {
		T2 string `json:"t2"`
	}

	type UserRes struct {
		ID         uint
		Username   string  `json:"username" gorm:"index"`
		Password   string  `json:"password"`
		TestStruct TestRes `json:"testStruct"  gorm:"embedded"`
	}

	gbtest.C(t, func(t *gbtest.T) {
		users := make([]*User, 0)
		output := make([]UserRes, 0)
		handler := New(cache.New())

		for i := 0; i < 3; i++ {
			users = append(users, &User{
				Model:      gorm.Model{ID: 1},
				Username:   "ghostbb",
				Password:   "123456",
				TestStruct: Test{T1: "", T2: "test2"},
			})
			output = append(output, UserRes{
				ID:       1,
				Username: "ghostbb",
				Password: "123456",
				TestStruct: TestRes{
					T2: "test2",
				},
			})
		}

		t.Assert(handler.doFormat(users), `[{"ID":1,"username":"ghostbb","password":"123456","testStruct":{"t2":"test2"}},{"ID":1,"username":"ghostbb","password":"123456","testStruct":{"t2":"test2"}},{"ID":1,"username":"ghostbb","password":"123456","testStruct":{"t2":"test2"}}]`)
	})
}

func Test_doFormatRecursion(t *testing.T) {
	type Info struct {
		Phone   string
		Country string
	}
	type Friend struct {
		RealName string
		Info     *Info   `gorm:"embedded"`
		Friend1  *Friend `gorm:"embedded"`
		Friend2  *Friend `gorm:"embedded"`
		Friend3  *Friend
	}
	type User struct {
		Username string
		Password string
		Info     Info   `gorm:"embedded"`
		Info2    *Info  `gorm:"embedded"`
		Info3    *Info  `gorm:"embedded"`
		Info4    Info   `gorm:"embedded"`
		Friend   Friend `dbcache:"true"`
	}
	gbtest.C(t, func(t *gbtest.T) {
		handler := New(cache.New())
		result := handler.doFormat(&User{
			Username: "ghostbb",
			Password: "123456",
			Info:     Info{Phone: "0900123456"},
			Info2:    nil,
			Info3:    &Info{Phone: "", Country: "taiwan"},
			Info4:    Info{Phone: ""},
			Friend: Friend{
				RealName: "golang",
				Info:     nil,
				Friend1: &Friend{
					RealName: "goland",
					Info:     &Info{Phone: "123456", Country: ""},
					Friend2: &Friend{
						RealName: "",
						Info:     &Info{Phone: ""},
					},
					Friend3: &Friend{
						RealName: "ghostbb123456",
						Info:     &Info{Phone: "123456", Country: ""},
					},
				},
				Friend2: nil,
			},
		})
		t.Assert(result, `{"Username":"ghostbb","Password":"123456","Info":{"Phone":"0900123456"},"Info3":{"Country":"taiwan"},"Friend":{"RealName":"golang","Friend1":{"RealName":"goland","Info":{"Phone":"123456"}}}}`)
		// {
		//     Username: "ghostbb",
		//     Password: "123456",
		//     Info:     {
		//         Phone: "0900123456",
		//     },
		//     Info3:    {
		//         Country: "taiwan",
		//     },
		//     Friend:   {
		//         RealName: "golang",
		//         Friend1:  {
		//             RealName: "goland",
		//             Info:     {
		//                 Phone: "123456",
		//             },
		//         },
		//     },
		// }
	})
}
