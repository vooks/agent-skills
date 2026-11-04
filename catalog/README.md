# vooks-catalog for Go

Official Go client for the [Vooks](https://www.vooks.com) Public Catalog API.

```bash
go get github.com/vooks/agent-skills/catalog
```

```go
client := catalog.New()
page, err := client.ListTitles(ctx, catalog.TitleQuery{Q: "dinosaurs", Age: 4, Limit: 5})
for _, t := range page.Data {
    fmt.Println(t.Title, t.URL)
}
```

Reads need no credential. A token only raises the rate limit and unlocks exports.
No dependencies outside the standard library.

`Titles` follows the cursor for you, and refuses a cursor that does not advance rather
than looping. `APIError` carries the API's stable `Code`; branch on that, not the
message. `RetryAfter` reads both forms RFC 9110 allows and rounds a future date up,
because a date names when a retry becomes acceptable.

Full guide: https://www.vooks.com/developers
