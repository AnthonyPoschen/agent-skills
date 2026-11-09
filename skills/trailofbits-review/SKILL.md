---
name: trailofbits-review
description: Manual pointer to pinned Trail of Bits rust-review and differential-review. Use for genos-engine (Rust, unsafe GPU/FFI) or a security-sensitive PR diff.
disable-model-invocation: true
metadata:
  opencode/slash: "true"
  opencode/autoinvoke: "false"
---

# Trail of Bits review

This skill is a router. It does not contain the Trail of Bits audit. Those
plugins are CC-BY-SA-4.0 and stay in `trailofbits/skills`. Do not paste their
prompts, skill text, or agent instructions into this repository.

## When to use it

- `rust-review`: Rust security review, including genos-engine code that uses
  `unsafe`, GPU code, or FFI.
- `differential-review`: a security-sensitive pull request diff.

Invoke this skill by name. Do not load it for ordinary edits. Each audit is
large and spends a lot of tokens.

## Steps

The pin below is copied from `third-party/pins.json` in the agent-skills
repository. Use this copy. A consumer checkout, plugin install, or
`npx skills add` does not include that file.

Pinned commit: `e4dd0bb5b15535f7c5a0834187dffba6ac37e4ff`

Install both plugins at that commit. This writes them on the machine that
runs the command. It does not add them to the MIT agent-skills tree.

```bash
npx skills add trailofbits/skills@e4dd0bb5b15535f7c5a0834187dffba6ac37e4ff --skill rust-review --skill differential-review --full-depth
```

Then run the installed upstream plugin. Do not follow `main`. Attribute the
result to Trail of Bits under CC-BY-SA-4.0.
