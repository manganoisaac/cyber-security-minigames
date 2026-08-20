#!/usr/bin/env bash
# Adds a new HTTP handler method to an existing minigame, optionally backed by
# a new service method, following the email_client/password pattern.
#
# Usage:
#   ./scripts/new_handler.sh                      interactive (uses fzf)
#   ./scripts/new_handler.sh --game password --name Score \
#       --method POST --route /password/score --service yes
#
# Interactive mode uses fzf for the game, the HTTP method, and whether to add
# a service method; handler name and route path are typed. Flags supply the
# same values and skip the prompts entirely (handy for scripting and tests).
#
# Run from the repo root. The script only edits files inside
# internal/handlers/<game>/ — it does not touch cmd/server/main.go or game.go.
# It runs gofmt on the edited files, but NOT make generate; run that yourself
# if the new method renders a new templ component.

set -euo pipefail

GAME=""
NAME=""
METHOD=""
ROUTE=""
SERVICE=""

usage() {
	cat <<'EOF'
usage: scripts/new_handler.sh [--game GAME] [--name NAME] [--method VERB]
       [--route /path] [--service yes|no]

Adds a new HTTP handler method to an existing minigame under
internal/handlers/. With no flags, prompts interactively.

  --game    target game directory under internal/handlers (skips fzf)
  --name    handler method name, e.g. Score (skips the read prompt)
  --method  HTTP verb: GET, POST, PUT or DELETE (skips fzf)
  --route   mount path, e.g. /password/score (skips the read prompt)
  --service yes adds a I*Service method that the handler delegates to;
            no keeps all logic in the handler (skips fzf)
EOF
}

while [[ $# -gt 0 ]]; do
	case "$1" in
		--game) GAME=${2:-}; shift 2 ;;
		--name) NAME=${2:-}; shift 2 ;;
		--method) METHOD=${2:-}; shift 2 ;;
		--route) ROUTE=${2:-}; shift 2 ;;
		--service) SERVICE=${2:-}; shift 2 ;;
		-h | --help) usage; exit 0 ;;
		*) echo "error: unknown argument: $1" >&2; usage; exit 1 ;;
	esac
done

if [[ ! -f go.mod ]]; then
	echo "error: run this script from the repo root (go.mod not found)" >&2
	exit 1
fi

# ---------------------------------------------------------------------------
# 1. Pick the game package
# ---------------------------------------------------------------------------
GAMES=$(find internal/handlers -mindepth 1 -maxdepth 1 -type d -printf '%f\n' | sort)
if [[ -z "$GAMES" ]]; then
	echo "error: no game packages found under internal/handlers/" >&2
	exit 1
fi

if [[ -z "$GAME" ]]; then
	if ! command -v fzf >/dev/null 2>&1; then
		echo "error: fzf is required for interactive game selection" >&2
		echo "       pass --game <dir> to skip the picker" >&2
		exit 1
	fi
	GAME=$(printf '%s\n' "$GAMES" | fzf --prompt 'Game: ' --height 40% \
		--border-label 'new_handler.sh' \
		--header 'which minigame gets the new handler?')
fi

if [[ -z "$GAME" ]]; then
	echo "aborted: no game selected" >&2
	exit 1
fi
HANDLER_DIR="internal/handlers/$GAME"
if [[ ! -d "$HANDLER_DIR" ]]; then
	echo "error: $HANDLER_DIR is not a directory" >&2
	usage
	exit 1
fi

