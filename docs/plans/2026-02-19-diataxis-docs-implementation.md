# Diataxis Documentation Restructure Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Restructure all documentation into the Diataxis model (tutorials, how-to, reference, explanation) with Hugo-compatible `_index.md` files, eliminating duplication and filling content gaps.

**Architecture:** Split existing README.md and docs/configuration.md into 15 focused pages across 4 Diataxis categories. Each page has one purpose. README becomes a slim landing page. Old monolithic docs/configuration.md is deleted after migration.

**Tech Stack:** Markdown, Hugo-compatible front matter, Diataxis documentation framework

**Source files to read before starting:**
- `README.md` — current content to migrate
- `docs/configuration.md` — current content to migrate
- `docs/plans/2026-02-19-diataxis-docs-design.md` — approved design
- `internal/server/server.go` — tool registrations (for reference/tools.md)
- `internal/server/handlers.go` — handler implementations (for reference/tools.md)
- `internal/gw2api/client.go` — struct definitions and API methods
- `internal/cache/manager.go` — cache TTL constants (for reference/caching.md)
- `Makefile` — build/dev commands (for how-to/contribute.md)

---

### Task 1: Create directory structure and section index pages

**Agent:** `diataxis-docs:diataxis-orchestrator`

**Files:**
- Create: `docs/_index.md`
- Create: `docs/tutorials/_index.md`
- Create: `docs/how-to/_index.md`
- Create: `docs/reference/_index.md`
- Create: `docs/explanation/_index.md`

**Step 1: Create all directories**

```bash
mkdir -p docs/tutorials docs/how-to docs/reference docs/explanation
```

**Step 2: Write docs/_index.md**

Root documentation index. Brief overview of the project and links to each Diataxis section. 5-10 lines. Must mention that the docs follow the Diataxis model.

**Step 3: Write docs/tutorials/_index.md**

Section index. Explain that tutorials are learning-oriented walkthroughs for newcomers. List the pages in the section (getting-started.md). 5-10 lines.

**Step 4: Write docs/how-to/_index.md**

Section index. Explain that how-to guides are task-oriented recipes for specific goals. List the pages (configure-mcp-clients.md, add-a-new-tool.md, contribute.md). 5-10 lines.

**Step 5: Write docs/reference/_index.md**

Section index. Explain that reference pages are information-oriented specs. List the pages (tools.md, api-scopes.md, caching.md, configuration.md). 5-10 lines.

**Step 6: Write docs/explanation/_index.md**

Section index. Explain that explanation pages provide understanding and context. List the pages (architecture.md, design-decisions.md). 5-10 lines.

**Step 7: Commit**

```bash
git add docs/tutorials/ docs/how-to/ docs/reference/ docs/explanation/ docs/_index.md
git commit -m "docs: add diataxis directory structure with section indexes"
```

---

### Task 2: Write reference/tools.md — consolidated tool reference

**Agent:** `diataxis-docs:doc-reference-gen`

**Files:**
- Create: `docs/reference/tools.md`
- Read: `internal/server/server.go` (tool registrations and descriptions)
- Read: `README.md` (current tool tables, lines 103-174)
- Read: `docs/configuration.md` (current tool tables, lines 162-211, and examples, lines 213-289)

**Step 1: Write docs/reference/tools.md**

This is the single source of truth for all tools. Must include ALL 37 tools (34 original + 3 new composite tools). For each tool:
- Tool name
- Auth requirement (Yes / No / Optional)
- Description
- Parameters with types and whether required/optional
- Example invocation JSON

