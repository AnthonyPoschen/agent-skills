# Per-agent TMPDIR (proposal)

All agents share one box. A fixed `/tmp/prbody.md` written by one agent was posted as another agent's PR body (k8s#87).

- `scripts/agent-tmp.sh` (source it): makes `/tmp/agents/<name>` (mode 700), exports `TMPDIR`/`AGENT_TMPDIR`, offers `agent_mktemp`, `agent_mktemp_d`, `pr_body_file`. Name is `$AGENT_NAME`, else `$GROK_AGENT_ID`, else `$USER-$$`; unsafe chars become `_`.
- `scripts/check-fixed-tmp.sh`: fails on literal `/tmp/<name>` in scripts/skills/tests. Allowed: `/tmp/agents`, `mktemp`, or a line marked `# tmp-ok`. Run with `npm run check-tmp`; suggest adding to the pre-commit hook.
- `tests/agent-tmp-test.sh`: two agents get distinct 700 dirs; 200 parallel `agent_mktemp` never collide; names are sanitised.

Follow-ups outside this repo (not changed live): `/home/box/shared/log-run.sh` has no fixed /tmp today; fleet-rules.md could point to `pr_body_file`; agent runtimes could export `AGENT_NAME` so the helper is automatic.
