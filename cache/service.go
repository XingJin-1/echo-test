package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

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

func main() {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c := cache.New(5*time.Minute, 10*time.Minute)

	keyString := "key"
	valueString := "value"
	// Set the value of the key "foo" to "bar", with the default expiration time
	c.Add(keyString, valueString, cache.DefaultExpiration)

	// Set the value of the key "baz" to 42, with no expiration time
	// (the item won't be removed until it is re-set, or removed using
	// c.Delete("baz")
	c.Set("baz", 42, cache.NoExpiration)

	// Get the string associated with the key "foo" from the cache

	if val, time, found := c.GetWithExpiration(keyString); found {
		fmt.Println("print from Main function Value: ", val)
		fmt.Println("print from Main function Time: ", time)
	} else {
		c.Add(keyString, valueString, cache.DefaultExpiration)
	}

	basicAuth := "Basic eGluZ2ppbjpKeDE1NTI4MjUwMjI3IQ=="
	url := "https://gam.intra.infineon.com/rest/ad/users"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", basicAuth)
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)

	if err == nil && res.StatusCode == 200 {
		fmt.Println("No Errors")
		info := make([]*UserInfo, 1)

		if err := json.NewDecoder(res.Body).Decode(&info); err == nil {
			fmt.Println("print from Main function Info: ", info[0])
		}

	} else {
		fmt.Println("Error occured ")
	}

	/*
		// Want performance? Store pointers!
		c.Set("foo", &MyStruct, cache.DefaultExpiration)
		if x, found := c.Get("foo"); found {
			foo := x.(*MyStruct)
			// ...
		}
	*/

}
