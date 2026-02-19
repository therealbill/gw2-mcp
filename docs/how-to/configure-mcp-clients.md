---
title: Configure MCP Clients
---

# How to Configure MCP Clients

**Goal**: Connect your MCP client to the GW2 MCP Server so you can query Guild Wars 2 data from your AI assistant.

## Prerequisites

Before starting, you should have:

- A GW2 API key created at [Guild Wars 2 API Key Management](https://account.arena.net/applications) -- see [API Key Scopes](../reference/api-scopes/) for which permissions to enable
- **Either** [Docker](https://docs.docker.com/get-docker/) installed **or** the `gw2-mcp` binary downloaded/built -- see [Getting Started](../tutorials/getting-started/) for installation steps

## Claude Desktop

Open Claude Desktop and navigate to **Settings > Developer > Edit Config**. This opens the `claude_desktop_config.json` file.

### Option A: Docker

```json
{
  "mcpServers": {
    "gw2-mcp": {
      "command": "docker",
      "args": ["run", "--rm", "-i", "-e", "GW2_API_KEY", "alyxpink/gw2-mcp:v1"],
      "env": {
        "GW2_API_KEY": "YOUR_KEY_HERE"
      }
    }
  }
}
```

### Option B: Direct binary

```json
{
  "mcpServers": {
    "gw2-mcp": {
      "command": "/path/to/gw2-mcp",
      "env": {
        "GW2_API_KEY": "YOUR_KEY_HERE"
      }
    }
  }
}
```

Replace `/path/to/gw2-mcp` with the absolute path to the binary (for example, `C:\\Users\\you\\bin\\gw2-mcp.exe` on Windows or `/usr/local/bin/gw2-mcp` on Linux/macOS).

### Verify it works

1. Restart Claude Desktop after saving the config file.
2. Open a new conversation and look for the hammer icon in the input area -- this confirms MCP tools are loaded.
3. Ask Claude: "What is the current Guild Wars 2 game build number?"

Claude should call the `get_game_build` tool and return a numeric build ID.

## Claude Code

Create or edit the `.mcp.json` file in your project root.

### Option A: Docker

```json
{
  "mcpServers": {
    "gw2-mcp": {
      "command": "docker",
      "args": ["run", "--rm", "-i", "-e", "GW2_API_KEY", "alyxpink/gw2-mcp:v1"],
      "env": {
        "GW2_API_KEY": "YOUR_KEY_HERE"
      }
    }
  }
}
```

### Option B: Direct binary

```json
{
  "mcpServers": {
    "gw2-mcp": {
      "command": "/path/to/gw2-mcp",
      "env": {
        "GW2_API_KEY": "YOUR_KEY_HERE"
      }
    }
  }
}
```

### Verify it works

1. Restart Claude Code or reload the window after saving `.mcp.json`.
2. Ask Claude Code to search the Guild Wars 2 wiki for "Mystic Coin".

Claude Code should invoke the `wiki_search` tool and return wiki results.

## LM Studio

### Auto-install

Click the badge on the [GW2 MCP Server README](https://github.com/AlyxPink/gw2-mcp) to install automatically:

[![Install in LM Studio](https://files.lmstudio.ai/deeplink/mcp-install-light.svg)](https://lmstudio.ai/install-mcp?name=gw2-mcp&config=eyJjb21tYW5kIjoiZG9ja2VyIiwiYXJncyI6WyJydW4iLCItLXJtIiwiLWkiLCJhbHl4cGluay9ndzItbWNwOnYxIl19)

This registers the Docker-based server. After installing, open LM Studio's MCP server settings and add your `GW2_API_KEY` to the environment variables.

### Manual configuration

Add the server manually in LM Studio with:

- **Command**: `docker`
- **Args**: `run --rm -i -e GW2_API_KEY=YOUR_KEY_HERE alyxpink/gw2-mcp:v1`

Or point to the binary directly:

- **Command**: `/path/to/gw2-mcp`
- **Environment variable**: `GW2_API_KEY=YOUR_KEY_HERE`

## Other MCP Clients (Cursor, Windsurf, etc.)

Any MCP client that supports **stdio** transport can use the GW2 MCP Server. The configuration pattern is the same across all clients.

### Docker

```json
{
  "mcpServers": {
    "gw2-mcp": {
      "command": "docker",
      "args": ["run", "--rm", "-i", "-e", "GW2_API_KEY", "alyxpink/gw2-mcp:v1"],
      "env": {
        "GW2_API_KEY": "YOUR_KEY_HERE"
      }
    }
  }
}
```

### Direct binary

```json
{
  "mcpServers": {
    "gw2-mcp": {
      "command": "/path/to/gw2-mcp",
      "env": {
        "GW2_API_KEY": "YOUR_KEY_HERE"
      }
    }
  }
}
```

Consult your client's documentation for where to place this JSON. Common locations:

| Client | Config location |
|--------|-----------------|
| Cursor | `.cursor/mcp.json` in your project root |
| Windsurf | MCP settings in the Windsurf preferences panel |
| VS Code (Copilot) | `.vscode/mcp.json` in your project root |

The key details are always the same: set `command` to either `docker` or the binary path, pass `GW2_API_KEY` through the `env` block, and use stdio transport.

## Troubleshooting

### Problem: Tools do not appear in the client

**Symptom**: The MCP client connects but no GW2 tools are listed.

**Cause**: The server process failed to start or crashed immediately.

**Solution**:

1. Run the server manually in a terminal to check for errors:
   ```bash
   docker run --rm -i -e GW2_API_KEY="YOUR_KEY_HERE" alyxpink/gw2-mcp:v1
   ```
2. If using the binary, confirm the path is correct and the file is executable.
3. Verify Docker is running if using the Docker configuration.

### Problem: "GW2_API_KEY environment variable not configured"

**Symptom**: Tools that require authentication return this error message.

**Cause**: The API key is not being passed to the server process.

**Solution**: Confirm the `env` block in your client config contains `GW2_API_KEY` with a valid key. For Docker configs, both the `args` array (`-e GW2_API_KEY`) and the `env` block are required -- the args flag tells Docker to forward the variable, and the env block sets its value.

### Problem: Authorization error from the GW2 API

**Symptom**: A tool returns a permissions or authorization error.

**Cause**: The API key is missing a required scope for that tool.

**Solution**: Check [API Key Scopes](../reference/api-scopes/) for the scopes each tool requires, then update your key at [Guild Wars 2 API Key Management](https://account.arena.net/applications).

### Problem: Docker image not found

**Symptom**: `docker: Error response from daemon: manifest for alyxpink/gw2-mcp:v1 not found`

**Cause**: Docker cannot pull the image.

**Solution**: Try the GitHub Container Registry mirror instead:
```json
"args": ["run", "--rm", "-i", "-e", "GW2_API_KEY", "ghcr.io/alyxpink/gw2-mcp:v1"]
```

## See also

- [Configuration reference](../reference/configuration/) -- environment variables and server startup options
- [API Key Scopes](../reference/api-scopes/) -- required permissions for each tool
- [Getting Started tutorial](../tutorials/getting-started/) -- full installation and first-query walkthrough
