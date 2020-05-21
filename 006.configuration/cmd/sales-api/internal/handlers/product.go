package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"garagesale/006.configuration/internal/product"
)

// List give all products from the database as a list
func List(w http.ResponseWriter, req *http.Request) {
	// var br BookRepo
	getdata := product.GetAllData()
	data, err := json.Marshal(getdata)
	if err != nil {
		fmt.Print("error in getting data from database", err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(w, "%s\n", data)

}
