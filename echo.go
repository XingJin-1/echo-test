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

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Pre(middleware.NonWWWRedirect())
	//e.Pre(middleware.WWWRedirect())
	//e.Pre(middleware.NonWWWRedirect())

	// Setup proxy
	/*
		url1, err := url.Parse("http://localhost:8080")
		if err != nil {
			e.Logger.Fatal(err)
		}
		url2, err := url.Parse("http://localhost:4200")
		if err != nil {
			e.Logger.Fatal(err)
		}
		targets := []*middleware.ProxyTarget{
			{
				URL: url1,
			},
			{
				URL: url2,
			},
		}
		e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
	*/

	// Routes
	e.GET("/", sidecar)

	// Start server
	e.Logger.Fatal(e.Start(":8082"))
}

// Handler
func sidecar(c echo.Context) error {
	var token string
	var authHeader string

	//url1, _ := url.Parse("http://localhost:8080")
	url2, _ := url.Parse("http://localhost:4200")
	//proxy1 := httputil.NewSingleHostReverseProxy(url1)
	proxy2 := httputil.NewSingleHostReverseProxy(url2)
	req := c.Request()
	res := c.Response().Writer

	if val, ok := c.Request().Header["Authorization"]; ok {
		fmt.Printf("Authozation header provided: \n")
		authHeader = val[0]
		if strings.Contains(authHeader, "Basic") {
			// if basic auth
			fmt.Printf("Basic Authorization header provided! A JWT token will be created by the basic auth sidecar\n")
			// TODO: GAM endpoint to get claims
			// TODO: cache to prevent from calling GAM all the time
			token = JWT()
			// TODO: Attach header with the redirect
			//c.Redirect(302, "http://localhost:4200")
			req.Host = url2.Host
			req.URL.Host = url2.Host
			req.URL.Scheme = url2.Scheme
			proxy2.ServeHTTP(res, req)
		} else {
			//JWT token is present
			fmt.Printf("JWT Authorization header provided! \n")
			//Decode the JWT token
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
						// 		 attach header with the redirect
						//c.Redirect(302, "http://localhost:4200")
						req.Host = url2.Host
						req.URL.Host = url2.Host
						req.URL.Scheme = url2.Scheme
						proxy2.ServeHTTP(res, req)
					} else {
						fmt.Println("Dealing with JWT token issued by the iam sidecar! The issuer is: ", val, "The request will be redirected to the iam side car.")
						// TODO: append JWT to the header of the request
						c.Redirect(302, "http://localhost:8080")
						/*
							req.Host = url1.Host
							req.URL.Host = url1.Host
							req.URL.Scheme = url1.Scheme
							proxy1.ServeHTTP(res, req)
						*/
					}
				} else {
					fmt.Println("No Issuer provided!")
				}
			}
		}
	} else {
		//forward to the iam side car
		fmt.Printf("No Authorization header provided! The request will be redirected to the iam side car.\n")
		c.Redirect(302, "http://localhost:8080")
	}

	return c.String(http.StatusOK, token)
}
