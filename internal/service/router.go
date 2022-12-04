package service

import (
	"github.com/jmoiron/sqlx"
	"warehouse-service/internal/pkg/address"
	"warehouse-service/internal/pkg/delivery"
	"warehouse-service/internal/pkg/ingredient"
	usedIngredient "warehouse-service/internal/pkg/used_ingredient"
	"warehouse-service/internal/pkg/warehouse"
	"warehouse-service/internal/pkg/warehouse_ingredient"
	addressHandlers "warehouse-service/internal/service/handlers/address"
	deliveryHandlers "warehouse-service/internal/service/handlers/delivery"
	ingredientHandlers "warehouse-service/internal/service/handlers/ingredient"
	usedIngredientHandlers "warehouse-service/internal/service/handlers/used_ingredient"
	warehouseHandlers "warehouse-service/internal/service/handlers/warehouse"
	warehouseIngredientHandlers "warehouse-service/internal/service/handlers/warehouse_ingredient"
	"warehouse-service/internal/service/helpers"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router(db *sqlx.DB) chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			helpers.CtxLog(s.log),
			helpers.CtxAddressesQuery(address.NewQuery(db)),
			helpers.CtxDeliveriesQuery(delivery.NewQuery(db)),
			helpers.CtxIngredientsQuery(ingredient.NewQuery(db)),
			helpers.CtxUsedIngredientsQuery(usedIngredient.NewQuery(db)),
			helpers.CtxWarehousesQuery(warehouse.NewQuery(db)),
			helpers.CtxWarehouseIngredientsQuery(warehouse_ingredient.NewQuery(db)),
		),
	)
	r.Route("/integrations/warehouse-service", func(r chi.Router) {
		r.Route("/addresses", func(r chi.Router) {
			r.Get("/", addressHandlers.GetAddressList)
		})
		r.Route("/address", func(r chi.Router) {
			r.Post("/", addressHandlers.CreateAddress)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", addressHandlers.GetAddress)
				r.Put("/", addressHandlers.UpdateAddress)
				r.Delete("/", addressHandlers.DeleteAddress)
			})
		})
		r.Route("/deliveries", func(r chi.Router) {
			r.Get("/", deliveryHandlers.GetDeliveryList)
		})
		r.Route("/delivery", func(r chi.Router) {
			r.Post("/", deliveryHandlers.CreateDelivery)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", deliveryHandlers.GetDelivery)
				r.Put("/", deliveryHandlers.UpdateDelivery)
				r.Delete("/", deliveryHandlers.DeleteDelivery)
			})
		})
		r.Route("/ingredients", func(r chi.Router) {
			r.Get("/", ingredientHandlers.GetIngredientList)
		})
		r.Route("/ingredient", func(r chi.Router) {
			r.Post("/", ingredientHandlers.CreateIngredient)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", ingredientHandlers.GetIngredient)
				r.Put("/", ingredientHandlers.UpdateIngredient)
				r.Delete("/", ingredientHandlers.DeleteIngredient)
			})
		})
		r.Route("/used/ingredients", func(r chi.Router) {
			r.Get("/", usedIngredientHandlers.GetUsedIngredientList)
		})
		r.Route("/used/ingredient", func(r chi.Router) {
			r.Post("/", usedIngredientHandlers.CreateUsedIngredient)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", usedIngredientHandlers.GetUsedIngredient)
				r.Put("/", usedIngredientHandlers.UpdateUsedIngredient)
				r.Delete("/", usedIngredientHandlers.DeleteUsedIngredient)
			})
		})
		r.Route("/warehouses", func(r chi.Router) {
			r.Get("/", warehouseHandlers.GetWarehouseList)
		})
		r.Route("/warehouse", func(r chi.Router) {
			r.Post("/", warehouseHandlers.CreateWarehouse)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", warehouseHandlers.GetWarehouse)
				r.Put("/", warehouseHandlers.UpdateWarehouse)
				r.Delete("/", warehouseHandlers.DeleteWarehouse)
			})
		})
		r.Route("/warehouse/ingredients", func(r chi.Router) {
			r.Get("/", warehouseIngredientHandlers.GetWarehouseIngredientList)
		})
		r.Route("/warehouse/ingredient", func(r chi.Router) {
			r.Post("/", warehouseIngredientHandlers.CreateWarehouseIngredient)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", warehouseIngredientHandlers.GetWarehouseIngredient)
				r.Put("/", warehouseIngredientHandlers.UpdateWarehouseIngredient)
				r.Delete("/", warehouseIngredientHandlers.DeleteWarehouseIngredient)
			})
		})
	})

	return r
}
