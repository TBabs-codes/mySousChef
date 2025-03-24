package main

import (
	"net/http"

	"github.com/TBabs-codes/mySousChef/auth"
	"github.com/TBabs-codes/mySousChef/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwt, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Error finding JSON web token")
			return
		}

		userID, err := auth.ValidateJWT(jwt, cfg.JWTSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Error validating JSON web token")
			return
		}

		user, err := cfg.DB.ReturnUserFromID(r.Context(), userID.String())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error retrieving User from ID")
			return
		}

		handler(w, r, user)
	}
}