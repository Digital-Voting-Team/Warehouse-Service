package handlers

import (
	"net/http"
	"strconv"
	"warehouse-service/internal/pkg/delivery"
	"warehouse-service/internal/pkg/ingredient"
	"warehouse-service/internal/pkg/warehouse"
	"warehouse-service/internal/pkg/warehouse_ingredient"
	"warehouse-service/internal/service/helpers"
	requests "warehouse-service/internal/service/requests/warehouse_ingredient"
	"warehouse-service/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetWarehouseIngredientList(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewGetWarehouseIngredientListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	warehouseIngredientsQ := helpers.WarehouseIngredientsQuery(r)
	applyFilters(warehouseIngredientsQ, request)
	warehouseIngredients, err := warehouseIngredientsQ.Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get warehouse ingredients")
		ape.Render(w, problems.InternalError())
		return
	}

	ingredients, err := helpers.IngredientsQuery(r).FilterById(getIngredientIds(warehouseIngredients)...).Select()
	deliveries, err := helpers.DeliveriesQuery(r).FilterById(getDeliveryIds(warehouseIngredients)...).Select()
	warehouses, err := helpers.WarehousesQuery(r).FilterById(getWarehouseIds(warehouseIngredients)...).Select()

	response := resources.WarehouseIngredientListResponse{
		Data:     newWarehouseIngredientsList(warehouseIngredients),
		Links:    helpers.GetOffsetLinks(r, request.OffsetPageParams),
		Included: newWarehouseIngredientIncluded(ingredients, deliveries, warehouses),
	}
	ape.Render(w, response)
}

func applyFilters(q warehouse_ingredient.Query, request requests.GetWarehouseIngredientListRequest) {
	q.Page(request.OffsetPageParams)

	if len(request.FilterIngredientId) > 0 {
		q.FilterByIngredientId(request.FilterIngredientId...)
	}

	if len(request.FilterWarehouseId) > 0 {
		q.FilterByWarehouseId(request.FilterWarehouseId...)
	}

	if len(request.FilterQuantity) > 0 {
		q.FilterByQuantity(request.FilterQuantity...)
	}

	if len(request.FilterOrigin) > 0 {
		q.FilterByOrigin(request.FilterOrigin...)
	}
	if len(request.FilterPrice) > 0 {
		q.FilterByPrice(request.FilterPrice...)
	}

	if len(request.FilterExpirationDate) > 0 {
		q.FilterByExpirationDate(request.FilterExpirationDate...)
	}

	if len(request.FilterDeliveryId) > 0 {
		q.FilterByDeliveryId(request.FilterDeliveryId...)
	}

}

func newWarehouseIngredientsList(warehouseIngredients []warehouse_ingredient.WarehouseIngredient) []resources.WarehouseIngredient {
	result := make([]resources.WarehouseIngredient, len(warehouseIngredients))
	for i, currentWarehouseIngredient := range warehouseIngredients {
		result[i] = resources.WarehouseIngredient{
			Key: resources.NewKeyInt64(currentWarehouseIngredient.Id, resources.WAREHOUSE_INGREDIENT),
			Attributes: resources.WarehouseIngredientAttributes{
				Quantity:       currentWarehouseIngredient.Quantity,
				Origin:         currentWarehouseIngredient.Origin,
				Price:          currentWarehouseIngredient.Price,
				ExpirationDate: currentWarehouseIngredient.ExpirationDate,
			},
			Relationships: resources.WarehouseIngredientRelationships{
				Ingredient: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(currentWarehouseIngredient.IngredientId, 10),
						Type: resources.INGREDIENT,
					},
				},
				Delivery: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(currentWarehouseIngredient.DeliveryId, 10),
						Type: resources.DELIVERY,
					},
				},
				Warehouse: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(currentWarehouseIngredient.WarehouseId, 10),
						Type: resources.WAREHOUSE,
					},
				},
			},
		}
	}
	return result
}

func getIngredientIds(warehouseIngredients []warehouse_ingredient.WarehouseIngredient) []int64 {
	ingredientIDs := make([]int64, len(warehouseIngredients))
	for i := 0; i < len(warehouseIngredients); i++ {
		ingredientIDs[i] = warehouseIngredients[i].IngredientId
	}
	return ingredientIDs
}

func getDeliveryIds(warehouseIngredients []warehouse_ingredient.WarehouseIngredient) []int64 {
	deliveryIDs := make([]int64, len(warehouseIngredients))
	for i := 0; i < len(warehouseIngredients); i++ {
		deliveryIDs[i] = warehouseIngredients[i].DeliveryId
	}
	return deliveryIDs
}

func getWarehouseIds(warehouseIngredients []warehouse_ingredient.WarehouseIngredient) []int64 {
	warehouseIDs := make([]int64, len(warehouseIngredients))
	for i := 0; i < len(warehouseIngredients); i++ {
		warehouseIDs[i] = warehouseIngredients[i].WarehouseId
	}
	return warehouseIDs
}

func newWarehouseIngredientIncluded(ingredients []ingredient.Ingredient, deliveries []delivery.Delivery, warehouses []warehouse.Warehouse) resources.Included {
	result := resources.Included{}
	for _, item := range ingredients {
		resource := newIngredientModel(item)
		result.Add(&resource)
	}
	for _, item := range deliveries {
		resource := newDeliveryModel(item)
		result.Add(&resource)
	}
	for _, item := range warehouses {
		resource := newWarehouseModel(item)
		result.Add(&resource)
	}
	return result
}

func newIngredientModel(ingredient ingredient.Ingredient) resources.Ingredient {
	return resources.Ingredient{
		Key: resources.NewKeyInt64(ingredient.Id, resources.INGREDIENT),
		Attributes: resources.IngredientAttributes{
			Name: ingredient.Name,
		},
	}
}

func newDeliveryModel(delivery delivery.Delivery) resources.Delivery {
	return resources.Delivery{
		Key: resources.NewKeyInt64(delivery.Id, resources.DELIVERY),
		Attributes: resources.DeliveryAttributes{
			Date:  delivery.Date,
			Price: delivery.Price,
		},
		Relationships: resources.DeliveryRelationships{
			Destination: resources.Relation{
				Data: &resources.Key{
					ID:   strconv.FormatInt(delivery.DestinationId, 10),
					Type: resources.WAREHOUSE,
				},
			},
			Source: resources.Relation{
				Data: &resources.Key{
					ID:   strconv.FormatInt(delivery.SourceId, 10),
					Type: resources.WAREHOUSE,
				},
			},
		},
	}
}

func newWarehouseModel(warehouse warehouse.Warehouse) resources.Warehouse {
	return resources.Warehouse{
		Key: resources.NewKeyInt64(warehouse.Id, resources.WAREHOUSE),
		Attributes: resources.WarehouseAttributes{
			CafeId:   warehouse.CafeId,
			Capacity: warehouse.Capacity,
		},
		Relationships: resources.WarehouseRelationships{
			Address: resources.Relation{
				Data: &resources.Key{
					ID:   strconv.FormatInt(warehouse.AddressId, 10),
					Type: resources.ADDRESS,
				},
			},
		},
	}
}
