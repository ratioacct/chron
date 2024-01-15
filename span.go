package chron

import (
	"encoding/json"
	"fmt"

	"github.com/ratioacct/chron/dura"
	"time"
)

type Span interface {
	Start() Chron
	End() Chron
	Duration() dura.Time
	Comparer
}

type Comparer interface {
	Before(Span) bool
	After(Span) bool
	Contains(Span) bool
}

type Interval struct {
	start Chron
	end   Chron
	d     dura.Time
}

func NewInterval(start Chron, d dura.Time) *Interval {
	if d.Duration() == time.Nanosecond && d.Years() == 0 && d.Months() == 0 && d.Days() == 0 {
		return &Interval{start: start, end: start, d: dura.Nano}
	}
	return &Interval{
		start: start,
		end:   start.Increment(d).Decrement(dura.Nano),
		d:     d,
	}
}

func (s Interval) Contains(t Span) bool {
	return !s.Before(t) && !s.After(t)
}

func (s Interval) Before(t Span) bool {
	return s.End().AsTime().Before(t.Start().AsTime())
}

func (s Interval) After(t Span) bool {
	return s.Start().AsTime().After(t.End().AsTime())
}

func (s Interval) Duration() dura.Time {
	return s.d
}

func (s Interval) Start() Chron {
	return s.start
}

func (s Interval) End() Chron {
	return s.end
}

func (s Interval) String() string {
	return fmt.Sprintf("start:%s, end:%s, len:%s", s.start, s.end, s.d)
}

func (s Interval) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"start":"%s","end":"%s","len":"%s"}`, s.start, s.end, s.d)), nil
}

func (s *Interval) UnmarshalJSON(b []byte) error {
	var start, end, len string
	if _, err := fmt.Sscanf(string(b), `{"start":"%s","end":"%s","len":"%s"}`, &start, &end, &len); err != nil {
		return fmt.Errorf("unmarshalling chron.Interval: %w", err)
	}
	type tmp struct {
		Start Chron     `json:"start"`
		End   Chron     `json:"end"`
		Len   dura.Time `json:"len"`
	}
	var t tmp
	if err := json.Unmarshal(b, &t); err != nil {
		return fmt.Errorf("unmarshalling chron.Interval: %w", err)
	}
	s.start = t.Start
	s.end = t.End
	s.d = t.Len
	return nil
}
