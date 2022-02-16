package main

import (
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

func JWT() string {

	userID := "XingJin"

	token, err := CreateToken(userID)

	if err == nil {
		println("Token: ", token)
	} else {
		println("Error while generating tokens")
	}
	return token
}

func CreateToken(userid string) (string, error) {
	var err error
	// TODO: get signing keys 	Openshift secrete
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file

	//atClaims := jwt.MapClaims{}
	//atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	atClaims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   "Jin.Xing@infineon.com",
			Audience:  "rddl",
			Id:        "test-id",
			Issuer:    "Managed-Identity-And-Access-Management-For-Infineon",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
		AuthTime:      time.Now().Unix(),
		EmailVerified: true,
		Name:          "Xing Jin (IFAG IT DSA RD PE)",
		GivenName:     "Xing",
		FamilyName:    "Jin",
		Email:         "Jin.Xing@infineon.com",
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
