---
title: Guild Lookup
weight: 20
---

# How to Look Up Guild Information

**Goal**: Find a guild by name and view its public info, or access detailed guild data if you are the leader.

**Time**: Approximately 5 minutes

## Prerequisites

- **For public info** (name, tag, level): No API key required.
- **For detailed data** (members, stash, treasury, etc.): An API key with the **guilds** scope, and you must be the **guild leader**.

## Steps

### 1. Search for a guild by name

> Ask your AI: "Find the guild called Mighty Teapot"

Your assistant uses the `search_guild` tool to search for guilds matching that name. The result is one or more guild IDs -- long identifier strings the GW2 API uses internally. You do not need to remember these IDs; the AI uses them in the next step automatically.

### 2. Get public guild info

Once the AI has the guild ID from the search, it calls the `get_guild` tool to retrieve public information:

- **Guild name** -- the full display name
- **Tag** -- the short tag shown in brackets next to member names (e.g., `[TEA]`)
- **Level** -- the guild's current level (max 69)

This works for any guild in the game, regardless of whether you are a member.

> Ask your AI: "What's the guild tag for Mighty Teapot?"

### 3. Get detailed guild data (guild leader only)

If you are the guild leader and your API key has the **guilds** scope, you can access internal guild data. Ask for specific detail types:

> Ask your AI: "Show me the guild members for Mighty Teapot"

The AI calls `get_guild_details` with the type set to `members`. The available detail types are:

| Detail type | What it shows |
|-------------|---------------|
| `members` | All guild members with ranks and join dates |
| `log` | Recent guild activity (invites, kicks, upgrades, stash changes) |
| `ranks` | Guild rank names and permission settings |
| `stash` | Items stored in the guild vault |
| `storage` | Decorations and other guild-level storage |
| `treasury` | Items donated toward guild upgrades |
| `teams` | PvP team rosters |
| `upgrades` | Completed and in-progress guild upgrades |

Example prompts:

> Ask your AI: "What's in our guild stash?"

> Ask your AI: "Show the guild activity log"

> Ask your AI: "What upgrades has our guild completed?"

## Verify it works

Test with a known guild name:

> Ask your AI: "Find the guild called ArenaNet and show me their tag and level"

You should see the guild's public info -- name, tag, and level -- returned without needing an API key. If you then try to access members or stash for a guild you do not lead, the AI will report an authorization error. That is expected.

## Troubleshooting

### Problem: No results for a guild name
**Symptom**: The search returns no matching guilds.
**Cause**: The guild name must be an exact match. Partial names and abbreviations do not work with the guild search API.
**Solution**: Use the guild's full display name, not its tag. If you only know the tag, try searching the wiki or asking in-game.

### Problem: Authorization error on detailed data
**Symptom**: Public info works, but requesting members, stash, or other details returns a permissions error.
**Cause**: Detailed guild data requires both a **guilds** scope API key and guild leader status. Regular members cannot access this data through the API.
**Solution**: Confirm your API key has the **guilds** scope at [Guild Wars 2 API Key Management](https://account.arena.net/applications). If you are not the guild leader, you can only view public info (name, tag, level).

### Problem: Wrong guild returned
**Symptom**: The search returns a guild with the right name but it is not the one you expected.
**Cause**: Multiple guilds can share the same name in GW2. The search returns any guilds matching the name.
**Solution**: Check the guild tag to confirm you have the right one. If multiple results come back, ask the AI to show the tag for each.

## See also

- [Know Your Account](../tutorials/account-overview/) -- view your account's guild memberships and other account data
- [Tools reference](../reference/tools/) -- full specification for `search_guild`, `get_guild`, and `get_guild_details`
