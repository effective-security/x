package format_test

import (
	"testing"
	"time"

	"github.com/effective-security/x/format"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseStringTime(t *testing.T) {
	t1 := format.ParseStringTime("2023-01-01T12:00:00Z")
	require.False(t, t1.IsZero())
	assert.Equal(t, "2023-01-01T12:00:00Z", t1.UTC().Format(time.RFC3339))

	t2 := format.ParseStringTime("")
	assert.True(t, t2.IsZero())

	t3 := format.ParseStringTime("2023-01-02 12:01:01")
	assert.Equal(t, "2023-01-02T12:01:01Z", t3.UTC().Format(time.RFC3339))

	t4 := format.ParseStringTime("2023-01-02")
	assert.Equal(t, "2023-01-02T00:00:00Z", t4.UTC().Format(time.RFC3339))

	t5 := format.ParseStringTime("2023-01-02T12:01:02+03:00")
	assert.Equal(t, "2023-01-02T09:01:02Z", t5.UTC().Format(time.RFC3339))

	t6 := format.ParseStringTime(t5.UTC().Format(time.RFC3339Nano))
	assert.Equal(t, "2023-01-02T09:01:02Z", t6.UTC().Format(time.RFC3339))

	t7 := format.ParseStringTime("2024-02-12T17:07:02.123Z")
	assert.Equal(t, "2024-02-12T17:07:02Z", t7.UTC().Format(time.RFC3339))
	t8 := format.ParseStringTime("2024-02-12T17:07:02.123+07:00")
	assert.Equal(t, "2024-02-12T17:07:02Z", t8.UTC().Format(time.RFC3339))
}

func TestParseTime(t *testing.T) {
	t1 := format.ParseTime("2023-01-01T12:00:00Z")
	require.False(t, t1.IsZero())
	assert.Equal(t, "2023-01-01T12:00:00Z", t1.UTC().Format(time.RFC3339))

	t1u := format.ParseTime(t1.Unix())
	assert.Equal(t, "2023-01-01T12:00:00Z", t1u.UTC().Format(time.RFC3339))

	t2 := format.ParseTime(time.Now())
	assert.False(t, t2.IsZero())

	t3 := format.ParseTime(nil)
	assert.True(t, t3.IsZero())
}

func TestFormatTime(t *testing.T) {
	defer func() {
		format.DefaultTimePrintFormat = "UTC"
		format.NowFunc = time.Now
	}()

	format.DefaultTimePrintFormat = "UTC"
	assert.Equal(t, "2023-01-01T12:00:00Z", format.Time("2023-01-01T12:00:00Z"))
	assert.Equal(t, "", format.Time(nil))
	assert.Equal(t, "never", format.Time(time.Time{}))
	assert.Equal(t, "2023-01-01 12:00:00", format.Time(time.Date(2023, 1, 1, 12, 0, 0, 0, time.Local)))
	assert.Equal(t, "2023-01-01 12:00:00", format.Time(time.Date(2023, 1, 1, 12, 0, 0, 0, time.Local).Unix()))

	format.DefaultTimePrintFormat = "Local"
	//	assert.Equal(t, "2023-01-01 12:00:00", format.Time("2023-01-01T12:00:00Z"))
	assert.Equal(t, "2023-01-01 12:00:00", format.Time(time.Date(2023, 1, 1, 12, 0, 0, 0, time.Local)))
	assert.Equal(t, "2023-01-01 12:00:00", format.Time(time.Date(2023, 1, 1, 12, 0, 0, 0, time.Local).Unix()))

	format.NowFunc = func() time.Time { return time.Date(2024, 3, 4, 12, 59, 5, 0, time.UTC) }
	format.DefaultTimePrintFormat = "Ago"
	assert.Equal(t, "428 days ago", format.Time("2023-01-01T12:00:00Z"))
	assert.Equal(t, "never", format.Time(time.Time{}))
	assert.Equal(t, "just now", format.Time(time.Date(2024, 3, 4, 12, 59, 5, 0, time.UTC)))
	assert.Equal(t, "1 mins ago", format.Time(time.Date(2024, 3, 4, 12, 58, 5, 0, time.UTC)))
	assert.Equal(t, "3 hours ago", format.Time(time.Date(2024, 3, 4, 9, 0, 0, 0, time.UTC)))
	assert.Equal(t, "an hour ago", format.Time(time.Date(2024, 3, 4, 11, 0, 0, 0, time.UTC)))
	assert.Equal(t, "6 hours ago", format.Time(time.Date(2024, 3, 4, 6, 58, 4, 0, time.UTC)))
	assert.Equal(t, "yesterday", format.Time(time.Date(2024, 3, 3, 6, 58, 4, 0, time.UTC)))
	assert.Equal(t, "427 days ago", format.Time(time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)))
	assert.Equal(t, "in 1 days", format.Time(time.Date(2024, 3, 6, 12, 0, 0, 0, time.UTC).Unix()))
	assert.Equal(t, "in 62 days", format.Time(time.Date(2024, 5, 6, 12, 0, 0, 0, time.UTC).Unix()))
	assert.Equal(t, "just now", format.Time(time.Date(2024, 3, 4, 12, 59, 5, 3, time.UTC).Unix()))
	assert.Equal(t, "in a few seconds", format.Time(time.Date(2024, 3, 4, 12, 59, 6, 3, time.UTC).Unix()))
	assert.Equal(t, "in 3 minutes", format.Time(time.Date(2024, 3, 4, 13, 02, 6, 3, time.UTC).Unix()))
	assert.Equal(t, "in 1 hours", format.Time(time.Date(2024, 3, 4, 14, 01, 6, 3, time.UTC).Unix()))
}

