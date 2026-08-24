# File and folder naming research

## Conclusion

Treat a path as a sentence. Each folder narrows the area. The filename names the
thing or job in that area. A reader should be able to guess one likely path
before opening source.

Do not repeat a parent folder name without a reason. `Enemy/Enemy.cpp` is fine
when it is the canonical implementation of `Enemy`, or when the language and
project use a matching type and implementation file. Otherwise the filename
should add the missing information: `Enemy/spawn.cpp`, `Enemy/render.cpp`, or
`Enemy/behaviour.cpp`.

This is a semantic rule, not a ban on matching names. A path segment earns its
place when it gives the reader information that the other segments do not.

Do not use file length as a split rule. A file can be large and still have a
clear job. Split when the path and filename can name a more useful responsibility
for a reader, not because a line counter reaches a fixed limit.

## What the sources support

- The [Google C++ Style Guide](https://google.github.io/styleguide/cppguide.html#General_Naming_Rules)
  says a name should describe its purpose to a new reader without repeating
  information that is already obvious from immediate context. Its
  [file naming rule](https://google.github.io/styleguide/cppguide.html#File_Names)
  also says filenames should be specific, follow the project convention, and
  explicitly allows a matching `foo_bar.h` and `foo_bar.cc` pair for `FooBar`.
  Matching names are therefore normal when they identify a type and its primary
  implementation.
- Go's official [module layout guidance](https://go.dev/doc/modules/layout)
  uses `internal/auth/auth.go` and `cmd/api-server/main.go`. The folder states
  the package or command; `main.go` has a clear entry point role. For server
  projects, it recommends keeping server logic in `internal` and commands in
  `cmd`.
- The [Go standard library HTTP package](https://github.com/golang/go/tree/master/src/net/http)
  uses filenames such as `request.go`, `response.go`, `server.go`, and
  `transport.go`. The directory supplies the broad `http` context and each file
  names a narrower responsibility.
- [Kubernetes core API v1](https://github.com/kubernetes/kubernetes/tree/master/staging/src/k8s.io/api/core/v1)
  uses `objectreference.go`, `resource.go`, `taint.go`, and `toleration.go`
  inside a broad versioned API folder. Its
  [`pkg/kubelet`](https://github.com/kubernetes/kubernetes/tree/master/pkg/kubelet)
  directory also has `kubelet.go` alongside `kubelet_pods.go`; the first is a
  useful example of a canonical package file, while the second adds its subject.
- [LLVM's IR implementation directory](https://github.com/llvm/llvm-project/tree/main/llvm/lib/IR)
  uses category folders and files such as `Instruction.cpp`, `Instructions.cpp`,
  and `Type.cpp`, rather than repeating `IR` in every filename.
- [Godot's 2D scene code](https://github.com/godotengine/godot/tree/master/scene/2d)
  uses `scene/2d/node_2d.cpp`: the folders add category and dimension, while the
  file names the component.

## Linux kernel and Kubernetes layouts

The Linux kernel and Kubernetes both use paths as navigation, but neither uses
"one object, one file" as a blanket rule.

### Linux kernel

- [`kernel/sched/`](https://github.com/torvalds/linux/tree/master/kernel/sched)
  collects CPU scheduler work. Its siblings include `core.c`, `fair.c`, `rt.c`,
  `deadline.c`, `idle.c`, and `topology.c`. The header of
  [`core.c`](https://github.com/torvalds/linux/blob/master/kernel/sched/core.c)
  calls it "Core kernel CPU scheduler code." The folder gives the subsystem;
  each file names a scheduler responsibility or policy.
- [`drivers/net/ethernet/intel/ice/`](https://github.com/torvalds/linux/tree/master/drivers/net/ethernet/intel/ice)
  narrows a driver from the general driver tree to vendor and device family.
  It has `ice_main.c`, `ice_txrx.c`, `ice_ethtool.c`, `ice_ptp.c`, and many
  matching headers. [`ice_main.c`](https://github.com/torvalds/linux/blob/master/drivers/net/ethernet/intel/ice/ice_main.c)
  carries the driver module setup, so its repeated `ice` name is a useful
  canonical driver entry file. `ice_txrx.c` and `ice_ptp.c` add their separate
  jobs rather than repeating the folder name alone.

### Kubernetes

- [`pkg/controller/deployment/`](https://github.com/kubernetes/kubernetes/tree/master/pkg/controller/deployment)
  groups one controller. It then separates `deployment_controller.go`,
  `rolling.go`, `recreate.go`, `progress.go`, and `sync.go` by controller job.
  For example, [`rolling.go`](https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/deployment/rolling.go)
  starts with the rolling rollout logic. The path predicts where to find a
  deployment concern without claiming every method needs its own file.
- [`pkg/apis/core/`](https://github.com/kubernetes/kubernetes/tree/master/pkg/apis/core)
  has dedicated files for concepts with their own focused work, including
  `objectreference.go`, `resource.go`, `taint.go`, and `toleration.go`. It also
  has [`types.go`](https://github.com/kubernetes/kubernetes/blob/master/pkg/apis/core/types.go),
  which intentionally holds closely related core API schema types such as
  `Volume` and `VolumeSource`. This is a useful counterexample to splitting
  every named type into its own file.

## Implications for the proposed rule

The rule should improve discovery, not maximize folders or files.

- Give a distinct boundary its own file when a reader will look for that
  boundary by name, or it has separate behavior, ownership, or reasons to
  change. Kubernetes' `taint.go` and the deployment controller's `rolling.go`
  are good examples.
- Keep related concepts together when they form one contract or schema and a
  reader benefits from seeing them together. Kubernetes' core `types.go` is a
  good example.
- Let a repeated name stand only when it tells the reader that the file is the
  folder's primary implementation or is part of an established language
  convention. Linux's `ice_main.c`, Kubernetes' `kubelet.go`, and conventional
  C++ header and implementation pairs qualify. A second `Enemy/Enemy.cpp` that
  only contains spawning or rendering does not.
- Add the job to the filename when the folder already names the object:
  `Enemy/spawn.cpp`, `Enemy/render.cpp`, or `Enemy/behaviour.cpp`. Avoid a
  catch all `Enemy/Enemy.cpp` that accumulates every concern for that object.

## Pstack's deslop guidance

Pstack points to `deslop` as a Cursor Team Kit skill rather than shipping it
itself. The [deslop skill](https://github.com/cursor/plugins/blob/main/cursor-team-kit/skills/deslop/SKILL.md)
asks for minimal edits and consistency with the surrounding codebase. It covers
unneeded comments, unusual defensive checks, unnecessary casts, and deep
nesting. It has no file or folder naming rule.

Pstack's [unslop skill](https://github.com/cursor/plugins/blob/main/pstack/skills/unslop/SKILL.md)
asks for plain words and short sentences. It rejects em dashes and hyphens used
as dashes, not normal compound words. For this library's prose, a useful rule is:

> Write plain English. Prefer short sentences and ordinary words. Do not use an
> em dash or a made up hyphenated phrase when a normal sentence says it better.

This does not affect required kebab case skill folder names.

## Rule to carry into the coding skill

> Choose paths for discovery. Each folder must narrow the area, and each file
> must name the thing or job inside it. Avoid repeating a parent name unless the
> file is that folder's canonical entry point or a matching type or module file
> is an established local convention.

Treat a file that grows large as a prompt to inspect its responsibilities, not a
violation. Keep it intact when the path and filename still give readers the best
starting point. Split it when a child responsibility has a clear name and a
reader would look for it separately.
