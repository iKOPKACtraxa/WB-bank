package internalhttp

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/iKOPKACtraxa/wb-bank/internal/storage"
)

type Server struct {
	HTTPServer *http.Server
	Storage    storage.Storage
}

// NewServer returns a new Server object.
func NewServer(s storage.Storage, hostPort string) *Server {
	mux := http.NewServeMux()
	server := &Server{
		HTTPServer: &http.Server{
			Addr:    hostPort,
			Handler: mux,
		},
		Storage: s,
	}

	mux.HandleFunc("/set", server.set)
	mux.HandleFunc("/get", server.get)
	mux.HandleFunc("/transfer", server.transfer)
	return server
}

// Start starts Server.
func (s *Server) Start() error {
	err := s.HTTPServer.ListenAndServe()
	return err
}

// set sets new account with balance.
func (s *Server) set(w http.ResponseWriter, r *http.Request) {
	account := storage.AccountID(r.FormValue("account"))
	balanceInt, err := strconv.Atoi(r.FormValue("balance"))
	if err != nil {
		_, err := w.Write([]byte(fmt.Sprintf("Ошибка: %v\n", err)))
		if err != nil {
			log.Println(err)
		}
		return
	}
	balance := storage.AccountBalance(balanceInt)
	err = s.Storage.Set(account, balance)
	if err != nil {
		_, err := w.Write([]byte(fmt.Sprintf("Ошибка: %v\n", err)))
		if err != nil {
			log.Println(err)
		}
		return
	}
	_, err = w.Write([]byte(fmt.Sprintf("Аккаунт %v добавлен с балансом %v \n", account, balance)))
	if err != nil {
		log.Println(err)
	}
	s.Storage.Print()
}

// get gets balance from account.
func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	account := storage.AccountID(r.FormValue("account"))
	balance, err := s.Storage.Get(account)
	if err != nil {
		_, err := w.Write([]byte(fmt.Sprintf("Ошибка: %v\n", err)))
		if err != nil {
			log.Println(err)
		}
		return
	}
	_, err = w.Write([]byte(fmt.Sprintf("Аккаунт %v, баланс %v\n", account, balance)))
	if err != nil {
		log.Println(err)
	}
}

// transfer transfers balance from one account to another.
func (s *Server) transfer(w http.ResponseWriter, r *http.Request) {
	accountFrom := storage.AccountID(r.FormValue("accountfrom"))
	accountTo := storage.AccountID(r.FormValue("accountto"))
	balanceInt, err := strconv.Atoi(r.FormValue("balance"))
	if err != nil {
		_, err := w.Write([]byte(fmt.Sprintf("Ошибка: %v\n", err)))
		if err != nil {
			log.Println(err)
		}
		return
	}
	balanceToTransfer := storage.AccountBalance(balanceInt)
	errTansfer := s.Storage.Transfer(accountFrom, accountTo, balanceToTransfer)
	if errTansfer != nil {
		_, err := w.Write([]byte(fmt.Sprintf("Операция невозможна, ошибка: %v\n", errTansfer)))
		if err != nil {
			log.Println(err)
		}
		return
	}
	_, err = w.Write([]byte(fmt.Sprintln("Операция успешна, балансы обновлены")))
	if err != nil {
		log.Println(err)
	}
	s.Storage.Print()
}
