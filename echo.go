package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	URLIAM, _      = url.Parse("http://localhost:8080")
	URLUPSTREAM, _ = url.Parse("http://localhost:4200")
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", sidecar)

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
func sidecar(c echo.Context) error {
	var token string
	var authHeader string

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
			// TODO: GAM endpoint to get claims
			// TODO: cache to prevent from calling GAM all the time store it in memory
			// go-cache
			token = JWT(authHeader)
			bearerToken := "Bearer " + token
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
						res.Header().Add("Authorization", bearerToken)
						// Attach header with the redirect
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
		//c.Redirect(302, "http://localhost:8080")
		req.Host = URLIAM.Host
		req.URL.Host = URLIAM.Host
		req.URL.Scheme = URLIAM.Scheme
		// Attach header with the redirect
		proxyIAM.ServeHTTP(res, req)
	}

	return c.String(http.StatusOK, token)
}
