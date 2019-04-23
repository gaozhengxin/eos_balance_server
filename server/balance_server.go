package main

import (
	"../dao"
	"../transfer_tracker"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main () {
	port := flag.String("port","7000","port")
	dbpath := flag.String("dbpath", "./data", "database path")
	reinit := flag.Bool("reinit", false, "reinit")
	flag.Parse()
	path := "0.0.0.0:" + *port
	defer dao.Close()
	go tracker.Run(dbpath,reinit)
	http.HandleFunc("/get_balance", GetBalance)
	go http.ListenAndServe(path, nil)
	fmt.Printf("service is running on %s\n", path)
	select{}
}

type Resp struct {
	Code string `json:"code"`
	Msg string `json:"Msg,omitempty"`
	Balance string `json:"balance,omitempty"`
}

func GetBalance (writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	userKey, ok := request.Form["user_key"]
	var result Resp
	if !ok {
		result.Code = "401"
		result.Msg = "no user_key"
	} else {
		result.Code = "200"
		result.Balance = dao.GetBalance(userKey[0])
	}
	if err := json.NewEncoder(writer).Encode(result); err != nil {
		log.Fatal(err)
	}
}
