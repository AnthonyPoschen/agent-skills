# Variables And Magic Values Examples

[Back to rules](../SKILL.md#variables-and-state)

## Bad: Mutable global runtime state

```go
var db *sql.DB

func initDB() error {
    conn, err := sql.Open("postgres", os.Getenv("DB_DSN"))
    if err != nil {
        return err
    }
    db = conn
    return nil
}

func loadUser(ctx context.Context, id string) (User, error) {
    return queryUser(ctx, db, id)
}
```

## Good: Create in main, inject dependency

```go
type Store struct {
    DB *sql.DB
}

func main() {
    db := mustOpenDB()
    store := Store{DB: db}
    _ = runServer(store)
}

func (s Store) loadUser(ctx context.Context, id string) (User, error) {
    return queryUser(ctx, s.DB, id)
}
```

## Bad: Repeated magic literals

```go
if status == "ready" {
    return "ready"
}
if prev == "ready" {
    return "ready"
}
```

## Good: Shared named constant

```go
const statusReady = "ready"

if status == statusReady {
    return statusReady
}
if prev == statusReady {
    return statusReady
}
```
