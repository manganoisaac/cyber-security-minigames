#!/usr/bin/env python3
# /// script
# requires-python = ">=3.11"
# dependencies = [
#     "langgraph>=1.2",
#     "tree-sitter>=0.25",
# ]
# ///
"""Fill a templ component body using tree-sitter + a LangGraph pipeline.

The model is given ONLY the component and the description as context,
then produces the component body markup. tree-sitter (the templ grammar)
is used to (a) extract the component signature / imports / body region so
only the body is replaced, and (b) validate that the regenerated template
still parses before anything is written.

The LLM is the model powering this opencode session, invoked through the
`opencode run` CLI so no separate API key is needed.

Usage:
  uv run scripts/fill_component.py views/components/<feature>/<name>.templ "<description>" \
    [--force] [--model opencode/deepseek-v4-flash-free] [--max-tries 3] [--dry-run]

Run from the repo root. Requires: git + a C compiler (gcc/clang) on first run,
so the templ grammar can be built into scripts/build/templ.so (cached after).
"""

from __future__ import annotations

import argparse
import ctypes
import json
import os
import re
import shutil
import subprocess
import sys
import warnings
from pathlib import Path
from typing import Any, TypedDict

from tree_sitter import Language, Parser

from typing_extensions import NotRequired

REPO_ROOT = Path(__file__).resolve().parent.parent
BUILD_DIR = REPO_ROOT / "scripts" / "build"
TEMPL_SO = BUILD_DIR / "templ.so"
CSS_DIR = REPO_ROOT / "public" / "css" / "features"
INDEX_CSS = REPO_ROOT / "public" / "css" / "index.css"

TEMPL_GRAMMAR_URL = "https://github.com/vrischmann/tree-sitter-templ"

DEFAULT_MODEL = "opencode/deepseek-v4-flash-free"

# ---------------------------------------------------------------------------
# tree-sitter (templ grammar) loading
# ---------------------------------------------------------------------------


def _lib_from_so(so_path: str | Path) -> Language:
    """Load the templ grammar .so with py-tree-sitter."""
    lib = ctypes.CDLL(str(so_path))
    lib.tree_sitter_templ.restype = ctypes.c_void_p
    with warnings.catch_warnings():
        warnings.simplefilter("ignore")
        return Language(lib.tree_sitter_templ())


def _build_templ_so() -> Path:
    """Compile the templ tree-sitter grammar (parser.c + scanner.c) into a .so."""
    git = shutil.which("git")
    cc = shutil.which("cc") or shutil.which("gcc") or shutil.which("clang")
    if git is None:
        sys.exit("error: 'git' is required to fetch the templ grammar on first run")
    if cc is None:
        sys.exit(
            "error: no C compiler found (cc/gcc/clang); the templ tree-sitter grammar "
            "needs to be compiled once"
        )

    BUILD_DIR.mkdir(parents=True, exist_ok=True)
    src_dir = BUILD_DIR / "_grammar-src"
    if src_dir.exists():
        shutil.rmtree(src_dir)

    print(f"[grammar] cloning {TEMPL_GRAMMAR_URL} ...")
    subprocess.run(
        [git, "clone", "--depth", "1", "--quiet", TEMPL_GRAMMAR_URL, str(src_dir)],
        check=True,
    )
    print(f"[grammar] compiling with {cc} -> {TEMPL_SO} ...")
    subprocess.run(
        [cc, "-shared", "-fPIC", "-Isrc", "-o", str(TEMPL_SO), "src/parser.c", "src/scanner.c"],
        check=True,
        cwd=src_dir,
    )
    shutil.rmtree(src_dir)
    return TEMPL_SO


def load_templ_language() -> Language:
    """Load the templ grammar, building it from source if not already available."""
    so_path: Path | None = None
    env_so = Path(os.environ.get("TREESITTER_TEMPL_SO", ""))
    if env_so.is_file():
        so_path = env_so
    elif TEMPL_SO.is_file():
        so_path = TEMPL_SO
    if so_path is None:
        so_path = _build_templ_so()
    return _lib_from_so(so_path)


def parse_templ(lang: Language, source: bytes):
    return Parser(lang).parse(source)


# ---------------------------------------------------------------------------
# tree-sitter helpers
# ---------------------------------------------------------------------------


def _named_children(node) -> list[Any]:
    return list(node.named_children)


