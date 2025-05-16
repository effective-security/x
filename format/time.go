package format

import (
	"fmt"
	"strings"
	"time"
)

var (
	// DefaultTimePrintFormat is the default time format: UTC, Local, Ago
	DefaultTimePrintFormat = "UTC"

	// DefaultTimeFormatZone is the default time format
	DefaultTimeFormatZone = "2006-01-02T15:04:05.999+07:00"
	// DefaultTimeFormatUTC is the default time format
	DefaultTimeFormatUTC = "2006-01-02T15:04:05.999Z"

	// DefaultTimeTruncate is the default time to truncate as Postgres time precision is default to 6
	// However, JavaScript and AWS accept time milliseconds only, 3 digits, so we truncate to 3
	DefaultTimeTruncate = time.Millisecond

	NowFunc = time.Now
)

// ParseStringTime returns Time from a string
func ParseStringTime(val string) time.Time {
	if val == "" {
		return time.Time{}
	}

	var t time.Time
	switch len(val) {
	case len(DefaultTimeFormatZone):
		t, _ = time.Parse(DefaultTimeFormatZone, val)
	case len(DefaultTimeFormatUTC):
		t, _ = time.Parse(DefaultTimeFormatUTC, val)
	case len(time.RFC3339):
		t, _ = time.Parse(time.RFC3339, val)
	case len(time.DateTime):
		t, _ = time.Parse(time.DateTime, val)
	case len(time.DateOnly):
		t, _ = time.Parse(time.DateOnly, val)
	default:
		t, _ = time.Parse(time.RFC3339Nano, val)
	}
	return t.Truncate(DefaultTimeTruncate)
}

// ParseTime returns time from any type
func ParseTime(val any) time.Time {
	switch t := val.(type) {
	case int64:
		return time.Unix(t, 0)
	case string:
		return ParseStringTime(t)
	case time.Time:
		return t
	case *time.Time:
		if t == nil {
			return time.Time{}
		}
		return *t
	}
	return time.Time{}
}

// Time returns formatted time
func Time(val any) string {
	if strings.EqualFold(DefaultTimePrintFormat, "Local") {
		return LocalTime(val)
	} else if strings.EqualFold(DefaultTimePrintFormat, "Ago") {
		return TimeAgo(val)
	}

	switch t := val.(type) {
	case int64:
		return time.Unix(t, 0).Format("2006-01-02 15:04:05")
	case string:
		return t
	case *time.Time:
		if t == nil {
			return "never"
		}
		return t.Format("2006-01-02 15:04:05")
	case time.Time:
		if t.IsZero() {
			return ""
		}
		return t.Format("2006-01-02 15:04:05")
	}
	return ""
}

// LocalTime returns local time
func LocalTime(val any) string {
	t := ParseTime(val)
	return t.Local().Format("2006-01-02 15:04:05")
}

// TimeAgo returns elapsed time since
func TimeAgo(val any) string {
	if val == nil {
		return "never"
	}
	t := ParseTime(val).UTC()
	now := NowFunc()

	if t.IsZero() {
		return "never"
	}

	if t.After(now) {
		diff := t.Sub(now)
		if diff < time.Minute {
			return "in a few seconds"
		}
		if diff < time.Hour {
			minutes := diff / time.Minute
			return fmt.Sprintf("in %d minutes", minutes)
		}
		if diff < 24*time.Hour {
			hours := diff / time.Hour
			return fmt.Sprintf("in %d hours", hours)
		}
		days := diff / (24 * time.Hour)
		return fmt.Sprintf("in %d days", days)
	}

	ago := now.Sub(t) / time.Second * time.Second
	if ago < time.Minute {
		return "just now"
	}
	if ago < time.Hour {
		minutes := ago / time.Minute
		return fmt.Sprintf("%d mins ago", minutes)
	}
	if ago < 2*time.Hour {
		return "an hour ago"
	}
	if ago < 24*time.Hour {
		hours := ago / time.Hour
		return fmt.Sprintf("%d hours ago", hours)
	}
	if ago < 48*time.Hour {
		return "yesterday"
	}
	days := ago / (24 * time.Hour)
	return fmt.Sprintf("%d days ago", days)
}

// TimesElapsed returns formatted times and elapsed time
func TimesElapsed(from, to any) (created string, updated string, elapsed string) {
	createdTime := ParseTime(from)
	updatedTime := ParseTime(to)

	created = Time(createdTime)
	updated = Time(updatedTime)

	if !createdTime.IsZero() && !updatedTime.IsZero() {
		elapsedTime := updatedTime.Sub(createdTime) / time.Second * time.Second
		elapsed = elapsedTime.String()
	}
	return
}
