package chron

import (
	"time"

	"database/sql/driver"
	"fmt"
	"github.com/ratioacct/chron/dura"
	"reflect"
	"strings"
)

type Second struct {
	time.Time
}

func NewSecond(year int, month time.Month, day, hour, min, sec int) Second {
	return Second{time.Date(year, month, day, hour, min, sec, 0, time.UTC)}
}

func ThisSecond() Second {
	return Now().AsSecond()
}

func SecondOf(t time.Time) Second {
	t = t.UTC()
	return Second{t.Truncate(time.Second)}
}

func (s Second) AsYear() Year      { return YearOf(s.Time) }
func (s Second) AsMonth() Month    { return MonthOf(s.Time) }
func (s Second) AsDay() Day        { return DayOf(s.Time) }
func (s Second) AsHour() Hour      { return HourOf(s.Time) }
func (s Second) AsMinute() Minute  { return MinuteOf(s.Time) }
func (s Second) AsSecond() Second  { return s }
func (s Second) AsMilli() Milli    { return MilliOf(s.Time) }
func (s Second) AsMicro() Micro    { return MicroOf(s.Time) }
func (s Second) AsChron() Chron    { return TimeOf(s.Time) }
func (s Second) AsTime() time.Time { return s.Time }

func (s Second) Increment(l dura.Time) Chron {
	return Chron{s.AddDate(l.Years(), l.Months(), l.Days()).Add(l.Duration())}
}

func (s Second) Decrement(l dura.Time) Chron {
	return Chron{s.AddDate(-1*l.Years(), -1*l.Months(), -1*l.Days()).Add(-1 * l.Duration())}
}

func (s Second) AddN(n int) Second {
	return Second{s.Add(time.Duration(int(time.Second) * n))}
}

// / span.Time implementation
func (s Second) Start() Chron {
	return s.AsChron()
}

func (s Second) End() Chron {
	return s.AddN(1).Decrement(dura.Nano)
}

func (s Second) Contains(t Span) bool {
	return !s.Before(t) && !s.After(t)
}

func (s Second) Before(t Span) bool {
	return s.End().AsTime().Before(t.Start().AsTime())
}

func (s Second) After(t Span) bool {
	return s.Start().AsTime().After(t.End().AsTime())
}

func (s Second) Duration() dura.Time {
	return dura.Second
}

func (s Second) AddYears(y int) Second {
	return s.Increment(dura.Years(y)).AsSecond()
}

func (s Second) AddMonths(m int) Second {
	return s.Increment(dura.Months(m)).AsSecond()
}

func (s Second) AddDays(d int) Second {
	return s.Increment(dura.Days(d)).AsSecond()
}

func (s Second) AddHours(h int) Second {
	return s.Increment(dura.Hours(h)).AsSecond()
}

func (s Second) AddMinutes(m int) Second {
	return s.Increment(dura.Mins(m)).AsSecond()
}

func (s Second) AddSeconds(secs int) Second {
	return s.AddN(secs)
}

func (s Second) AddMillis(m int) Milli {
	return s.AsMilli().AddN(m)
}

func (s Second) AddMicros(m int) Micro {
	return s.AsMicro().AddN(m)
}

func (s Second) AddNanos(n int) Chron {
	return s.AsChron().AddN(n)
}

func (s *Second) Scan(value interface{}) error {
	if value == nil {
		*s = ZeroValue().AsSecond()
		return nil
	}
	if t, ok := value.(time.Time); ok {
		*s = SecondOf(t)
		return nil
	}
	return fmt.Errorf("unsupported Scan, storing %s into type *chron.Day", reflect.TypeOf(value))
}

func (s Second) Value() (driver.Value, error) {
	// todo: error check the range.
	return s.Time, nil
}

func (s *Second) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	st := strings.Trim(string(data), `"`)
	t, err := Parse(st)
	*s = SecondOf(t)
	return err
}
