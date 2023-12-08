package entity

import "errors"

type Order struct {
	ID            string
	Investor      *Investor
	Asset         *Asset
	Shares        int
	PendingShares int
	Price         float64
	OrderType     string
	Status        string
	Transactions  []*Transaction
}

func NewOrder(
	orderId string,
	investor *Investor,
	asset *Asset,
	shares int,
	price float64,
	orderType string,
) *Order {
	return &Order{
		ID:            orderId,
		Investor:      investor,
		Asset:         asset,
		Shares:        shares,
		PendingShares: shares,
		Price:         price,
		OrderType:     orderType,
		Status:        "OPEN",
		Transactions:  []*Transaction{},
	}
}

func (o *Order) Close() error {
	if o.PendingShares <= 0 {
		o.Status = "CLOSED"
		return nil
	}

	return errors.New("cannot close order with pending shares")
}

func (o *Order) RemovePendingShares(qnt int) {
	o.PendingShares -= qnt
}
