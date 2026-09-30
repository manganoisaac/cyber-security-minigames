# Caught in the Net

**Caught in the Net** is a browser-based cybersecurity learning experience developed for Griffith University. Players take the role of a private investigator and work through a connected cybercrime case while completing six interactive cybersecurity mini-games.

The application is a single Go web server using [templ](https://templ.guide/) for server-rendered HTML, [htmx](https://htmx.org/) for partial page interactions, and vanilla JavaScript/CSS for game-specific behaviour. It uses no database and is packaged to run as a single Docker container.

## Activities

The final application contains six activities:

1. **Passwords & Social Engineering** — investigate personal information and use discovered clues to reconstruct a predictable password.
2. **Online Privacy** — analyse social-media posts and identify information that contributes to a person's digital footprint.
3. **Phishing** — identify a phishing email, then analyse the warning signs within the message.
4. **Encryption** — use a simplified five-letter key puzzle to demonstrate the relationship between plaintext, ciphertext and an encryption key.
5. **Network Defence** — deploy complementary security controls against different threat types and demonstrate defence in depth.
6. **Digital Investigation** — examine evidence, construct a timeline and compare investigative theories to determine what happened.

The full story can be played from **Start Investigation**, while **Chapter Select** allows individual activities to be opened directly. Each activity is connected to story/comic content and a post-game **What Have We Learned?** page.

---

## Running the handover package

### Requirements

For normal use, only Docker is required.

- **Windows / macOS:** install and run [Docker Desktop](https://www.docker.com/products/docker-desktop/).
- **Linux:** install Docker Engine and the Docker Compose plugin.

You do **not** need Go, templ, Node.js or any other development tools just to run the application.

### Start the application

1. Extract the supplied ZIP file.
2. Open a terminal in the extracted `cyber-security-minigames` folder — the folder containing `docker-compose.yml` and `Dockerfile`.
3. Run:

```bash
docker compose up --build
```

4. When the server is running, open:

```text
http://localhost:3000
```

The first build can take a few minutes because Docker may need to download the required base images and build the application. Later starts are usually much faster.

To stop the application, press **Ctrl + C** in the terminal. You can also stop/remove the container with:

```bash
docker compose down
```

For later runs, if the source has not changed, this is normally enough:

```bash
docker compose up
```

### If `docker` is not recognised

If the terminal reports that `docker` is not recognised as a command, Docker Desktop is either not installed or its command-line tools are not currently available.

Install/start Docker Desktop, wait until Docker reports that it is running, then close and reopen the terminal and check:

```bash
docker --version
docker compose version
```

---

## Changing the browser-facing port

The container listens on port `3000` internally. If port `3000` is already in use on the host machine, copy `.env.example` to `.env` and change `HOST_PORT`.

Example:

```text
HOST_PORT=9090
```

Then start the application normally and open:

```text
http://localhost:9090
```

The Compose configuration maps `${HOST_PORT:-3000}` on the host to port `3000` in the container.

---

## Player state and reset behaviour

There is **no database**. Player progress is stored only in server memory for the current browser session.

- A random `HttpOnly` session cookie named `cyber_session` identifies the browser session.
- Session data expires after **30 minutes of inactivity**.
- Restarting the server/container clears all in-memory session state.
- Returning to the landing page, story introduction or Chapter Select resets all minigame progress for that browser session so a new player does not inherit the previous player's state.
- Privacy and Phishing are also designed to restart cleanly on a full page refresh, making them convenient for repeated classroom or event use.

No player accounts or persistent personal data are stored by the application.

---

## Technology stack

- **Go 1.26.5** — HTTP server and server-side game logic
- **templ v0.3.1020** — type-safe HTML templates
- **htmx** — server-driven partial page updates
- **Vanilla JavaScript** — client-side minigame behaviour where appropriate
- **CSS** — game-specific and shared styling
- **Docker / Docker Compose** — build and deployment

The application runs as a single container and does not require a separate frontend service or database.

---

## Project structure

```text
cmd/server/
  main.go                     Application entry point and game registration

data/
  *.go                        Static game content and instructional data

internal/
  handlers/
    comic/                    Story comic routes
    password/                 Password activity handlers/state/service
    privacy/                  Online Privacy handlers/state/service
    email_client/             Phishing handlers/state/service
    encryption/               Encryption page handler
    network_defense/          Network Defence page handler
    digitalinvestigation/     Digital Investigation page handlers
    router.go                 Core routes, static assets and game registration
  session/                    In-memory browser session management

views/
  layouts/                    Shared page layout
  pages/                      Full-page templ views
  pages/comic/                Story/comic pages
  pages/whwl/                 What Have We Learned? pages
  components/                 Reusable UI and game components

public/
  css/                        Shared and game-specific styles
  js/                         htmx and game-specific JavaScript
  images/                     Application/game artwork and UI assets
  fonts/                      Local font assets
  data/                       Static browser-side data files

types/                        Shared Go data types

Dockerfile                    Multi-stage production image build
docker-compose.yml            Local/run configuration
.env.example                  Optional host-port configuration
go.mod / go.sum               Go module/dependency definitions
Makefile                      Convenience development commands
```

---

## templ source and generated files

The `.templ` files are the source templates. Files ending in `_templ.go` are generated Go code and should **not** be edited manually.

The handover package includes the generated `_templ.go` files, but the Docker build also runs:

```bash
templ generate
```

before compiling the application. This ensures generated code matches the `.templ` source during a Docker build.

After manually editing a `.templ` file outside Docker, regenerate the Go output with:

```bash
go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate
```

---

## Development

### Docker development workflow

Docker Compose can rebuild the application automatically when source files are changed:

```bash
docker compose watch
```

This requires Docker Compose v2.22 or newer. Check with:

```bash
docker compose version
```

The watch configuration rebuilds when application source changes and ignores development-only files such as `.git/`, `tmp/`, Markdown files and `.env`.

### Local Go development (optional)

Local development without Docker requires a compatible Go installation.

Regenerate templ files:

```bash
go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate
```

Run the server:

```bash
go run ./cmd/server
```

The server uses the `PORT` environment variable when supplied and defaults to `3000`.

For templ live reload:

```bash
go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate --watch --proxy="http://localhost:3000" --cmd="go run ./cmd/server" --open-browser=false
```

Then open:

```text
http://localhost:7331
```

Stop Docker first if it is already using port `3000`.

---

## Common Docker commands

| Task | Command |
| --- | --- |
| First run / rebuild after source changes | `docker compose up --build` |
| Start an existing build | `docker compose up` |
| Development auto-rebuild | `docker compose watch` |
| Build without starting | `docker compose build` |
| Stop and remove containers | `docker compose down` |
| Follow container logs | `docker compose logs -f` |

If `make` is available, equivalent shortcuts are provided in the `Makefile` (`make up`, `make watch`, `make rebuild`, `make down`, `make logs`).

---

## Troubleshooting

### `docker` is not recognised

Docker Desktop is not installed, is not running, or the terminal was opened before Docker was installed. Start/install Docker Desktop, then open a new terminal and run:

```bash
docker --version
docker compose version
```

### Cannot connect to the Docker daemon

Docker is installed but the Docker service/Desktop application is not running. Start Docker Desktop and retry.

### `http://localhost:3000` does not load

Check that the container is running:

```bash
docker compose ps
```

Then inspect logs if required:

```bash
docker compose logs -f
```

### Port 3000 is already in use

Set a different `HOST_PORT` in `.env`, for example:

```text
HOST_PORT=9090
```

Then open `http://localhost:9090`.

### Source changes are not appearing

If using normal `docker compose up`, Docker does not automatically rebuild the application. Stop it and run:

```bash
docker compose up --build
```

For active development, use:

```bash
docker compose watch
```

### A `.templ` change does not appear during local Go development

Regenerate the generated Go files before restarting the server:

```bash
go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate
```

---

## Deployment notes

The project is intentionally self-contained:

- one Go application;
- one Docker container;
- static assets included in the image;
- no database;
- no external frontend build step;
- no machine-specific filesystem paths required at runtime.

The repository/ZIP can therefore be extracted to any location on the target machine and built from the project root with Docker Compose.

The Dockerfile uses a multi-stage build: templ code is generated and the Go server is compiled in the build stage, then only the compiled server and `public/` assets are copied into the final runtime image.

---

## Source control

The final source is maintained in the project Git repository. Git is **not required** to run the handover ZIP, but future development can continue through normal feature branches and pull requests before merging changes into `main`.
