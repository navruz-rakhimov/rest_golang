package main

import (
	"context"
	"net/http"
)

func (app *application) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie("SESSTOKEN")
		if err != nil {
			err = app.writeJSON(w, http.StatusUnauthorized, nil, nil)
			return
		}

		isValid, err := app.IsAccessTokenValid(accessToken.Value)

		if !isValid || err != nil {
			err = app.writeJSON(w, http.StatusUnauthorized, nil, nil)
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		userId, err := app.GetCurrentUserId(accessToken.Value)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}

		ctx := context.WithValue(r.Context(), "user_id", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
