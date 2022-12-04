package helpers

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"net/http"
	"warehouse-service/internal/pkg/address"
	"warehouse-service/internal/pkg/delivery"
	"warehouse-service/internal/pkg/ingredient"
	usedIngredient "warehouse-service/internal/pkg/used_ingredient"
	"warehouse-service/internal/pkg/warehouse"
	warehouseIngredient "warehouse-service/internal/pkg/warehouse_ingredient"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	addressesQueryCtxKey
	warehousesQueryCtxKey
	deliveriesQueryCtxKey
	ingredientsQueryCtxKey
	warehouseIngredientsQueryCtxKey
	usedIngredientsQueryCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxAddressesQuery(entry address.Query) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, addressesQueryCtxKey, entry)
	}
}

func AddressesQuery(r *http.Request) address.Query {
	return r.Context().Value(addressesQueryCtxKey).(address.Query).New()
}

func CtxDeliveriesQuery(entry delivery.Query) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, deliveriesQueryCtxKey, entry)
	}
}

func DeliveriesQuery(r *http.Request) delivery.Query {
	return r.Context().Value(deliveriesQueryCtxKey).(delivery.Query).New()
}

func CtxWarehousesQuery(entry warehouse.Query) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, warehousesQueryCtxKey, entry)
	}
}

func WarehousesQuery(r *http.Request) warehouse.Query {
	return r.Context().Value(warehousesQueryCtxKey).(warehouse.Query).New()
}

func CtxIngredientsQuery(entry ingredient.Query) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ingredientsQueryCtxKey, entry)
	}
}

func IngredientsQuery(r *http.Request) ingredient.Query {
	return r.Context().Value(ingredientsQueryCtxKey).(ingredient.Query).New()
}

func CtxUsedIngredientsQuery(entry usedIngredient.Query) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, usedIngredientsQueryCtxKey, entry)
	}
}

func UsedIngredientsQuery(r *http.Request) usedIngredient.Query {
	return r.Context().Value(usedIngredientsQueryCtxKey).(usedIngredient.Query).New()
}

func CtxWarehouseIngredientsQuery(entry warehouseIngredient.Query) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, warehouseIngredientsQueryCtxKey, entry)
	}
}

func WarehouseIngredientsQuery(r *http.Request) warehouseIngredient.Query {
	return r.Context().Value(warehouseIngredientsQueryCtxKey).(warehouseIngredient.Query).New()
}
