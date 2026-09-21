# Comments And Docs Examples

[Back to rules](../references/standards.md#comments-and-docs)

Amounts in both samples are integer cents. `saveOrder` is idempotent for
`in.RequestID`, so a retry returns the existing order. The bad sample hides
that contract under narration. The good sample states it once.

## Bad: Comment spam narrates each line

```go
// processOrder processes an order.
func processOrder(ctx context.Context, in OrderInput) (*Receipt, error) {
    // check customer id
    if in.CustomerID == "" {
        // return invalid customer
        return nil, ErrInvalidCustomer
    }

    // check items
    if len(in.Items) == 0 {
        // return empty items
        return nil, ErrEmptyItems
    }

    // make subtotal variable
    subtotal := 0

    // loop items
    for _, item := range in.Items {
        // check quantity
        if item.Qty <= 0 {
            // return invalid qty
            return nil, ErrInvalidQty
        }

        // check price
        if item.UnitPriceCents < 0 {
            // return invalid price
            return nil, ErrInvalidPrice
        }

        // add line total
        subtotal += item.Qty * item.UnitPriceCents
    }

    // make discount variable
    discount := 0

    // check loyalty
    if in.IsLoyalCustomer {
        // set loyalty discount
        discount = subtotal / 10
    }

    // make tax variable
    tax := (subtotal - discount) * 7 / 100

    // make total variable
    total := subtotal - discount + tax

    // call repository save
    orderID, err := saveOrder(ctx, in, subtotal, discount, tax, total)
    if err != nil {
        // return save error
        return nil, err
    }

    // make receipt variable
    receipt := &Receipt{
        OrderID:  orderID,
        Subtotal: subtotal,
        Discount: discount,
        Tax:      tax,
        Total:    total,
    }

    // call notifier
    if err := sendConfirmation(ctx, in.CustomerID, orderID, total); err != nil {
        // return confirmation pending error
        return receipt, ErrConfirmationPending
    }

    // return receipt
    return receipt, nil
}
```

## Good: Comment the contract the code does not show

```go
// processOrder prices an order in integer cents and stores it.
// Loyalty discount and tax round down.
// Confirmation mail is a separate step. If the order is saved and the mail
// fails, the receipt is returned with ErrConfirmationPending. Retrying that
// error with the same RequestID must not create a second order.
func processOrder(ctx context.Context, in OrderInput) (*Receipt, error) {
    if in.CustomerID == "" {
        return nil, ErrInvalidCustomer
    }
    if len(in.Items) == 0 {
        return nil, ErrEmptyItems
    }

    subtotal := 0
    for _, item := range in.Items {
        if item.Qty <= 0 {
            return nil, ErrInvalidQty
        }
        if item.UnitPriceCents < 0 {
            return nil, ErrInvalidPrice
        }
        subtotal += item.Qty * item.UnitPriceCents
    }

    discount := 0
    if in.IsLoyalCustomer {
        discount = subtotal / 10
    }
    tax := (subtotal - discount) * 7 / 100
    total := subtotal - discount + tax

    orderID, err := saveOrder(ctx, in, subtotal, discount, tax, total)
    if err != nil {
        return nil, err
    }

    receipt := &Receipt{
        OrderID:  orderID,
        Subtotal: subtotal,
        Discount: discount,
        Tax:      tax,
        Total:    total,
    }
    if err := sendConfirmation(ctx, in.CustomerID, orderID, total); err != nil {
        return receipt, ErrConfirmationPending
    }
    return receipt, nil
}
```

Returning a bare error after `saveOrder` succeeds tells the caller the order
was not created. `ErrConfirmationPending` plus the receipt is the contract the
comment exists to state.

[Back to rules](../references/standards.md#comments-and-docs)
