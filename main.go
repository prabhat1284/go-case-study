package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Item struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    string `json:"price"`
}

var Items []Item

var ApiUrls = []string{
	"https://run.mocky.io/v3/c51441de-5c1a-4dc2-a44e-aab4f619926b",
	"https://run.mocky.io/v3/4ec58fbc-e9e5-4ace-9ff0-4e893ef9663c",
	"https://run.mocky.io/v3/e6c77e5c-aec9-403f-821b-e14114220148",
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", Home).Methods("GET")
	r.HandleFunc("/buy-item/{name}", getByName).Methods("GET")
	r.HandleFunc("/buy-item-qty/{name}&{quantity}", getByQuantity).Methods("GET")
	r.HandleFunc("/buy-item-qty-price/{name}&{quantity}&{price}", getByPrice).Methods("GET")
	fmt.Println("Server Started on 4000 port")
	log.Fatal(http.ListenAndServe(":4000", r))

}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to Food Aggregator</h1>"))
}

func getByName(w http.ResponseWriter, r *http.Request) {
	flag := 0
	for _, web := range ApiUrls {

		res, err := http.Get(web)

		if err != nil {
			panic(err)
		}

		dataBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		data := strClean(dataBytes)

		json.Unmarshal(data, &Items)

		params := mux.Vars(r)

		for _, item := range Items {

			if item.Name == params["name"] {
				flag = 1
				json.NewEncoder(w).Encode(item)
				return
			}
		}

	}

	if flag == 0 {
		json.NewEncoder(w).Encode("NOT_FOUND")
	}

}

func getByQuantity(w http.ResponseWriter, r *http.Request) {
	flag := 0
	for _, web := range ApiUrls {

		res, err := http.Get(web)

		if err != nil {
			panic(err)
		}

		dataBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		data := strClean(dataBytes)

		json.Unmarshal(data, &Items)

		params := mux.Vars(r)
		qty, _ := strconv.Atoi(params["quantity"])

		for _, item := range Items {
			fmt.Println(item)
			if item.Name == params["name"] && item.Quantity >= qty {
				flag = 1
				json.NewEncoder(w).Encode(item)
				return
			}
		}

	}

	if flag == 0 {
		json.NewEncoder(w).Encode("NOT_FOUND")
	}

}

func getByPrice(w http.ResponseWriter, r *http.Request) {
	flag := 0
	for _, web := range ApiUrls {

		res, err := http.Get(web)

		if err != nil {
			panic(err)
		}

		dataBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		data := strClean(dataBytes)

		json.Unmarshal(data, &Items)

		params := mux.Vars(r)
		qty, _ := strconv.Atoi(params["quantity"])

		for _, item := range Items {
			fmt.Println(item)
			if item.Name == params["name"] && item.Quantity >= qty && item.Price[1:] >= params["price"] {
				flag = 1
				json.NewEncoder(w).Encode(item)
				return
			}
		}

	}

	if flag == 0 {
		json.NewEncoder(w).Encode("NOT_FOUND")
	}

}

func strClean(dataBytes []byte) []byte {
	dat0 := strings.ReplaceAll(string(dataBytes), "itemId", "id")
	dat1 := strings.ReplaceAll(dat0, "itemName", "name")
	dat2 := strings.ReplaceAll(dat1, "productId", "id")
	dat3 := strings.ReplaceAll(dat2, "productName", "name")
	return []byte(dat3)
}
