// Exercise 7.11: Add additional handlers so that clients can create, read, update, and delete
// database entries. For example, a request of the form /update?item=socks&price=6 will
// update the price of an item in the inventory and report an error if the item does not exist or if
// the price is invalid. (Warning: this change introduces concurrent variable updates.)

package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/read", db.read)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

var dbMutex sync.Mutex

func (db database) list(w http.ResponseWriter, req *http.Request) {
	dbMutex.Lock()
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
	dbMutex.Unlock()
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	dbMutex.Lock()
	price, ok := db[item]
	dbMutex.Unlock()

	if ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	dbMutex.Lock()
	_, exists := db[item]
	dbMutex.Unlock()

	if exists {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "attemtping to create item that already exists %s\n", item)
		return
	}
	price := "0"
	if req.URL.Query().Has("price") {
		price = req.URL.Query().Get("price")
	}
	p, err := strconv.ParseInt(price, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "invalid price provided: %q\n", item)
		return
	}
	fmt.Fprintf(w, "Creating new entry %s with a price %s\n", item, price)
	dbMutex.Lock()
	db[item] = dollars(p)
	dbMutex.Unlock()
}

func (db database) read(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	dbMutex.Lock()
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "attemtping to create item that already exists %s\n", item)
	}
	dbMutex.Unlock()
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if req.URL.Query().Has("price") {
		price := req.URL.Query().Get("price")
		p, _ := strconv.ParseInt(price, 10, 64)

		dbMutex.Lock()
		defer dbMutex.Unlock()
		if _, ok := db[item]; ok {
			fmt.Fprintf(w, "updating existing item: %s to have price %s\n", item, price)
		} else {
			fmt.Fprintf(w, "no item %s found. Creating new with a price %s\n", item, price)
		}
		db[item] = dollars(p)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no price provided for item update: %q\n", item)
	}

}

func (db database) delete(w http.ResponseWriter, req *http.Request) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	item := req.URL.Query().Get("item")
	if _, ok := db[item]; ok {
		fmt.Fprintf(w, "deleting item %s from database", item)
		delete(db, item)
	} else {
		fmt.Fprintf(w, "no item %s found in database", item)
	}
}