SNAKE=$GAME
KEBAB=${GAME//_/-}
pascal() {
	printf '%s' "$1" | sed -E 's/(^|[-_])([[:alnum:]])/\u\2/g'
}
PASCAL=$(pascal "$KEBAB")
# lowerFirst: EmailClient -> emailClient (unexported handler struct)
PC_LOWER=$(printf '%s' "$PASCAL" | sed -E 's/^./\l&/')

HANDLERS="$HANDLER_DIR/handlers.go"
SERVICE_FILE="$HANDLER_DIR/service.go"
ROUTES="$HANDLER_DIR/routes.go"
for f in "$HANDLERS" "$ROUTES"; do
	if [[ ! -f "$f" ]]; then
		echo "error: $f not found - does the game follow the email_client/password scaffold?" >&2
		exit 1
	fi
done

# ---------------------------------------------------------------------------
# 2. Handler method name
# ---------------------------------------------------------------------------
if [[ -z "$NAME" ]]; then
	while [[ -z "$NAME" ]]; do
		read -r -p 'Handler method name (e.g. Score, Submit): ' NAME
	done
fi
NAME=$(printf '%s' "$NAME" | tr -cd '[:alnum:]')
if [[ -z "$NAME" ]]; then
	echo "error: handler name must contain at least one letter or digit" >&2
	exit 1
fi
NAME=${NAME^}

if grep -qE "^[[:space:]]*${NAME}\(" "$HANDLERS"; then
	echo "error: $NAME already exists in $HANDLERS" >&2
	exit 1
fi

# ---------------------------------------------------------------------------
# 3. HTTP method
# ---------------------------------------------------------------------------
if [[ -z "$METHOD" ]]; then
	if ! command -v fzf >/dev/null 2>&1; then
		echo "error: fzf is required for HTTP method selection" >&2
		echo "       pass --method <VERB> to skip the picker" >&2
		exit 1
	fi
	METHOD=$(printf 'GET\nPOST\nPUT\nDELETE\n' | fzf --prompt 'HTTP method: ' \
		--height 40% --border-label 'new_handler.sh' \
		--header 'which HTTP verb?' --default 'GET' 2>/dev/null || printf 'GET\n')
fi
METHOD=$(printf '%s' "$METHOD" | tr '[:lower:]' '[:upper:]')
case "$METHOD" in
	GET | POST | PUT | DELETE) ;;
	*)
		echo "error: invalid HTTP method: $METHOD (want GET, POST, PUT or DELETE)" >&2
		exit 1
		;;
esac

# ---------------------------------------------------------------------------
# 4. Route path
# ---------------------------------------------------------------------------
DEFAULT_ROUTE=("/$KEBAB/$(printf '%s' "$NAME" | tr '[:upper:]' '[:lower:]')")
if [[ -z "$ROUTE" ]]; then
	read -r -e -p "Route path [$DEFAULT_ROUTE]: " -i "$DEFAULT_ROUTE" ROUTE
	ROUTE=${ROUTE:-$DEFAULT_ROUTE}