Organize by category (Wiki, Account, Trading Post, Game Data, Wizard's Vault, Guilds, Game Metadata, Composite Tools). Use tables for the overview, then detailed parameter/example sections below.

Content sources:
- Tool names, descriptions, params from `internal/server/server.go` (arg struct jsonschema tags)
- Auth requirements from current README.md tool tables
- Example invocations from `docs/configuration.md` lines 213-289
- New composite tools: `get_item_by_name`, `get_item_recipe_by_name`, `get_tp_price_by_name`

Cross-links: Link to `../how-to/configure-mcp-clients.md` for setup, `api-scopes.md` for auth details.

**Step 2: Commit**

```bash
git add docs/reference/tools.md
git commit -m "docs: add consolidated tool reference"
```

---

### Task 3: Write reference/api-scopes.md — API key scopes reference

**Agent:** `diataxis-docs:doc-reference-gen`

**Files:**
- Create: `docs/reference/api-scopes.md`
- Read: `docs/configuration.md` (lines 7-22 for key creation, lines 291-315 for scopes table)

**Step 1: Write docs/reference/api-scopes.md**

Pure reference. Contains:
- Table mapping each authenticated tool to its required GW2 API key scopes (extracted from `docs/configuration.md` lines 296-315)
- List of all available scopes with descriptions
- Link to ArenaNet key management page: `https://account.arena.net/applications`

No advice, no tutorial steps. Just the facts.

Cross-links: Link to `tools.md` for tool details, `configuration.md` for how the key is configured, `../how-to/configure-mcp-clients.md` for setup instructions.

**Step 2: Commit**

```bash
git add docs/reference/api-scopes.md
git commit -m "docs: add API scopes reference"
```

---

### Task 4: Write reference/caching.md — cache TTL reference

**Agent:** `diataxis-docs:doc-reference-gen`

**Files:**
- Create: `docs/reference/caching.md`
- Read: `README.md` (lines 186-196 for caching strategy)
- Read: `internal/cache/manager.go` (TTL constants)

**Step 1: Write docs/reference/caching.md**

Pure reference table. Contains:
- Table of all cache categories with their TTL values (read actual constants from `internal/cache/manager.go`)
- Brief description of cache behavior: in-memory only, lost on process exit, per-API-key isolation via SHA-256 hash

No advice on tuning or optimization. Just the spec.

Cross-links: Link to `../explanation/architecture.md` for why caching works this way, `../explanation/design-decisions.md` for rationale.

**Step 2: Commit**

```bash
git add docs/reference/caching.md
git commit -m "docs: add caching TTL reference"
```

---

### Task 5: Write reference/configuration.md — environment and startup reference

**Agent:** `diataxis-docs:doc-reference-gen`

**Files:**
- Create: `docs/reference/configuration.md`
- Read: `docs/configuration.md` (lines 64-78 for API key behavior, lines 316-321 for security, lines 325-379 for troubleshooting)

**Step 1: Write docs/reference/configuration.md**

Pure reference. Contains:
- Environment variables table (`GW2_API_KEY`, `LOG_LEVEL` if applicable)
- Startup behavior: what happens with/without API key
- Security notes (extracted from `docs/configuration.md` lines 316-321): key hashing, no persistent storage, stdio-only, HTTPS
- Troubleshooting section: error messages and their meanings (extracted from `docs/configuration.md` lines 325-379)

Cross-links: Link to `api-scopes.md` for scope details, `../how-to/configure-mcp-clients.md` for setup.

**Step 2: Commit**

```bash
git add docs/reference/configuration.md
git commit -m "docs: add configuration and troubleshooting reference"
```

---

### Task 6: Write how-to/configure-mcp-clients.md

**Agent:** `diataxis-docs:doc-howto-writer`

**Files:**
- Create: `docs/how-to/configure-mcp-clients.md`
- Read: `docs/configuration.md` (lines 82-160 for client configs)

**Step 1: Write docs/how-to/configure-mcp-clients.md**

Task-oriented guide. Sections for each MCP client:
1. Claude Desktop (Docker + binary)
2. Claude Code (Docker + binary)
3. LM Studio (auto-install badge + manual)
4. Other MCP clients (generic stdio pattern)

Each section: numbered steps, JSON config blocks, variation for Docker vs binary. Assumes the reader already has the binary or Docker image (installation is in the tutorial).

Cross-links: Link to `../reference/configuration.md` for env var details, `../reference/api-scopes.md` for which scopes to enable, `../tutorials/getting-started.md` for first-time setup.

**Step 2: Commit**

```bash
git add docs/how-to/configure-mcp-clients.md
git commit -m "docs: add MCP client configuration how-to"
```

---

### Task 7: Write how-to/contribute.md

**Agent:** `diataxis-docs:doc-howto-writer`

**Files:**
- Create: `docs/how-to/contribute.md`
- Read: `README.md` (lines 211-244 for development/contributing)
- Read: `Makefile` (all targets)

**Step 1: Write docs/how-to/contribute.md**

Task-oriented guide for contributors. Sections:
1. Development setup (clone, install tools via `make tools`, build via `make build`)
2. Running tests (`make test`, `go test ./...`)
3. Linting and formatting (`make lint`, `make format`)
4. Commit conventions (conventional commits: feat/fix/docs/refactor)
5. Pull request workflow (fork, branch, changes, tests, lint, PR)
6. Code standards (gofumpt, golangci-lint, unit tests for core functionality)

Reference all Makefile targets. Assumes Go 1.24+ is installed.

Cross-links: Link to `add-a-new-tool.md` for the specific workflow of adding tools, `../explanation/architecture.md` for understanding the codebase.

**Step 2: Commit**

```bash
git add docs/how-to/contribute.md
git commit -m "docs: add contributor how-to guide"
```

---

### Task 8: Write how-to/add-a-new-tool.md

**Agent:** `diataxis-docs:doc-howto-writer`

**Files:**
- Create: `docs/how-to/add-a-new-tool.md`
- Read: `internal/server/server.go` (tool registration pattern)
- Read: `internal/server/handlers.go` (handler pattern)
- Read: `internal/gw2api/client.go` (API client pattern)

**Step 1: Write docs/how-to/add-a-new-tool.md**

Developer recipe. Walk through adding a hypothetical tool step-by-step:

1. Define the response struct in `internal/gw2api/client.go`
2. Add the API method in `internal/gw2api/client.go`
3. Define the args struct in `internal/server/server.go`
4. Write the handler in `internal/server/handlers.go`
5. Register the tool in `registerTools()` in `internal/server/server.go`
6. Add tests in `internal/server/handlers_test.go`
7. Build and verify (`make build && go test ./...`)

Use a concrete example (e.g., a `get_titles` tool) with actual code snippets showing the pattern. Reference real existing tools as examples of the pattern.

Cross-links: Link to `../reference/tools.md` for existing tool specs, `../explanation/architecture.md` for why the code is structured this way, `contribute.md` for the PR workflow.

**Step 2: Commit**

```bash
git add docs/how-to/add-a-new-tool.md
git commit -m "docs: add how-to guide for adding new tools"
```

---

### Task 9: Write tutorials/getting-started.md

**Agent:** `diataxis-docs:doc-tutorial-writer`

**Files:**
- Create: `docs/tutorials/getting-started.md`
- Read: `README.md` (installation section)
- Read: `docs/configuration.md` (installation + client config sections)

**Step 1: Write docs/tutorials/getting-started.md**

Learning-oriented tutorial. One golden path, no alternatives. Target: get from zero to querying GW2 data in ~15 minutes (not 90 min — this is an MCP server, not a framework).

Sections:
1. **What you'll build** — brief description of outcome
2. **Prerequisites** — Docker (recommended path), Claude Desktop or Claude Code
3. **Step 1: Get the server running** — `docker pull` + verify
4. **Step 2: Create a GW2 API key** — walk through ArenaNet site
5. **Step 3: Configure your MCP client** — Claude Desktop config (one client only, golden path)
6. **Step 4: Make your first query** — ask Claude to check your wallet
7. **Step 5: Explore game data** — look up an item by name, check TP prices
8. **Checkpoint: What you've accomplished** — summary
9. **Next steps** — links to how-to guides and reference

Each step has a checkpoint: "You should see..." verification.

Cross-links: Link to `../how-to/configure-mcp-clients.md` for other clients, `../reference/tools.md` for all available tools.

**Step 2: Commit**

```bash
git add docs/tutorials/getting-started.md
git commit -m "docs: add getting started tutorial"
```

---

### Task 10: Write explanation/architecture.md

**Agent:** `diataxis-docs:doc-explanation-writer`

**Files:**
- Create: `docs/explanation/architecture.md`
- Read: `README.md` (lines 198-208 for architecture section)
- Read: all files in `internal/` (package structure)
- Read: `internal/cache/manager.go` (caching layer)

**Step 1: Write docs/explanation/architecture.md**

Understanding-oriented. Explain:
1. **Package structure** — what each package does and why it's separate (`server/`, `gw2api/`, `wiki/`, `cache/`)
2. **Request flow** — how an MCP tool call flows: MCP client → stdio → server → handler → gw2api client → GW2 API, with caching at the gw2api layer
3. **Caching architecture** — why in-memory, how TTLs are chosen, API key hashing for cache isolation
4. **Tool registration pattern** — how `mcp.AddTool` connects arg structs, handlers, and the MCP protocol
5. **Wiki integration** — how wiki search + infobox parsing enables composite tools

Prose, not bullet lists. Explain the "why" behind each decision. This is not a reference — it's understanding.

Cross-links: Link to `../reference/caching.md` for exact TTL values, `design-decisions.md` for deeper rationale, `../how-to/add-a-new-tool.md` for the practical workflow.

**Step 2: Commit**

```bash
git add docs/explanation/architecture.md
git commit -m "docs: add architecture explanation"
```

---

### Task 11: Write explanation/design-decisions.md

**Agent:** `diataxis-docs:doc-explanation-writer`

**Files:**
- Create: `docs/explanation/design-decisions.md`
- Read: `internal/gw2api/client.go` (struct patterns)
- Read: `internal/server/handlers.go` (composite tool patterns)

**Step 1: Write docs/explanation/design-decisions.md**

Understanding-oriented. Explain key decisions:
1. **json.RawMessage for variable structures** — 20+ item subtypes, character equipment/skills/specializations. Typed Go structs add no value when the LLM reads JSON directly. Trade-off: no compile-time validation, but maximum flexibility.
2. **Typed structs where appropriate** — CraftingDiscipline (3 simple fields, commonly queried), ColorComponent (shared across materials), AchievementTier/Reward/Bit (stable schema). When to use typed vs raw.
3. **Composite tools (wiki + API)** — why `get_item_by_name` exists: reduces LLM context consumption, eliminates multi-step chaining errors, collapses 2-3 tool calls into 1.
4. **In-memory cache only** — no Redis, no disk. Process lifetime matches session lifetime. Simplicity over persistence.
5. **stdio-only transport** — no HTTP server, no network ports. Security by design: the server is a subprocess, not a network service.
6. **Excluding character bags from CharacterInfo** — inventory handled by dedicated `get_inventory` tool to avoid massive response payloads.

Cross-links: Link to `architecture.md` for the structural context, `../reference/caching.md` for TTL specs, `../reference/tools.md` for the composite tool specs.

**Step 2: Commit**

```bash
git add docs/explanation/design-decisions.md
git commit -m "docs: add design decisions explanation"
```

---

### Task 12: Rewrite README.md as slim landing page

**Agent:** None (direct implementation)

**Files:**
- Modify: `README.md`

**Step 1: Rewrite README.md**

Replace the current ~250-line README with a ~50-line landing page:

1. Title + LM Studio badge (keep existing badges)
2. **Prominent docs site link** — large, visible link to the hosted documentation (placeholder URL for now: `https://therealbill.github.io/gw2-mcp/`)
3. One-paragraph description (keep existing first paragraph)
4. Quick start (3 steps: Docker run command, API key setup, Claude Desktop JSON config)
5. Feature highlights (5-6 bullet points, not the full list)
6. Documentation links table:

| Section | Description |
|---------|-------------|
| [Getting Started](docs/tutorials/getting-started.md) | Install, configure, and make your first query |
| [How-To Guides](docs/how-to/) | Configure clients, add tools, contribute |
| [Reference](docs/reference/) | Tools, API scopes, caching, configuration |
| [Architecture](docs/explanation/) | System design and decision rationale |

7. License + Acknowledgments (keep existing)

**Step 2: Commit**

```bash
git add README.md
git commit -m "docs: rewrite README as slim landing page with docs site link"
```

---

### Task 13: Delete old docs/configuration.md

**Agent:** None (direct implementation)

**Files:**
- Delete: `docs/configuration.md`

**Step 1: Verify all content migrated**

Before deleting, verify every section of `docs/configuration.md` has a new home:
- Prerequisites/API key creation → `reference/api-scopes.md`
- Installation (Docker, binary) → `tutorials/getting-started.md`
- API key configuration → `reference/configuration.md`
- MCP client configs → `how-to/configure-mcp-clients.md`
- Tool tables → `reference/tools.md`
- Example invocations → `reference/tools.md`
- API scopes table → `reference/api-scopes.md`
- Security notes → `reference/configuration.md`
- Troubleshooting → `reference/configuration.md`

**Step 2: Delete and commit**

```bash
git rm docs/configuration.md
git commit -m "docs: remove old monolithic configuration.md (content migrated to diataxis structure)"
```

---

### Task 14: Cross-link validation

**Agent:** `diataxis-docs:doc-crosslink-validator`

**Files:**
- Read: all files in `docs/` recursively

**Step 1: Validate all cross-links**

Check every markdown link in every docs file:
- All relative links resolve to existing files
- No broken links
- No orphaned pages (every page is linked from at least one other page)
- Each `_index.md` lists all pages in its directory
- README links to docs sections correctly

**Step 2: Fix any broken links**

If any links are broken, fix them.

**Step 3: Verify no content duplication**

Scan for duplicated content (same tool tables, same config blocks) across files. Each piece of information should live in exactly one place.

**Step 4: Commit fixes if any**

```bash
git add docs/
git commit -m "docs: fix cross-links and remove duplication"
```

---

### Task 15: Final verification

**Agent:** None (direct implementation)

**Step 1: Verify Hugo compatibility**

Check that every directory under `docs/` has an `_index.md`:
- `docs/_index.md`
- `docs/tutorials/_index.md`
- `docs/how-to/_index.md`
- `docs/reference/_index.md`
- `docs/explanation/_index.md`

**Step 2: Verify build still works**

```bash
make build
go test ./...
```

Expected: both pass (docs changes don't affect code).

**Step 3: Verify old file is gone**

```bash
test ! -f docs/configuration.md && echo "OK: old config deleted"
```

**Step 4: List final structure**

```bash
find docs/ -name "*.md" | sort
```

Expected output:
```
docs/_index.md
docs/explanation/_index.md
docs/explanation/architecture.md
docs/explanation/design-decisions.md
docs/how-to/_index.md
docs/how-to/add-a-new-tool.md
docs/how-to/configure-mcp-clients.md
docs/how-to/contribute.md
docs/plans/2026-02-19-diataxis-docs-design.md
docs/plans/2026-02-19-diataxis-docs-implementation.md
docs/reference/_index.md
docs/reference/api-scopes.md
docs/reference/caching.md
docs/reference/configuration.md
docs/reference/tools.md
docs/tutorials/_index.md
docs/tutorials/getting-started.md
```
