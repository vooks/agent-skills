# Working in this repo

This repo holds the Vooks agent surface descriptors: the skills an agent can load, and
the MCP servers it can connect to. It holds no application code. The implementation
lives in a private repo; what is here is the public contract.

## The one rule

**Nothing here may describe a capability that does not exist.** Every URL, price, filter
and tool name in these files is checked against the live API before it is published. If
you are adding a claim you cannot verify with curl, do not add it.

## Verify before you commit

```bash
# Every plan price named in skills/vooks-pricing/SKILL.md must match this
curl -s https://www.vooks.com/api/v1/plans

# Every filter named in skills/vooks-catalog/SKILL.md must work
curl -s "https://www.vooks.com/api/v1/titles?q=dinosaurs&age=4&limit=3"

# Both MCP servers in mcp.json must answer, and their tool names must match
curl -s -X POST https://www.vooks.com/mcp \
  -H 'content-type: application/json' -H 'accept: application/json' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}'
```

## Layout, and why it is this shape

```
plugin.json              Agent Plugins manifest
mcp.json                 MCP servers, discovered only at the plugin root
skills/
  vooks-catalog/SKILL.md finding titles
  vooks-pricing/SKILL.md what it costs and how schools buy
AGENTS.md                this file
```

Agent Plugins clients discover skills **only** under `skills/<skill-name>/SKILL.md` and
MCP servers **only** through `mcp.json` at the root. A flat layout publishes a plugin
that exposes nothing, so do not rearrange these paths.

The directory name under `skills/` must equal the `name` in that skill's frontmatter.

## What does not belong here

- Application code, tests, build config. This is a descriptor repo.
- Anything secret. There is nothing to protect here and nothing should arrive.
- Skills for capabilities that are planned rather than shipped.

## Where the real documentation is

- Developer guide: https://www.vooks.com/developers
- OpenAPI 3.1: https://www.vooks.com/openapi.json
- Agent index: https://www.vooks.com/llms.txt
- Authentication: https://www.vooks.com/auth.md
