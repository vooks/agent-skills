// Package catalog is the official Go client for the Vooks Public Catalog API:
// a library of ad-free animated read-along storybooks, audiobooks and ebooks
// for children roughly ages 2 to 8.
//
// Reads need no credential. A token only raises the rate limit from 600 to
// 3,000 requests a minute and unlocks catalog exports.
//
//	client := catalog.New()
//	page, err := client.ListTitles(ctx, catalog.TitleQuery{Q: "dinosaurs", Age: 4, Limit: 5})
//
// No dependencies outside the standard library. A dependency in a published
// SDK is a dependency in every project that imports it.
package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	Version        = "1.0.0"
	BaseURL        = "https://www.vooks.com/api/v1"
	SandboxBaseURL = "https://www.vooks.com/api/v1/sandbox"
)

// APIError is an HTTP error response, carrying the API's stable
// machine-readable Code. Branch on that, not on the message.
//
// A transport failure, where the request never got an answer, returns the
// underlying error instead: there is no status or code to report, so wrapping
// it would invent one.
type APIError struct {
	Status            int               `json:"-"`
	Code              string            `json:"code"`
	Title             string            `json:"title"`
	Detail            string            `json:"detail"`
	Errors            []ValidationError `json:"errors"`
	RetryAfterSeconds *int              `json:"-"`
}

type ValidationError struct {
	Parameter string `json:"parameter"`
	Detail    string `json:"detail"`
}

func (e *APIError) Error() string {
	if e.Detail != "" {
		return e.Detail
	}
	if e.Title != "" {
		return e.Title
	}
	return fmt.Sprintf("the Vooks API answered %d", e.Status)
}

// RetryAfter reads either form RFC 9110 allows: a count of seconds, or an
// HTTP-date converted to seconds from now.
//
// A future date rounds UP. The date names the moment a retry becomes
// acceptable, so rounding short of it would tell the caller to retry before
// the server was ready. Returns nil when the header is absent, unparseable,
// or already past, because "the server said nothing" and "retry immediately"
// are different answers and guessing the second is what makes a client hammer
// a struggling API.
func RetryAfter(header string, now time.Time) *int {
	value := strings.TrimSpace(header)
	if value == "" {
		return nil
	}
	if seconds, err := strconv.Atoi(value); err == nil && seconds >= 0 {
		return &seconds
	}
	when, err := http.ParseTime(value)
	if err != nil {
		return nil
	}
	delta := when.Sub(now)
	if delta <= 0 {
		return nil
	}
	seconds := int(math.Ceil(delta.Seconds()))
	return &seconds
}

type Client struct {
	BaseURL     string
	AccessToken string
	UserAgent   string
	HTTP        *http.Client
}

