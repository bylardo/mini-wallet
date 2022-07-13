package controllers

import (
	"fmt"
	"net/http"

	"miniwallet.co.id/workers"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Mini Wallet API")
}

func HandleInitWallet(w http.ResponseWriter, r *http.Request) {
	workers.DoInitWallet(w, r)
}

func EnableWallet(w http.ResponseWriter, r *http.Request) {
	workers.DoEnableWallet(w, r)
}

func DisableWallet(w http.ResponseWriter, r *http.Request) {
	workers.DoDisableWallet(w, r)
}

func ViewWallet(w http.ResponseWriter, r *http.Request) {
	workers.ViewWallet(w, r)
}

func DepositMoney(w http.ResponseWriter, r *http.Request) {
	workers.DepositMoney(w, r)
}

func WithdrawMoney(w http.ResponseWriter, r *http.Request) {
	workers.WithdrawMoney(w, r)
}