def _component_decl(source: bytes, lang: Language):
    """Find the single templ component declaration in source."""
    tree = parse_templ(lang, source)
    decl = None
    for node in tree.root_node.named_children:
        if node.type == "component_declaration":
            if decl is not None:
                sys.exit("error: file contains more than one templ component")
            decl = node
    if decl is None:
        sys.exit("error: no templ component found in file")
    block = next(
        (c for c in _named_children(decl) if c.type == "component_block"), None
    )
    if block is None:
        sys.exit("error: could not locate the component body block")
    name_node = next(
        (c for c in _named_children(decl) if c.type == "component_identifier"), None
    )
    params_node = next(
        (c for c in _named_children(decl) if c.type == "parameter_list"), None
    )
    return {
        "tree": tree,
        "decl": decl,
        "block": block,
        "name": source[name_node.start_byte : name_node.end_byte].decode() if name_node else "",
        "params": (
            source[params_node.start_byte : params_node.end_byte].decode()
            if params_node
            else "()"
        ),
        "signature": source[decl.start_byte : block.start_byte].strip(),
    }


def _body_is_empty(block) -> bool:
    """True if the block contains only comments/whitespace between its braces."""

    def walk(node):
        for c in node.named_children:
            if c.type != "comment":
                return False
            if not walk(c):
                return False
        return True

    return walk(block)


def _collect_errors(node, source: bytes, acc: list[str]) -> None:
    if node.is_missing:
        snippet = source[min(node.start_byte, len(source) - 1) : node.end_byte + 8]
        acc.append(f"missing {node.type} near {snippet!r}")
    if node.type == "ERROR":
        snippet = source[node.start_byte : node.end_byte][:96]
        acc.append(f"syntax error near: {snippet.decode('utf-8', 'replace')!r}")
    for c in node.children:
        _collect_errors(c, source, acc)


# ---------------------------------------------------------------------------
# prompt / LLM
# ---------------------------------------------------------------------------


def build_prompt(
    signature: str,
    package: str,
    imports: str,
    description: str,
    feedback: str | None,
) -> str:
    conv = """\
You are an expert templ + htmx developer. Fill in the BODY of a templ component.
You have only the component and its description as context; do not use any tools
and do not try to read the repository. Respond with text only.

CONVENTIONS (must follow):
- Output ONLY the body markup: HTML-like elements and templ control flow
  (if / else / for / switch). Do NOT repeat the `templ` signature, do NOT add
  the enclosing braces, and do NOT emit a trailing closing brace.
- Classic BEM naming: `Block__element--modifier`. The Block is the kebab-case
  component name (e.g. email-phrase, modal). Reuse the existing class names.
- The component receives all data as parameters; never reference files,
  globals, or other components.
- Interactivity goes through htmx attributes in the markup (hx-get / hx-post /
  hx-target / hx-swap / hx-trigger), not JavaScript. Fragments are served by the
  server; state lives in the HTML.
- templ syntax: an attribute embedding Go uses `attr={ expr }`; text
  interpolation uses `{ expr }`; class attribute lists use `class={ "a", "b" }`
  and conditional classes use `templ.KV("class", cond)`.
- If an existing parameter provides the content, render it — do not invent
  unrelated content.

RESPONSE FORMAT (strict):
1. A fenced code block tagged templ (```templ ... ```) containing exactly the
   body markup.
2. Optional line `IMPORT: pkg/path` (comma separated) for any Go import the
   body needs (e.g. strconv) that is not already in the component header.
3. Line `BEM: a b c` listing every BEM class name used (space separated)."""
    comp_lines = [f"Component: {signature} {{", "    // body to fill", "}"]
    if package:
        comp_lines.insert(0, f"Package: {package}")
    if imports:
        comp_lines.insert(1, f"Imports: {imports}")
    parts = [
        conv,
        "Component:",
        "\n".join(comp_lines),
        f"Description: {description}",
    ]
    if feedback:
        parts.append(
            f"The previous attempt failed templ validation with:\n{feedback}\nFix those errors and retry."
        )
    return "\n\n".join(parts)


