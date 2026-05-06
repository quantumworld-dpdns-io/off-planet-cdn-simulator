# mcp-server

The MCP (Mission Control Protocol) server exposes a tool-call HTTP API used by
AI agents and dashboards to inspect and operate the Off-Planet CDN system.

## Running

```
MCP_SERVER_PORT=8084 ./mcp-server
```

## API

`POST /tools/call` accepts:

```json
{ "tool": "<tool_name>", "input": { ... } }
```

## Available tools

| Tool | Description |
|------|-------------|
| `cache_status` | Returns the fill ratio, pinned-object count, and top cached objects for a given node (`node_id`). |
| `generate_preload_plan` | Produces a prioritized list of URLs to prefetch onto a node before the next communication blackout. |
| `inspect_node` | Returns live diagnostic information for an edge node: status, capacity, and usage. |
| `simulate_eviction` | Predicts which cached objects would be evicted to free a requested amount of space on a node. |
| `summarize_incident` | Generates a human-readable summary of a CDN incident from audit logs and telemetry. |
