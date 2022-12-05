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

func CreateAddress(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewCreateAddressRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var resultAddress address.Address

	currentAddress := address.Address{
		Building:   request.Data.Attributes.Building,
		Street:     request.Data.Attributes.Street,
		City:       request.Data.Attributes.City,
		District:   request.Data.Attributes.District,
		Region:     request.Data.Attributes.Region,
		PostalCode: request.Data.Attributes.PostalCode,
	}

	resultAddress, err = helpers.AddressesQuery(r).Insert(currentAddress)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create address")
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
