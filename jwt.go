package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	jwt.StandardClaims
	AuthTime      int64  "json:\"auth_time,omitempty\""
	EmailVerified bool   "json:\"email_verified,omitempty\""
	Name          string "json:\"name,omitempty\""
	GivenName     string "json:\"given_name,omitempty\""
	FamilyName    string "json:\"family_name,omitempty\""
	Email         string "json:\"email,omitempty\""
	PiSri         string "json:\"pi.sri,omitempty\""
	SHash         string "json:\"s_hash,omitempty\""
}

type UserInfo struct {
	ID                     string
	UserID                 string
	Login                  string
	FirstName              string
	LastName               string
	Department             string
	GlobalID               string
	Email                  string
	DistinguishedName      string
	Domain                 string
	AccountType            string
	ExternalSystemType     string
	Removable              string
	Access                 bool
	Status                 string
	ThumbnailPictureBase64 string
}

func JWTTokenGeneration(basicAuth string) string {
	userInfo, err := verifyCredentials(basicAuth)
	if err != nil {
		log.Fatalln(err)
	}
	token, err := createToken(userInfo)
	if err != nil {
		log.Fatalln(err)
	}
	println("Token: ", token)

	return token
}

func verifyCredentials(basicAuth string) (*UserInfo, error) {
	// curl -X GET https://gam.intra.infineon.com/rest/ad/users -H "accept: application/json" -H "Authorization: Basic eGluZ2ppbjpKeDE1NTI4MjUwMjI3IQ=="
	url := "https://gam.intra.infineon.com/rest/ad/users"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", basicAuth)
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)

	if err == nil && res.StatusCode == 200 {
		fmt.Println("No Errors")
		info := make([]*UserInfo, 1)
		if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
			return nil, errors.New("response body of the GAM request is malformed")
		} else {
			return info[0], nil
		}
	} else {
		fmt.Println("provided basic authentication did not pass the verification")
		return nil, errors.New("provided basic authentication did not pass the verification")
	}
}

func createToken(userinfo *UserInfo) (string, error) {
	var err error
	// TODO: get signing keys from Openshift secrete
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file

	//atClaims := jwt.MapClaims{}
	//atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	atClaims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   userinfo.Email,
			Audience:  "miam",
			Id:        "test-id",
			Issuer:    "Managed-Identity-And-Access-Management-For-Infineon",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		},
		AuthTime:      time.Now().Unix(),
		EmailVerified: true,
		Name:          userinfo.Login,
		GivenName:     userinfo.FirstName,
		FamilyName:    userinfo.LastName,
		Email:         userinfo.Email,
		PiSri:         "un8EQ2xr1HCTQVInLk0DL7z4Btg..Al2g",
		SHash:         "6eQETHDIvtTEss0dNZFWWw",
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
