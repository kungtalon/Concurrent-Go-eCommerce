package main

import (
	"log"
	"net/http"
	"sync"
)

var sum int64 = 0

// count of product to be sold
var productNum int64 = 10000

// mutual exclusive lock
var mutex sync.Mutex

func GetOneProduct() bool {
	mutex.Lock()
	defer mutex.Unlock()
	// check whether the count of product has exceeded storage
	if sum < productNum {
		sum += 1
		return true
	}
	return false
}

func GetProduct(rw http.ResponseWriter, req *http.Request) {
	if GetOneProduct() {
		rw.Write([]byte("true"))
		return
	}
	rw.Write([]byte("false"))
}

func main() {
	http.HandleFunc("/getOne", GetProduct)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("Err: ", err)
	}
}
