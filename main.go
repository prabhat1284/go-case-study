package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

// Create struct with name Item
type Item struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    string `json:"price"`
}

// Create a cache with a default expiration time of 5 minutes, purges expired items every 10 minutes
var c = cache.New(5*time.Minute, 10*time.Minute)

// Suppliers API
var apiUrls = []string{
	"https://run.mocky.io/v3/c51441de-5c1a-4dc2-a44e-aab4f619926b",
	"https://run.mocky.io/v3/4ec58fbc-e9e5-4ace-9ff0-4e893ef9663c",
	"https://run.mocky.io/v3/e6c77e5c-aec9-403f-821b-e14114220148",
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/food-aggregator", foodAggregator).Methods("GET")
	r.HandleFunc("/buy-item/{name}", getByName).Methods("GET")
	r.HandleFunc("/buy-item-qty/{name}&{quantity}", getByQuantity).Methods("GET")
	r.HandleFunc("/buy-item-qty-price/{name}&{quantity}&{price}", getByPrice).Methods("GET")
	r.HandleFunc("/show-summary", showSummary).Methods("GET")
	r.HandleFunc("/fast-buy-item/{name}", getFastByName).Methods("GET")

	fmt.Println("Server Started on 4000 port")
	log.Fatal(http.ListenAndServe(":4000", r))

}

//foodAggregator
func foodAggregator(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to Food Aggregator</h1>"))
}

// getByName : API to check the item
func getByName(w http.ResponseWriter, r *http.Request) {

	flag := 0
	params := mux.Vars(r)
	name := params["name"]

	result, err := Suppliers(name)

	if err != nil {
		panic(err)
	}

	if len(name) == 0 {
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}

	if len(result) > 0 {
		for i := 0; i < len(result); i++ {
			if result[i].Name == name {
				flag = 1
				json.NewEncoder(w).Encode(result[i])
				return
			}
		}
	}

	if flag == 0 {
		json.NewEncoder(w).Encode("NOT_FOUND")
	}

}

// Handler function to check by name and quantity
func getByQuantity(w http.ResponseWriter, r *http.Request) {

	flag := 0
	params := mux.Vars(r)
	name := params["name"]
	quantity, _ := strconv.Atoi(params["quantity"])

	if len(name) == 0 || quantity <= 0 {
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}

	result, err := Suppliers(name)

	if err != nil {
		panic(err)
	}

	if len(result) > 0 {
		for i := 0; i < len(result); i++ {
			if result[i].Name == name && result[i].Quantity >= quantity {
				flag = 1
				json.NewEncoder(w).Encode(result[i])
				return
			}
		}
	}

	if flag == 0 {
		json.NewEncoder(w).Encode("NOT_FOUND")
	}

}

// getByPrice :- to check by name, quantity and price
func getByPrice(w http.ResponseWriter, r *http.Request) {

	flag := 0
	params := mux.Vars(r)
	name := params["name"]
	quantity, _ := strconv.Atoi(params["quantity"])
	price := params["price"]

	if len(name) == 0 || len(price) == 0 || quantity <= 0 {
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}

	cacheData, found := c.Get(name)
	if found {
		json.NewEncoder(w).Encode(cacheData)
		return
	}

	result, err := Suppliers(name)

	if err != nil {
		panic(err)
	}

	if len(result) > 0 {
		for i := 0; i < len(result); i++ {
			if result[i].Name == name && result[i].Quantity >= quantity && result[i].Price[1:] == price {
				flag = 1
				c.Set(result[i].Name, result[i], cache.NoExpiration)
				json.NewEncoder(w).Encode(result[i])
				return
			}
		}
	}

	if flag == 0 {
		json.NewEncoder(w).Encode("NOT_FOUND")
	}

}

func showSummary(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(c.Items())
}

// fast-buy-item API :- to check the item by name
func getFastByName(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	name := params["name"]

	if len(name) == 0 {
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}

	c1 := make(chan Item)
	c2 := make(chan Item)
	c3 := make(chan Item)

	go FastSuppliers(name, c1, apiUrls[0])
	go FastSuppliers(name, c2, apiUrls[1])
	go FastSuppliers(name, c3, apiUrls[2])

	select {
	case result1 := <-c1:
		json.NewEncoder(w).Encode(result1)
	case result2 := <-c2:
		json.NewEncoder(w).Encode(result2)
	case result3 := <-c3:
		json.NewEncoder(w).Encode(result3)
	}

}

// Suppliers: API to buy items from suppliers
func Suppliers(name string) ([]Item, error) {

	for _, url := range apiUrls {
		res, err := http.Get(url)

		if err != nil {
			panic(err)
		}

		dataBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		data := strClean(dataBytes)
		item := []Item{}

		json.Unmarshal(data, &item)

		for _, product := range item {

			if product.Name == name {
				return item, nil
			}
		}
	}
	return nil, nil
}

// FastSuppliers: API to buy items from suppliers
func FastSuppliers(name string, c chan Item, url string) {

	res, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	dataBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	data := strClean(dataBytes)
	item := []Item{}

	json.Unmarshal(data, &item)

	for _, product := range item {

		if product.Name == name {
			c <- product
		}
	}

}

// to make the key values uniform
func strClean(dataBytes []byte) []byte {
	dat0 := strings.ReplaceAll(string(dataBytes), "itemId", "id")
	dat1 := strings.ReplaceAll(dat0, "itemName", "name")
	dat2 := strings.ReplaceAll(dat1, "productId", "id")
	dat3 := strings.ReplaceAll(dat2, "productName", "name")
	return []byte(dat3)
}
