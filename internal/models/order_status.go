package models

type OrderStatus int

const (
	Received  OrderStatus = 0
	Preparing OrderStatus = 1
	Ready     OrderStatus = 2
	Delivered OrderStatus = 3
)
