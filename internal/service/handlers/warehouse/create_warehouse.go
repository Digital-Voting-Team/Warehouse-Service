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

func CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewCreateWarehouseRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	dataWarehouse := warehouse.Warehouse{
		CafeId:    request.Data.Attributes.CafeId,
		Capacity:  request.Data.Attributes.Capacity,
		AddressId: cast.ToInt64(request.Data.Relationships.Address.Data.ID),
	}

	var resultWarehouse warehouse.Warehouse
	relateAddress, err := helpers.AddressesQuery(r).FilterById(dataWarehouse.AddressId).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get address")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	resultWarehouse, err = helpers.WarehousesQuery(r).Insert(dataWarehouse)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create warehouse")
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
