package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
	"github.com/theskyinflames/cqrs-eda/pkg/events"
)

// PurchaseProductCmd is a command
type PurchaseProductCmd struct {
	ID uuid.UUID
}

// PurchaseProductName is self-described
var PurchaseProductName = "purchase.product"

// Name implements the Command interface
func (cmd PurchaseProductCmd) Name() string {
	return PurchaseProductName
}

// PurchaseProduct is a command handler
type PurchaseProduct struct {
	pr ProductsRepository
}

// NewPurchaseProduct is a constructor
func NewPurchaseProduct(pr ProductsRepository) PurchaseProduct {
	return PurchaseProduct{pr: pr}
}

// Handle implements CommandHandler interface
func (ch PurchaseProduct) Handle(ctx context.Context, cmd cqrs.Command) ([]events.Event, error) {
	co, ok := cmd.(PurchaseProductCmd)
	if !ok {
		return nil, NewInvalidCommandError(PurchaseProductName, cmd.Name())
	}

	p, err := ch.pr.FindByID(ctx, co.ID)
	if err != nil {
		return nil, err
	}

	if err := p.Purchase(); err != nil {
		return nil, err
	}

	if err := ch.pr.UpdateAvailable(ctx, p); err != nil {
		return nil, err
	}

	return p.Events(), nil
}
