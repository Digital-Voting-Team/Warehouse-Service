package handlers

import (
	"github.com/spf13/cast"
	"net/http"
	"strconv"
	"warehouse-service/internal/pkg/delivery"
	"warehouse-service/internal/service/helpers"
	requests "warehouse-service/internal/service/requests/delivery"
	"warehouse-service/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func UpdateDelivery(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewUpdateDeliveryRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	currentDelivery, err := helpers.DeliveriesQuery(r).FilterById(request.DeliveryID).Get()
	if currentDelivery == nil {
		helpers.Log(r).WithError(err).Info("did not found delivery to update")
		ape.Render(w, problems.NotFound())
		return
	}

	//userId := r.Context().Value("userId").(int64)
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//_, _, deliveryId, err := helpers.GetIdsForGivenUser(r, userId)
	//if err != nil {
	//	helpers.Log(r).WithError(err).Info("wrong relations")
	//	ape.RenderErr(w, problems.InternalError())
	//	return
	//}
	//if *accessLevel != resources.Admin && deliveryId != delivery.ID {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	newDelivery := delivery.Delivery{
		Price:         request.Data.Attributes.Price,
		Date:          request.Data.Attributes.Date,
		SourceId:      cast.ToInt64(request.Data.Relationships.Source.Data.ID),
		DestinationId: cast.ToInt64(request.Data.Relationships.Destination.Data.ID),
	}

	relateDestination, err := helpers.WarehousesQuery(r).FilterById(newDelivery.DestinationId).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get new destination")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	relateSource, err := helpers.WarehousesQuery(r).FilterById(newDelivery.SourceId).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get new source")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	var resultDelivery delivery.Delivery
	resultDelivery, err = helpers.DeliveriesQuery(r).FilterById(currentDelivery.Id).Update(newDelivery)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to update delivery")
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
