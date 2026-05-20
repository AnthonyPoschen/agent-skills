---
name: makefile
description: >
  Use when creating, reviewing, or editing Makefiles, make targets, make
  recipes, task runners, or documentation that tells users to run make commands.
  Always use this skill when a make target accepts user-provided variables such
  as tokens, usernames, branch names, paths, namespaces, environments, or other
  inputs. Ensures required Make variables are validated with clear error
  descriptions before commands run.
---

# Makefile

## Goal

Write Makefiles that fail before doing work when required user input is
missing. A target that needs a user-provided variable must explain what the
variable is for in the Make error, close to the target that needs it.

## Required Variable Guard

When adding a target that reads user input from Make variables, include this
helper once near the top of the Makefile if an equivalent helper is not already
present:

```make
# https://stackoverflow.com/a/10858332
check_defined = \
    $(strip $(foreach 1,$1, \
        $(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
    $(if $(value $1),, \
        $(error Undefined $1$(if $2, ($2))$(if $(value @), \
		required by target `$@`)))
```

Use the helper at the start of each target that requires input:

```make
.PHONY: init_cluster
init_cluster:
	@:$(call check_defined, GIT_USERNAME, username used for Git operations)
	@:$(call check_defined, GIT_TOKEN, access token with required repository permissions)
	@:$(call check_defined, GIT_BRANCH, Git branch to reconcile)
	@:$(call check_defined, KUSTOMIZATION_PATH, path from repo root to the kustomization)
	flux create source git k8s --url=https://github.com/AnthonyPoschen/k8s.git --username="$(GIT_USERNAME)" --password="$(GIT_TOKEN)" --branch="$(GIT_BRANCH)" -n flux-system
	flux create kustomization k8s --source=k8s -n flux-system --path="$(KUSTOMIZATION_PATH)" --interval=1m
```

## Rules

- Validate every required user-provided Make variable before the first command
  that uses it.
- Give each `check_defined` call a short, concrete description of what the
  variable is for, not just "required".
- Keep validation inside the target that needs the variable unless several
  targets share the exact same required inputs.
- Prefer explicit variable names that do not shadow common environment
  variables. For example, use `KUSTOMIZATION_PATH`, `CONFIG_PATH`, or
  `MANIFEST_PATH` instead of `PATH`.
- Do not check optional variables. Document their defaults near the assignment
  or in the target help text.
- Use `$(VAR)` in Make recipes. Avoid `${VAR}` unless the command specifically
  needs shell-style expansion after Make runs.
- Quote variable expansions in shell commands when values may contain special
  characters or spaces: `"$(VAR)"`.
- Keep `.PHONY` declarations for command targets.
- Do not print secrets. It is fine to validate that `GIT_TOKEN` or another
  secret variable is present, but do not echo the value.

## Shared Guards

If multiple targets require the same group of variables, extract a small guard
macro with the same descriptive checks:

```make
check_flux_git_inputs = \
	$(call check_defined, GIT_USERNAME, username used for Git operations) \
	$(call check_defined, GIT_TOKEN, access token with required repository permissions) \
	$(call check_defined, GIT_BRANCH, Git branch to reconcile)

.PHONY: init_cluster
init_cluster:
	@:$(call check_flux_git_inputs)
	@:$(call check_defined, KUSTOMIZATION_PATH, path from repo root to the kustomization)
	flux create source git k8s --url="$(GIT_REPO_URL)" --username="$(GIT_USERNAME)" --password="$(GIT_TOKEN)" --branch="$(GIT_BRANCH)" -n flux-system
```

Keep shared guard macros narrow. Do not force unrelated targets to provide
variables they do not use.

## Review Checklist

- Every required input variable has a `check_defined` call before use.
- Each error description tells the user what value to provide and why.
- Secret variables are validated but never echoed.
- Target names are declared `.PHONY`.
- Variable names avoid collisions with important shell environment variables
  such as `PATH`, `HOME`, `SHELL`, and `USER`.
- Commands are not run through unnecessary semicolon chains; each recipe line is
  readable and fails clearly.
