package requests

import (
	"encoding/json"
	"net/http"
	"warehouse-service/internal/service/helpers"
	"warehouse-service/resources"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type CreateWarehouseRequest struct {
	Data resources.Warehouse
}

func NewCreateWarehouseRequest(r *http.Request) (CreateWarehouseRequest, error) {
	var request CreateWarehouseRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *CreateWarehouseRequest) validate() error {
	return helpers.MergeErrors(validation.Errors{
		"/data/attributes/cafe_id":    validation.Validate(&r.Data.Attributes.CafeId, validation.Required),
		"/data/attributes/capacity":   validation.Validate(&r.Data.Attributes.Capacity, validation.Required),
		"/data/attributes/address_id": validation.Validate(&r.Data.Relationships.Address, validation.Required),
	}).Filter()
}
