package handlers

import (
	"github.com/Digital-Voting-Team/warehouse-service/internal/pkg/warehouse_ingredient"
	"github.com/Digital-Voting-Team/warehouse-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/warehouse-service/internal/service/requests/warehouse_ingredient"
	"github.com/Digital-Voting-Team/warehouse-service/resources"
	"github.com/spf13/cast"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func UpdateWarehouseIngredient(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewUpdateWarehouseIngredientRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	currentWarehouseIngredient, err := helpers.WarehouseIngredientsQuery(r).FilterById(request.WarehouseIngredientID).Get()
	if currentWarehouseIngredient == nil {
		helpers.Log(r).WithError(err).Info("did not found warehouseIngredient to update")
		ape.Render(w, problems.NotFound())
		return
	}

	//userId := r.Context().Value("userId").(int64)
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//_, _, warehouseIngredientId, err := helpers.GetIdsForGivenUser(r, userId)
	//if err != nil {
	//	helpers.Log(r).WithError(err).Info("wrong relations")
	//	ape.RenderErr(w, problems.InternalError())
	//	return
	//}
	//if *accessLevel != resources.Admin && warehouseIngredientId != warehouseIngredient.ID {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	newWarehouseIngredient := warehouse_ingredient.WarehouseIngredient{
		Quantity:       request.Data.Attributes.Quantity,
		Origin:         request.Data.Attributes.Origin,
		Price:          request.Data.Attributes.Price,
		ExpirationDate: request.Data.Attributes.ExpirationDate,
		IngredientId:   cast.ToInt64(request.Data.Relationships.Ingredient.Data.ID),
		WarehouseId:    cast.ToInt64(request.Data.Relationships.Warehouse.Data.ID),
		DeliveryId:     cast.ToInt64(request.Data.Relationships.Delivery.Data.ID),
	}

	relateIngredient, err := helpers.IngredientsQuery(r).FilterById(newWarehouseIngredient.IngredientId).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get ingredient")
		ape.RenderErr(w, problems.NotFound())
		return
	}
	relateDelivery, err := helpers.DeliveriesQuery(r).FilterById(newWarehouseIngredient.DeliveryId).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get delivery")
		ape.RenderErr(w, problems.NotFound())
		return
	}
	relateWarehouse, err := helpers.WarehousesQuery(r).FilterById(newWarehouseIngredient.WarehouseId).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get warehouse")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	var resultWarehouseIngredient warehouse_ingredient.WarehouseIngredient
	resultWarehouseIngredient, err = helpers.WarehouseIngredientsQuery(r).FilterById(currentWarehouseIngredient.Id).Update(newWarehouseIngredient)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to update warehouseIngredient")
		ape.RenderErr(w, problems.InternalError())
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
			Key: resources.NewKeyInt64(resultWarehouseIngredient.Id, resources.WAREHOUSE_INGREDIENT),
			Attributes: resources.WarehouseIngredientAttributes{
				Quantity:       resultWarehouseIngredient.Quantity,
				Origin:         resultWarehouseIngredient.Origin,
				Price:          resultWarehouseIngredient.Price,
				ExpirationDate: resultWarehouseIngredient.ExpirationDate,
			},
			Relationships: resources.WarehouseIngredientRelationships{
				Ingredient: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(resultWarehouseIngredient.IngredientId, 10),
						Type: resources.INGREDIENT,
					},
				},
				Delivery: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(resultWarehouseIngredient.DeliveryId, 10),
						Type: resources.DELIVERY,
					},
				},
				Warehouse: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(resultWarehouseIngredient.WarehouseId, 10),
						Type: resources.WAREHOUSE,
					},
				},
			},
		},
		Included: includes,
	}
	ape.Render(w, result)
}
