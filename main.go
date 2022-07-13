package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	c "miniwallet.co.id/controllers"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", c.HandleHome).Methods("GET")
	router.HandleFunc("/api/v1/init", c.HandleInitWallet).Methods("POST")
	router.HandleFunc("/api/v1/wallet", c.EnableWallet).Methods("POST")
	router.HandleFunc("/api/v1/wallet", c.DisableWallet).Methods("PATCH")
	router.HandleFunc("/api/v1/wallet", c.ViewWallet).Methods("GET")
	router.HandleFunc("/api/v1/wallet/deposits", c.DepositMoney).Methods("POST")
	router.HandleFunc("/api/v1/wallet/withdrawals", c.WithdrawMoney).Methods("POST")
	log.Fatal(http.ListenAndServe(":1991", router))
}
