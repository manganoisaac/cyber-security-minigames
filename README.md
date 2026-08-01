# GU Cyber Minigames

A suite of browser-based cybersecurity mini-games for Griffith University Open Days.
It's a single Go web server that renders HTML with [templ](https://templ.guide) and adds
interactivity with [htmx](https://htmx.org) — no database and no separate frontend build.
The whole thing runs as one container.

---

## What you need

Just **Docker Desktop**. No Go, no Node, nothing else - for both running *and* developing the app.

- **Windows / macOS:** install Docker Desktop → https://www.docker.com/products/docker-desktop/
- **Linux:** install Docker Engine + the Docker Compose plugin.

Make sure Docker Desktop is **running** before you start.

(There's one *optional* power-user dev mode further down that needs Go installed. You can ignore it unless confident,
everything works with Docker alone.)

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

You do **not** need Go installed for this, everything is built inside the container.

---

## Everyday commands

| I want to…                     | Command                        |
| ------------------------------ | ------------------------------ |
| Start the app                  | `docker compose up`            |
| Develop (auto-rebuild on save) | `docker compose watch`         |
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
  images/hero-banner.svg  hero banner image
go.mod / go.sum      Go dependencies
```

Two things worth knowing:

- **`.templ` files are the source; `*_templ.go` files are generated from them.** Never edit a
  `*_templ.go` file directly — edit the matching `.templ` and regenerate (see below).
- There is **no database**. Nothing about players is stored. That's deliberate: it keeps privacy
  obligations minimal and makes the app trivial to deploy anywhere.

---

## Developing (auto-rebuild on save) — only Docker

When you're **editing code**, use `docker compose watch` instead of `docker compose up`:

```bash
docker compose watch
```

Now edit any `.templ`, `.go`, or file under `public/`, save, and Docker automatically rebuilds and
restarts the app. Refresh **http://localhost:3000** to see the change. Stop with **Ctrl + C**.

This needs **only Docker** — Docker Desktop watches your files on the host side and rebuilds, which
works reliably on Windows, Mac, and Linux. (It needs Docker Compose v2.22 or newer; check with
`docker compose version`.) The first rebuild after a change takes a few seconds; thanks to build
caching, later ones are faster.

> **Tip:** keep the `docker compose watch` terminal open while you work. When you save, watch the
> log — you'll see it rebuild, then "listening on port 3000" — then refresh the browser.

### Optional: instant hot reload with auto browser refresh (needs Go)

If you're doing heavy front-end work and want *sub-second* reload with the browser refreshing on its
own (no manual refresh), there's a faster loop, but it needs **Go installed**
(https://go.dev/dl/ — 1.26+). Stop Docker first (`docker compose down`), since both use port 3000,
then from the repo root:

```bash
go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate --watch --proxy="http://localhost:3000" --cmd="go run ./cmd/server" --open-browser=false
```

Open **http://localhost:7331** (that's templ's reload proxy, not 3000). Shortcuts for the same
command: `make dev` (macOS/Linux, or Windows via `winget install ezwinports.make`), or `dev` on
Windows (runs `dev.cmd`). Wait for `listening on port 3000` in the terminal before refreshing.

> Live reload only injects itself into pages with a closing `</body>` tag, so keep that in the layout.

If you edited a `.templ` and want to regenerate the Go code once without running any watcher (needs
Go): `go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate`.

---

## Adding a new activity / mini-game

Each activity is a page (or a few pages) plus its handlers:

1. Add a `.templ` under `views/pages/` for the activity's screen(s).
2. Add the routes / logic in `internal/handlers/`.
3. Register the route where the server wires up handlers in `cmd/server/main.go`.
4. Run `make generate` (or use `make dev`) and check it in the browser.

Keep each activity self-contained so it can be played on its own — that's a core design
requirement.

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
- **`docker compose watch` says it's not a valid command** — your Docker Compose is older than
  v2.22. Update Docker Desktop. In the meantime you can still develop with `docker compose up --build`
  after each change.
- **`docker compose watch` isn't picking up my saves** — make sure you started it with `watch`
  (not `up`), and give it a few seconds — it rebuilds on each save. Editing a file it's told to
  ignore (`*.md`, `.env`, `tmp/`, `.git/`) won't trigger a rebuild, by design.
- **Edited a `.templ`/`.go` but nothing changed** — if you're using `docker compose up`, it doesn't
  auto-rebuild; use `docker compose watch`, or stop and run `docker compose up --build`.
- **Old code still showing after `git pull`** — force a rebuild: `docker compose up --build`.

Optional Go hot-reload path only:

- **Live dev server won't load at localhost:7331** — the app behind the proxy couldn't start,
  almost always because **Docker is still running and holding port 3000**. Stop it
  (`docker compose down`) and start the dev server again. Also wait for `listening on port 3000`
  in the terminal before refreshing.
- **`go` is not recognised / `make` is not recognised** — you don't have Go (or make) installed.
  This optional path needs Go; the rest of the project doesn't. Use `docker compose watch` instead,
  or install Go from https://go.dev/dl/ (and open a fresh terminal afterwards).

---

## Deploying

Deployment to the client's server is coordinated by the Integration / Deployment lead. Because the
app is a single container with no database, deploying is essentially "build the image and run it on
the target host."
