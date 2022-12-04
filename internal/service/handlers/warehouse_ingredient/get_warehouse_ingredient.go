package handlers

import (
	"net/http"
	"strconv"
	"warehouse-service/internal/service/helpers"
	requests "warehouse-service/internal/service/requests/warehouse_ingredient"
	"warehouse-service/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetWarehouseIngredient(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewGetWarehouseIngredientRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	warehouseIngredient, err := helpers.WarehouseIngredientsQuery(r).FilterById(request.WarehouseIngredientID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get warehouse ingredient from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	if warehouseIngredient == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	relateIngredient, err := helpers.IngredientsQuery(r).FilterById(warehouseIngredient.IngredientId).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get ingredient")
		ape.RenderErr(w, problems.NotFound())
		return
	}
	relateDelivery, err := helpers.DeliveriesQuery(r).FilterById(warehouseIngredient.DeliveryId).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get delivery")
		ape.RenderErr(w, problems.NotFound())
		return
	}
	relateWarehouse, err := helpers.WarehousesQuery(r).FilterById(warehouseIngredient.WarehouseId).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get warehouse")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	var includes resources.Included
	includes.Add(&resources.Ingredient{
		Key: resources.NewKeyInt64(relateIngredient.Id, resources.INGREDIENT),
		Attributes: resources.IngredientAttributes{
			Name: relateIngredient.Name,
		},
	})
	includes.Add(&resources.Delivery{
		Key: resources.NewKeyInt64(relateDelivery.Id, resources.DELIVERY),
		Attributes: resources.DeliveryAttributes{
			Price: relateDelivery.Price,
			Date:  relateDelivery.Date,
		},
		Relationships: resources.DeliveryRelationships{
			Destination: resources.Relation{
				Data: &resources.Key{
					ID:   strconv.FormatInt(relateDelivery.DestinationId, 10),
					Type: resources.WAREHOUSE,
				},
			},
			Source: resources.Relation{
				Data: &resources.Key{
					ID:   strconv.FormatInt(relateDelivery.SourceId, 10),
					Type: resources.WAREHOUSE,
				},
			},
		},
	})
	includes.Add(&resources.Warehouse{
		Key: resources.NewKeyInt64(relateWarehouse.Id, resources.WAREHOUSE),
		Attributes: resources.WarehouseAttributes{
			CafeId:   relateWarehouse.CafeId,
			Capacity: relateWarehouse.Capacity,
		},
		Relationships: resources.WarehouseRelationships{
			Address: resources.Relation{
				Data: &resources.Key{
					ID:   strconv.FormatInt(relateWarehouse.AddressId, 10),
					Type: resources.ADDRESS,
				},
			},
		},
	})

	result := resources.WarehouseIngredientResponse{
		Data: resources.WarehouseIngredient{
			Key: resources.NewKeyInt64(warehouseIngredient.Id, resources.WAREHOUSE_INGREDIENT),
			Attributes: resources.WarehouseIngredientAttributes{
				Quantity:       warehouseIngredient.Quantity,
				Origin:         warehouseIngredient.Origin,
				Price:          warehouseIngredient.Price,
				ExpirationDate: warehouseIngredient.ExpirationDate,
			},
			Relationships: resources.WarehouseIngredientRelationships{
				Ingredient: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(warehouseIngredient.IngredientId, 10),
						Type: resources.INGREDIENT,
					},
				},
				Delivery: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(warehouseIngredient.DeliveryId, 10),
						Type: resources.DELIVERY,
					},
				},
				Warehouse: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(warehouseIngredient.WarehouseId, 10),
						Type: resources.WAREHOUSE,
					},
				},
			},
		},
		Included: includes,
	}

	ape.Render(w, result)
}
