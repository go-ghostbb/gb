package gbutil_test

import (
	"bytes"
	gbtype "ghostbb.io/gb/container/gb_type"
	"ghostbb.io/gb/frame/g"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	gbstr "ghostbb.io/gb/text/gb_str"
	gbmeta "ghostbb.io/gb/util/gb_meta"
	gbutil "ghostbb.io/gb/util/gb_util"
	"testing"
)

func Test_Dump(t *testing.T) {
	type CommonReq struct {
		AppId      int64  `json:"appId" v:"required" in:"path" des:"应用Id" sum:"应用Id Summary"`
		ResourceId string `json:"resourceId" in:"query" des:"资源Id" sum:"资源Id Summary"`
	}
	type SetSpecInfo struct {
		StorageType string   `v:"required|in:CLOUD_PREMIUM,CLOUD_SSD,CLOUD_HSSD" des:"StorageType"`
		Shards      int32    `des:"shards 分片数" sum:"Shards Summary"`
		Params      []string `des:"默认参数(json 串-ClickHouseParams)" sum:"Params Summary"`
	}
	type CreateResourceReq struct {
		CommonReq
		gbmeta.Meta `path:"/CreateResourceReq" method:"POST" tags:"default" sum:"CreateResourceReq sum"`
		Name        string
		CreatedAt   *gbtime.Time
		SetMap      map[string]*SetSpecInfo
		SetSlice    []SetSpecInfo
		internal    string
	}
	req := &CreateResourceReq{
		CommonReq: CommonReq{
			AppId:      12345678,
			ResourceId: "tdchqy-xxx",
		},
		Name:      "john",
		CreatedAt: gbtime.Now(),
		SetMap: map[string]*SetSpecInfo{
			"test1": {
				StorageType: "ssd",
				Shards:      2,
				Params:      []string{"a", "b", "c"},
			},
			"test2": {
				StorageType: "hssd",
				Shards:      10,
				Params:      []string{},
			},
		},
		SetSlice: []SetSpecInfo{
			{
				StorageType: "hssd",
				Shards:      10,
				Params:      []string{"h"},
			},
		},
	}
	gbtest.C(t, func(t *gbtest.T) {
		gbutil.Dump(map[int]int{
			100: 100,
		})
		gbutil.Dump(req)
		gbutil.Dump(true, false)
		gbutil.Dump(make(chan int))
		gbutil.Dump(func() {})
		gbutil.Dump(nil)
		gbutil.Dump(gbtype.NewInt(1))
	})
}

func Test_Dump_Map(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		buffer := bytes.NewBuffer(nil)
		m := g.Map{
			"k1": g.Map{
				"k2": "v2",
			},
		}
		gbutil.DumpTo(buffer, m, gbutil.DumpOption{})
		t.Assert(buffer.String(), `{
    "k1": {
        "k2": "v2",
    },
}`)
	})
}

func TestDumpWithType(t *testing.T) {
	type CommonReq struct {
		AppId      int64  `json:"appId" v:"required" in:"path" des:"应用Id" sum:"应用Id Summary"`
		ResourceId string `json:"resourceId" in:"query" des:"资源Id" sum:"资源Id Summary"`
	}
	type SetSpecInfo struct {
		StorageType string   `v:"required|in:CLOUD_PREMIUM,CLOUD_SSD,CLOUD_HSSD" des:"StorageType"`
		Shards      int32    `des:"shards 分片数" sum:"Shards Summary"`
		Params      []string `des:"默认参数(json 串-ClickHouseParams)" sum:"Params Summary"`
	}
	type CreateResourceReq struct {
		CommonReq
		gbmeta.Meta `path:"/CreateResourceReq" method:"POST" tags:"default" sum:"CreateResourceReq sum"`
		Name        string
		CreatedAt   *gbtime.Time
		SetMap      map[string]*SetSpecInfo `v:"required" des:"配置Map"`
		SetSlice    []SetSpecInfo           `v:"required" des:"配置Slice"`
		internal    string
	}
	req := &CreateResourceReq{
		CommonReq: CommonReq{
			AppId:      12345678,
			ResourceId: "tdchqy-xxx",
		},
		Name:      "john",
		CreatedAt: gbtime.Now(),
		SetMap: map[string]*SetSpecInfo{
			"test1": {
				StorageType: "ssd",
				Shards:      2,
				Params:      []string{"a", "b", "c"},
			},
			"test2": {
				StorageType: "hssd",
				Shards:      10,
				Params:      []string{},
			},
		},
		SetSlice: []SetSpecInfo{
			{
				StorageType: "hssd",
				Shards:      10,
				Params:      []string{"h"},
			},
		},
	}
	gbtest.C(t, func(t *gbtest.T) {
		gbutil.DumpWithType(map[int]int{
			100: 100,
		})
		gbutil.DumpWithType(req)
		gbutil.DumpWithType([][]byte{[]byte("hello")})
	})
}

