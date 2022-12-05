package handlers

import (
	"github.com/Digital-Voting-Team/warehouse-service/internal/pkg/ingredient"
	"github.com/Digital-Voting-Team/warehouse-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/warehouse-service/internal/service/requests/ingredient"
	"github.com/Digital-Voting-Team/warehouse-service/resources"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetIngredientList(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewGetIngredientListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	ingredientsQ := helpers.IngredientsQuery(r)
	applyFilters(ingredientsQ, request)
	currentIngredient, err := ingredientsQ.Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get ingredient")
		ape.Render(w, problems.InternalError())
		return
	}

	response := resources.IngredientListResponse{
		Data:  newIngredientsList(currentIngredient),
		Links: helpers.GetOffsetLinks(r, request.OffsetPageParams),
	}
	ape.Render(w, response)
}

func applyFilters(q ingredient.Query, request requests.GetIngredientListRequest) {
	q.Page(request.OffsetPageParams)

	if len(request.FilterName) > 0 {
		q.FilterByName(request.FilterName...)
	}

}

func newIngredientsList(ingredientes []ingredient.Ingredient) []resources.Ingredient {
	result := make([]resources.Ingredient, len(ingredientes))
	for i, currentIngredient := range ingredientes {
		result[i] = resources.Ingredient{
			Key: resources.NewKeyInt64(currentIngredient.Id, resources.INGREDIENT),
			Attributes: resources.IngredientAttributes{
				Name: currentIngredient.Name,
			},
		}
	}
	return result
}
