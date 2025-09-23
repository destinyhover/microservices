package domain

import "time"

type Payment struct {
	ID         int64  `json:"id"`
	CustomerID int64  `json:"customer_id"`
	Status     string `json:"status"`
	OrderID    int64  `json:"order_id"`
	TotalPrice int32  `json:"total_price"`
	CreatedAt  int64  `json:"created_at"`
}

func NewPayment(id int64, cusid int64, orderID int64, totalPrice int32) Payment {
	return Payment{ID: id, CustomerID: cusid, Status: "Pending",
		OrderID: orderID, TotalPrice: totalPrice, CreatedAt: time.Now().Unix()}
}
