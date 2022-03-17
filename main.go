package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/patrickmn/go-cache"
)

var (
	URLIAM, _      = url.Parse("http://localhost:8080")
	URLUPSTREAM, _ = url.Parse("http://localhost:4200")
)

var (
	MemCache = cache.New(5*time.Minute, 10*time.Minute)
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", handler)

	targetsOauth2 := []*middleware.ProxyTarget{
		{
			URL: URLIAM,
		},
	}
	groupOauth2 := e.Group("/oauth2")
	groupOauth2.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targetsOauth2)))

	// Start server
	e.Logger.Fatal(e.Start(":8082"))
}

// Handler.
func handler(c echo.Context) error {
	var token string
	var authHeader string
	var bearerToken string

	proxyIAM := httputil.NewSingleHostReverseProxy(URLIAM)
	proxyUpstream := httputil.NewSingleHostReverseProxy(URLUPSTREAM)

	req := c.Request()
	res := c.Response().Writer

	if val, ok := c.Request().Header["Authorization"]; ok {
		// use logging library
		fmt.Printf("Authozation header provided: \n")
		authHeader = val[0]
		if strings.Contains(authHeader, "Basic") {
			// If basic auth
			fmt.Printf("Basic Authorization header provided! A JWT token will be created by the basic auth sidecar\n")
			// Check whether the corresponding token is already in the cache
			if val, time, found := MemCache.GetWithExpiration(authHeader); found {
				fmt.Println("The basic auth is already in the cache", val, "the token will be expired at: ", time)
				bearerToken = val.(string)
			} else {
				// verify the basic auth and generate JWT token
				token = JWTTokenGeneration(authHeader)
				bearerToken = "Bearer " + token
				// Add the token to the cache
				err := MemCache.Add(authHeader, bearerToken, cache.DefaultExpiration)
				if err != nil {
					log.Fatalf("ERROR: %v", err)
				}
			}
			req.Host = URLUPSTREAM.Host
			req.URL.Host = URLUPSTREAM.Host
			req.URL.Scheme = URLUPSTREAM.Scheme
			// Attach header with the redirect
			res.Header().Add("Authorization", bearerToken)
			proxyUpstream.ServeHTTP(res, req)

		} else {
			// JWT token is present
			fmt.Printf("JWT Authorization header provided! \n")
			// Decode the JWT token
			bearerToken := authHeader
			authHeader = strings.Replace(authHeader, "Bearer ", "", 1)
			claims := jwt.MapClaims{}
			_, err := jwt.ParseWithClaims(authHeader, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(""), nil
			})
			if err != nil {
				if val, ok := claims["iss"]; ok {
					if claims.VerifyIssuer("Managed-Identity-And-Access-Management-For-Infineon", false) {
						// The JWT token is issued by the basic auth sidecar
						fmt.Println("Dealing with JWT token issued by the basic auth sidecar! the issuer is: ", val)
						// TODO: verify the JWT token
						req.Host = URLUPSTREAM.Host
						req.URL.Host = URLUPSTREAM.Host
						req.URL.Scheme = URLUPSTREAM.Scheme
						// Attach header with the redirect
						res.Header().Add("Authorization", bearerToken)
						proxyUpstream.ServeHTTP(res, req)
					} else {
						fmt.Println("Dealing with JWT token issued by the iam sidecar! The issuer is: ", val, "The request will be redirected to the iam side car.")
						req.Host = URLIAM.Host
						req.URL.Host = URLIAM.Host
						req.URL.Scheme = URLIAM.Scheme
						// Attach header with the redirect
						res.Header().Add("Authorization", bearerToken)
						proxyIAM.ServeHTTP(res, req)
						// forward the url with oauth2
					}
				} else {
					fmt.Println("No Issuer provided!")
				}
			}
		}
	} else {
		// forward to the iam side car
		fmt.Printf("No Authorization header provided! The request will be redirected to the iam side car.\n")
		req.Host = URLIAM.Host
		req.URL.Host = URLIAM.Host
		req.URL.Scheme = URLIAM.Scheme
		proxyIAM.ServeHTTP(res, req)
	}

	return c.String(http.StatusOK, token)
}
