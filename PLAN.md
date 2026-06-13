# Markdown Output Fix and Library Update

## Context

The user reports that when a step has `output_type: "markdown"`, only the first few lines of the rendered markdown are visible ("Project README" and "Overview"). The rest of the content appears to be missing or not scrollable. The user also wants all charmbracelet libraries updated to their latest v2 versions.

## Current Issues

1. **Library versions**: Using v1 import paths (`github.com/charmbracelet/bubbletea`, `github.com/charmbracelet/lipgloss`, etc.). The latest stable versions are v2 under `charm.land/` import paths.
2. **Markdown rendering**: Glamour produces ANSI-styled output (~19 lines for the example). The viewport only shows ~7-10 lines at a time. The user needs to scroll with `PgUp`/`PgDown` to see the rest, but the viewport's default keymap handling is not wired up correctly.
3. **v2 API changes**: `tea.View` returns a `tea.View` struct instead of `string`. `tea.Cursor` is explicit. `viewport` has updated APIs.

## Latest Versions

- `charm.land/bubbletea/v2@v2.0.7`
- `charm.land/lipgloss/v2@v2.0.4`
- `charm.land/bubbles/v2@v2.1.0`
- `charm.land/glamour/v2@v2.0.1`

## Approach

1. **Update imports**: Change all `github.com/charmbracelet/bubbletea` to `charm.land/bubbletea/v2`, `github.com/charmbracelet/lipgloss` to `charm.land/lipgloss/v2`, etc.
2. **Update API usage**: 
   - `View()` returns `tea.View` instead of `string`
   - `tea.Place` / `tea.NewView` for compositing
   - `Cursor()` method for cursor positioning
   - Update `viewport` to v2 API
3. **Fix scroll handling**: Add explicit `pgup`/`pgdown`/`home`/`end` cases in the key switch that directly call viewport methods.
4. **Verify markdown rendering**: Ensure glamour v2 renders correctly and viewport scrolls properly.

## Files to modify

- `go.mod` / `go.sum` — update to v2 module paths
- `main.go` — update import and `tea.NewProgram` call
- `app.go` — update imports, `View()` signature, scroll handling, viewport usage
- `session.go` — no charmbracelet imports, no changes needed
- `workflow.go` — no charmbracelet imports, no changes needed
- `runner.go` — no charmbracelet imports, no changes needed
- `app_test.go` — update imports if needed

## Reuse

- Existing application logic (model, update, key handling)
- Existing `renderMarkdown` function (glamour v2 API may differ slightly)

## Steps

- [ ] Update go.mod with v2 module paths
- [ ] Update all imports in main.go, app.go, app_test.go
- [ ] Update `View()` to return `tea.View`
- [ ] Add explicit `pgup`/`pgdown`/`home`/`end` key handling
- [ ] Fix viewport creation and scroll handling for v2 API
- [ ] Test build and verify markdown rendering

## Verification

Build the app, run `./tui-workflow examples/markdown.json`, press `r` to run the step, then use `PgDown` to scroll through the rendered markdown output.
