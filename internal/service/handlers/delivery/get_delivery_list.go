package handlers

import (
	"github.com/Digital-Voting-Team/warehouse-service/internal/pkg/delivery"
	"github.com/Digital-Voting-Team/warehouse-service/internal/pkg/warehouse"
	"github.com/Digital-Voting-Team/warehouse-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/warehouse-service/internal/service/requests/delivery"
	"github.com/Digital-Voting-Team/warehouse-service/resources"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetDeliveryList(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewGetDeliveryListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	deliveriesQ := helpers.DeliveriesQuery(r)
	applyFilters(deliveriesQ, request)
	deliveries, err := deliveriesQ.Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get deliveries")
		ape.Render(w, problems.InternalError())
		return
	}
	destinations, err := helpers.WarehousesQuery(r).FilterById(getDestinationIds(deliveries)...).Select()
	sources, err := helpers.WarehousesQuery(r).FilterById(getSourceIds(deliveries)...).Select()

	response := resources.DeliveryListResponse{
		Data:     newDeliveriesList(deliveries),
		Links:    helpers.GetOffsetLinks(r, request.OffsetPageParams),
		Included: newDeliveryIncluded(destinations, sources),
	}
	ape.Render(w, response)
}

func applyFilters(q delivery.Query, request requests.GetDeliveryListRequest) {
	q.Page(request.OffsetPageParams)

	if len(request.FilterSourceId) > 0 {
		q.FilterBySourceId(request.FilterSourceId...)
	}

	if len(request.FilterDestinationId) > 0 {
		q.FilterByDestinationId(request.FilterDestinationId...)
	}

	if len(request.FilterPrice) > 0 {
		q.FilterByPrice(request.FilterPrice...)
	}

	if len(request.FilterDate) > 0 {
		q.FilterByDate(request.FilterDate...)
	}
}

func newDeliveriesList(deliveries []delivery.Delivery) []resources.Delivery {
	result := make([]resources.Delivery, len(deliveries))
	for i, currentDelivery := range deliveries {
		result[i] = resources.Delivery{
			Key: resources.NewKeyInt64(currentDelivery.Id, resources.DELIVERY),
			Attributes: resources.DeliveryAttributes{
				Date:  currentDelivery.Date,
				Price: currentDelivery.Price,
			},
			Relationships: resources.DeliveryRelationships{
				Destination: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(currentDelivery.DestinationId, 10),
						Type: resources.WAREHOUSE,
					},
				},
				Source: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(currentDelivery.SourceId, 10),
						Type: resources.WAREHOUSE,
					},
				},
			},
		}
	}
	return result
}

func getDestinationIds(deliveries []delivery.Delivery) []int64 {
	destinationIDs := make([]int64, len(deliveries))
	for i := 0; i < len(deliveries); i++ {
		destinationIDs[i] = deliveries[i].DestinationId
	}
	return destinationIDs
}

func getSourceIds(deliveries []delivery.Delivery) []int64 {
	sourceIDs := make([]int64, len(deliveries))
	for i := 0; i < len(deliveries); i++ {
		sourceIDs[i] = deliveries[i].SourceId
	}
	return sourceIDs
}

func newDeliveryIncluded(destinations []warehouse.Warehouse, sources []warehouse.Warehouse) resources.Included {
	result := resources.Included{}
	for _, item := range destinations {
		resource := newWarehouseModel(item)
		result.Add(&resource)
	}
	for _, item := range sources {
		resource := newWarehouseModel(item)
		result.Add(&resource)
	}
	return result
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
