# Design: Diataxis Documentation Restructure

## Problem

The current documentation (README.md + docs/configuration.md) mixes tutorial, how-to, reference, and explanation content. Tool reference tables are duplicated across both files. Key content is missing entirely: no getting-started tutorial, no contributor guide for adding tools, no architecture explanation.

## Audience

Primary: end users installing and using the MCP server with Claude, LM Studio, etc.
Secondary: developers contributing to or extending the codebase.

## Approach

Full Diataxis model with subdirectories per category. Hugo-compatible with `_index.md` files for each section. README becomes a slim landing page with a prominent link to the hosted docs site.

## File Structure

```
README.md                              # Slim landing page (~50 lines)

docs/
├── _index.md                          # Docs root index
├── tutorials/
│   ├── _index.md                      # Tutorials section index
│   └── getting-started.md             # NEW: Install → configure → first query
├── how-to/
│   ├── _index.md                      # How-to section index
│   ├── configure-mcp-clients.md       # EXTRACTED from docs/configuration.md
│   ├── add-a-new-tool.md             # NEW: Developer recipe for adding tools
│   └── contribute.md                  # EXPANDED from README Contributing section
├── reference/
│   ├── _index.md                      # Reference section index
│   ├── tools.md                       # CONSOLIDATED from README + configuration.md
│   ├── api-scopes.md                  # EXTRACTED from docs/configuration.md
│   ├── caching.md                     # EXTRACTED from README
│   └── configuration.md               # EXTRACTED from docs/configuration.md
└── explanation/
    ├── _index.md                      # Explanation section index
    ├── architecture.md                # EXPANDED from README Architecture section
    └── design-decisions.md            # NEW: Why json.RawMessage, composites, etc.
```

## Content Migration Map

### Current docs/configuration.md splits into:
- **how-to/configure-mcp-clients.md**: Claude Desktop, Claude Code, LM Studio, Docker configs, other clients
- **reference/tools.md**: Tool tables (merged with README tables, single source of truth)
- **reference/api-scopes.md**: Scope-to-tool mapping, key creation instructions
- **reference/configuration.md**: Environment variables, startup behavior, security notes
- Troubleshooting content moves to reference/configuration.md

### Current README.md splits into:
- **README.md**: Rewritten as slim landing page with prominent docs site link
- **reference/tools.md**: Tool reference tables (consolidated)
- **reference/caching.md**: TTL reference table
- **explanation/architecture.md**: Package structure, data flow
- **how-to/contribute.md**: Development standards, contributing workflow

### New content to write:
- **tutorials/getting-started.md**: Step-by-step walkthrough (~90 min)
- **how-to/add-a-new-tool.md**: Developer recipe with concrete example
- **explanation/design-decisions.md**: Rationale for key technical choices

### Deleted after migration:
- **docs/configuration.md** (old monolithic file)

## README Rewrite

The README (~50 lines) contains:
1. Title + badges (LM Studio, license, build)
2. Prominent link to the hosted documentation site
3. One-paragraph description
4. Quick start (3 steps: Docker, API key, configure client)
5. Feature highlights (short bullet list)
6. Documentation section links table
7. License + Acknowledgments

## Content Guidelines

### Tutorials (tutorials/)
Learning-oriented. Second person, imperative steps, one golden path. Checkpoints to verify progress. No alternatives or options.

### How-to (how-to/)
Task-oriented. Assumes existing knowledge. Numbered steps, covers variations. Each page solves one specific problem.

### Reference (reference/)
Information-oriented. Tables, complete specs, no advice. Austere, accurate, comprehensive.

### Explanation (explanation/)
Understanding-oriented. Prose, context, alternatives considered, trade-offs. References other sections.

### _index.md pages
Brief (5-10 lines) section overview with list of pages in the section.

### Cross-linking
Each page links to related pages in other categories. Tutorials reference the full specs. How-to guides link to relevant reference pages. Explanations provide context for why reference specs are the way they are.

## Implementation Agent Assignments

Each task uses the appropriate diataxis subagent:

| Task | Agent |
|------|-------|
| Section index pages (_index.md) | diataxis-docs:diataxis-orchestrator |
| tutorials/getting-started.md | diataxis-docs:doc-tutorial-writer |
| how-to/configure-mcp-clients.md | diataxis-docs:doc-howto-writer |
| how-to/add-a-new-tool.md | diataxis-docs:doc-howto-writer |
| how-to/contribute.md | diataxis-docs:doc-howto-writer |
| reference/tools.md | diataxis-docs:doc-reference-gen |
| reference/api-scopes.md | diataxis-docs:doc-reference-gen |
| reference/caching.md | diataxis-docs:doc-reference-gen |
| reference/configuration.md | diataxis-docs:doc-reference-gen |
| explanation/architecture.md | diataxis-docs:doc-explanation-writer |
| explanation/design-decisions.md | diataxis-docs:doc-explanation-writer |
| Cross-link validation | diataxis-docs:doc-crosslink-validator |

## Verification

1. All content from current README.md and docs/configuration.md is accounted for (nothing lost)
2. No duplication between files
3. Each page follows its diataxis type rules
4. Cross-links are valid
5. Hugo _index.md files present in every directory
6. README prominently links to docs site
7. `docs/configuration.md` (old) is deleted
