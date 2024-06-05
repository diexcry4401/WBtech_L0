package server

import (
	memcache "consumer/pkg/memcache"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func ServerStart(c *memcache.Cache) error {

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		HundlerGet(w, r, c)
	}) // GET
	fmt.Println("Server start listening on port 3001")
	err := http.ListenAndServe("localhost:3001", nil)

	if err != nil {
		fmt.Println(err)
		return err
	}
	select {}
}

func HundlerGet(w http.ResponseWriter, r *http.Request, c *memcache.Cache) {
	if r.Method == http.MethodGet {

		getOrder(w, r, c)
	} else {
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
	}

}

func getOrder(w http.ResponseWriter, r *http.Request, c *memcache.Cache) {

	idOrderStr := strings.TrimPrefix(r.URL.Path, "/api/")
	idOrder, err := strconv.Atoi(idOrderStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	fmt.Println("order")
	order, found := c.Get(fmt.Sprintf("%d", idOrder))
	fmt.Println(order)
	if !found {
		http.Error(w, "Order not found in cache", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
