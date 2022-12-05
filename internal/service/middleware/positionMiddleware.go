package middleware

import (
	staffRes "github.com/Digital-Voting-Team/staff-service/resources"
	"github.com/Digital-Voting-Team/warehouse-service/internal/service/helpers"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"net/http"
)

func CheckManagerPosition() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessLevel := r.Context().Value("accessLevel").(staffRes.AccessLevel)
			if accessLevel < staffRes.Manager {
				helpers.Log(r).Info("insufficient user permissions")
				ape.RenderErr(w, problems.Forbidden())
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
