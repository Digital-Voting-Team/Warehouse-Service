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

func GetAddressList(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewGetAddressListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	addressesQ := helpers.AddressesQuery(r)
	applyFilters(addressesQ, request)
	currentAddress, err := addressesQ.Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get address")
		ape.Render(w, problems.InternalError())
		return
	}

	response := resources.AddressListResponse{
		Data:  newAddressesList(currentAddress),
		Links: helpers.GetOffsetLinks(r, request.OffsetPageParams),
	}
	ape.Render(w, response)
}

func applyFilters(q address.Query, request requests.GetAddressListRequest) {
	q.Page(request.OffsetPageParams)

	if len(request.FilterBuildingNumber) > 0 {
		q.FilterByBuilding(request.FilterBuildingNumber...)
	}

	if len(request.FilterStreet) > 0 {
		q.FilterByStreet(request.FilterStreet...)
	}

	if len(request.FilterCity) > 0 {
		q.FilterByCity(request.FilterCity...)
	}

	if len(request.FilterDistrict) > 0 {
		q.FilterByDistrict(request.FilterDistrict...)
	}
	if len(request.FilterRegion) > 0 {
		q.FilterByRegion(request.FilterRegion...)
	}

	if len(request.FilterPostalCode) > 0 {
		q.FilterByPostalCode(request.FilterPostalCode...)
	}

}

func newAddressesList(addresses []address.Address) []resources.Address {
	result := make([]resources.Address, len(addresses))
	for i, currentAddress := range addresses {
		result[i] = resources.Address{
			Key: resources.NewKeyInt64(currentAddress.Id, resources.ADDRESS),
			Attributes: resources.AddressAttributes{
				Building:   currentAddress.Building,
				Street:     currentAddress.Street,
				City:       currentAddress.City,
				District:   currentAddress.District,
				Region:     currentAddress.Region,
				PostalCode: currentAddress.PostalCode,
			},
		}
	}
	return result
}
