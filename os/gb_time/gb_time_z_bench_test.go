package gbtime_test

import (
	gbtime "ghostbb.io/gb/os/gb_time"
	"testing"
	"time"
)

func Benchmark_Timestamp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gbtime.Timestamp()
	}
}

func Benchmark_TimestampMilli(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gbtime.TimestampMilli()
	}
}

func Benchmark_TimestampMicro(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gbtime.TimestampMicro()
	}
}

func Benchmark_TimestampNano(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gbtime.TimestampNano()
	}
}

func Benchmark_StrToTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gbtime.StrToTime("2018-02-09T20:46:17.897Z")
	}
}

func Benchmark_StrToTime_Format(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gbtime.StrToTime("2018-02-09 20:46:17.897", "Y-m-d H:i:su")
	}
}

func Benchmark_StrToTime_Layout(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gbtime.StrToTimeLayout("2018-02-09T20:46:17.897Z", time.RFC3339)
	}
}

func Benchmark_ParseTimeFromContent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gbtime.ParseTimeFromContent("2018-02-09T20:46:17.897Z")
	}
}

func Benchmark_NewFromTimeStamp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gbtime.NewFromTimeStamp(1542674930)
	}
}

func Benchmark_Date(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gbtime.Date()
	}
}

func Benchmark_Datetime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gbtime.Datetime()
	}
}
