# AGENTS.md

Single Go web server (module `cyber-security-minigames`, Go 1.26.5) rendering HTML with templ + htmx. No database, no frontend build step.

## Commands
- `make generate` — regenerates `*_templ.go` from `.templ` sources. **Run after every `.templ` edit**; `go build` silently ignores `.templ` changes.
- `go build ./...` — the only verification step; there are no tests in this repo.
- `make dev` — templ watch + live-reload proxy: browser on `:7331`, app on `:3000` (needs Go, Docker stopped). Docker alternative: `docker compose watch`.

## Architecture
- Entrypoint `cmd/server/main.go` is a thin shell. **Routes and static-file serving are wired in `internal/handlers/games.go` (`Routes()`)** — register new routes there, not in main.go.
- Game content data lives in `data/` (`EmailContent`), shared DTO structs in `types/` (`EmailItemProps`: Title/Subject/Content/IsMalware), shared helpers in `utils/` (`textmanipulation.TruncateChars`).
- `views/*.templ` are the source of truth; sibling `*_templ.go` files are generated — never edit them by hand.
- Static assets served from `./public` under `/public/...` via `http.FileServer`.

## Conventions
- templ imports are one per line (`import "cyber-security-minigames/data"`), no grouped blocks. Use `templ.KV("class", cond)` for conditional classes and `{ expr }` for attribute values.
- Routes use Go 1.22 ServeMux patterns (`GET /email/{id}`); read params with `r.PathValue("id")`. ServeMux `{id}` matches a single segment.
- `<body hx-boost="true">` is set in `views/layouts/base.templ`. htmx is loaded from `public/js/htmx.js`.
- Email-client pattern: `GET /email-client` renders the page; clicking an item does `GET /email/{id}` which re-renders the whole `EmailList` fragment, swapped via `hx-target="#email-list-container" hx-swap="outerHTML"` (highlight + details in one round trip, selection defaulting to index 0).

## Workflow
- Feature branch + PR to Isaac for review; never commit to `main`.
