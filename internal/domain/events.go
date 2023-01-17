package domain

import (
	"github.com/theskyinflames/cqrs-eda/pkg/events"
)

// ProductPurchasedEventName is self-described
const ProductPurchasedEventName = "product.purchased"

// ProductPurchasedEvent is an event
type ProductPurchasedEvent struct {
	events.EventBasic
}

// NewProductPurchasedEvent is a constructor
func NewProductPurchasedEvent(p Product) ProductPurchasedEvent {
	return ProductPurchasedEvent{
		EventBasic: events.NewEventBasic(p.ID(), ProductPurchasedEventName, nil),
	}
}