def call_opencode(prompt: str, model: str, timeout: int = 300) -> str:
    """Run the model via `opencode run --format json` and return its text."""
    cmd = ["opencode", "run", "--format=json", "-m", model, prompt]
    proc = subprocess.run(
        cmd,
        cwd=REPO_ROOT,
        capture_output=True,
        text=True,
        timeout=timeout,
    )
    if proc.returncode != 0:
        raise RuntimeError(
            f"opencode run failed ({proc.returncode}): {(proc.stderr or proc.stdout)[-800:]}"
        )
    chunks: list[str] = []
    for line in proc.stdout.splitlines():
        line = line.strip()
        if not line:
            continue
        try:
            event = json.loads(line)
        except json.JSONDecodeError:
            continue
        if event.get("type") == "done":
            break
        part = event.get("part")
        if isinstance(part, dict) and isinstance(part.get("text"), str) and part["text"]:
            chunks.append(part["text"])
    return "".join(chunks).strip()


# ---------------------------------------------------------------------------
# generation output parsing
# ---------------------------------------------------------------------------

FENCE_RE = re.compile(r"```(?:templ)?\s*(.*?)```", re.S)
CLASS_ATTR_RE = re.compile(r'class="([^"]+)"')
CLASS_LIST_RE = re.compile(r"class=\{(.*?)\}", re.S)
KV_RE = re.compile(r'templ\.KV\(\s*"([^"]+)"')
BEM_TOKEN_RE = re.compile(
    r"^[a-z0-9]+(?:-[a-z0-9]+)*"
    r"(?:__[a-z0-9]+(?:-[a-z0-9]+)*)?"
    r"(?:--[a-z0-9]+(?:-[a-z0-9]+)*)?$"
)


def _extract_classes(body: str) -> list[str]:
    found: list[str] = []
    for m in CLASS_ATTR_RE.finditer(body):
        found.extend(m.group(1).split())
    for m in CLASS_LIST_RE.finditer(body):
        found.extend(q for q in re.findall(r'"([^"}]+)"', m.group(1)))
    found.extend(m.group(1) for m in KV_RE.finditer(body))
    return sorted({c for c in found if BEM_TOKEN_RE.match(c)})


def parse_generation(text: str) -> dict[str, Any]:
    """Split the model reply into body markup, needed imports and BEM classes."""
    m = FENCE_RE.search(text)
    if not m:
        return {"generated": None, "gen_error": "model reply missing a ```templ ...``` code fence"}
    body = m.group(1).strip("\n")
    imports: list[str] = []
    bem: list[str] = []
    for line in text.splitlines():
        ls = line.strip()
        low = ls.lower()
        if low.startswith("import:"):
            imports = [
                p.strip().strip('"')
                for p in re.split(r"[,\s]+", ls[len("import:") :])
                if p.strip()
            ]
        elif low.startswith("bem:"):
            bem = [c for c in ls[len("bem:") :].split() if c.strip()]
    classes = bem if bem else _extract_classes(body)
    return {"generated": body, "imports": imports, "classes": classes}


# ---------------------------------------------------------------------------
# file splicing / imports / css
# ---------------------------------------------------------------------------


def _kebab(name: str) -> str:
    return re.sub(r"(?<!^)(?=[A-Z])", "-", name).lower()


def _upsert_imports(source: bytes, need: list[str]) -> bytes:
    """Add Go import paths that the body needs but the file lacks."""
    need = [p.strip().strip('"') for p in need if p.strip().strip('"')]
    if not need:
        return source
    text = source.decode("utf-8")

    present: set[str] = set(re.findall(r'^import\s+"([^"]+)"\s*$', text, re.M))
    grouped = re.search(r"import\s*\((.*?)\)", text, re.S)
    if grouped:
        present |= set(re.findall(r'"([^"]+)"', grouped.group(1)))
    add = [p for p in need if p not in present]
    if not add:
        return source

    if grouped:
        closing = text.find(")", grouped.start(1))
        text = text[:closing] + "".join(f'\t"{p}"\n' for p in add) + text[closing:]
    else:
        bare = list(re.finditer(r"^\s*import\s+\"[^\"]+\"\s*$", text, re.M))
        if bare:
            at = text.index("\n", bare[-1].start()) + 1
            text = text[:at] + "".join(f'import "{p}"\n' for p in add) + text[at:]
        else:
            pkg = re.search(r"package\s+\w+", text)
            if pkg is None:
                return source
            text = text[: pkg.end()] + "\n\n" + "".join(f'import "{p}"\n' for p in add) + text[pkg.end() :]
    return text.encode("utf-8")


