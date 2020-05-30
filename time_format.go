package log

import "time"

const TimeFormatLayout = "2006-01-02 15:04:05.000"

const (
	digits01 = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
	digits10 = "0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999"
)

// FormatTime format time.Time as layout "2006-01-02 15:04:05.000".
func FormatTime(t time.Time) string {
	year, month, day := t.Date()
	if year < 1 || year > 9999 {
		return t.Format(TimeFormatLayout)
	}
	hour, min, sec := t.Clock()
	millisecond := t.Nanosecond() / 1e6

	year100 := year / 100
	year1 := year % 100
	millisecond100 := millisecond / 100
	millisecond1 := millisecond % 100

	var result [23]byte
	result[0], result[1], result[2], result[3] = digits10[year100], digits01[year100], digits10[year1], digits01[year1]
	result[4] = '-'
	result[5], result[6] = digits10[month], digits01[month]
	result[7] = '-'
	result[8], result[9] = digits10[day], digits01[day]
	result[10] = ' '
	result[11], result[12] = digits10[hour], digits01[hour]
	result[13] = ':'
	result[14], result[15] = digits10[min], digits01[min]
	result[16] = ':'
	result[17], result[18] = digits10[sec], digits01[sec]
	result[19] = '.'
	result[20], result[21], result[22] = digits01[millisecond100], digits10[millisecond1], digits01[millisecond1]
	return string(result[:])
}
