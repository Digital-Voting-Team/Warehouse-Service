package handlers

import (
	"github.com/Digital-Voting-Team/warehouse-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/warehouse-service/internal/service/requests/warehouse"
	"github.com/Digital-Voting-Team/warehouse-service/resources"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetWarehouse(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewGetWarehouseRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	warehouse, err := helpers.WarehousesQuery(r).FilterById(request.WarehouseID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get warehouse from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	if warehouse == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	relateAddress, err := helpers.AddressesQuery(r).FilterById(warehouse.AddressId).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get address")
		ape.RenderErr(w, problems.NotFound())
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
		},
		Included: includes,
	}

	ape.Render(w, result)
}
