package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Grocery struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

var groceries = []Grocery{
	{Name: "Almond Milk", Quantity: 2},
	{Name: "Apple", Quantity: 6},
}

func AllGroceries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(groceries)
}

func SingleGrocery(w http.ResponseWriter, r *http.Request) {
	// Path params aise nikal sakte hain
	vars := mux.Vars(r)
	name := vars["name"]

	// _ ki jagah index likhte agar uska kahi use hota.
	for _, grocery := range groceries {
		if grocery.Name == name {
			// This chains marshalling and writing to w.
			// Could have done a json.Marashal on the object
			// written that object to w.
			// i guess
			json.NewEncoder(w).Encode(grocery)
		}
	}
}

func GroceriesToBuy(w http.ResponseWriter, r *http.Request) {
	// how you read request body
	body, _ := ioutil.ReadAll(r.Body)
	var grocery Grocery
	// marshall the body into object
	json.Unmarshal(body, &grocery)
	groceries = append(groceries, grocery)

	json.NewEncoder(w).Encode(groceries)
}

func UpdateGrocery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	for index, item := range groceries {
		if item.Name == name {
			// This is how you delete an entry from a slice in go :P
			groceries = append(groceries[:index], groceries[index+1:]...)
			var updatedGrocery Grocery
			json.NewDecoder(r.Body).Decode(&updatedGrocery)
			groceries = append(groceries, updatedGrocery)
			json.NewEncoder(w).Encode(updatedGrocery)
			return
		}
	}

}

func DeleteGrocery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	for index, item := range groceries {
		if item.Name == name {
			// removed the entry
			groceries = append(groceries[:index], groceries[index+1:]...)
			json.NewEncoder(w).Encode(item)
			return
		}
	}

}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/allgroceries", AllGroceries).Methods("GET")
	r.HandleFunc("/groceries/{name}", SingleGrocery).Methods("GET")
	r.HandleFunc("/groceries", GroceriesToBuy).Methods("POST")
	r.HandleFunc("/groceries/{name}", UpdateGrocery).Methods("PUT") // ----> To update a grocery
	r.HandleFunc("/groceries/{name}", DeleteGrocery).Methods("DELETE")

	http.ListenAndServe(":8080", r)
}