def css_for(block: str, classes: list[str]) -> str:
    rules = [f".{block} {{", "", "}"]
    for c in sorted(set(classes)):
        if c == block or not c.startswith(block):
            continue
        rules += ["", f".{c} {{", "", "}"]
    return "\n".join(rules) + "\n"


# ---------------------------------------------------------------------------
# LangGraph state + nodes
# ---------------------------------------------------------------------------


class FillState(TypedDict, total=False):
    path: str
    description: str
    force: bool
    model: str
    max_tries: int
    dry_run: bool
    _lang: NotRequired[Any]

    source_bytes: bytes
    name: str
    package: str
    imports: str
    signature: str
    decl_start: int
    body_start: int
    body_end: int

    generated: NotRequired[str | None]
    gen_error: NotRequired[str | None]
    imports_needed: NotRequired[list[str]]
    classes: NotRequired[list[str]]
    feedback: NotRequired[str | None]
    attempts: NotRequired[int]

    error: NotRequired[str | None]
    out_template: NotRequired[str]
    out_css: NotRequired[str]
    css_written: NotRequired[bool]
    css_path: NotRequired[str]
    import_wired: NotRequired[bool]


def node_extract(state: FillState) -> dict[str, Any]:
    lang = state["_lang"]
    source = Path(state["path"]).read_bytes()
    comp = _component_decl(source, lang)
    block = comp["block"]

    if not _body_is_empty(block) and not state.get("force"):
        return {
            "error": (
                "component body is not empty (only whitespace/comments allowed). "
                "Re-run with --force to overwrite the existing body."
            )
        }

    pkg_node = next(
        (
            c
            for c in comp["tree"].root_node.named_children
            if c.type == "package_clause"
        ),
        None,
    )
    package = ""
    imports = ""
    if pkg_node is not None:
        package = source[pkg_node.start_byte : pkg_node.end_byte].decode()
    import_nodes = [
        c
        for c in comp["tree"].root_node.named_children
        if c.type == "import_declaration"
    ]
    imports = "\n".join(
        source[n.start_byte : n.end_byte].decode() for n in import_nodes
    )

    return {
        "source_bytes": source,
        "name": comp["name"],
        "package": package,
        "imports": imports,
        "signature": comp["signature"],
        "decl_start": comp["decl"].start_byte,
        "body_start": block.start_byte,
        "body_end": block.end_byte,
        "generated": None,
        "feedback": None,
        "attempts": 0,
        "error": None,
    }


def node_generate(state: FillState) -> dict[str, Any]:
    feedback = state.get("feedback")
    prompt = build_prompt(
        state["signature"],
        state.get("package", ""),
        state.get("imports", ""),
        state["description"],
        feedback,
    )
    print(f"[generate] calling {state['model']} (attempt {state.get('attempts', 0) + 1}) ...")
    text = call_opencode(prompt, state["model"])
    parsed = parse_generation(text)
    return {
        "generated": parsed["generated"],
        "gen_error": parsed.get("gen_error") if parsed.get("generated") is None else None,
        "imports_needed": parsed.get("imports", []),
        "classes": parsed.get("classes", []),
        "attempts": state.get("attempts", 0) + 1,
    }


def node_validate(state: FillState) -> dict[str, Any]:
    lang = state["_lang"]
    src = state["source_bytes"]
    generated = state.get("generated")

    if generated is None:
        feedback = state.get("gen_error") or "model returned no body"
    elif not generated.strip():
        feedback = "model returned an empty body"
    else:
        pre = src[: state["body_start"] + 1]
        post = src[state["body_end"] - 1 :]
        rebuilt = pre + ("\n" + generated + "\n").encode("utf-8") + post
        tree = parse_templ(lang, rebuilt)
        errors: list[str] = []
        _collect_errors(tree.root_node, rebuilt, errors)
        feedback = "\n".join(errors[:6]) if errors else None

    if feedback is not None and state.get("attempts", 0) >= state.get("max_tries", 3):
        return {
            "feedback": feedback,
            "error": (
                f"could not produce a valid templ body after "
                f"{state['attempts']} attempts: {feedback}"
            ),
        }
    return {"feedback": feedback}


def _route_after_extract(state: FillState) -> str:
    return "end" if state.get("error") else "generate"


