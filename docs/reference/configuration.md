---
title: Configuration
---

# Configuration

Environment variables, startup behavior, security properties, and error messages for the GW2 MCP Server.

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `GW2_API_KEY` | No | Guild Wars 2 API key. Enables authenticated tools that access account-specific data. Created at [account.arena.net/applications](https://account.arena.net/applications). See [API Key Scopes](api-scopes/) for the permissions each tool requires. |

`GW2_API_KEY` is the only environment variable the server reads.

## Startup Behavior

The server reads `GW2_API_KEY` from the environment once at startup and passes it to the API client. The key is not re-read during the server's lifetime.

### With `GW2_API_KEY` set

1. The server starts and registers all 37 tools.
2. Both authenticated and unauthenticated tools are available.
3. The server logs its version, commit hash, and build date at startup.

### Without `GW2_API_KEY`

1. The server logs a warning to stderr: `GW2_API_KEY environment variable not set; authenticated endpoints will be unavailable`
2. The server starts and registers all 37 tools.
3. Unauthenticated tools function normally.
4. Authenticated tools return the error: `GW2_API_KEY environment variable not configured`

## Communication

The server communicates exclusively over stdio (standard input/output). It does not open network ports or HTTP listeners.

## Security

- **Single-read API key.** `GW2_API_KEY` is read from the process environment at startup. It is never accepted as a tool parameter.
- **Hashed cache keys.** The API key is hashed with SHA-256. Only the first 8 bytes of the hash are used as a cache key prefix. The raw API key is not stored in the cache.
- **In-memory cache only.** All cached data is held in process memory. No data is written to disk. All cache entries are lost when the server process exits.
- **No network listeners.** The server uses stdio transport only. It does not bind to any port or start any HTTP server.
- **HTTPS transmission.** The API key is sent to the GW2 API (`api.guildwars2.com`) as an `Authorization: Bearer` header over HTTPS.

## Troubleshooting

### Error Messages

| Error Message | Source | Meaning |
|---------------|--------|---------|
| `GW2_API_KEY environment variable not configured` | Server | An authenticated tool was called but `GW2_API_KEY` was not set at startup. |
| `Cannot connect to the Docker daemon` | Docker | The Docker daemon is not running. The server cannot start in Docker mode without it. |
| `spawn gw2-mcp ENOENT` | MCP Client | The MCP client cannot find the `gw2-mcp` binary at the configured path. |
| `API request failed with status 401` | Server | The GW2 API rejected the API key. The key is invalid or has been deleted. |
| `API request failed with status 403` | Server | The API key is valid but lacks the required permission scopes for the requested endpoint. See [API Key Scopes](api-scopes/) for scope requirements per tool. |
| `invalid transaction type "...": must be one of current/buys, current/sells, history/buys, history/sells` | Server | The `type` parameter passed to `get_tp_transactions` is not one of the four accepted values. |
| `invalid direction "...": must be "coins" or "gems"` | Server | The `direction` parameter passed to `get_gem_exchange` is not `coins` or `gems`. |

### Startup Warning

The following log line is emitted to stderr at `WARN` level when `GW2_API_KEY` is absent:

```
GW2_API_KEY environment variable not set; authenticated endpoints will be unavailable
```

This is informational. The server continues to start and unauthenticated tools remain functional.

## See Also

- [API Key Scopes](api-scopes/) -- permissions required by each authenticated tool
- [How to Configure MCP Clients](../how-to/configure-mcp-clients/) -- setup instructions for Claude Desktop, Claude Code, and other clients
