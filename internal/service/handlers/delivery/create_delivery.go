package handlers

import (
	model "github.com/Digital-Voting-Team/warehouse-service/internal/pkg/delivery"
	"github.com/Digital-Voting-Team/warehouse-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/warehouse-service/internal/service/requests/delivery"
	"github.com/Digital-Voting-Team/warehouse-service/resources"
	"github.com/spf13/cast"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func CreateDelivery(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewCreateDeliveryRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	delivery := model.Delivery{
		Price:         request.Data.Attributes.Price,
		Date:          request.Data.Attributes.Date,
		SourceId:      cast.ToInt64(request.Data.Relationships.Source.Data.ID),
		DestinationId: cast.ToInt64(request.Data.Relationships.Destination.Data.ID),
	}

	var resultDelivery model.Delivery
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

	resultDelivery, err = helpers.DeliveriesQuery(r).Insert(delivery)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create delivery")
		ape.RenderErr(w, problems.InternalError())
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
			Key: resources.NewKeyInt64(resultDelivery.Id, resources.DELIVERY),
			Attributes: resources.DeliveryAttributes{
				Price: resultDelivery.Price,
				Date:  resultDelivery.Date,
			},
			Relationships: resources.DeliveryRelationships{
				Destination: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(resultDelivery.DestinationId, 10),
						Type: resources.WAREHOUSE,
					},
				},
				Source: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(resultDelivery.SourceId, 10),
						Type: resources.WAREHOUSE,
					},
				},
			},
		},
		Included: includes,
	}
	ape.Render(w, result)
}