fi
if [[ "$ROUTE" != /* ]]; then
	echo "error: route must start with '/': $ROUTE" >&2
	exit 1
fi

if grep -qF "\"$METHOD $ROUTE\"" "$ROUTES"; then
	echo "error: \"$METHOD $ROUTE\" is already registered in $ROUTES" >&2
	exit 1
fi

# ---------------------------------------------------------------------------
# 5. Service method: yes (handler delegates) or no (all logic in the handler)
# ---------------------------------------------------------------------------
if [[ -z "$SERVICE" ]]; then
	if ! command -v fzf >/dev/null 2>&1; then
		echo "error: fzf is required for the service-method choice" >&2
		echo "       pass --service yes|no to skip the picker" >&2
		exit 1
	fi
	SERVICE=$(printf 'yes\nno\n' | fzf --prompt 'Add a service method? ' \
		--height 40% --border-label 'new_handler.sh' \
		--header 'yes = handler delegates to a new I*Service method' \
		--default 'yes' 2>/dev/null || printf 'yes\n')
fi
case "$SERVICE" in
	yes | no) ;;
	*)
		echo "error: --service must be yes or no, got: $SERVICE" >&2
		exit 1
		;;
esac

# ---------------------------------------------------------------------------
# Apply edits: handlers.go, then service.go (yes only), then routes.go
# ---------------------------------------------------------------------------
TAB=$'\t'

INSERT=$'\t'"${NAME}(w http.ResponseWriter, r *http.Request)"
H_LINE=$(grep -nE '^[[:space:]]*Page\(' "$HANDLERS" | tail -1 | cut -d: -f1 || true)
if [[ -z "$H_LINE" ]]; then
	echo "error: could not find the Page method in the $NAME handler interface" >&2
	echo "add this line to I${PASCAL}Handler manually:" >&2
	echo "  ${INSERT}" >&2
else
	sed -i "${H_LINE}a\\
${INSERT}" "$HANDLERS"
fi

if [[ "$SERVICE" == yes ]]; then
	if [[ ! -f "$SERVICE_FILE" ]]; then
		echo "error: $SERVICE_FILE missing; add the service method manually" >&2
	else
		S_INSERT=$'\t'"${NAME}(st I${PASCAL}GameState) templ.Component"
		S_LINE=$(grep -nF "Page(st I${PASCAL}GameState) templ.Component" "$SERVICE_FILE" | head -1 | cut -d: -f1 || true)
		if [[ -z "$S_LINE" ]]; then
			echo "error: could not find the Page signature in I${PASCAL}Service" >&2
			echo "add this line manually:" >&2
			echo "  ${S_INSERT}" >&2
		else
			sed -i "${S_LINE}a\\
${S_INSERT}" "$SERVICE_FILE"
		fi
		cat >> "$SERVICE_FILE" <<EOF

func (s *Service) ${NAME}(st I${PASCAL}GameState) templ.Component {
	return nil
}
EOF
	fi
fi

if [[ "$SERVICE" == yes ]]; then
	cat >> "$HANDLERS" <<EOF

// ${NAME} handles ${METHOD} ${ROUTE}, delegating to the service.
func (h *${PC_LOWER}Handler) ${NAME}(w http.ResponseWriter, r *http.Request) {
	h.svc.${NAME}(stateForRequest(r)).Render(r.Context(), w)
}
EOF
else
	cat >> "$HANDLERS" <<EOF

// ${NAME} handles ${METHOD} ${ROUTE}.
func (h *${PC_LOWER}Handler) ${NAME}(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
EOF
fi

R_INSERT=$'\t'"mux.HandleFunc(\"$METHOD $ROUTE\", h.${NAME})"
R_LINE=$(grep -nF 'mux.HandleFunc' "$ROUTES" | tail -1 | cut -d: -f1 || true)
if [[ -z "$R_LINE" ]]; then
	echo "error: could not find a mux.HandleFunc call in $ROUTES" >&2
	echo "add this line inside Register manually:" >&2
	echo "  ${R_INSERT}" >&2
else
	sed -i "${R_LINE}a\\
${R_INSERT}" "$ROUTES"
fi

# ---------------------------------------------------------------------------
# gofmt the touched files
# ---------------------------------------------------------------------------
if command -v gofmt >/dev/null 2>&1; then
	gofmt -w "$HANDLERS" "$ROUTES"
	[[ "$SERVICE" == yes ]] && gofmt -w "$SERVICE_FILE"
fi

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
echo
echo "Added (${GAME}):"
echo "  ${METHOD} ${ROUTE}  ->  h.${NAME}"
echo
echo "Edited:"
echo "  $HANDLERS"
echo "  $ROUTES"
[[ "$SERVICE" == yes ]] && echo "  $SERVICE_FILE"
echo
if [[ "$SERVICE" == yes ]]; then
	echo "Next steps:"
	echo "  1. Implement the ${NAME} stub in $SERVICE_FILE (it returns nil now)"
	echo "  2. go build ./..."
	echo "  3. make generate  (only if the method renders a new templ component)"
else
	echo "Next steps:"
	echo "  1. Implement the ${NAME} body in $HANDLERS (returns 501 not implemented now)"
	echo "  2. go build ./..."
	echo "  3. make generate  (only if the method renders a new templ component)"
fi
