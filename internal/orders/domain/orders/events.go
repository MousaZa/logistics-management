package orders

type OrderEvent struct {
	Topic     string
	Payload   any
	Published bool
}

type OrderPlacedEvent struct {
	OrderUUID     string
	ProductsUUIDs []string
}
