package catalog

import (
	"net/http"
	"testing"
	"time"
)

func TestRetryAfterPlainSeconds(t *testing.T) {
	if got := RetryAfter("120", time.Now()); got == nil || *got != 120 {
		t.Fatalf("want 120, got %v", got)
	}
}

func TestRetryAfterZeroIsARealAnswer(t *testing.T) {
	// The server named the number itself. Different from saying nothing.
	if got := RetryAfter("0", time.Now()); got == nil || *got != 0 {
		t.Fatalf("want 0, got %v", got)
	}
}

func TestRetryAfterAbsentOrJunkIsNil(t *testing.T) {
	// nil, not 0: "said nothing" and "retry now" are different answers, and
	// guessing the second is what makes a client hammer a struggling API.
	for _, header := range []string{"", "soon", "   "} {
		if got := RetryAfter(header, time.Now()); got != nil {
			t.Fatalf("%q: want nil, got %v", header, *got)
		}
	}
}

func TestRetryAfterHTTPDateRoundsUp(t *testing.T) {
	// An HTTP-date carries whole seconds but now does not, so a header naming
	// the very next second can be a fraction away. A date names when a retry
	// becomes acceptable, so rounding short of it is permission to retry early.
	whole := time.Date(2026, 11, 4, 8, 30, 0, 0, time.UTC)
	header := whole.Format(http.TimeFormat)
	cases := []struct {
		offset time.Duration
		want   int
	}{
		{400 * time.Millisecond, 1},
		{time.Millisecond, 1},
		{1400 * time.Millisecond, 2},
		{90 * time.Second, 90},
	}
	for _, c := range cases {
		got := RetryAfter(header, whole.Add(-c.offset))
		if got == nil || *got != c.want {
			t.Fatalf("offset %v: want %d, got %v", c.offset, c.want, got)
		}
	}
}

func TestRetryAfterPastDateIsNil(t *testing.T) {
	whole := time.Date(2026, 11, 4, 8, 30, 0, 0, time.UTC)
	for _, at := range []time.Time{whole, whole.Add(5 * time.Second)} {
		if got := RetryAfter(whole.Format(http.TimeFormat), at); got != nil {
			t.Fatalf("want nil, got %v", *got)
		}
	}
}

func TestQueryDropsEmptyValues(t *testing.T) {
	got := TitleQuery{Q: "dinosaurs", Age: 4, Topic: "", Limit: 0}.values().Encode()
	if got != "age=4&q=dinosaurs" {
		t.Fatalf("got %q", got)
	}
}

func TestAPIErrorMessageFallsBack(t *testing.T) {
	if (&APIError{Status: 502}).Error() != "the Vooks API answered 502" {
		t.Fatal("bare status message wrong")
	}
	if (&APIError{Status: 400, Title: "Bad"}).Error() != "Bad" {
		t.Fatal("title fallback wrong")
	}
	if (&APIError{Status: 400, Title: "Bad", Detail: "age must be a number"}).Error() != "age must be a number" {
		t.Fatal("detail should win")
	}
}
