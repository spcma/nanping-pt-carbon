package valueobject

import (
	"errors"
	"fmt"
	"time"
)

// TimeRange time range value object
type TimeRange struct {
	StartTime time.Time
	EndTime   time.Time
}

// NewTimeRange creates a time range
func NewTimeRange(startTime, endTime time.Time) (*TimeRange, error) {
	if startTime.After(endTime) {
		return nil, errors.New("start time cannot be later than end time")
	}

	return &TimeRange{
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

// Contains checks if time is within range
func (tr *TimeRange) Contains(t time.Time) bool {
	return !t.Before(tr.StartTime) && !t.After(tr.EndTime)
}

// Overlaps checks if overlaps with another time range
func (tr *TimeRange) Overlaps(other *TimeRange) bool {
	return tr.StartTime.Before(other.EndTime) && other.StartTime.Before(tr.EndTime)
}

// Duration gets the duration
func (tr *TimeRange) Duration() time.Duration {
	return tr.EndTime.Sub(tr.StartTime)
}

// String string representation
func (tr *TimeRange) String() string {
	return fmt.Sprintf("[%s, %s]", tr.StartTime.Format(time.RFC3339), tr.EndTime.Format(time.RFC3339))
}

// Equal checks if equal
func (tr *TimeRange) Equal(other *TimeRange) bool {
	if other == nil {
		return false
	}
	return tr.StartTime.Equal(other.StartTime) && tr.EndTime.Equal(other.EndTime)
}
