@echo off
REM Live-reload dev server for Windows (needs Go installed).
REM Stop Docker first (docker compose down) - both use port 3000.
REM Open http://localhost:7331 once it prints "listening on port 3000".
go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate --watch --proxy="http://localhost:3000" --cmd="go run ./cmd/server" --open-browser=false
