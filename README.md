# Vooks agent skills

Agent-facing descriptors for [Vooks](https://www.vooks.com): a library of ad-free
animated read-along storybooks, audiobooks and ebooks for children roughly ages 2 to 8.

Every read on the Vooks Public Catalog API is open. No key, no registration, no form.

## Skills

| Skill | For |
|---|---|
| [`vooks-catalog`](skills/vooks-catalog/SKILL.md) | Finding titles by age, topic, language or free text |
| [`vooks-pricing`](skills/vooks-pricing/SKILL.md) | What a subscription costs, the teacher rate, the trial, school licensing |

## MCP servers

Both are Streamable HTTP, both public, and both listed in the
[MCP Registry](https://registry.modelcontextprotocol.io) under `com.vooks`.

| Server | URL |
|---|---|
| `com.vooks/catalog` | `https://www.vooks.com/mcp` |
| `com.vooks/docs` | `https://www.vooks.com/docs/mcp` |

```bash
curl -X POST https://www.vooks.com/mcp \
  -H 'content-type: application/json' -H 'accept: application/json' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}'
```

## Clients

| Ecosystem | Install |
|---|---|
| npm | `npx @vooks/catalog search dinosaurs --age 4` |
| PyPI | `pip install vooks-catalog` |
| Go | `go get github.com/vooks/agent-skills/catalog` |

The Go client lives in [`catalog/`](catalog/) in this repo.

## Reading on

- [Developer guide](https://www.vooks.com/developers)
- [OpenAPI 3.1](https://www.vooks.com/openapi.json)
- [Agent index](https://www.vooks.com/llms.txt)
- [Authentication](https://www.vooks.com/auth.md)

The implementation is not open source. These descriptors are published so agents can
discover and use the public API.
