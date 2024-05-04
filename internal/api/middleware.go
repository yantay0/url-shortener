package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"

	"github.com/yantay0/url-shortener/internal/model"
	"github.com/yantay0/url-shortener/internal/storage"
	"github.com/yantay0/url-shortener/internal/validator"
)

func (app *App) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, model.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		// Split the authorization header into two parts. The first part should be
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		token := headerParts[1]
		v := validator.New()
		if model.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.Storage.Users.GetForToken(storage.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, storage.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}
		// Call the contextSetUser() helper to add the user information to the request
		// context.
		r = app.contextSetUser(r, user)
		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
func (app *App) rateLimit(next http.Handler) http.Handler {
	// Define a client struct to hold the rate limiter and last seen time for reach client
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	// Declare a mutex and a map to hold pointers to a client struct.
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Launch a background goroutine which removes old entries from the clients map once every
	// minute.
	go func() {
		for {
			time.Sleep(time.Minute)

			// Lock the mutex to prevent any rate limiter checks from happening while the cleanup
			// is taking place.
			mu.Lock()

			// Loop through all clients. if they haven't been seen within the last three minutes,
			// then delete the corresponding entry from the clients map.
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			// Importantly, unlock the mutex when the cleanup is complete.
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only carry out the check if rate limited is enabled.
		if app.Config.Limiter.Enabled {
			// Use the realip.FromRequest function to get the client's real IP address.
			ip := realip.FromRequest(r)

			// Lock the mutex to prevent this code from being executed concurrently.
			mu.Lock()

			// Check to see if the IP address already exists in the map. If it doesn't,
			// then initialize a new rate limiter and add the IP address and limiter to the map.
			if _, found := clients[ip]; !found {
				// Use the requests-per-second and burst values from the app.config struct.
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.Config.Limiter.RPS), app.Config.Limiter.Burst)}
			}

			// Update the last seen time for the client.
			clients[ip].lastSeen = time.Now()

			// Call the limiter.Allow() method on the rate limiter for the current IP address.
			// If the request isn't allowed, unlock the mutex and send a 429 Too Many Requests
			// response.
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			// Very importantly, unlock the mutex before calling the next handler in the chain.
			// Notice that we DON'T use defer to unlock the mutex, as that would mean that the mutex
			// isn't unlocked until all handlers downstream of this middleware have also returned.
			mu.Unlock()
		}
		next.ServeHTTP(w, r)
	})
}

func (app *App) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
