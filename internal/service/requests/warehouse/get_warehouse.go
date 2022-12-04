package requests

import (
	"github.com/go-chi/chi"
	"github.com/spf13/cast"

	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

type GetWarehouseRequest struct {
	WarehouseID int64 `url:"-"`
}

func NewGetWarehouseRequest(r *http.Request) (GetWarehouseRequest, error) {
	request := GetWarehouseRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.WarehouseID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