func Test_Dump_Slashes(t *testing.T) {
	type Req struct {
		Content string
	}
	req := &Req{
		Content: `{"name":"john", "age":18}`,
	}
	gbtest.C(t, func(t *gbtest.T) {
		gbutil.Dump(req)
		gbutil.Dump(req.Content)

		gbutil.DumpWithType(req)
		gbutil.DumpWithType(req.Content)
	})
}

// https://github.com/gogf/gf/issues/1661
func Test_Dump_Issue1661(t *testing.T) {
	type B struct {
		ba int
		bb string
	}
	type A struct {
		aa int
		ab string
		cc []B
	}
	gbtest.C(t, func(t *gbtest.T) {
		var q1 []A
		var q2 []A
		q2 = make([]A, 0)
		q1 = []A{{aa: 1, ab: "1", cc: []B{{ba: 1}, {ba: 2}, {ba: 3}}}, {aa: 2, ab: "2", cc: []B{{ba: 1}, {ba: 2}, {ba: 3}}}}
		for _, q1v := range q1 {
			x := []string{"11", "22"}
			for _, iv2 := range x {
				ls := q1v
				for i := range ls.cc {
					sj := iv2
					ls.cc[i].bb = sj
				}
				q2 = append(q2, ls)
			}
		}
		buffer := bytes.NewBuffer(nil)
		gbutil.DumpTo(buffer, q2, gbutil.DumpOption{})
		t.Assert(buffer.String(), `[
    {
        aa: 1,
        ab: "1",
        cc: [
            {
                ba: 1,
                bb: "22",
            },
            {
                ba: 2,
                bb: "22",
            },
            {
                ba: 3,
                bb: "22",
            },
        ],
    },
    {
        aa: 1,
        ab: "1",
        cc: [
            {
                ba: 1,
                bb: "22",
            },
            {
                ba: 2,
                bb: "22",
            },
            {
                ba: 3,
                bb: "22",
            },
        ],
    },
    {
        aa: 2,
        ab: "2",
        cc: [
            {
                ba: 1,
                bb: "22",
            },
            {
                ba: 2,
                bb: "22",
            },
            {
                ba: 3,
                bb: "22",
            },
        ],
    },
    {
        aa: 2,
        ab: "2",
        cc: [
            {
                ba: 1,
                bb: "22",
            },
            {
                ba: 2,
                bb: "22",
            },
            {
                ba: 3,
                bb: "22",
            },
        ],
    },
]`)
	})
}

func Test_Dump_Cycle_Attribute(t *testing.T) {
	type Abc struct {
		ab int
		cd *Abc
	}
	abc := Abc{ab: 3}
	abc.cd = &abc
	gbtest.C(t, func(t *gbtest.T) {
		buffer := bytes.NewBuffer(nil)
		g.DumpTo(buffer, abc, gbutil.DumpOption{})
		t.Assert(gbstr.Contains(buffer.String(), "cycle"), true)
	})
}

func Test_DumpJson(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		var jsonContent = `{"a":1,"b":2}`
		gbutil.DumpJson(jsonContent)
	})
}
