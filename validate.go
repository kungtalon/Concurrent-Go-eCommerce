package main

import (
	"errors"
	"fmt"
	"jzmall/common"
	"jzmall/distributed"
	"net/http"
)

var hostArray = []string{"127.0.0.1", "127.0.0.1"}

var localHost = "127.0.0.1"

var port = "8081"

var hashConsistent *distributed.Consistent

var accessControl = &distributed.AccessControl{SourcesArray: make(map[int]interface{})}

func Check(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("Done check!")
}

// Auth unified authentication filter, every uri registered with this handle
// needs to be validated before any other handle
func Auth(rw http.ResponseWriter, r *http.Request) error {
	// add cookie based authentication
	err := ValidateCookie(r)
	if err != nil {
		return err
	}
	fmt.Println("Done Auth!")
	return nil
}

func ValidateCookie(r *http.Request) error {
	uid, err1 := r.Cookie("uid")
	uidStr, err2 := r.Cookie("sign")
	if err1 != nil || err2 != nil {
		return errors.New("Failed to get userid from cookie")
	}
	decoded, err := common.DePwdCode(uidStr.Value)
	if err != nil {
		return errors.New("Error when decoding encoded userid...")
	}
	if string(decoded) != uid.Value {
		return errors.New("Invalid user information found in cookies! Logged out...")
	}
	return nil
}

func main() {
	hashConsistent = distributed.NewConsistent()
	accessControl.SetHosts(localHost, port)
	accessControl.SetConsistentHash(hashConsistent)
	// add node with consistent hash algo
	for _, v := range hostArray {
		// add hostIP to hash ring
		hashConsistent.Add(v)
	}
	filter := common.NewFiler()
	filter.RegisterFilterUri("/check", Auth)
	http.HandleFunc("/check", filter.Handle(Check))
	http.ListenAndServe("localhost:8083", nil)
}
