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

func GetUsedIngredientList(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewGetUsedIngredientListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	usedIngredientsQ := helpers.UsedIngredientsQuery(r)
	applyFilters(usedIngredientsQ, request)
	currentUsedIngredient, err := usedIngredientsQ.Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get used ingredient")
		ape.Render(w, problems.InternalError())
		return
	}

	response := resources.UsedIngredientListResponse{
		Data:  newUsedIngredientsList(currentUsedIngredient),
		Links: helpers.GetOffsetLinks(r, request.OffsetPageParams),
	}
	ape.Render(w, response)
}

func applyFilters(q usedIngredient.Query, request requests.GetUsedIngredientListRequest) {
	q.Page(request.OffsetPageParams)

	if len(request.FilterName) > 0 {
		q.FilterByName(request.FilterName...)
	}

	if len(request.FilterQuantity) > 0 {
		q.FilterByQuantity(request.FilterQuantity...)
	}

	if len(request.FilterOrigin) > 0 {
		q.FilterByOrigin(request.FilterOrigin...)
	}

	if len(request.FilterPrice) > 0 {
		q.FilterByPrice(request.FilterPrice...)
	}
	if len(request.FilterDeletionDate) > 0 {
		q.FilterByDeletionDate(request.FilterDeletionDate...)
	}

	if len(request.FilterReason) > 0 {
		q.FilterByReason(request.FilterReason...)
	}
}

func newUsedIngredientsList(usedIngredients []usedIngredient.UsedIngredient) []resources.UsedIngredient {
	result := make([]resources.UsedIngredient, len(usedIngredients))
	for i, currentUsedIngredient := range usedIngredients {
		result[i] = resources.UsedIngredient{
			Key: resources.NewKeyInt64(currentUsedIngredient.Id, resources.USED_INGREDIENT),
			Attributes: resources.UsedIngredientAttributes{
				DeletionDate: currentUsedIngredient.DeletionDate,
				Name:         currentUsedIngredient.Name,
				Origin:       currentUsedIngredient.Origin,
				Price:        currentUsedIngredient.Price,
				Quantity:     currentUsedIngredient.Quantity,
				Reason:       currentUsedIngredient.Reason,
			},
		}
	}
	return result
}
