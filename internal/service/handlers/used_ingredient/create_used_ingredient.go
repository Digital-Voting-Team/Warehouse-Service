package handlers

import (
	"net/http"
	usedIngredient "warehouse-service/internal/pkg/used_ingredient"
	"warehouse-service/internal/service/helpers"
	requests "warehouse-service/internal/service/requests/used_ingredient"
	"warehouse-service/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func CreateUsedIngredient(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewCreateUsedIngredientRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var resultUsedIngredient usedIngredient.UsedIngredient

	currentUsedIngredient := usedIngredient.UsedIngredient{
		Name:         request.Data.Attributes.Name,
		Quantity:     request.Data.Attributes.Quantity,
		Origin:       request.Data.Attributes.Origin,
		Price:        request.Data.Attributes.Price,
		DeletionDate: request.Data.Attributes.DeletionDate,
		Reason:       request.Data.Attributes.Reason,
	}

	resultUsedIngredient, err = helpers.UsedIngredientsQuery(r).Insert(currentUsedIngredient)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create used ingredient")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.UsedIngredientResponse{
		Data: resources.UsedIngredient{
			Key: resources.NewKeyInt64(resultUsedIngredient.Id, resources.USED_INGREDIENT),
			Attributes: resources.UsedIngredientAttributes{
				DeletionDate: resultUsedIngredient.DeletionDate,
				Name:         resultUsedIngredient.Name,
				Origin:       resultUsedIngredient.Origin,
				Price:        resultUsedIngredient.Price,
				Quantity:     resultUsedIngredient.Quantity,
				Reason:       resultUsedIngredient.Reason,
			},
		},
	}
	ape.Render(w, result)
}
