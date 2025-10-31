package util

import (
	"time"
)

const (
	ProgramFmt = "2006-01-02 15:04:05.000"

	DefaultFmt = "2006-01-02 15:04:05"

	YMDFmt = "2006-01-02"

	YMNoSepFmt = "200601"

	YMDNoSepFmt = "20060102"

	TimeFmt = "15:04:05"
)

var cnLoc, _ = time.LoadLocation("Asia/Shanghai")

func UtcNow() time.Time {
	return time.Now().UTC()
}

func Now() time.Time {
	return CnNow()
}

func CnNow() time.Time {
	return time.Now().In(cnLoc)
}

func Parse2MillsInCn(val string) (int64, error) {
	return parse2ts(val, cnLoc, 13)
}

func Parse2SecInCn(val string) (int64, error) {
	return parse2ts(val, cnLoc, 10)
}

func _parseWithLoc(val string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(DefaultFmt, val, loc)
}

func parse2ts(val string, loc *time.Location, num uint8) (int64, error) {
	t, err := _parseWithLoc(val, loc)
	if err != nil {
		return 0, err
	}

	if num == 10 {
		return t.Unix(), nil
	} else if num == 13 {
		return t.UnixMilli(), nil
	} else {
		return t.UnixNano(), nil
	}
}
