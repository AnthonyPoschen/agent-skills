# Go Layout

Use this reference with `./file-organization.md` when starting a Go application,
placing Go packages, or reorganizing Go files.

## Start With Go Conventions

Follow the existing Go module layout when it is coherent. Otherwise use the
standard Go package model and the category based layout used by large Go projects
such as Kubernetes. Create only the categories the project needs.

```text
cmd/<binary>/main.go     executable entrypoint
internal/<area>/         private application code
pkg/<name>/              deliberate reusable library code
api/<version>/           API types and schemas, when needed
test/                    integration and end to end support, when needed
hack/                    developer tooling, when needed
```

`internal` is the usual home for supporting application packages. Go enforces
its import boundary. Use `pkg` only for code that is deliberately reusable and
worth supporting as an import surface. Do not create empty `api`, `pkg`, `test`,
or `hack` directories because another project has them.

## Keep Commands As Commands

`cmd/<binary>/main.go` is the process entrypoint. It reads configuration,
constructs dependencies, starts the application, and handles top level shutdown
and errors. Do not put feature code, database access, HTTP handlers, or domain
logic under `cmd/<binary>` merely because the program started there.

When a command grows beyond startup work, move the feature or adapter into a
purposeful package outside `cmd`. Do not replace one dumping ground with a broad
`internal/app` package.

## Name Packages And Files For Discovery

Use package names that name the area of work. Within a package, use a canonical
file for its main type or component when that helps navigation, then name related
files by the part they own:

```text
kubelet/
  kubelet.go       main Kubelet type and core work
  kubelet_pods.go  pod related Kubelet work
  kubelet_stats.go stats related Kubelet work
```

Keep related schema types together when they form one contract. Use a separate
file when its type or behavior is an independently searchable concept. Do not
force one Go type, method, or interface into its own file.

Go package names should be short and specific. Avoid `util`, `utils`, `common`,
`helpers`, `manager`, and broad `app` packages when a domain or role name exists.

## When The Existing Layout Is Weak

Do not keep placing new code under `cmd/<binary>` or another dumping ground just
because existing code lives there. Put new work in the correct conventional
package when the task can do so cleanly. Do not move unrelated code only to make
the repository match a target layout. Suggest a broader migration when it would
help, and perform it only when the user asks or the current change needs it.
