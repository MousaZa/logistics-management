package orders

type OrderEvent struct {
	Topic     string
	Payload   any
	Published bool
}

type EventLineItem struct {
	ProductUUID string
	Quantity    int
}

type OrderPlacedEvent struct {
	OrderUUID string
	LineItems []EventLineItem
}

type OrderCanceledEvent struct {
	OrderUUID string
}
