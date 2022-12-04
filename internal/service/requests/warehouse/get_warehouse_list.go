package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

type GetWarehouseListRequest struct {
	pgdb.OffsetPageParams
	FilterCafeId    []int64 `filter:"cafe_id"`
	FilterAddressId []int64 `filter:"address_id"`
	FilterCapacity  []int64 `filter:"capacity"`
}

func NewGetWarehouseListRequest(r *http.Request) (GetWarehouseListRequest, error) {
	var request GetWarehouseListRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}
