#!/usr/bin/env bash
# Scaffolds a blank new minigame: handler package, templ page, CSS, and wiring.
#
# Usage:  ./scripts/new_game.sh <game-name> [--route /path]
#
#   <game-name>   kebab- or snake-case, e.g. password-breach / password_breach
#   --route       mount path (defaults to /<game-name>)
#
# Run from the repo root. The script creates the handler package, page,
# feature CSS, registers the CSS import, and adds the game to the explicit
# games list in cmd/server/main.go. It does NOT run `make generate` — run it
# yourself after scaffolding.

set -euo pipefail

if [[ $# -lt 1 ]]; then
	echo "usage: new_game.sh <game-name> [--route /path]" >&2
	exit 1
fi

if [[ ! -f go.mod ]]; then
	echo "error: run this script from the repo root (go.mod not found)" >&2
	exit 1
fi

NAME=$(printf '%s' "$1" | tr '[:upper:]' '[:lower:]')

ROUTE=""
if [[ ${2:-} == "--route" ]]; then
	ROUTE="${3:-}"
	if [[ -z "$ROUTE" || "$ROUTE" != /* ]]; then
		echo "error: --route must be a path starting with '/'" >&2
		exit 1
	fi
fi

SNAKE=$(printf '%s' "$NAME" | tr '-' '_')
KEBAB=$(printf '%s' "$NAME" | tr '_' '-')
ROUTE=${ROUTE:-"/$KEBAB"}

# pascal: hello-world -> HelloWorld
pascal() {
	printf '%s' "$1" | sed -E 's/(^|[-_])([[:alnum:]])/\u\2/g'
}
PASCAL=$(pascal "$NAME")
# lowerFirst: HelloWorld -> helloWorld (for the unexported handler struct)
PC_LOWER=$(printf '%s' "$PASCAL" | sed -E 's/^./\l&/')
# title: HelloWorld -> Hello World (for the page <title>)
TITLE=$(printf '%s' "$PASCAL" | sed -E 's/([[:lower:][:digit:]])([[:upper:]])/\1 \2/g')

MODULE=$(sed -n 's/^module[[:space:]]\+//p' go.mod | tr -d '\r')
if [[ -z "$MODULE" ]]; then
	echo "error: could not determine module name from go.mod" >&2
	exit 1
fi

HANDLER_DIR="internal/handlers/$SNAKE"
PAGE_TEMPL="views/pages/$SNAKE.templ"
CSS_FILE="public/css/features/$KEBAB/$KEBAB.css"

if [[ -e "$HANDLER_DIR" ]]; then
	echo "error: $HANDLER_DIR already exists" >&2
	exit 1
fi
if [[ -e "$PAGE_TEMPL" ]]; then
	echo "error: $PAGE_TEMPL already exists" >&2
	exit 1
fi
if [[ -e "$CSS_FILE" ]]; then
	echo "error: $CSS_FILE already exists" >&2
	exit 1
fi

# ---------------------------------------------------------------------------
# Handler package
# ---------------------------------------------------------------------------
mkdir -p "$HANDLER_DIR"

cat > "$HANDLER_DIR/routes.go" <<EOF
package $SNAKE

import "net/http"

// Register wires up the $KEBAB page and all of its htmx endpoints,
// delegating to the given handler.
func Register(mux *http.ServeMux, h I${PASCAL}Handler) {
	mux.HandleFunc("GET $ROUTE", h.Page)
}
EOF

cat > "$HANDLER_DIR/handlers.go" <<EOF
package $SNAKE

import "net/http"

// I${PASCAL}Handler is the HTTP handler contract for the $KEBAB minigame.
type I${PASCAL}Handler interface {
	Page(w http.ResponseWriter, r *http.Request)
}

// ${PC_LOWER}Handler implements I${PASCAL}Handler, delegating to the service.
type ${PC_LOWER}Handler struct {
	svc I${PASCAL}Service
}

// NewHandler wires the service into the HTTP layer.
func NewHandler(svc I${PASCAL}Service) I${PASCAL}Handler {
	return &${PC_LOWER}Handler{svc: svc}
}

// Page renders the $KEBAB minigame page.
func (h *${PC_LOWER}Handler) Page(w http.ResponseWriter, r *http.Request) {
	h.svc.Page(stateForRequest(r)).Render(r.Context(), w)
}
EOF

cat > "$HANDLER_DIR/service.go" <<EOF
package $SNAKE

import (
	"sync"

	"$MODULE/views/pages"

	"github.com/a-h/templ"
)

// I${PASCAL}Service owns the $KEBAB minigame's orchestration. Methods receive
// only the parameters they need; the HTTP handlers decode requests and render
// the returned components.
type I${PASCAL}Service interface {
	Page(st I${PASCAL}GameState) templ.Component
}

// Service implements I${PASCAL}Service.
type Service struct{}

// newService is the private constructor for Service.
func newService() *Service {
	return &Service{}
}

var (
	serviceOnce   sync.Once
	serviceGlobal I${PASCAL}Service
)

// ${PASCAL}Service returns the app-wide $KEBAB service singleton, constructed
// once on first use.
func ${PASCAL}Service() I${PASCAL}Service {
	serviceOnce.Do(func() { serviceGlobal = newService() })
	return serviceGlobal
}

func (s *Service) Page(st I${PASCAL}GameState) templ.Component {
	return pages.${PASCAL}Page()
}
EOF

cat > "$HANDLER_DIR/state.go" <<EOF
package $SNAKE

import (
	"net/http"
	"sync"

	"$MODULE/internal/session"
)

// stateKey is the name of this minigame's state slot inside the shared session.
const stateKey = "$KEBAB"

// I${PASCAL}GameState is the state contract for the $KEBAB minigame. Only
// Reset is shared with other games; every other operation is game-specific.
type I${PASCAL}GameState interface {
	session.IState
}

// gameState is the per-user state for the $KEBAB minigame. The server owns
// every piece of state; the HTML rendered to the client is just a projection.
type gameState struct {
	mu sync.Mutex
}

// newGameState returns the concrete state; callers receive it through
// I${PASCAL}GameState.
func newGameState() *gameState {
	return &gameState{}
}

// stateFor returns the $KEBAB minigame's state for a session, creating it on
// first use. The state is stored in the shared session so it survives
// switching to other minigames and back.
func stateFor(s session.ISession) I${PASCAL}GameState {
	if st, ok := s.State(stateKey); ok {
		return st.(I${PASCAL}GameState)
	}
	gs := newGameState()
	s.SetState(stateKey, gs)
	return gs
}

// stateForRequest returns the $KEBAB minigame's state for the request's
// session, as resolved into the context by session.Middleware.
func stateForRequest(r *http.Request) I${PASCAL}GameState {
	return stateFor(session.SessionFrom(r))
}

// Reset clears this minigame's per-run state.
func (s *gameState) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
}
EOF

cat > "$HANDLER_DIR/game.go" <<EOF
package $SNAKE

import (
	"net/http"

	"$MODULE/internal/handlers"
)

// game adapts the $KEBAB minigame to handlers.Game so it can be registered
// alongside the other minigames.
type game struct {
	handler I${PASCAL}Handler
}

// NewGame returns the $KEBAB minigame as a handlers.Game, ready to register
// its page and htmx endpoints onto the app mux.
func NewGame(h I${PASCAL}Handler) handlers.Game {
	return &game{handler: h}
}

// Register mounts the $KEBAB minigame's routes, delegating to Register.
func (g *game) Register(mux *http.ServeMux) {
	Register(mux, g.handler)
}
EOF

# ---------------------------------------------------------------------------
# Templ page
# ---------------------------------------------------------------------------
cat > "$PAGE_TEMPL" <<EOF
package pages

import "$MODULE/views/layouts"

templ ${PASCAL}Page() {
	@layouts.Base("$TITLE") {
		<main class="$KEBAB">
		</main>
	}
}
EOF

# ---------------------------------------------------------------------------
# CSS
# ---------------------------------------------------------------------------
mkdir -p public/css/features/$KEBAB

cat > "$CSS_FILE" <<EOF
.$KEBAB {
  --min-height: 100vh;
  min-height: var(--min-height);
  background-color: var(--bg);
  display: flex;
  flex-direction: column;
  align-items: center;
}
EOF

if ! grep -q "@import \"./features/$KEBAB/$KEBAB.css\";" public/css/index.css; then
	printf '@import "./features/%s/%s.css";\n' "$KEBAB" "$KEBAB" >> public/css/index.css
fi

# ---------------------------------------------------------------------------
# Wire into the explicit games list in cmd/server/main.go
# ---------------------------------------------------------------------------
EMAIL_IMPORT='internal/handlers/email_client"'
NEW_IMPORT="$MODULE/internal/handlers/$SNAKE"
EMAIL_ENTRY='email_client.NewGame(email_client.NewHandler(email_client.EmailService())),'
GAME_ENTRY="$SNAKE.NewGame($SNAKE.NewHandler($SNAKE.${PASCAL}Service())),"

if ! grep -q "$NEW_IMPORT" cmd/server/main.go; then
	if grep -q "$EMAIL_IMPORT" cmd/server/main.go; then
		sed -i "\|$EMAIL_IMPORT|a\\
	\"$NEW_IMPORT\"" cmd/server/main.go
	else
		echo "warning: could not find the email_client import in cmd/server/main.go" >&2
		echo "add this import manually:" >&2
		echo "  \"$NEW_IMPORT\"" >&2
	fi
fi

if ! grep -qF "$GAME_ENTRY" cmd/server/main.go; then
	if grep -qF "$EMAIL_ENTRY" cmd/server/main.go; then
		sed -i "\|$EMAIL_ENTRY|a\\
		$GAME_ENTRY" cmd/server/main.go
	else
		echo "warning: could not find the email_client entry in the games list" >&2
		echo "add this line to the games slice in cmd/server/main.go manually:" >&2
		echo "  $GAME_ENTRY" >&2
	fi
fi

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
echo
echo "Created:"
echo "  $HANDLER_DIR/routes.go"
echo "  $HANDLER_DIR/handlers.go"
echo "  $HANDLER_DIR/service.go"
echo "  $HANDLER_DIR/state.go"
echo "  $HANDLER_DIR/game.go"
echo "  $PAGE_TEMPL"
echo "  $CSS_FILE"
echo
echo "Wired:"
echo "  public/css/index.css"
echo "  cmd/server/main.go"
echo
echo "Next steps:"
echo "  1. make generate   (regenerate *_templ.go)"
echo "  2. Open /$KEBAB to see the blank page"
echo "  3. Add a chapter-select link in views/pages/chapter_select.templ"
