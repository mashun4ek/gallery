package middleware

import (
	"net/http"
	"strings"

	"github.com/mashun4ek/webdevcalhoun/gallery/context"
	"github.com/mashun4ek/webdevcalhoun/gallery/models"
)

type User struct {
	models.UserService
}

func (mw *User) RequireUserMiddleware(next http.Handler) http.HandlerFunc {
	return mw.RUMFn(next.ServeHTTP)
}

func (mw *User) RUMFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// if static asset or image we will not need to look up the current user
		if strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/images/") {
			next(w, r)
			return
		}
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next(w, r)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next(w, r)
	})

}

// RequireUser assumes that User has already been run, otherwise it won't work correctly
type RequireUser struct {
	User
}

// RequireUserMiddleware assumes that User has already been run, otherwise it won't work correctly
func (mw *RequireUser) RequireUserMiddleware(next http.Handler) http.HandlerFunc {
	return mw.RUMFn(next.ServeHTTP)
}

// RUMFn assumes that User has already been run, otherwise it won't work correctly
func (mw *RequireUser) RUMFn(next http.HandlerFunc) http.HandlerFunc {
	return mw.User.RUMFn(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	})
}
