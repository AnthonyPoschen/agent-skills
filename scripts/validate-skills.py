#!/usr/bin/env python3
"""Validate the repository's Agent Skills layout."""

from __future__ import annotations

import re
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SKILLS_DIR = ROOT / "skills"
NAME_RE = re.compile(r"^[a-z0-9]+(?:-[a-z0-9]+)*$")
MD_LINK_RE = re.compile(r"\[[^\]]+\]\(([^):#]+\.md)(?:#[^)]+)?\)")


def fail(errors: list[str], message: str) -> None:
    errors.append(message)


def parse_frontmatter(path: Path, errors: list[str]) -> dict[str, str]:
    text = path.read_text(encoding="utf-8")
    lines = text.splitlines()

    if not lines or lines[0].strip() != "---":
        fail(errors, f"{path}: missing opening frontmatter delimiter")
        return {}

    try:
        end = lines[1:].index("---") + 1
    except ValueError:
        fail(errors, f"{path}: missing closing frontmatter delimiter")
        return {}

    fields: dict[str, str] = {}
    current_key = ""
    for line in lines[1:end]:
        if not line.strip():
            continue
        if line.startswith((" ", "\t")):
            if current_key:
                fields[current_key] = f"{fields[current_key]} {line.strip()}"
            continue
        if ":" not in line:
            fail(errors, f"{path}: invalid frontmatter line: {line}")
            continue
        key, value = line.split(":", 1)
        current_key = key.strip()
        fields[current_key] = value.strip().strip("\"'")

    return fields


def validate_markdown_links(skill_dir: Path, skill_md: Path, errors: list[str]) -> None:
    text = skill_md.read_text(encoding="utf-8")
    for match in MD_LINK_RE.finditer(text):
        link = match.group(1)
        if "://" in link or link.startswith("/"):
            continue
        target = (skill_dir / link).resolve()
        if not target.exists():
            fail(errors, f"{skill_md}: missing linked markdown file: {link}")


def validate_skill(skill_dir: Path, names: set[str], errors: list[str]) -> None:
    skill_md = skill_dir / "SKILL.md"
    folder_name = skill_dir.name

    if not NAME_RE.fullmatch(folder_name):
        fail(errors, f"{skill_dir}: folder name must be lowercase kebab-case")

    if not skill_md.exists():
        fail(errors, f"{skill_dir}: missing SKILL.md")
        return

    fields = parse_frontmatter(skill_md, errors)
    name = fields.get("name", "")
    description = fields.get("description", "")

    if not name:
        fail(errors, f"{skill_md}: missing frontmatter field: name")
    elif not NAME_RE.fullmatch(name):
        fail(errors, f"{skill_md}: name must be lowercase kebab-case")
    elif name != folder_name:
        fail(errors, f"{skill_md}: name must match folder name '{folder_name}'")
    elif name in names:
        fail(errors, f"{skill_md}: duplicate skill name '{name}'")
    else:
        names.add(name)

    if not description:
        fail(errors, f"{skill_md}: missing frontmatter field: description")

    validate_markdown_links(skill_dir, skill_md, errors)


def main() -> int:
    errors: list[str] = []

    if not SKILLS_DIR.exists():
        fail(errors, f"{SKILLS_DIR}: missing skills directory")
    else:
        skill_dirs = [path for path in sorted(SKILLS_DIR.iterdir()) if path.is_dir()]
        if not skill_dirs:
            fail(errors, f"{SKILLS_DIR}: no skills found")

        names: set[str] = set()
        for skill_dir in skill_dirs:
            validate_skill(skill_dir, names, errors)

    if errors:
        print("Skill validation failed:", file=sys.stderr)
        for error in errors:
            print(f"- {error}", file=sys.stderr)
        return 1

    print("Skill validation passed.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
