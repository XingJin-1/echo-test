package main

import (
	"encoding/json"
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

func JWT(basicAuth string) string {

	userInfo, err := VerifyCredentials(basicAuth)

	token, err := CreateToken(userInfo)

	if err == nil {
		println("Token: ", token)
	} else {
		println("Error while generating tokens")
	}

	return token
}

func VerifyCredentials(basicAuth string) (*UserInfo, error) {
	// TODO: somehow get the username or email
	url := "https://gam.intra.infineon.com/rest/ad/users?username=XingJin"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", basicAuth)
	res, _ := client.Do(req)

	info := make([]*UserInfo, 1)
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return nil, err
	}
	return info[0], nil
}

func CreateToken(userinfo *UserInfo) (string, error) {
	var err error
	// TODO: get signing keys 	Openshift secrete
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
