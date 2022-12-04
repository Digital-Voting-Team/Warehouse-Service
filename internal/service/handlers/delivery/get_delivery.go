package handlers

import (
	"net/http"
	"strconv"
	"warehouse-service/internal/service/helpers"
	requests "warehouse-service/internal/service/requests/delivery"
	"warehouse-service/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetDelivery(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewGetDeliveryRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	delivery, err := helpers.DeliveriesQuery(r).FilterById(request.DeliveryID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get delivery from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	if delivery == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	relateDestination, err := helpers.WarehousesQuery(r).FilterById(delivery.DestinationId).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get destination")
		ape.RenderErr(w, problems.NotFound())
		return
	}
	relateSource, err := helpers.WarehousesQuery(r).FilterById(delivery.SourceId).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get source")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	var includes resources.Included
	includes.Add(&resources.Warehouse{
		Key: resources.NewKeyInt64(relateDestination.Id, resources.WAREHOUSE),
		Attributes: resources.WarehouseAttributes{
			CafeId:   relateDestination.CafeId,
			Capacity: relateDestination.Capacity,
		},
		Relationships: resources.WarehouseRelationships{
			Address: resources.Relation{
				Data: &resources.Key{
					ID:   strconv.FormatInt(relateDestination.AddressId, 10),
					Type: resources.ADDRESS,
				},
			},
		},
	})
	includes.Add(&resources.Warehouse{
		Key: resources.NewKeyInt64(relateSource.Id, resources.WAREHOUSE),
		Attributes: resources.WarehouseAttributes{
			CafeId:   relateSource.CafeId,
			Capacity: relateSource.Capacity,
		},
		Relationships: resources.WarehouseRelationships{
			Address: resources.Relation{
				Data: &resources.Key{
					ID:   strconv.FormatInt(relateSource.AddressId, 10),
					Type: resources.ADDRESS,
				},
			},
		},
	})

	result := resources.DeliveryResponse{
		Data: resources.Delivery{
			Key: resources.NewKeyInt64(delivery.Id, resources.DELIVERY),
			Attributes: resources.DeliveryAttributes{
				Price: delivery.Price,
				Date:  delivery.Date,
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
		},
		Included: includes,
	}
	ape.Render(w, result)
}
