#!/usr/bin/env python3
"""Check Cursor and Copilot plugin manifests and their skill paths."""

from __future__ import annotations

import json
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
PLUGIN_NAME = "apstack"

CURSOR_PLUGIN = ROOT / ".cursor-plugin" / "plugin.json"
CURSOR_MARKETPLACE = ROOT / ".cursor-plugin" / "marketplace.json"
COPILOT_PLUGIN = ROOT / "plugin.json"
COPILOT_MARKETPLACE = ROOT / ".github" / "plugin" / "marketplace.json"

AGENT_PLUGINS_SCHEMA = "https://agent-plugins.org/schemas/1.0.0/plugin.schema.json"


def fail(errors: list[str], message: str) -> None:
    errors.append(message)


def load_json(path: Path, errors: list[str]) -> dict:
    if not path.is_file():
        fail(errors, f"{path.relative_to(ROOT)}: missing manifest")
        return {}
    try:
        data = json.loads(path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as exc:
        fail(errors, f"{path.relative_to(ROOT)}: invalid JSON: {exc}")
        return {}
    if not isinstance(data, dict):
        fail(errors, f"{path.relative_to(ROOT)}: manifest must be a JSON object")
        return {}
    return data


def skill_dirs(skills_root: Path) -> set[str]:
    if not skills_root.is_dir():
        return set()
    return {
        child.name
        for child in skills_root.iterdir()
        if child.is_dir() and (child / "SKILL.md").is_file()
    }


def as_path_list(value: object, label: str, errors: list[str]) -> list[str]:
    if isinstance(value, str):
        return [value]
    if isinstance(value, list) and all(isinstance(item, str) for item in value):
        return list(value)
    fail(errors, f"{label}: skills must be a string or array of strings")
    return []


def resolve_skill_path(base: Path, raw: str, label: str, errors: list[str]) -> Path | None:
    if raw.startswith("/") or ".." in Path(raw).parts:
        fail(errors, f"{label}: skill path must be relative and must not contain '..': {raw}")
        return None
    return (base / raw).resolve()


def require_name(data: dict, path: Path, errors: list[str]) -> None:
    name = data.get("name")
    if name != PLUGIN_NAME:
        fail(errors, f"{path.relative_to(ROOT)}: name must be '{PLUGIN_NAME}', got {name!r}")


def check_listed_skills(
    paths: list[str],
    base: Path,
    label: str,
    repo_skills: set[str],
    errors: list[str],
) -> None:
    found: set[str] = set()
    for raw in paths:
        target = resolve_skill_path(base, raw, label, errors)
        if target is None:
            continue
        try:
            display = target.relative_to(ROOT)
        except ValueError:
            display = target
        if not target.exists():
            fail(errors, f"{label}: skill path does not exist: {raw}")
            continue
        if target.is_file():
            if target.name != "SKILL.md":
                fail(errors, f"{label}: skill file must be SKILL.md: {raw}")
                continue
            found.add(target.parent.name)
            continue
        if not target.is_dir():
            fail(errors, f"{label}: skill path is not a file or directory: {raw}")
            continue
        children = skill_dirs(target)
        if not children:
            fail(errors, f"{label}: skill directory has no SKILL.md children: {display}")
        found.update(children)

    missing = repo_skills - found
    extra = found - repo_skills
    if missing:
        fail(
            errors,
            f"{label}: skills/ not fully covered: {', '.join(sorted(missing))}",
        )
    if extra:
        fail(errors, f"{label}: listed skills outside skills/: {', '.join(sorted(extra))}")


def validate_cursor(errors: list[str], repo_skills: set[str]) -> None:
    plugin = load_json(CURSOR_PLUGIN, errors)
    marketplace = load_json(CURSOR_MARKETPLACE, errors)
    if not plugin or not marketplace:
        return

    require_name(plugin, CURSOR_PLUGIN, errors)
    if marketplace.get("name") != PLUGIN_NAME:
        fail(
            errors,
            f"{CURSOR_MARKETPLACE.relative_to(ROOT)}: name must be '{PLUGIN_NAME}'",
        )
    owner = marketplace.get("owner")
    if not isinstance(owner, dict) or not owner.get("name"):
        fail(errors, f"{CURSOR_MARKETPLACE.relative_to(ROOT)}: owner.name is required")

    plugins = marketplace.get("plugins")
    if not isinstance(plugins, list) or not plugins:
        fail(errors, f"{CURSOR_MARKETPLACE.relative_to(ROOT)}: plugins must be a non-empty array")
        return

    entry = next((item for item in plugins if isinstance(item, dict) and item.get("name") == PLUGIN_NAME), None)
    if entry is None:
        fail(errors, f"{CURSOR_MARKETPLACE.relative_to(ROOT)}: missing plugin entry '{PLUGIN_NAME}'")
        return

    source = entry.get("source")
    if not isinstance(source, str) or not source:
        fail(errors, f"{CURSOR_MARKETPLACE.relative_to(ROOT)}: plugin source must be a path string")
        return
    if source.startswith("/") or ".." in Path(source).parts:
        fail(errors, f"{CURSOR_MARKETPLACE.relative_to(ROOT)}: plugin source must stay inside the repo")
        return

    plugin_root = (ROOT / source).resolve()
    manifest = plugin_root / ".cursor-plugin" / "plugin.json"
    if manifest.resolve() != CURSOR_PLUGIN.resolve():
        fail(
            errors,
            f"{CURSOR_MARKETPLACE.relative_to(ROOT)}: source {source!r} does not resolve to {CURSOR_PLUGIN.relative_to(ROOT)}",
        )

    paths = as_path_list(plugin.get("skills"), str(CURSOR_PLUGIN.relative_to(ROOT)), errors)
    check_listed_skills(paths, plugin_root, str(CURSOR_PLUGIN.relative_to(ROOT)), repo_skills, errors)


def validate_copilot(errors: list[str], repo_skills: set[str]) -> None:
    plugin = load_json(COPILOT_PLUGIN, errors)
    marketplace = load_json(COPILOT_MARKETPLACE, errors)
    if not plugin or not marketplace:
        return

    require_name(plugin, COPILOT_PLUGIN, errors)
    if plugin.get("$schema") != AGENT_PLUGINS_SCHEMA:
        fail(
            errors,
            f"{COPILOT_PLUGIN.relative_to(ROOT)}: $schema must be {AGENT_PLUGINS_SCHEMA}",
        )
    if "skills" in plugin:
        fail(
            errors,
            f"{COPILOT_PLUGIN.relative_to(ROOT)}: Agent Plugins 1.0 rejects a skills path field; skills load from skills/",
        )

    if marketplace.get("name") != PLUGIN_NAME:
        fail(errors, f"{COPILOT_MARKETPLACE.relative_to(ROOT)}: name must be '{PLUGIN_NAME}'")
    owner = marketplace.get("owner")
    if not isinstance(owner, dict) or not owner.get("name"):
        fail(errors, f"{COPILOT_MARKETPLACE.relative_to(ROOT)}: owner.name is required")

    plugins = marketplace.get("plugins")
    if not isinstance(plugins, list) or not plugins:
        fail(errors, f"{COPILOT_MARKETPLACE.relative_to(ROOT)}: plugins must be a non-empty array")
        return

    entry = next((item for item in plugins if isinstance(item, dict) and item.get("name") == PLUGIN_NAME), None)
    if entry is None:
        fail(errors, f"{COPILOT_MARKETPLACE.relative_to(ROOT)}: missing plugin entry '{PLUGIN_NAME}'")
        return

    source = entry.get("source")
    if not isinstance(source, str) or not source:
        fail(errors, f"{COPILOT_MARKETPLACE.relative_to(ROOT)}: plugin source must be a path string")
        return
    if source.startswith("/") or ".." in Path(source).parts:
        fail(errors, f"{COPILOT_MARKETPLACE.relative_to(ROOT)}: plugin source must stay inside the repo")
        return

    plugin_root = (ROOT / source).resolve()
    manifest = plugin_root / "plugin.json"
    if not manifest.is_file() or manifest.resolve() != COPILOT_PLUGIN.resolve():
        fail(
            errors,
            f"{COPILOT_MARKETPLACE.relative_to(ROOT)}: source {source!r} does not resolve to {COPILOT_PLUGIN.relative_to(ROOT)}",
        )

    # Agent Plugins 1.0 always loads skills/ at the plugin root. The marketplace
    # skills field is the explicit path for legacy Copilot clients.
    implicit = plugin_root / "skills"
    if not implicit.is_dir():
        fail(errors, f"{COPILOT_PLUGIN.relative_to(ROOT)}: Agent Plugins skills directory is missing: skills/")
    else:
        check_listed_skills(
            ["skills"],
            plugin_root,
            f"{COPILOT_PLUGIN.relative_to(ROOT)} (implicit skills/)",
            repo_skills,
            errors,
        )

    paths = as_path_list(
        entry.get("skills"),
        str(COPILOT_MARKETPLACE.relative_to(ROOT)),
        errors,
    )
    check_listed_skills(
        paths,
        plugin_root,
        str(COPILOT_MARKETPLACE.relative_to(ROOT)),
        repo_skills,
        errors,
    )


def main() -> int:
    errors: list[str] = []
    repo_skills = skill_dirs(ROOT / "skills")
    if not repo_skills:
        fail(errors, "skills/: no skill directories with SKILL.md")

    validate_cursor(errors, repo_skills)
    validate_copilot(errors, repo_skills)

    if errors:
        print("Plugin manifest validation failed:", file=sys.stderr)
        for error in errors:
            print(f"- {error}", file=sys.stderr)
        return 1

    print(
        f"Plugin manifest validation passed ({len(repo_skills)} skills in Cursor and Copilot manifests)."
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
