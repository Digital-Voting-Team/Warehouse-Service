package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"warehouse-service/internal/service/helpers"
	"warehouse-service/resources"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/urlval"
)

type UpdateWarehouseRequest struct {
	WarehouseID int64 `url:"-" json:"-"`
	Data        resources.Warehouse
}

func NewUpdateWarehouseRequest(r *http.Request) (UpdateWarehouseRequest, error) {
	request := UpdateWarehouseRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.WarehouseID = cast.ToInt64(chi.URLParam(r, "id"))

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *UpdateWarehouseRequest) validate() error {
	return helpers.MergeErrors(validation.Errors{
		"/data/attributes/cafe_id":    validation.Validate(&r.Data.Attributes.CafeId, validation.Required),
		"/data/attributes/capacity":   validation.Validate(&r.Data.Attributes.Capacity, validation.Required),
		"/data/attributes/address_id": validation.Validate(&r.Data.Relationships.Address, validation.Required),
	}).Filter()
}
