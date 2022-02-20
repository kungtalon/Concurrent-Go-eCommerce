package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"jzmall/common"
	"jzmall/datamodels"
	"jzmall/distributed"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var hostArray = []string{"127.0.0.1", "172.23.240.1"}

var localHost = ""

// internal ip of productNum controll service
var GetOneIp = "127.0.0.1"

var GetOnePort = "8084"

var port = "8083"

var hashConsistent *distributed.Consistent

var rabbitmqValidate *distributed.RabbitMQ

var blackList = &distributed.BlackList{Lookup: make(map[int]bool)}

var accessControl = &distributed.AccessControl{SourcesArray: make(map[int]time.Time), BlackList: blackList}

func CheckUserRight(rw http.ResponseWriter, r *http.Request) {
	right := accessControl.GetDistributedRight(r)
	if !right {
		rw.Write([]byte("false"))
		return
	}
	rw.Write([]byte("true"))
	return
}

// Check checks whether the product has sold out
func Check(rw http.ResponseWriter, r *http.Request) {
	var (
		failure = []byte("false")
		success = []byte("true")
	)

	log.Println("Doing check...")
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["productID"]) <= 0 || len(queryForm["productID"][0]) <= 0 {
		log.Println("Wrong query!")
		rw.Write(failure)
		return
	}
	productIdStr := queryForm["productID"][0]
	userCookie, err := r.Cookie("uid")
	if err != nil {
		log.Println("Product id conversion error...")
		rw.Write(failure)
		return
	}

	// 1. Distributed Authentication
	right := accessControl.GetDistributedRight(r)
	if right == false {
		log.Println("User check failed...")
		rw.Write(failure)
		return
	}
	// 2. Get the control of product number, avoid oversale
	hostUrl := "http://" + GetOneIp + ":" + GetOnePort + "/getOne"
	responseValidate, validateBody, err := accessControl.GetUrl(hostUrl, r)
	if err != nil {
		log.Println("Failed to get response from access control")
		rw.Write(failure)
		return
	}
	// check the count control service statuscode
	if responseValidate.StatusCode != 200 || string(validateBody) != "true" {
		log.Println("Response is not expected...")
		rw.Write(failure)
		return
	}

	// place order now!
	productId, err := strconv.ParseUint(productIdStr, 10, 64)
	if err != nil {
		log.Println("fail to convert product id: " + productIdStr)
		rw.Write(failure)
		return
	}
	userId, err := strconv.ParseUint(userCookie.Value, 10, 64)
	if err != nil {
		rw.Write(failure)
		return
	}
	message := datamodels.NewMessage(uint(userId), uint(productId))
	// data type conversion
	byteMessage, err := json.Marshal(message)
	if err != nil {
		rw.Write(failure)
		return
	}

	err = rabbitmqValidate.PublishSimple(string(byteMessage))
	if err != nil {
		rw.Write(failure)
		return
	}
	rw.Write(success)
	return
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
	// add node with consistent hash algo
	for _, v := range hostArray {
		// add hostIP to hash ring
		hashConsistent.Add(v)
	}

	localIp, err := common.GetIntranceIp()
	if err != nil {
		log.Println(err)
	}
	localHost = localIp
	log.Println("LocalHost = " + localHost)

	accessControl.SetHosts(localHost, port)
	accessControl.SetConsistentHash(hashConsistent)

	rabbitmqValidate = distributed.NewRabbitMQSimple(common.AMQP_QUEUE_NAME)
	defer rabbitmqValidate.Destroy()

	http.Handle("/html/", http.StripPrefix("/html", http.FileServer(http.Dir("./frontend/web/htmlProductShow"))))
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir(common.CDN_DOMAIN_URL))))

	filter := common.NewFiler()
	filter.RegisterFilterUri("/check", Auth)
	filter.RegisterFilterUri("/checkRight", Auth)
	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/checkRight", filter.Handle(CheckUserRight))
	http.ListenAndServe("localhost:8083", nil)
}