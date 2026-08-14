## Rules

### CSS Organization

- All CSS lives in `public/css/`.
- `index.css` contains **only** `@import` statements (Google Fonts first, then `./themes.css`, then feature files).
- `themes.css` contains only `:root` CSS variables, `@font-face` declarations, and global element/utility styles.
- Feature-specific tokens may be scoped on the feature's container selector in its own `features/*.css` file (e.g. `--details-width` on `.email-client__list`); the shared `--size-*` scale stays in `themes.css`.
- Each feature/component gets its own file in `public/css/features/<feature>.css` (e.g. `email-details.css`, `chapter-select.css`). A page and a component may share a block prefix only if they also share the component name (e.g. `email_client.templ`); otherwise the page gets its own block (`email-client-page`).
- Image URLs inside `public/css/features/*` must use `../../images/...`; in `themes.css` use `../fonts/...`.
- New CSS must be wired into `index.css` via `@import` — never added inline to templates or another file.
- Use the BEM naming convention (classic `Block__element--modifier`) for all CSS classes, e.g. `.email-item__text`, `.chapter-select--1`.

### File Naming

- Go source files (`.go`), templ source files (`.templ`), and their generated `*_templ.go` files use `snake_case` — multi-word names joined with underscores. Underscores are the de-facto Go convention (`lowercase_with_underscores.go`) and match templ's generated `_templ.go` suffix, e.g. `email_client.go`, `chapter_select.templ`, `email_list_templ.go`.
- CSS files use kebab-case (hyphens) so filenames match the BEM classes they style, e.g. `features/email-details.css` styles `.email-details__*`, `features/chapter-select.css` styles `.chapter-select`.
- JS asset files use kebab-case to match conventional web-asset naming and the CSS they complement, e.g. `public/js/htmx-extensions.js`.

### Sizing Tokens

- Recurring sizes use semantic `:root` tokens in `themes.css`: `--size-xs` (0.6rem), `--size-s` (0.8rem), `--size-m` (1rem), `--size-l` (1.2rem), `--size-xl` (1.5rem), `--size-2xl` (1.7rem), `--size-3xl` (2rem).
- Use `var(--size-*)` for matching `font-size` and spacing (`margin`/`padding`/`gap`) values that fit the scale.
- All other sizes (e.g. `2.5rem`, `3rem`, `4rem`, `5.5rem`, `px` dimensions, and `vw`/`vh`/`%` values) must be declared as scoped variables on their feature's container selector in `features/*.css` — no hardcoded values in rule bodies.
- Use `calc()` when a size is derived from or related to another size (e.g. `calc(var(--size-m) * 2)`) instead of hardcoding the result.

### Templ Components

- Components receive data as typed parameters (prop drilling); never read package globals (e.g. `data.Emails`) from inside a component — pass the data in from the caller.
- Pass page-scoped data explicitly down the tree; reserve `context` for genuinely app-wide data (user/session/theme).
- Use `context` sparingly — Go's `context` is untyped and errors surface at runtime; prefer explicit params.
- Layouts wrap content with `{ children... }`; components that render slots accept a `templ.Component` parameter.

### htmx

- Keep the hypermedia model: the server returns HTML fragments, not JSON — state lives in the HTML.
- Express interactivity in markup with `hx-*` attributes (`hx-get`/`hx-post`, `hx-target`, `hx-swap`, `hx-trigger`); avoid hand-written JS.
- htmx requests should render only the affected component/fragment and swap it in (use `templ.WithFragments` + `templ.Fragment` when a page has multiple fragments).
- Use the `HX-Request` header to distinguish fragment requests from full page loads.
- Use `hx-boost` on links/forms for smooth navigation; wrap interactive inputs in a `<form>` for progressive enhancement (no-JS fallback).
- Integrate any client scripting through htmx events (e.g. `htmx:beforeSwap`, `htmx:afterSwap`) rather than ad-hoc DOM manipulation.
