# Contributing to tilt-tui

tilt-tui is a small TUI built with [Bubbletea](https://github.com/charmbracelet/bubbletea) that talks to Tilt's local API. It's intentionally minimal — the codebase fits in a handful of files and the scope is narrow. Keep that in mind when proposing changes.

## Before opening a PR

Check if there's already an issue or discussion for what you want to do. For anything beyond a clear bug fix, open an issue first. It avoids situations where you put in work on something that won't be merged.

## Setting up

```bash
git clone https://github.com/jmeiracorbal/tt
cd tt
go build ./...
```

You don't need a running Kubernetes cluster or Tilt to work on the UI. The mock server in `../mock-server` serves a fake `/api/view` response on `localhost:10350` so you can iterate without any infrastructure:

```bash
# terminal 1
cd ../mock-server && go run .

# terminal 2
go run . 
```

The mock cycles resources through `building`, `ok`, and `error` states on different intervals, which is enough to exercise most UI paths.

## Where things live

```
internal/api/   HTTP client and types that mirror Tilt's /api/view response
internal/ui/    Bubbletea model, update loop, rendering
main.go         Entry point, flag parsing
```

The API types in `internal/api/types.go` map directly to Tilt's JSON schema. If Tilt changes its API, that's where it needs to be updated.

## What's in scope

- Bug fixes
- Keybinding improvements
- Rendering issues (wrapping, alignment, color)
- Performance on large resource lists or long logs

## What's probably not in scope

- Support for non-Tilt backends
- Configuration files or persistent state
- Mouse support (Bubbletea supports it but it adds complexity with little payoff for a keyboard-first tool)

If you're unsure, open an issue before building.

## Code style

Run `go vet ./...` before submitting. No linter config beyond that. Keep functions short, avoid abstractions that only have one callsite. Comments only when the code can't explain itself.

## Commits

Prefer small, focused commits. One logical change per commit. The message should say what changed and, if it's not obvious, why.
