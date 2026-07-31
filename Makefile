# Convenience shortcuts.
#   `make up`, `make rebuild`, `make down`, `make logs`  -> only need Docker.
#   `make dev`, `make generate`                          -> need Go installed locally.
#
# Windows without `make`? Just run the command after the ":" for the target you want.

.PHONY: up rebuild build down logs generate dev

# Run the app in Docker. Open http://localhost:3000
up:
	docker compose up

# Rebuild the image, then run. Use this after pulling new code.
rebuild:
	docker compose up --build

# Build the image without running.
build:
	docker compose build

# Stop and remove the containers.
down:
	docker compose down

# Follow the logs.
logs:
	docker compose logs -f

# Regenerate *_templ.go from *.templ (needed after editing a .templ, if not using `make dev`).
generate:
	go run github.com/a-h/templ/cmd/templ generate

# Live-reload dev server (needs Go). Edits to .templ/.go refresh the browser automatically.
# Open http://localhost:7331  (the app itself runs on :3000; 7331 is templ's reload proxy).
dev:
	go run github.com/a-h/templ/cmd/templ generate \
		--watch \
		--proxy="http://localhost:3000" \
		--cmd="go run ./cmd/server" \
		--open-browser=false
