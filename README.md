# GU Cyber Minigames

A suite of browser-based cybersecurity mini-games for Griffith University Open Days.
It's a single Go web server that renders HTML with [templ](https://templ.guide) and adds
interactivity with [htmx](https://htmx.org) — no database and no separate frontend build.
The whole thing runs as one container.

---

## What you need

Just **Docker Desktop**. No Go, no Node, nothing else.

- **Windows / macOS:** install Docker Desktop → https://www.docker.com/products/docker-desktop/
- **Linux:** install Docker Engine + the Docker Compose plugin.

Make sure Docker Desktop is **running** before you start.

---

## Quick start

```bash
# 1. Clone the repo
git clone <REPO_URL>
cd <REPO_FOLDER>

# 2. Start it
docker compose up
```

Then open **http://localhost:3000** in your browser.

The first run takes a minute or two while it downloads the base image and builds. Later runs
are much faster. Press **Ctrl + C** in the terminal to stop it.

You do **not** need Go installed for this — everything is built inside the container.

---

## Everyday commands

| I want to…                     | Command                        |
| ------------------------------ | ------------------------------ |
| Start the app                  | `docker compose up`            |
| Start after pulling new code   | `docker compose up --build`    |
| Stop it                        | `Ctrl + C` (or `docker compose down` in another terminal) |
| See the logs                   | `docker compose logs -f`       |

If you have `make`, there are shortcuts: `make up`, `make rebuild`, `make down`, `make logs`.

> **Pulled new changes and something looks off?** Rebuild with `docker compose up --build`.
> Docker reuses the previous build until you tell it to rebuild.

---

## Changing the port

If port 3000 is already in use on your machine, copy the example env file and pick another:

```bash
cp .env.example .env
# then edit .env, e.g. HOST_PORT=9090  ->  open http://localhost:9090
```

> This only changes the port on *your* machine. The app itself still listens on `:3000` inside
> the container.

### Making the port configurable (optional, recommended for deployment)

`main.go` currently hardcodes `:3000`. For deployment flexibility (the client's server may want a
specific port), have it read the `PORT` environment variable with `:3000` as the default:

```go
port := os.Getenv("PORT")
if port == "" {
    port = "3000"
}
fmt.Println("listening on port " + port)
http.ListenAndServe(":"+port, nil)
```

The Docker setup already passes `PORT`, so once this change is in, the port can be set without
touching code. Nothing breaks if you leave it as-is — everything defaults to 3000 either way.

---

## Project layout

```
cmd/server/          main.go — starts the web server
internal/handlers/   games.go — request handlers (routes / game logic)
views/
  layouts/           base.templ — shared page shell
  pages/             landing.templ etc. — individual pages
  *_templ.go         GENERATED from the .templ files — do not edit by hand
public/
  css/styles.css     styles
  js/htmx.js         htmx library
go.mod / go.sum      Go dependencies
```

Two things worth knowing:

- **`.templ` files are the source; `*_templ.go` files are generated from them.** Never edit a
  `*_templ.go` file directly — edit the matching `.templ` and regenerate (see below).
- There is **no database**. Nothing about players is stored. That's deliberate: it keeps privacy
  obligations minimal and makes the app trivial to deploy anywhere.

---

## Active development (live reload)

`docker compose up` is great for *running* the app. But if you're **editing code** and want the
browser to refresh automatically, use the live-reload dev server. This one needs **Go installed**
(https://go.dev/dl/ — version 1.26 or newer):

```bash
make dev
# or, without make:
go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate \
    --watch --proxy="http://localhost:3000" --cmd="go run ./cmd/server" --open-browser=false
```

Then open **http://localhost:7331** — note that's **7331**, not 3000. That's templ's live-reload
proxy. Edit a `.templ` or `.go` file, save, and the page refreshes on its own.

> Live reload only injects itself into pages with a closing `</body>` tag, so keep that in the layout.

If you edited a `.templ` and are **not** using live reload, regenerate the Go code once:

```bash
make generate
# or: go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate
```

---

## Adding a new activity / mini-game

Each activity is a page (or a few pages) plus its handlers:

1. Add a `.templ` under `views/pages/` for the activity's screen(s).
2. Add the routes / logic in `internal/handlers/`.
3. Register the route where the server wires up handlers in `cmd/server/main.go`.
4. Run `make generate` (or use `make dev`) and check it in the browser.

Keep each activity self-contained so it can be played on its own — that's a core design
requirement (the "Previously…" recap model).

---

## Git workflow (reminder)

Per the team agreement: work on a **feature branch**, push, and open a **pull request** to the
lead developer (Isaac) for review before it's merged to `main`. Don't commit straight to `main`.

---

## Troubleshooting

- **"Cannot connect to the Docker daemon" / errors on `docker compose`** — Docker Desktop isn't
  running. Start it and retry.
- **Page won't load / "connection refused" at localhost:3000** — confirm the container is up
  (`docker compose ps`). The app binds to `:3000` inside the container; if that ever changes in
  `cmd/server/main.go`, update the right-hand port in `docker-compose.yml` and the `EXPOSE`/health
  check to match.
- **Port 3000 already in use** — set `HOST_PORT` in `.env` (see "Changing the port").
- **Edited a `.templ` but nothing changed** — regenerate: `make generate`, or run with `make dev`.
  If you're on the Docker path, rebuild: `docker compose up --build`.
- **Old code still showing after `git pull`** — force a rebuild: `docker compose up --build`.

---

## Deploying

Deployment to the client's server is coordinated by the Integration / Deployment lead. Because the
app is a single container with no database, deploying is essentially "build the image and run it on
the target host." Check with Jeremy before pushing anything to a live server.