def _route_after_validate(state: FillState) -> str:
    if state.get("error"):
        return "end"
    if state.get("feedback") is not None:
        return "generate"
    return "write"


def node_write(state: FillState) -> dict[str, Any]:
    src = state["source_bytes"]
    generated = state["generated"]
    indented = "\n".join(
        ("\t" + line) if line.strip() else "" for line in generated.splitlines()
    )
    pre = src[: state["body_start"] + 1]
    post = src[state["body_end"] - 1 :]
    body = ("\n" + indented + "\n").encode("utf-8")
    new_src = pre + body + post
    new_src = _upsert_imports(new_src, state.get("imports_needed", []))

    path = Path(state["path"])
    if not state.get("dry_run"):
        print(f"[write] splicing body into {path}")
        path.write_bytes(new_src)
    out_template = new_src.decode("utf-8")
    if state.get("dry_run"):
        print("[dry-run] generated body:\n" + generated.rstrip("\n") + "\n")

    # --- CSS -----------------------------------------------------------------
    block = _kebab(state["name"])
    classes = state.get("classes", [])
    css = css_for(block, classes)
    css_path = CSS_DIR / f"{block}.css"
    css_written = False
    import_wired = False
    if classes and not state.get("dry_run"):
        if css_path.exists():
            print(f"[css] {css_path} already exists, leaving it untouched")
        else:
            CSS_DIR.mkdir(parents=True, exist_ok=True)
            css_path.write_text(css)
            css_written = True
            print(f"[css] wrote {css_path}")
            if INDEX_CSS.exists():
                index = INDEX_CSS.read_text()
                if f'@import "./features/{block}.css";' not in index:
                    INDEX_CSS.write_text(
                        index.rstrip("\n") + f'\n@import "./features/{block}.css";\n'
                    )
                    import_wired = True
                    print(f"[css] registered import in {INDEX_CSS}")
    elif classes and state.get("dry_run"):
        print(f"[dry-run] would write {css_path}")
        print(css)

    return {
        "out_template": out_template,
        "out_css": css,
        "css_written": css_written,
        "css_path": str(css_path),
        "import_wired": import_wired,
    }


def build_graph():
    """Build the FillState pipeline once."""
    from langgraph.graph import END, START, StateGraph

    graph = StateGraph(FillState)
    graph.add_node("extract", node_extract)
    graph.add_node("generate", node_generate)
    graph.add_node("validate", node_validate)
    graph.add_node("write", node_write)

    graph.add_edge(START, "extract")
    graph.add_conditional_edges(
        "extract",
        _route_after_extract,
        {"generate": "generate", "end": END},
    )
    graph.add_edge("generate", "validate")
    graph.add_conditional_edges(
        "validate",
        _route_after_validate,
        {"write": "write", "generate": "generate", "end": END},
    )
    graph.add_edge("write", END)
    return graph.compile()


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("component", help="path to the .templ component file")
    parser.add_argument("description", help="what the component should render")
    parser.add_argument("--force", action="store_true", help="overwrite a non-empty body")
    parser.add_argument("--model", default=DEFAULT_MODEL, help="opencode model id")
    parser.add_argument("--max-tries", type=int, default=3, help="max generate/validate attempts")
    parser.add_argument("--dry-run", action="store_true", help="print changes without writing files")
    args = parser.parse_args()

    component = REPO_ROOT / args.component
    if not component.is_file():
        parser.error(f"component file not found: {component}")

    global _LANG
    _LANG = load_templ_language()

    graph = build_graph()
    state: FillState = {
        "path": str(component),
        "description": args.description,
        "force": args.force,
        "model": args.model,
        "max_tries": args.max_tries,
        "dry_run": args.dry_run,
        "_lang": _LANG,
    }
    result = graph.invoke(state)
    if result.get("error"):
        print(f"\nerror: {result['error']}", file=sys.stderr)
        return 1

    print(f"\nDone. Component body filled in {component}")
    if result.get("css_written"):
        print(f"CSS scaffold written: {result['css_path']}")
    if result.get("css_path") and not result.get("css_written") and result.get("classes"):
        print(f"CSS already existed, not touched: {result['css_path']}")
    if result.get("import_wired"):
        print("CSS @import registered in public/css/index.css")
    print("Run `make generate` to compile the template, then check the component in the browser.")
    return 0


_LANG: Language | None = None

if __name__ == "__main__":
    sys.exit(main())