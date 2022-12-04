package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/urlval"
)

type DeleteWarehouseRequest struct {
	WarehouseID int64 `url:"-"`
}

func NewDeleteWarehouseRequest(r *http.Request) (DeleteWarehouseRequest, error) {
	request := DeleteWarehouseRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.WarehouseID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
