package gbtime_test

import (
	"ghostbb.io/gb/frame/g"
	gbtime "ghostbb.io/gb/os/gb_time"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
	"time"
)

func Test_TimestampStr(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.AssertGT(len(gbtime.TimestampMilliStr()), 0)
		t.AssertGT(len(gbtime.TimestampMicroStr()), 0)
		t.AssertGT(len(gbtime.TimestampNanoStr()), 0)
	})
}

func Test_Nanosecond(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		nanos := gbtime.TimestampNano()
		timeTemp := time.Unix(0, nanos)
		t.Assert(nanos, timeTemp.UnixNano())
	})
}

func Test_Microsecond(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		micros := gbtime.TimestampMicro()
		timeTemp := time.Unix(0, micros*1e3)
		t.Assert(micros, timeTemp.UnixNano()/1e3)
	})
}

func Test_Millisecond(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		millis := gbtime.TimestampMilli()
		timeTemp := time.Unix(0, millis*1e6)
		t.Assert(millis, timeTemp.UnixNano()/1e6)
	})
}

func Test_Second(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		s := gbtime.Timestamp()
		timeTemp := time.Unix(s, 0)
		t.Assert(s, timeTemp.Unix())
	})
}

func Test_Date(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		t.Assert(gbtime.Date(), time.Now().Format("2006-01-02"))
	})
}

func Test_Datetime(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		datetime := gbtime.Datetime()
		timeTemp, err := gbtime.StrToTime(datetime, "Y-m-d H:i:s")
		if err != nil {
			t.Error("test fail")
		}
		t.Assert(datetime, timeTemp.Time.Format("2006-01-02 15:04:05"))
	})
	gbtest.C(t, func(t *gbtest.T) {
		timeTemp, err := gbtime.StrToTime("")
		t.Assert(err, nil)
		t.AssertLT(timeTemp.Unix(), 0)
		timeTemp, err = gbtime.StrToTime("2006-01")
		t.AssertNE(err, nil)
		t.Assert(timeTemp, nil)
	})
}

func Test_ISO8601(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		iso8601 := gbtime.ISO8601()
		t.Assert(iso8601, gbtime.Now().Format("c"))
	})
}

func Test_RFC822(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		rfc822 := gbtime.RFC822()
		t.Assert(rfc822, gbtime.Now().Format("r"))
	})
}

func Test_StrToTime(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		// Correct datetime string.
		var testDateTimes = []string{
			"2006-01-02 15:04:05",
			"2006/01/02 15:04:05",
			"2006.01.02 15:04:05.000",
			"2006.01.02 - 15:04:05",
			"2006.01.02 15:04:05 +0800 CST",
			"2006-01-02T20:05:06+05:01:01",
			"2006-01-02T14:03:04Z01:01:01",
			"2006-01-02T15:04:05Z",
			"02-jan-2006 15:04:05",
			"02/jan/2006 15:04:05",
			"02.jan.2006 15:04:05",
			"02.jan.2006:15:04:05",
		}

		for _, item := range testDateTimes {
			timeTemp, err := gbtime.StrToTime(item)
			t.AssertNil(err)
			t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), "2006-01-02 15:04:05")
		}

		// Correct date string,.
		var testDates = []string{
			"2006.01.02",
			"2006.01.02 00:00",
			"2006.01.02 00:00:00.000",
		}

		for _, item := range testDates {
			timeTemp, err := gbtime.StrToTime(item)
			t.AssertNil(err)
			t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), "2006-01-02 00:00:00")
		}

		// Correct time string.
		var testTimes = g.MapStrStr{
			"16:12:01":     "15:04:05",
			"16:12:01.789": "15:04:05.000",
		}

		for k, v := range testTimes {
			time1, err := gbtime.StrToTime(k)
			t.AssertNil(err)
			time2, err := time.ParseInLocation(v, k, time.Local)
			t.AssertNil(err)
			t.Assert(time1.Time, time2)
		}

		// formatToStdLayout
		var testDateFormats = []string{
			"Y-m-d H:i:s",
			"\\T\\i\\m\\e Y-m-d H:i:s",
			"Y-m-d H:i:s\\",
			"Y-m-j G:i:s.u",
			"Y-m-j G:i:su",
		}

		var testDateFormatsResult = []string{
			"2007-01-02 15:04:05",
			"Time 2007-01-02 15:04:05",
			"2007-01-02 15:04:05",
			"2007-01-02 15:04:05.000",
			"2007-01-02 15:04:05.000",
		}

		for index, item := range testDateFormats {
			timeTemp, err := gbtime.StrToTime(testDateFormatsResult[index], item)
			if err != nil {
				t.Error("test fail")
			}
			t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05.000"), "2007-01-02 15:04:05.000")
		}

		// 异常日期列表
		var testDatesFail = []string{
			"2006.01",
			"06..02",
		}

		for _, item := range testDatesFail {
			_, err := gbtime.StrToTime(item)
			if err == nil {
				t.Error("test fail")
			}
		}

		// test err
		_, err := gbtime.StrToTime("2006-01-02 15:04:05", "aabbccdd")
		if err == nil {
			t.Error("test fail")
		}
	})
}

