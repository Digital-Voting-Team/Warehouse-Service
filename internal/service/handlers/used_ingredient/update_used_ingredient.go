package handlers

import (
	usedIngredient "github.com/Digital-Voting-Team/warehouse-service/internal/pkg/used_ingredient"
	"github.com/Digital-Voting-Team/warehouse-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/warehouse-service/internal/service/requests/used_ingredient"
	"github.com/Digital-Voting-Team/warehouse-service/resources"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func UpdateUsedIngredient(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewUpdateUsedIngredientRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	currentUsedIngredient, err := helpers.UsedIngredientsQuery(r).FilterById(request.UsedIngredientID).Get()
	if currentUsedIngredient == nil {
		helpers.Log(r).WithError(err).Info("did not found usedIngredient to update")
		ape.Render(w, problems.NotFound())
		return
	}

	//userId := r.Context().Value("userId").(int64)
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//_, _, usedIngredientId, err := helpers.GetIdsForGivenUser(r, userId)
	//if err != nil {
	//	helpers.Log(r).WithError(err).Info("wrong relations")
	//	ape.RenderErr(w, problems.InternalError())
	//	return
	//}
	//if *accessLevel != resources.Admin && usedIngredientId != usedIngredient.ID {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	newUsedIngredient := usedIngredient.UsedIngredient{
		Name:         request.Data.Attributes.Name,
		Quantity:     request.Data.Attributes.Quantity,
		Origin:       request.Data.Attributes.Origin,
		Price:        request.Data.Attributes.Price,
		DeletionDate: request.Data.Attributes.DeletionDate,
		Reason:       request.Data.Attributes.Reason,
	}

	var resultUsedIngredient usedIngredient.UsedIngredient
	resultUsedIngredient, err = helpers.UsedIngredientsQuery(r).FilterById(currentUsedIngredient.Id).Update(newUsedIngredient)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to update usedIngredient")
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
