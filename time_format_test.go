package log

import (
	"testing"
	"time"
)

func TestFormatTime(t *testing.T) {
	now := time.Now()
	have := FormatTime(now)
	want := now.Format(TimeFormatLayout)
	if len(want) < 23 {
		want += TimeFormatLayout[len(want):]
	}
	if have != want {
		t.Errorf("have:%s, want:%s", have, want)
		return
	}
}
