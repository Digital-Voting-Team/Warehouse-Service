package handlers

import (
	"net/http"
	"warehouse-service/internal/service/helpers"
	requests "warehouse-service/internal/service/requests/warehouse"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func DeleteWarehouse(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewDeleteWarehouseRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	warehouse, err := helpers.WarehousesQuery(r).FilterById(request.WarehouseID).Get()
	if warehouse == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	err = helpers.WarehousesQuery(r).Delete(request.WarehouseID)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to delete warehouse")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
