package handlers

import (
	"github.com/Digital-Voting-Team/warehouse-service/internal/pkg/address"
	"github.com/Digital-Voting-Team/warehouse-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/warehouse-service/internal/service/requests/address"
	"github.com/Digital-Voting-Team/warehouse-service/resources"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func UpdateAddress(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewUpdateAddressRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	currentAddress, err := helpers.AddressesQuery(r).FilterById(request.AddressID).Get()
	if currentAddress == nil {
		helpers.Log(r).WithError(err).Info("did bot found address to update")
		ape.Render(w, problems.NotFound())
		return
	}

	//userId := r.Context().Value("userId").(int64)
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//_, _, addressId, err := helpers.GetIdsForGivenUser(r, userId)
	//if err != nil {
	//	helpers.Log(r).WithError(err).Info("wrong relations")
	//	ape.RenderErr(w, problems.InternalError())
	//	return
	//}
	//if *accessLevel != resources.Admin && addressId != address.ID {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	newAddress := address.Address{
		Building:   request.Data.Attributes.Building,
		Street:     request.Data.Attributes.Street,
		City:       request.Data.Attributes.City,
		District:   request.Data.Attributes.District,
		Region:     request.Data.Attributes.Region,
		PostalCode: request.Data.Attributes.PostalCode,
	}

	var resultAddress address.Address
	resultAddress, err = helpers.AddressesQuery(r).FilterById(currentAddress.Id).Update(newAddress)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to update address")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.AddressResponse{
		Data: resources.Address{
			Key: resources.NewKeyInt64(resultAddress.Id, resources.ADDRESS),
			Attributes: resources.AddressAttributes{
				Building:   resultAddress.Building,
				Street:     resultAddress.Street,
				City:       resultAddress.City,
				District:   resultAddress.District,
				Region:     resultAddress.Region,
				PostalCode: resultAddress.PostalCode,
			},
		},
	}
	ape.Render(w, result)
}
