# Functions Examples

[Back to rules](../SKILL.md#functions)

## Good: Short Linear Phases Stay Inline

```go
func buildInvoice(in InvoiceInput, now time.Time) (Invoice, error) {
    if in.CustomerID == "" {
        return Invoice{}, ErrInvalidCustomer
    }
    if len(in.Lines) == 0 {
        return Invoice{}, ErrEmptyInvoice
    }

    subtotal := 0
    for _, line := range in.Lines {
        if line.Qty <= 0 {
            return Invoice{}, ErrInvalidQty
        }
        subtotal += line.Qty * line.UnitPrice
    }

    discount := 0
    if in.DiscountCode != "" {
        discount = lookupDiscount(in.DiscountCode, subtotal)
    }

    tax := (subtotal - discount) * taxRate / taxScale
    total := subtotal - discount + tax

    return Invoice{
        CustomerID: in.CustomerID,
        DueAt:      now.Add(invoiceDueAfter),
        Subtotal:   subtotal,
        Discount:   discount,
        Tax:        tax,
        Total:      total,
    }, nil
}
```

## Bad: One Helper Per Small Phase

```go
func buildInvoice(in InvoiceInput, now time.Time) (Invoice, error) {
    if err := validateInvoiceInput(in); err != nil {
        return Invoice{}, err
    }

    subtotal := computeSubtotal(in.Lines)
    discount := computeDiscount(in.DiscountCode, subtotal)
    tax := computeTax(subtotal, discount)

    return makeInvoice(in.CustomerID, now, subtotal, discount, tax), nil
}
```

## Good: Extract Cohesive Substantial Phase Groups

```go
func convertReport(in RawReport) (Report, error) {
    clean, err := sanitizeRawReport(in)
    if err != nil {
        return Report{}, err
    }

    groups := groupReportRows(clean.Rows)
    summaries := computeReportSummaries(groups)

    return Report{
        Metadata:  clean.Metadata,
        Groups:    groups,
        Summaries: summaries,
    }, nil
}
```

`sanitizeRawReport` is a good extraction when it groups related validation,
normalization, status mapping, and amount parsing. Avoid separate helpers for
each tiny sanitize step unless one becomes substantial or reusable.

## Good: Extract Shared Traversal Behavior

```go
func anyBoundCode(devices []Device, codes []Code, query Query) bool {
    for _, dev := range devices {
        for _, code := range codes {
            if queryCode(dev, code, query) {
                return true
            }
        }
    }
    return false
}

func down(devices []Device, action Action) bool {
    return anyBoundCode(devices, action.Codes, QueryDown)
}

func pressed(devices []Device, action Action) bool {
    return anyBoundCode(devices, action.Codes, QueryPressed)
}
```

Even if each caller is short, these functions share behavior: iterate attached
devices and bound codes, apply one query, and return on the first match.
Extracting that behavior prevents future callers from drifting on ordering,
filtering, early-return behavior, or edge handling.

## Good: Keep Public Names, Share Repeated Bodies

```zig
const ButtonQuery = enum {
    pressed,
    released,
};

fn buttonTransition(
    previous: ButtonState,
    current: ButtonState,
    query: ButtonQuery,
) bool {
    return switch (query) {
        .pressed => previous == .up and current == .down,
        .released => previous == .down and current == .up,
    };
}

pub fn pressed(self: *const MouseDevice, code: InputCode) bool {
    const idx = mouseIndex(code) orelse return false;
    return buttonTransition(self.prev_buttons[idx], self.buttons[idx], .pressed);
}

pub fn released(self: *const MouseDevice, code: InputCode) bool {
    const idx = mouseIndex(code) orelse return false;
    return buttonTransition(self.prev_buttons[idx], self.buttons[idx], .released);
}
```

The public methods remain because they are the consumer-facing API. The helper
exists because the bodies differ only by expected states, so centralizing the
transition logic reduces the chance that one method drifts.

## Good: Pass Concrete Varying State Into Shared Helpers

```zig
pub fn down(self: *const KeyboardDevice, code: InputCode) bool {
    return keyStateIs(self.keys[0..], code, .down) orelse false;
}

pub fn pressed(self: *const KeyboardDevice, code: InputCode) bool {
    const previous_up = keyStateIs(
        self.prev_keys[0..],
        code,
        .up,
    ) orelse return false;
    const current_down = keyStateIs(
        self.keys[0..],
        code,
        .down,
    ) orelse return false;
    return previous_up and current_down;
}

pub fn released(self: *const KeyboardDevice, code: InputCode) bool {
    const previous_down = keyStateIs(
        self.prev_keys[0..],
        code,
        .down,
    ) orelse return false;
    const current_up = keyStateIs(
        self.keys[0..],
        code,
        .up,
    ) orelse return false;
    return previous_down and current_up;
}

fn keyStateIs(
    keys: []const ButtonState,
    code: InputCode,
    state: ButtonState,
) ?bool {
    const idx: usize = @intFromEnum(code);
    if (idx >= keys.len) return null;
    return keys[idx] == state;
}
```

This helper centralizes the shared policy: translate the code to an index,
reject incompatible codes, and compare against the requested state. The callers
still spell out the transition they need by passing `self.keys`, `self.prev_keys`,
`.up`, and `.down` directly.

## Bad: Add A Mode Enum When Callers Already Have The Data

```zig
const KeyFrame = enum {
    current,
    previous,
};

fn keyStateIs(
    self: *const KeyboardDevice,
    frame: KeyFrame,
    code: InputCode,
    state: ButtonState,
) ?bool {
    const keys = switch (frame) {
        .current => self.keys,
        .previous => self.prev_keys,
    };
    const idx: usize = @intFromEnum(code);
    if (idx >= keys.len) return null;
    return keys[idx] == state;
}
```

The enum adds another decision table without adding policy. It makes the helper
less general and forces readers to map `.current` and `.previous` back to fields
that the caller could have passed directly.

## Good: Hoist Shared Preparation From Branches

```zig
pub fn prevAxis1d(self: *const GamepadDevice, code: InputCode) ?Axis1d {
    const left = applyDeadzone2d(self.prev_left_stick, self.left_stick_deadzone);
    const right = applyDeadzone2d(
        self.prev_right_stick,
        self.right_stick_deadzone,
    );
    const left_trigger = applyDeadzone(
        self.prev_left_trigger_value,
        self.left_trigger_deadzone,
    );
    const right_trigger = applyDeadzone(
        self.prev_right_trigger_value,
        self.right_trigger_deadzone,
    );

    return switch (code) {
        .gamepad_left_trigger => left_trigger,
        .gamepad_right_trigger => right_trigger,
        .gamepad_left_stick_up => positive(left.y),
        .gamepad_left_stick_down => positive(-left.y),
        .gamepad_left_stick_left => positive(-left.x),
        .gamepad_left_stick_right => positive(left.x),
        .gamepad_right_stick_up => positive(right.y),
        .gamepad_right_stick_down => positive(-right.y),
        .gamepad_right_stick_left => positive(-right.x),
        .gamepad_right_stick_right => positive(right.x),
        else => null,
    };
}
```

The switch now expresses which prepared value each code selects. The repeated
deadzone normalization is named once, which makes future changes to preparation
logic less error-prone.

## Good: Centralize Repeated Object Configuration

```go
type Button struct {
    Label    string
    Disabled bool
    Tone     Tone
    OnClick  func()
}

func (b *Button) ConfigureDanger(label string, onClick func()) {
    b.Label = label
    b.Disabled = false
    b.Tone = ToneDanger
    b.OnClick = onClick
}

func deleteButton(onClick func()) Button {
    var button Button
    button.ConfigureDanger("Delete", onClick)
    return button
}

func removeButton(onClick func()) Button {
    var button Button
    button.ConfigureDanger("Remove", onClick)
    return button
}
```

The fields may stay public for literals, tests, or low-level access. The helper
exists because multiple callers share a configuration behavior that may need new
defaults, validation, or derived state later.

## Good: Parent Retains Orchestration

```go
func reconcileAccount(ctx context.Context, id string, deps Deps) error {
    state, err := loadReconciliationState(ctx, id, deps.DB)
    if err != nil {
        return err
    }

    adjustments, err := computeLedgerAdjustments(state)
    if err != nil {
        return err
    }

    if len(adjustments) == 0 {
        return nil
    }

    return saveLedgerAdjustments(ctx, deps.DB, adjustments)
}
```

This parent is still useful because it owns ordering, error handling, and the
empty-adjustment policy. If the parent only forwarded calls without policy or
branching, reconsider whether the caller should compose the steps directly.
