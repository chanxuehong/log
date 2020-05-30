package log

import (
	"testing"
	"time"
)

func BenchmarkFormatTime(b *testing.B) {
	now := time.Now()
	var str string
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str = FormatTime(now)
	}
	_ = str
}

func BenchmarkStdTimeFormat(b *testing.B) {
	now := time.Now()
	var str string
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str = now.Format(TimeFormatLayout)
	}
	_ = str
}
