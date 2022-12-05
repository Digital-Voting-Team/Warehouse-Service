package handlers

import (
	"github.com/Digital-Voting-Team/warehouse-service/internal/pkg/address"
	"github.com/Digital-Voting-Team/warehouse-service/internal/pkg/warehouse"
	"github.com/Digital-Voting-Team/warehouse-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/warehouse-service/internal/service/requests/warehouse"
	"github.com/Digital-Voting-Team/warehouse-service/resources"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetWarehouseList(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewGetWarehouseListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	warehousesQ := helpers.WarehousesQuery(r)
	applyFilters(warehousesQ, request)
	warehouses, err := warehousesQ.Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get warehouses")
		ape.Render(w, problems.InternalError())
		return
	}
	addresses, err := helpers.AddressesQuery(r).FilterById(getAddressIds(warehouses)...).Select()

	response := resources.WarehouseListResponse{
		Data:     newWarehousesList(warehouses),
		Links:    helpers.GetOffsetLinks(r, request.OffsetPageParams),
		Included: newWarehouseIncluded(addresses),
	}
	ape.Render(w, response)
}

func applyFilters(q warehouse.Query, request requests.GetWarehouseListRequest) {
	q.Page(request.OffsetPageParams)

	if len(request.FilterCapacity) > 0 {
		q.FilterByCapacity(request.FilterCapacity...)
	}

	if len(request.FilterAddressId) > 0 {
		q.FilterByAddressId(request.FilterAddressId...)
	}

	if len(request.FilterCafeId) > 0 {
		q.FilterByCafeId(request.FilterCafeId...)
	}
}

func newWarehousesList(warehouses []warehouse.Warehouse) []resources.Warehouse {
	result := make([]resources.Warehouse, len(warehouses))
	for i, currentWarehouse := range warehouses {
		result[i] = resources.Warehouse{
			Key: resources.NewKeyInt64(currentWarehouse.Id, resources.WAREHOUSE),
			Attributes: resources.WarehouseAttributes{
				CafeId:   currentWarehouse.CafeId,
				Capacity: currentWarehouse.Capacity,
			},
			Relationships: resources.WarehouseRelationships{
				Address: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(currentWarehouse.AddressId, 10),
						Type: resources.ADDRESS,
					},
				},
			},
		}
	}
	return result
}

func getAddressIds(warehouses []warehouse.Warehouse) []int64 {
	addressIDs := make([]int64, len(warehouses))
	for i := 0; i < len(warehouses); i++ {
		addressIDs[i] = warehouses[i].AddressId
	}
	return addressIDs
}

func newWarehouseIncluded(addresses []address.Address) resources.Included {
	result := resources.Included{}
	for _, item := range addresses {
		resource := newAddressModel(item)
		result.Add(&resource)
	}
	return result
}

func newAddressModel(address address.Address) resources.Address {
	return resources.Address{
		Key: resources.NewKeyInt64(address.Id, resources.ADDRESS),
		Attributes: resources.AddressAttributes{
			Building:   address.Building,
			Street:     address.Street,
			City:       address.City,
			District:   address.District,
			Region:     address.Region,
			PostalCode: address.PostalCode,
		},
	}
}
