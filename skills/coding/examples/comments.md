# Comments And Docs Examples

[Back to rules](../SKILL.md#comments-and-docs)

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
        if item.UnitPrice < 0 {
            // return invalid price
            return nil, ErrInvalidPrice
        }

        // add line total
        subtotal += item.Qty * item.UnitPrice
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

    // call notifier
    if err := sendConfirmation(ctx, in.CustomerID, orderID, total); err != nil {
        // return notification error
        return nil, err
    }

    // return receipt
    return &Receipt{
        OrderID:  orderID,
        Subtotal: subtotal,
        Discount: discount,
        Tax:      tax,
        Total:    total,
    }, nil
}
```

## Good: One comment per phase, skimmable story

```go
// processOrder validates input, prices order, persists it, and returns receipt.
func processOrder(ctx context.Context, in OrderInput) (*Receipt, error) {
    // Validate request and fail fast on invalid input.
    if in.CustomerID == "" {
        return nil, ErrInvalidCustomer
    }
    if len(in.Items) == 0 {
        return nil, ErrEmptyItems
    }

    // Price line items and compute subtotal.
    subtotal := 0
    for _, item := range in.Items {
        if item.Qty <= 0 {
            return nil, ErrInvalidQty
        }
        if item.UnitPrice < 0 {
            return nil, ErrInvalidPrice
        }
        subtotal += item.Qty * item.UnitPrice
    }

    // Apply discount and tax to compute total.
    discount := 0
    if in.IsLoyalCustomer {
        discount = subtotal / 10
    }
    tax := (subtotal - discount) * 7 / 100
    total := subtotal - discount + tax

    // Persist order and notify customer.
    orderID, err := saveOrder(ctx, in, subtotal, discount, tax, total)
    if err != nil {
        return nil, err
    }
    if err := sendConfirmation(ctx, in.CustomerID, orderID, total); err != nil {
        return nil, err
    }

    // Return receipt payload.
    return &Receipt{
        OrderID:  orderID,
        Subtotal: subtotal,
        Discount: discount,
        Tax:      tax,
        Total:    total,
    }, nil
}
```

[Back to rules](../SKILL.md#comments-and-docs)
