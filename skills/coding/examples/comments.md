# Comments And Docs Examples

[Back to rules](../references/standards.md#comments-and-docs)

Amounts in both samples are integer cents. `saveOrder` is idempotent for
`in.RequestID`, so a retry returns the existing order. The bad sample hides
that contract under narration. The good sample moves it into a type and a test.

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

## Good: Types and tests carry the contract

```go
type Cents int64

func (c Cents) Percent(p int64) Cents { return c * Cents(p) / 100 }

type RequestID string

func processOrder(ctx context.Context, in OrderInput) (*Receipt, error) {
    if in.CustomerID == "" {
        return nil, ErrInvalidCustomer
    }
    if len(in.Items) == 0 {
        return nil, ErrEmptyItems
    }

    var subtotal Cents
    for _, item := range in.Items {
        if item.Qty <= 0 {
            return nil, ErrInvalidQty
        }
        if item.UnitPrice < 0 {
            return nil, ErrInvalidPrice
        }
        subtotal += Cents(item.Qty) * item.UnitPrice
    }

    var loyaltyDiscount Cents
    if in.IsLoyalCustomer {
        loyaltyDiscount = subtotal.Percent(10)
    }
    tax := (subtotal - loyaltyDiscount).Percent(7)
    receipt := &Receipt{
        Subtotal: subtotal,
        Discount: loyaltyDiscount,
        Tax:      tax,
        Total:    subtotal - loyaltyDiscount + tax,
    }

    orderID, err := saveOrderOnce(ctx, in.RequestID, in, receipt)
    if err != nil {
        return nil, err
    }
    receipt.OrderID = orderID

    if err := sendConfirmation(ctx, in.CustomerID, orderID, receipt.Total); err != nil {
        return receipt, ErrConfirmationPending
    }
    return receipt, nil
}
```

```go
func TestProcessOrderRetryAfterConfirmationFailureReusesOrder(t *testing.T) {
    store := newMemoryOrderStore()
    mail := &failingMailer{}
    svc := newOrderService(store, mail)
    in := OrderInput{
        RequestID:  "req-1",
        CustomerID: "cust-1",
        Items:      []Item{{Qty: 2, UnitPrice: Cents(1050)}},
    }

    first, err := svc.processOrder(context.Background(), in)
    if !errors.Is(err, ErrConfirmationPending) {
        t.Fatalf("first call err = %v, want ErrConfirmationPending", err)
    }
    second, _ := svc.processOrder(context.Background(), in)

    if second.OrderID != first.OrderID {
        t.Fatalf("retry created order %q, want %q", second.OrderID, first.OrderID)
    }
    if got := store.Count(); got != 1 {
        t.Fatalf("stored %d orders, want 1", got)
    }
    if first.Total != Cents(2247) {
        t.Fatalf("total = %d, want 2247", first.Total)
    }
}
```

`Cents` makes integer money a type the compiler checks, so no comment has to
say "amounts are cents". `saveOrderOnce` and the `RequestID` type name the
idempotency key. The test proves a retry after `ErrConfirmationPending` reuses
the order. If someone breaks that, the test fails. A comment would just go
stale.

[Back to rules](../references/standards.md#comments-and-docs)
