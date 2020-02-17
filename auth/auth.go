package auth

import (
	"log"
	"net/http"
)

// AuthenticationMiddleware Lookup middleware for checking adminToken in HTTP Header is valid
type AuthenticationMiddleware struct {
	tokenUsers map[string]string
}

// Populate initialise the auth muddleware and populate with the valid admin tokens
func (amw *AuthenticationMiddleware) Populate() {
	if amw.tokenUsers == nil {
        amw.tokenUsers = make(map[string]string)
    }
	amw.tokenUsers["example"] = "adminToken"
}

// Middleware Function to apply to router the check all HTTP requests for admin tokens header
func (amw *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("adminToken")
		if user, found := amw.tokenUsers[token]; found {
			log.Printf("Successfully Authenticated user %s\n", user)
			next.ServeHTTP(w, r)
		} else {
			log.Printf("Token %s is not a valid adminToken", token)
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func main() {
	log.Println("Run Auth")
}