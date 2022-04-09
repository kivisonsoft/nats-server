package client

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"net/http"
	"time"
)

var(
	nc *nats.Conn
)

func Start(log server.Logger) {
	nc = newNatsClient()
	r := mux.NewRouter()
	r.HandleFunc("/pb/{name}", postbackHandler)
	r.HandleFunc("/rpb/{name}", rejectPostbackHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:80",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Noticef("Client is ready")
	log.Fatalf("%v",srv.ListenAndServe())
}

func newNatsClient() *nats.Conn {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	return nc
}

type Message struct {
	Name string `json:"name"`
	Params map[string][]string `json:"params"`
}

func postbackHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("Name: %v, %v: %v \n", vars["name"], r.Method, r.URL.RawQuery)
	msg := &Message{
		Name:   vars["name"],
		Params: r.URL.Query(),
	}
	marshal, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("error:%v \n",err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error"))
		return
	}
	err = nc.Publish("postback-event", marshal)
	if err!=nil{
		fmt.Printf("error:%v \n",err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func rejectPostbackHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("Name: %v, %v: %v \n", vars["name"], r.Method, r.URL.RawQuery)
	msg := &Message{
		Name:   vars["name"],
		Params: r.URL.Query(),
	}
	marshal, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("error:%v \n",err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error"))
		return
	}
	err = nc.Publish("reject-postback-event", marshal)
	if err!=nil{
		fmt.Printf("error:%v \n",err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}