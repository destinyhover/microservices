package domain

import "time"

type OrderItem struct {
	ProductCode string  `json:"product_code"`
	UnitPrice   float32 `json:"unit_price"`
	Quantity    int32   `json:"quantity"`
}
type Order struct {
	ID         int64       `json:"id"`
	CustomerId int64       `json:"customer_id"`
	Status     string      `json:"status"`
	OrderItems []OrderItem `json:"order_items"`
	CreatedAt  int64       `json:"created_at"`
}

func NewOrder(customerId int64, orderItems []OrderItem) Order {
	return Order{
		CreatedAt:  time.Now().Unix(),
		CustomerId: customerId,
		Status:     "Pending",
		OrderItems: orderItems,
	}
}
func (o *Order) TotalPrice() float32 {
	var sum float32
	for _, v := range o.OrderItems {
		sum += v.UnitPrice
	}
	return sum
}
