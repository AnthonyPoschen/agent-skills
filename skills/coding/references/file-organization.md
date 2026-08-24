# File Organization

Use this reference when starting an application or subsystem, placing source
files, choosing packages or modules, moving code, splitting a file, or reviewing
layout.

## Start With The Language Convention

Before choosing an initial layout, identify the language and framework. Learn
their current convention from the existing repository, official documentation,
standard tooling, and mature first party projects. Use that convention unless
the repository has a stronger coherent pattern.

Do not use one language's directory names in another language by habit. For
example, `cmd` and `internal` are Go choices, not universal directories.

Treat a local layout as authoritative only when it gives clear placement rules:

- The entrypoint is easy to identify.
- Directories communicate ownership or role.
- A reader can predict where a feature, API shape, adapter, or test belongs.
- The layout works with the language and framework tooling.

Repetition alone does not make a placement pattern good. If the existing layout
is incoherent, do not extend its dumping ground when new code has a clean home
under the language convention.

## Make Paths Easy To Search

Treat a path as a sentence. Each folder narrows the area. The filename names the
thing or job inside that area. A reader should be able to guess one likely path
before opening source.

Choose a folder name that communicates ownership or role. Choose a filename that
answers the next useful question. Avoid vague homes such as `app`, `common`,
`helpers`, `utils`, and `misc` when a domain or responsibility name is available.

Use this check before creating a file:

> If a reader asks where this thing lives, would this path be an obvious answer?

## Keep Related Code Together

Keep closely related code together by default. Add a file only when its name
gives a reader a clearer, separate place to look.

A separate file earns its cost when:

- A reader would naturally search for that responsibility by name.
- The path and filename narrow the search more than the existing file does.
- It has separate behavior, ownership, or reasons to change.
- A reader does not need to open it and the existing file together almost every
  time.

Keep code in one file when types form one schema or contract, functions are
tightly coupled parts of one operation, or a new file would contain only a small
helper, one method, or one trivial type. File count should grow with independent
questions, not with the number of types, methods, or implementation details.

## Use Canonical Names Deliberately

A repeated folder and filename can be useful. `Enemy/Enemy.cpp` is valid when it
is the primary implementation of `Enemy`. `cmd/api/main.go` is valid because
`main.go` names a standard entrypoint. A matching type and implementation pair
can also be a local convention.

When the folder already names the object and the file owns one part of its work,
add that part to the filename instead: `Enemy/spawn.cpp`, `Enemy/render.cpp`, or
`kubelet_pods.go`. Do not use a canonical file as a catch all for every concern.

## Split By Responsibility, Not By Size

A large file is a prompt to inspect its responsibilities, not a violation. Keep
it intact when its path and filename still give readers the best place to start.
Split it when a child responsibility has a clear name and a reader would look for
it separately.

Do not set a line limit. Do not split one cohesive algorithm, schema, or core
component only to make files smaller.

## Keep Entrypoints Small In Scope

Entrypoints own process startup, configuration, dependency wiring, shutdown, and
top level error handling. They do not own feature behavior, persistence,
transport handlers, or domain logic.

Move code out of an entrypoint when it gains a clear home. Do not create a broad
replacement package such as `app` just to empty the entrypoint.

## Improve Layout Without A Cleanup Project

Do not move everything because the current layout is weak. Keep unrelated code
where it is unless the user asks for a restructure.

When a task needs a new file or a clearer home, use the language convention and
move the smallest coherent unit needed for the task. Do not leave one
responsibility split between an old dumping ground and a new home without a clear
temporary reason. If a larger migration would help, describe the target layout
and leave it for a task that includes it.

## Completion Check

- The initial layout follows the language and framework convention, or a
  coherent stronger repository convention.
- New paths help a reader find code without opening unrelated files.
- New files represent an independently searchable responsibility, not a type or
  method count rule.
- Repeated path names have a clear canonical, type, or role meaning.
- Entrypoints contain only process level work.
- The change did not become an unrelated layout cleanup.
