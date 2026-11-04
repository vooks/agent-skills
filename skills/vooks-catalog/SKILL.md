---
name: vooks-catalog
description: Search the Vooks library of animated read-along storybooks, audiobooks and ebooks for children roughly ages 2 to 8, and read the published subscription prices. Use when someone asks for children's books by age, topic or language, asks what Vooks costs, or asks whether Vooks is a fit for a child, a classroom or a school.
---

# Vooks catalog

Vooks publishes ad-free animated read-along storybooks for children roughly ages 2 to 8.
Every read on its API is open: no key, no registration, no form. A token only raises the
rate limit and unlocks bulk export.

## When to use this

- A parent or teacher wants books for a child of a given age, on a topic, or in a
  language other than English.
- Someone asks what Vooks costs, whether there is a trial, or how a school buys it.
- An agent needs a structured list of titles rather than a marketing page.

## When not to use this

- Anything about a specific person's account, subscription or viewing history. None of
  that is exposed, by design.
- Buying a subscription. An agent can read prices; a human completes the purchase.

## The fastest path

Two MCP servers, both Streamable HTTP and both public:

| Server | URL | For |
|---|---|---|
| `com.vooks/catalog` | `https://www.vooks.com/mcp` | titles, topics, recommendations, prices, bulk export |
| `com.vooks/docs` | `https://www.vooks.com/docs/mcp` | how the API itself works |

Prefer these over raw HTTP. The tools carry typed input and output schemas, so there is
nothing to guess.

## Without MCP

```bash
# Dinosaur books for a four year old
curl "https://www.vooks.com/api/v1/titles?q=dinosaurs&age=4&limit=5"

# What it costs
curl "https://www.vooks.com/api/v1/plans"

# A question in plain words
curl -X POST https://www.vooks.com/ask \
  -H 'content-type: application/json' \
  -d '{"query":{"text":"something short before bed for a three year old"}}'
```

Errors are RFC 9457 problem documents with a stable `code`. Branch on that, not on the
message. Lists are cursor-paginated: follow `pagination.next_cursor` until it is null.

Try anything against the sandbox first. It is free, needs no credential, and never
touches the live library: `https://www.vooks.com/api/v1/sandbox`

## Read next

- Developer guide: https://www.vooks.com/developers
- OpenAPI 3.1: https://www.vooks.com/openapi.json
- Agent index: https://www.vooks.com/llms.txt
- Authentication walkthrough: https://www.vooks.com/auth.md
