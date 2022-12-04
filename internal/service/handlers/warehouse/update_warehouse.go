package handlers

import (
	"github.com/spf13/cast"
	"net/http"
	"strconv"
	"warehouse-service/internal/pkg/warehouse"
	"warehouse-service/internal/service/helpers"
	requests "warehouse-service/internal/service/requests/warehouse"
	"warehouse-service/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func UpdateWarehouse(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewUpdateWarehouseRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	currentWarehouse, err := helpers.WarehousesQuery(r).FilterById(request.WarehouseID).Get()
	if currentWarehouse == nil {
		helpers.Log(r).WithError(err).Info("did not found warehouse to update")
		ape.Render(w, problems.NotFound())
		return
	}

	newWarehouse := warehouse.Warehouse{
		CafeId:    request.Data.Attributes.CafeId,
		Capacity:  request.Data.Attributes.Capacity,
		AddressId: cast.ToInt64(request.Data.Relationships.Address.Data.ID),
	}

	relateAddress, err := helpers.AddressesQuery(r).FilterById(newWarehouse.AddressId).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get new address")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	//userId := r.Context().Value("userId").(int64)
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//_, _, warehouseId, err := helpers.GetIdsForGivenUser(r, userId)
	//if err != nil {
	//	helpers.Log(r).WithError(err).Info("wrong relations")
	//	ape.RenderErr(w, problems.InternalError())
	//	return
	//}
	//if *accessLevel != resources.Admin && warehouseId != warehouse.ID {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	var resultWarehouse warehouse.Warehouse
	resultWarehouse, err = helpers.WarehousesQuery(r).FilterById(currentWarehouse.Id).Update(newWarehouse)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to update warehouse")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var includes resources.Included
	includes.Add(&resources.Address{
		Key: resources.NewKeyInt64(relateAddress.Id, resources.ADDRESS),
		Attributes: resources.AddressAttributes{
			Building:   relateAddress.Building,
			Street:     relateAddress.Street,
			City:       relateAddress.City,
			District:   relateAddress.District,
			Region:     relateAddress.Region,
			PostalCode: relateAddress.PostalCode,
		},
	})

	result := resources.WarehouseResponse{
		Data: resources.Warehouse{
			Key: resources.NewKeyInt64(resultWarehouse.Id, resources.WAREHOUSE),
			Attributes: resources.WarehouseAttributes{
				CafeId:   resultWarehouse.CafeId,
				Capacity: resultWarehouse.Capacity,
			},
			Relationships: resources.WarehouseRelationships{
				Address: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(resultWarehouse.AddressId, 10),
						Type: resources.ADDRESS,
					},
				},
			},
		},
		Included: includes,
	}
	ape.Render(w, result)
}
