package handlers

import (
	"net/http"
	"warehouse-service/internal/service/helpers"
	requests "warehouse-service/internal/service/requests/used_ingredient"
	"warehouse-service/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetUsedIngredient(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewGetUsedIngredientRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	usedIngredient, err := helpers.UsedIngredientsQuery(r).FilterById(request.UsedIngredientID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get used ingredient from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	if usedIngredient == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	result := resources.UsedIngredientResponse{
		Data: resources.UsedIngredient{
			Key: resources.NewKeyInt64(usedIngredient.Id, resources.USED_INGREDIENT),
			Attributes: resources.UsedIngredientAttributes{
				DeletionDate: usedIngredient.DeletionDate,
				Name:         usedIngredient.Name,
				Origin:       usedIngredient.Origin,
				Price:        usedIngredient.Price,
				Quantity:     usedIngredient.Quantity,
				Reason:       usedIngredient.Reason,
			},
		},
	}

	ape.Render(w, result)
}
