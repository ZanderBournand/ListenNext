package middlewares

import (
	"context"
	"main/tools"
	"net/http"
)

type authString string
type loginStypeString string

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth == "" {
			next.ServeHTTP(w, r)
			return
		}

		bearer := "Bearer "
		auth = auth[len(bearer):]

		validate, err := tools.JwtValidate(context.Background(), auth)
		if err != nil || !validate.Valid {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		customClaim, _ := validate.Claims.(*tools.JwtCustomClaim)

		ctx := context.WithValue(r.Context(), authString("auth"), customClaim)
		ctx = context.WithValue(ctx, authString("userID"), customClaim.ID)
		ctx = context.WithValue(ctx, loginStypeString("loginType"), customClaim.LoginType)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func CtxValue(ctx context.Context) *tools.JwtCustomClaim {
	raw, _ := ctx.Value(authString("auth")).(*tools.JwtCustomClaim)
	return raw
}

func CtxUserID(ctx context.Context) string {
	raw, _ := ctx.Value(authString("userID")).(string)
	return raw
}

func CtxLoginType(ctx context.Context) string {
	raw, _ := ctx.Value(loginStypeString("loginType")).(string)
	return raw
}
