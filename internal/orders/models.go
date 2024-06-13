package orders

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	OrderLines []*OrderLine       `bson:"order_lines,omitempty"`
	TotalPrice int32              `bson:"total_price,omitempty"`
	OrderDate  string             `bson:"order_date,omitempty"`
}

type OrderLine struct {
	BookId   string `bson:"book_id,omitempty"`
	Quantity int32  `bson:"quantity,omitempty"`
}

func NewOrder(orderLines []*OrderLine, totalPrice int32, orderDate string) *Order {
	return &Order{
		ID:         primitive.NewObjectID(),
		OrderLines: orderLines,
		TotalPrice: totalPrice,
		OrderDate:  orderDate,
	}
}

func NewOrderLine(bookId string, quantity int32) *OrderLine {
	return &OrderLine{
		BookId:   bookId,
		Quantity: quantity,
	}
}
