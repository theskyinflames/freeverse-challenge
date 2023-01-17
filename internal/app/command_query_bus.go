package app

import (
	"github.com/theskyinflames/cqrs-eda/pkg/bus"
	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
	"github.com/theskyinflames/cqrs-eda/pkg/helpers"
)

// BuildCommandQueryBus returns the command/query bus
func BuildCommandQueryBus(log cqrs.Logger, eventsBus bus.Bus, pr ProductsRepository) bus.Bus {
	chMw := cqrs.CommandHandlerMultiMiddleware(
		cqrs.ChEventMw(eventsBus),
		cqrs.ChErrMw(log),
	)

	purchaseProduct := chMw(NewPurchaseProduct(pr))
	productsQh := cqrs.QhErrMw(log)(NewProducts(pr))

	bus := bus.New()
	bus.Register(PurchaseProductName, helpers.BusChHandler(purchaseProduct))
	bus.Register(ProductsName, helpers.BusQhHandler(productsQh))
	return bus
}