func Test_ConvertZone(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		// 现行时间
		nowUTC := time.Now().UTC()
		testZone := "America/Los_Angeles"

		// 转换为洛杉矶时间
		t1, err := gbtime.ConvertZone(nowUTC.Format("2006-01-02 15:04:05"), testZone, "")
		if err != nil {
			t.Error("test fail")
		}

		// 使用洛杉矶时区解析上面转换后的时间
		laStr := t1.Time.Format("2006-01-02 15:04:05")
		loc, err := time.LoadLocation(testZone)
		t2, err := time.ParseInLocation("2006-01-02 15:04:05", laStr, loc)

		// 判断是否与现行时间匹配
		t.Assert(t2.UTC().Unix(), nowUTC.Unix())

	})

	// test err
	gbtest.C(t, func(t *gbtest.T) {
		// 现行时间
		nowUTC := time.Now().UTC()
		// t.Log(nowUTC.Unix())
		testZone := "errZone"

		// 错误时间输入
		_, err := gbtime.ConvertZone(nowUTC.Format("06..02 15:04:05"), testZone, "")
		if err == nil {
			t.Error("test fail")
		}
		// 错误时区输入
		_, err = gbtime.ConvertZone(nowUTC.Format("2006-01-02 15:04:05"), testZone, "")
		if err == nil {
			t.Error("test fail")
		}
		// 错误时区输入
		_, err = gbtime.ConvertZone(nowUTC.Format("2006-01-02 15:04:05"), testZone, testZone)
		if err == nil {
			t.Error("test fail")
		}
	})
}

func Test_ParseDuration(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		d, err := gbtime.ParseDuration("1d")
		t.AssertNil(err)
		t.Assert(d.String(), "24h0m0s")
	})
	gbtest.C(t, func(t *gbtest.T) {
		d, err := gbtime.ParseDuration("1d2h3m")
		t.AssertNil(err)
		t.Assert(d.String(), "26h3m0s")
	})
	gbtest.C(t, func(t *gbtest.T) {
		d, err := gbtime.ParseDuration("-1d2h3m")
		t.AssertNil(err)
		t.Assert(d.String(), "-26h3m0s")
	})
	gbtest.C(t, func(t *gbtest.T) {
		d, err := gbtime.ParseDuration("3m")
		t.AssertNil(err)
		t.Assert(d.String(), "3m0s")
	})
	// error
	gbtest.C(t, func(t *gbtest.T) {
		d, err := gbtime.ParseDuration("-1dd2h3m")
		t.AssertNE(err, nil)
		t.Assert(d.String(), "0s")
	})
}

func Test_ParseTimeFromContent(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		timeTemp := gbtime.ParseTimeFromContent("我是中文2006-01-02 15:04:05我也是中文", "Y-m-d H:i:s")
		t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), "2006-01-02 15:04:05")

		timeTemp1 := gbtime.ParseTimeFromContent("我是中文2006-01-02 15:04:05我也是中文")
		t.Assert(timeTemp1.Time.Format("2006-01-02 15:04:05"), "2006-01-02 15:04:05")

		timeTemp2 := gbtime.ParseTimeFromContent("我是中文02.jan.2006 15:04:05我也是中文")
		t.Assert(timeTemp2.Time.Format("2006-01-02 15:04:05"), "2006-01-02 15:04:05")

		// test err
		timeTempErr := gbtime.ParseTimeFromContent("我是中文", "Y-m-d H:i:s")
		if timeTempErr != nil {
			t.Error("test fail")
		}
	})

	gbtest.C(t, func(t *gbtest.T) {
		timeStr := "2021-1-27 9:10:24"
		t.Assert(gbtime.ParseTimeFromContent(timeStr, "Y-n-d g:i:s").String(), "2021-01-27 09:10:24")
	})
}

func Test_FuncCost(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		gbtime.FuncCost(func() {

		})
	})
}
