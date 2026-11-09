#!/usr/bin/env python3
"""Check that Trail of Bits plugins are pinned, not copied."""

from __future__ import annotations

import json
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
PINS = ROOT / "third-party" / "pins.json"
CLAUDE = ROOT / "third-party" / "claude" / "marketplace.json"
COPILOT = ROOT / "third-party" / "copilot" / "marketplace.json"
REQUIRED = ("name", "upstreamRepo", "path", "pinnedSha", "license", "attribution")
UPSTREAM = "https://github.com/trailofbits/skills"
LICENSE = "CC-BY-SA-4.0"


def fail(errors: list[str], message: str) -> None:
    errors.append(message)


def load(path: Path, errors: list[str]) -> dict:
    if not path.is_file():
        fail(errors, f"{path.relative_to(ROOT)}: missing")
        return {}
    try:
        data = json.loads(path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as exc:
        fail(errors, f"{path.relative_to(ROOT)}: invalid JSON: {exc}")
        return {}
    if not isinstance(data, dict):
        fail(errors, f"{path.relative_to(ROOT)}: expected an object")
        return {}
    return data


def main() -> int:
    errors: list[str] = []
    pins = load(PINS, errors)
    plugins = pins.get("plugins")
    if not isinstance(plugins, list) or len(plugins) != 2:
        fail(errors, "third-party/pins.json: plugins must list rust-review and differential-review")
        plugins = []

    by_name: dict[str, dict] = {}
    for item in plugins:
        if not isinstance(item, dict):
            fail(errors, "third-party/pins.json: each plugin must be an object")
            continue
        for key in REQUIRED:
            if not isinstance(item.get(key), str) or not item[key].strip():
                fail(errors, f"third-party/pins.json: {item.get('name')!r} missing {key}")
        name = item.get("name")
        if isinstance(name, str):
            by_name[name] = item
        if item.get("upstreamRepo") != UPSTREAM:
            fail(errors, f"third-party/pins.json: {name} upstreamRepo must be {UPSTREAM}")
        if item.get("license") != LICENSE:
            fail(errors, f"third-party/pins.json: {name} license must be {LICENSE}")
        sha = item.get("pinnedSha", "")
        if not isinstance(sha, str) or len(sha) != 40 or any(c not in "0123456789abcdef" for c in sha):
            fail(errors, f"third-party/pins.json: {name} pinnedSha must be 40 lowercase hex characters")
        path = item.get("path", "")
        if not isinstance(path, str) or path != f"plugins/{name}":
            fail(errors, f"third-party/pins.json: {name} path must be plugins/{name}")

    expected = {"rust-review", "differential-review"}
    if set(by_name) != expected:
        fail(errors, f"third-party/pins.json: plugin names must be {sorted(expected)}")

    shas = {item.get("pinnedSha") for item in by_name.values()}
    if len(shas) != 1:
        fail(errors, "third-party/pins.json: both plugins must share one pinnedSha")
    pinned = next(iter(shas), "")

    for skill_name in expected:
        copied = ROOT / "skills" / skill_name
        if copied.exists():
            fail(errors, f"skills/{skill_name}: upstream plugin must not be copied into this repo")

    check_marketplace(CLAUDE, "git-subdir", "https://github.com/trailofbits/skills.git", pinned, by_name, errors)
    check_marketplace(COPILOT, "github", "trailofbits/skills", pinned, by_name, errors)
    check_pointer_skill(pinned, errors)

    if errors:
        print("Third-party pin validation failed:", file=sys.stderr)
        for error in errors:
            print(f"- {error}", file=sys.stderr)
        return 1

    print("Third-party pin validation passed.")
    return 0


def check_pointer_skill(pinned: str, errors: list[str]) -> None:
    path = ROOT / "skills" / "trailofbits-review" / "SKILL.md"
    if not path.is_file():
        fail(errors, "skills/trailofbits-review/SKILL.md: missing")
        return
    text = path.read_text(encoding="utf-8")
    if pinned not in text:
        fail(errors, "skills/trailofbits-review/SKILL.md: must include pins.json pinnedSha")
    command = (
        f"npx skills add trailofbits/skills@{pinned} "
        "--skill rust-review --skill differential-review --full-depth"
    )
    if command not in text:
        fail(errors, "skills/trailofbits-review/SKILL.md: install command must use pins.json pinnedSha")
    for plugin in ("rust-review", "differential-review"):
        if plugin not in text:
            fail(errors, f"skills/trailofbits-review/SKILL.md: must name {plugin}")


def check_marketplace(
    path: Path,
    source_type: str,
    source_id: str,
    pinned: str,
    by_name: dict[str, dict],
    errors: list[str],
) -> None:
    data = load(path, errors)
    if not data:
        return
    label = str(path.relative_to(ROOT))
    if data.get("name") != "trailofbits-pins":
        fail(errors, f"{label}: name must be 'trailofbits-pins'")
    entries = data.get("plugins")
    if not isinstance(entries, list):
        fail(errors, f"{label}: plugins must be an array")
        return
    found: set[str] = set()
    id_key = "url" if source_type == "git-subdir" else "repo"
    for entry in entries:
        if not isinstance(entry, dict):
            fail(errors, f"{label}: plugin entry must be an object")
            continue
        name = entry.get("name")
        if not isinstance(name, str):
            fail(errors, f"{label}: plugin missing name")
            continue
        found.add(name)
        if entry.get("license") != LICENSE:
            fail(errors, f"{label}: {name} license must be {LICENSE}")
        source = entry.get("source")
        if not isinstance(source, dict):
            fail(errors, f"{label}: {name} source must be an object")
            continue
        if source.get("source") != source_type:
            fail(errors, f"{label}: {name} source.source must be {source_type}")
        if source.get(id_key) != source_id:
            fail(errors, f"{label}: {name} source.{id_key} must be {source_id}")
        if source.get("path") != f"plugins/{name}":
            fail(errors, f"{label}: {name} source.path must be plugins/{name}")
        if source.get("sha") != pinned or source.get("sha") != by_name.get(name, {}).get("pinnedSha"):
            fail(errors, f"{label}: {name} source.sha must match pins.json pinnedSha")
    if found != set(by_name):
        fail(errors, f"{label}: plugin names must match pins.json")


if __name__ == "__main__":
    raise SystemExit(main())
