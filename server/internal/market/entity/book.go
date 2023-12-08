package entity

import (
	"container/heap"
	"sync"
)

// Book represents a trading book.
type Book struct {
	Orders              []*Order
	Transactions        []*Transaction
	OrdersChannel       chan *Order
	OrdersChannelOutput chan *Order
	Wg                  *sync.WaitGroup
}

// NewBook creates a new Book instance.
func NewBook(orderChannel chan *Order, orderChannelOutput chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Orders:              []*Order{},
		Transactions:        []*Transaction{},
		OrdersChannel:       orderChannel,
		OrdersChannelOutput: orderChannelOutput,
		Wg:                  wg,
	}
}

// Trade executes trading logic for buy and sell orders.
func (b *Book) Trade() {
	buyOrders := NewOrderQueue()
	sellOrders := NewOrderQueue()

	heap.Init(buyOrders)
	heap.Init(sellOrders)

	for order := range b.OrdersChannel {
		switch order.OrderType {
		case "BUY":
			b.handleBuyOrder(order, sellOrders)
		case "SELL":
			b.handleSellOrder(order, buyOrders)
		}
	}
}

func (b *Book) handleBuyOrder(buyOrder *Order, sellOrders *OrderQueue) {
	sellOrders.Push(buyOrder)
	if sellOrders.Len() > 0 && sellOrders.Orders[0].Price <= buyOrder.Price {
		sellOrder := sellOrders.Pop().(*Order)
		b.executeTrade(buyOrder, sellOrder)
	}
}

func (b *Book) handleSellOrder(sellOrder *Order, buyOrders *OrderQueue) {
	if buyOrders.Len() > 0 && buyOrders.Orders[0].Price >= sellOrder.Price {
		buyOrder := buyOrders.Pop().(*Order)
		b.executeTrade(buyOrder, sellOrder)
	} else {
		buyOrders.Push(sellOrder)
	}
}

func (b *Book) executeTrade(buyOrder, sellOrder *Order) {
	transaction := NewTransaction(sellOrder, buyOrder, buyOrder.Shares, sellOrder.Price)

	b.AddTransaction(transaction, b.Wg)
	b.OrdersChannelOutput <- sellOrder
	b.OrdersChannelOutput <- buyOrder
}

func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	min := min(transaction.SellingOrder.PendingShares, transaction.BuyingOrder.PendingShares)

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -min)
	transaction.SellingOrder.RemovePendingShares(min)
	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.Asset.ID, -min)
	transaction.BuyingOrder.RemovePendingShares(min)

	transaction.CalculateTotal()
	transaction.CloseOrders()

	b.Transactions = append(b.Transactions, transaction)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