func TestLocalTime(t *testing.T) {
	t1 := format.LocalTime("2023-01-01T12:00:00Z")
	assert.NotEmpty(t, t1)

	t2 := format.LocalTime(time.Time{})
	assert.NotEmpty(t, t2)
}

func TestTimeAgo(t *testing.T) {
	now := time.Now()
	format.NowFunc = func() time.Time { return now }

	t1 := format.TimeAgo(now.Add(-time.Minute))
	assert.Equal(t, "1 mins ago", t1)

	t2 := format.TimeAgo(now.Add(-48 * time.Hour))
	assert.Equal(t, "2 days ago", t2)

	t3 := format.TimeAgo(now.Add(time.Minute))
	assert.Equal(t, "in 1 minutes", t3)
}

func TestTimeAg2(t *testing.T) {
	format.NowFunc = func() time.Time {
		return time.Date(2024, 5, 6, 7, 8, 9, 0, time.UTC)
	}
	defer func() {
		format.NowFunc = time.Now
	}()

	assert.Equal(t, "just now", format.TimeAgo("2024-05-06 07:08:09"))
	assert.Equal(t, "just now", format.TimeAgo(time.Date(2024, 5, 6, 7, 8, 9, 0, time.UTC)))
	assert.Equal(t, "7 hours ago", format.TimeAgo("2024-05-06"))
	assert.Equal(t, "7 hours ago", format.TimeAgo(time.Date(2024, 5, 6, 0, 0, 0, 0, time.UTC)))
	assert.Equal(t, "just now", format.TimeAgo(format.NowFunc()))
	assert.Equal(t, "in 2 minutes", format.TimeAgo(time.Date(2024, 5, 6, 7, 10, 9, 0, time.UTC)))
	assert.Equal(t, "2 hours ago", format.TimeAgo(time.Date(2024, 5, 6, 4, 10, 9, 0, time.UTC)))
	assert.Equal(t, "2 days ago", format.TimeAgo(time.Date(2024, 5, 4, 4, 10, 9, 0, time.UTC)))
}

func TestFormatTimesElapsed(t *testing.T) {
	from := "2023-01-01T12:00:00Z"
	to := "2023-01-01T14:00:00Z"

	created, updated, elapsed := format.TimesElapsed(from, to)
	assert.Equal(t, "2023-01-01 12:00:00", created)
	assert.Equal(t, "2023-01-01 14:00:00", updated)
	assert.Equal(t, "2h0m0s", elapsed)
}
