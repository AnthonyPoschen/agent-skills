---
name: renderdoc-gpu-debug
description: >
  Use when capturing or inspecting a Vulkan frame with RenderDoc: .rdc
  captures, GLSL or SPIR-V shader debugging, pipeline state, or pixel history
  via rdc-cli. Do not use for CSS, HTML, canvas, or browser rendering.
---

# RenderDoc GPU Debugging

Inspect Vulkan frame captures with host `rdc-cli` (RenderDoc's Python API).
D3D11, D3D12, and OpenGL captures use the same commands. This skill does not
include the upstream MCP server.

Source and MIT licence: [LICENSE-UPSTREAM](LICENSE-UPSTREAM).

## Where capture and replay run

Capture and replay run only on Anthony's GPU machine **zanven-pc** (RTX 4070),
which has host RenderDoc and `rdc-cli`. Cloud VMs and the bot box have no GPU.

If the current host is not zanven-pc, do not run `rdc`, `renderdoccmd`, or any
capture script. Hand the capture and replay steps to a session on zanven-pc.
Opening an existing `.rdc` is replay: it also belongs on that machine.

## Prerequisites

On zanven-pc, `rdc doctor` must pass before capture or inspection.

- `rdc-cli` is installed (`pip install rdc-cli` if the command is missing).
- `RENDERDOC_PYTHON_PATH` points at the RenderDoc Python module built for the
  host Python. On Linux that module is `renderdoc.so`. On Windows it is
  `renderdoc.pyd`, and `renderdoc.dll` sits in the same directory. The Windows
  names do not apply on Linux.
- `renderdoccmd` is on `PATH` (Windows binary name: `renderdoccmd.exe`).

Vulkan needs the RenderDoc implicit layer and
`ENABLE_VULKAN_RENDERDOC_CAPTURE=1`. On Linux, point the loader at the layer
JSON with `VK_ADD_LAYER_PATH` or the user `implicit_layer.d` directory. The
Windows registry key `HKCU\SOFTWARE\Khronos\Vulkan\ImplicitLayers` applies only
on Windows.

## Session

Open one capture, inspect, then close. A leaked daemon holds GPU memory.

```bash
rdc open path/to/capture.rdc
rdc status
rdc close
```

Use `--session name` when two captures must stay open.

## Capture

Run the app from the directory it uses to find assets:

```bash
cd /path/to/app && rdc capture -o capture.rdc -- ./app [app-args...]
```

Useful flags: `--frame N`, `--timeout S` (default 60), `--api-validation`,
`--ref-all-resources`, `--wait-for-exit`.

RenderDoc frame boundaries are present calls. A headless Vulkan app with no
swapchain cannot be captured. On zanven-pc, prefer an interactive present, then
`--trigger` plus `rdc capture-trigger`, then `renderdoccmd capture`.

## Frame exploration

Start small. Full event lists are too large.

```bash
rdc info --json
rdc stats --json
rdc passes
rdc draws --limit 20
rdc draws --pass "Shadow Pass" --json
rdc events --limit 50
rdc event 42 --json
rdc ls /textures -l
```

## Pipeline, shaders, and pixels

```bash
rdc pipeline EID --json          # or one section: vs ps om rs ia ds
rdc bindings EID --set 0 --json
rdc shader EID ps --source
rdc shader EID ps --target spirv
rdc shader EID ps --constants --json
rdc pixel X Y EID --json
rdc pick-pixel X Y EID --json
rdc debug pixel EID X Y --trace
rdc debug vertex EID VTXID --json
```

Export images next to the capture, then view the PNG with the image read tool.
Do not dump image bytes through the shell.

```bash
rdc rt EID -o captures/render_target.png
rdc texture RESID --mip 2 -o captures/mip2.png
```

Shader edit and replay, still on zanven-pc:

```bash
rdc shader-encodings --json
rdc shader EID ps --source -o captures/shader.frag
rdc shader-build captures/shader.frag --encoding GLSL --stage ps --json
rdc shader-replace EID ps --with SHADER_ID --json
rdc rt EID -o captures/after_edit.png
rdc shader-restore EID ps
```

Compare two captures:

```bash
rdc diff before.rdc after.rdc --shortstat
rdc diff before.rdc after.rdc --framebuffer --diff-output captures/diff.png
```

## Output size

1. `rdc draws --limit 20` and `rdc events --limit 50` while exploring.
2. Filter with `--pass` before listing every draw.
3. `rdc draws -q` returns event ids only.
4. Prefer the default table over `--json` until the set is small.
5. Ask for one pipeline section (`rdc pipeline EID ps`), not the full state.

## When a command fails

| Symptom | What to do on zanven-pc |
| --- | --- |
| `rdc` not found | `pip install rdc-cli` |
| `rdc doctor` fails | Fix `RENDERDOC_PYTHON_PATH` and the host RenderDoc module |
| Capture reports no swapchain | `rdc capture --trigger`, then `rdc capture-trigger` |
| Daemon not responding | `rdc status`, `rdc close`, `rdc open` again |
| Counters or pixel history empty | `rdc gpus --json` and `rdc counters --list`; the driver may omit them |

## References

Command flags and defaults: [commands-quick-ref](references/commands-quick-ref.md).

Pass, colour, shadow, and pixel recipes: [debugging-recipes](references/debugging-recipes.md).
