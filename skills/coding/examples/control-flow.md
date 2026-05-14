# Control Flow Examples

[Back to rules](../SKILL.md#control-flow)

## Bad: Nested happy path

```go
func process(user *User, req *Request) error {
    if user != nil {
        if req != nil {
            if req.Valid() {
                return runProcess(user, req)
            }
        }
    }
    return ErrInvalidInput
}
```

## Good: Guard clauses first

```go
func process(user *User, req *Request) error {
    if user == nil {
        return ErrInvalidInput
    }
    if req == nil {
        return ErrInvalidInput
    }
    if !req.Valid() {
        return ErrInvalidInput
    }

    return runProcess(user, req)
}
```