func New() *Client {
	return &Client{
		BaseURL:   BaseURL,
		UserAgent: "vooks-catalog-go/" + Version,
		HTTP:      &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) do(ctx context.Context, method, path string, body io.Reader, headers map[string]string, out any) error {
	req, err := http.NewRequestWithContext(ctx, method, strings.TrimRight(c.BaseURL, "/")+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("user-agent", c.UserAgent)
	if c.AccessToken != "" {
		req.Header.Set("authorization", "Bearer "+c.AccessToken)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		apiErr := &APIError{Status: res.StatusCode, Code: "unknown_error"}
		raw, _ := io.ReadAll(res.Body)
		// A non-2xx body is whatever the server sent. Unmarshalling a list or
		// a bare string into the struct fails, and the caller would lose both
		// the status and the typed error, so the status is set either way.
		_ = json.Unmarshal(raw, apiErr)
		if apiErr.Code == "" {
			apiErr.Code = "unknown_error"
		}
		apiErr.Status = res.StatusCode
		apiErr.RetryAfterSeconds = RetryAfter(res.Header.Get("retry-after"), time.Now())
		return apiErr
	}
	return json.NewDecoder(res.Body).Decode(out)
}

type TitleQuery struct {
	Q        string
	Age      int
	Topic    string
	Language string
	Sort     string
	Limit    int
	Cursor   string
}

func (q TitleQuery) values() url.Values {
	v := url.Values{}
	set := func(k, s string) {
		if s != "" {
			v.Set(k, s)
		}
	}
	set("q", q.Q)
	set("topic", q.Topic)
	set("language", q.Language)
	set("sort", q.Sort)
	set("cursor", q.Cursor)
	if q.Age > 0 {
		v.Set("age", strconv.Itoa(q.Age))
	}
	if q.Limit > 0 {
		v.Set("limit", strconv.Itoa(q.Limit))
	}
	return v
}

type Title struct {
	Slug            string   `json:"slug"`
	Title           string   `json:"title"`
	Subtitle        string   `json:"subtitle"`
	Author          string   `json:"author"`
	Summary         string   `json:"summary"`
	URL             string   `json:"url"`
	AgeLabel        string   `json:"age_label"`
	Language        string   `json:"language"`
	DurationSeconds int      `json:"duration_seconds"`
	Topics          []string `json:"topics"`
}

type Pagination struct {
	Limit      int    `json:"limit"`
	Total      int    `json:"total"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor"`
}

type TitlePage struct {
	Data       []Title    `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// ListTitles returns one page. Pagination.Total is the true count.
func (c *Client) ListTitles(ctx context.Context, q TitleQuery) (*TitlePage, error) {
	page := &TitlePage{}
	suffix := ""
	if encoded := q.values().Encode(); encoded != "" {
		suffix = "?" + encoded
	}
	return page, c.do(ctx, http.MethodGet, "/titles"+suffix, nil, nil, page)
}

// Titles calls yield for every title matching the query, following the cursor.
// Returning false from yield stops the walk. Prefer ListTitles when one page
// is enough: this makes a request per page.
func (c *Client) Titles(ctx context.Context, q TitleQuery, yield func(Title) bool) error {
	cursor := q.Cursor
	for {
		q.Cursor = cursor
		page, err := c.ListTitles(ctx, q)
		if err != nil {
			return err
		}
		for _, t := range page.Data {
			if !yield(t) {
				return nil
			}
		}
		next := page.Pagination.NextCursor
		// A page that hands back the cursor it was called with would be
		// fetched forever, yielding the same titles. Stopping silently would
		// look like the end of the catalog, which is a worse lie.
		if next != "" && next == cursor {
			return fmt.Errorf("the catalog returned the same pagination cursor it was given (%q), which would repeat this page forever", next)
		}
		if next == "" {
			return nil
		}
		cursor = next
	}
}

type TitleDetail struct {
	Title
	DescriptionLong string `json:"description_long"`
}

// GetTitle reads one title live, including the full description.
func (c *Client) GetTitle(ctx context.Context, slug string) (*TitleDetail, error) {
	out := &TitleDetail{}
	return out, c.do(ctx, http.MethodGet, "/titles/"+url.PathEscape(slug), nil, nil, out)
}

type Topic struct {
	Slug       string `json:"slug"`
	Name       string `json:"name"`
	TitleCount int    `json:"title_count"`
}

type TopicPage struct {
	Data       []Topic    `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// ListTopics returns the topic vocabulary TitleQuery.Topic accepts.
func (c *Client) ListTopics(ctx context.Context, limit int) (*TopicPage, error) {
	page := &TopicPage{}
	suffix := ""
	if limit > 0 {
		suffix = "?limit=" + strconv.Itoa(limit)
	}
	return page, c.do(ctx, http.MethodGet, "/topics"+suffix, nil, nil, page)
}

type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Plan struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Price     Price  `json:"price"`
	Interval  string `json:"interval"`
	Audience  string `json:"audience"`
	TrialDays int    `json:"trial_days"`
	SignupURL string `json:"signup_url"`
}

type Plans struct {
	Data   []Plan            `json:"data"`
	Quoted []json.RawMessage `json:"quoted"`
}

// ListPlans returns the published prices. Never answer a pricing question
// from memory: prices change and this is the published source.
func (c *Client) ListPlans(ctx context.Context) (*Plans, error) {
	out := &Plans{}
	return out, c.do(ctx, http.MethodGet, "/plans", nil, nil, out)
}
