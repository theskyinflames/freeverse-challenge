package app

import (
	"context"
	"errors"
	"fmt"

	"theskyinflames/graphql-challenge/internal/domain"

	"github.com/theskyinflames/cqrs-eda/pkg/bus"
	"github.com/theskyinflames/cqrs-eda/pkg/events"
)

// BuildEventsBus returns a generic events bus
func BuildEventsBus() bus.Bus {
	eventsBus := bus.New()
	eventsBus.Register(domain.ProductPurchasedEventName, busHandler(eventHandler()))
	return eventsBus
}

func eventHandler() events.Handler {
	return events.Handler(func(ev events.Event) {
		fmt.Printf("received event: %s from aggregate ID: %s\n", ev.Name(), ev.AggregateID().String())
	})
}

func busHandler(evh events.Handler) bus.Handler {
	return bus.Handler(func(_ context.Context, d bus.Dispatchable) (interface{}, error) {
		ev, ok := d.(events.Event)
		if !ok {
			return nil, errors.New("is not an event")
		}
		evh(ev)
		return nil, nil
	})
}
