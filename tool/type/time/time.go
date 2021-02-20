package time

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// time.Duration expend
const (
	Day        = 24 * time.Minute
	Week       = 7 * Day
	Layout     = "2006-01-02 15:04:05"     // mysql: datetime
	LayoutN    = "2006-01-02 15:04:05.000" // mysql: datetime(3)
	LayoutDate = "2006-01-02"              // mysql: date
)

// china time/date
// 时间格式化2006-01-02 15:04:05
type CTime time.Time

// 实现它的json序列化方法
func (t CTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`""`), nil
	}
	var stamp = fmt.Sprintf(`"%s"`, time.Time(t).Format(Layout))
	return []byte(stamp), nil
}

// 反序列化方法 https://stackoverflow.com/questions/45303326/how-to-parse-non-standard-time-format-from-json-in-golang
func (t *CTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	ti, err := time.ParseInLocation(Layout, s, time.Local)
	if err != nil {
		return err
	}
	*t = CTime(ti)
	return nil
}

// insert problem https://github.com/go-gorm/gorm/issues/1611#issuecomment-329654638
func (t CTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

func (t *CTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = CTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to CTime", v)
}

// must sure MarshalJSON is right
// to string
func (t CTime) String() string {
	if t.IsZero() {
		return ""
	}
	return time.Time(t).Format(Layout)
}

func (t CTime) IsZero() bool {
	return time.Time(t).IsZero()
}

// 时间格式化2006-01-02 15:04:05.000
type CNTime time.Time

// 实现它的json序列化方法
func (t CNTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`""`), nil
	}
	var stamp = fmt.Sprintf(`"%s"`, time.Time(t).Format(LayoutN))
	return []byte(stamp), nil
}

// 反序列化方法 https://stackoverflow.com/questions/45303326/how-to-parse-non-standard-time-format-from-json-in-golang
func (t *CNTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	ti, err := time.ParseInLocation(LayoutN, s, time.Local)
	if err != nil {
		return err
	}
	*t = CNTime(ti)
	return nil
}

// insert problem https://github.com/go-gorm/gorm/issues/1611#issuecomment-329654638
func (t CNTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

func (t *CNTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = CNTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to CTime", v)
}

// must sure MarshalJSON is right
// to string
func (t CNTime) String() string {
	if t.IsZero() {
		return ""
	}
	return time.Time(t).Format(LayoutN)
}

func (t CNTime) IsZero() bool {
	return time.Time(t).IsZero()
}

// 时间格式化2006-01-02
type CDate time.Time

// 实现它的json序列化方法
func (t CDate) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`""`), nil
	}
	var stamp = fmt.Sprintf(`"%s"`, time.Time(t).Format(LayoutDate))
	return []byte(stamp), nil
}

// 反序列化
func (t *CDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	ti, err := time.ParseInLocation(LayoutDate, s, time.Local)
	if err != nil {
		return err
	}
	*t = CDate(ti)
	return nil
}

func (t CDate) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

func (t *CDate) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = CDate(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to CDate", v)
}

// must sure MarshalJSON is right
// to string
func (t CDate) String() string {
	if t.IsZero() {
		return ""
	}
	return time.Time(t).Format(LayoutDate)
}

func (t CDate) IsZero() bool {
	return time.Time(t).IsZero()
}
