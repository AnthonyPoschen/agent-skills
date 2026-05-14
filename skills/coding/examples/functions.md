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
