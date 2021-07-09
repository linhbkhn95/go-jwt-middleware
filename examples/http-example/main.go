package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/auth0/go-jwt-middleware/validate/josev2"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(jwtmiddleware.ContextKey{})
	j, err := json.MarshalIndent(user, "", "\t")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}

	fmt.Fprintf(w, "This is an authenticated request")
	fmt.Fprintf(w, "Claim content:\n")
	fmt.Fprint(w, string(j))
})

func main() {
	// uncomment the below to use the caching key provider
	// u, err := url.Parse("https://<your-domain>")
	// if err != nil {
	// 	// we'll panic in order to fail fast
	// 	panic(err)
	// }

	// p := josev2.NewCachingJWKSProvider(*u, 5*time.Minute)

	keyFunc := func(ctx context.Context) (interface{}, error) {
		// our token must be signed using this data
		return []byte("secret"), nil
	}

	expectedClaimsFunc := func() jwt.Expected {
		// By setting up expected claims we are saying a token must
		// have the data we specify.
		return jwt.Expected{
			Issuer: "josev2-example",
		}
	}

	// setup the piece which will validate tokens
	validator, err := josev2.New(
		// p.KeyFunc, // uncomment this to use the caching key provider
		keyFunc,
		jose.HS256,
		josev2.WithExpectedClaims(expectedClaimsFunc),
	)
	if err != nil {
		// we'll panic in order to fail fast
		panic(err)
	}

	// setup the middleware
	m := jwtmiddleware.New(validator.ValidateToken)

	http.ListenAndServe("0.0.0.0:3000", m.CheckJWT(handler))
}
