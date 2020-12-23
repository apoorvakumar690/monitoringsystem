package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"stash.bms.bz/bms/monitoringsystem"
)

type handler struct {
	APM *monitoringsystem.Agent
}

func (h *handler) User(w http.ResponseWriter, req *http.Request) {
	// The call to StartTransaction must include the response writer and the
	// request.
	txn, _ := h.APM.StartWebTransaction("/users", w, req)
	h.APM.NoticeError(txn, errors.New("some thing went wrong"))
	h.APM.AddAttribute(txn, "feature", "iam.manage.user.r")
	defer h.APM.EndTransaction(txn, nil)
	time.Sleep(time.Second * 1)
	dataSegment, _ := h.APM.StartDataStoreSegment(txn, "Mongo", "find", "tblUsers")
	time.Sleep(20 * time.Millisecond)
	h.APM.EndSegment(dataSegment)
	segment, _ := h.APM.StartSegment(txn, "opt-service")
	time.Sleep(10 * time.Millisecond)
	h.APM.EndSegment(segment)
	externalSegment, _ := h.APM.StartExternalSegment(txn, "http://iam.bookmyhsow.com")
	time.Sleep(20 * time.Millisecond)
	h.APM.EndExternalSegment(externalSegment)
	time.Sleep(time.Second * 1)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	}); err != nil {
		fmt.Fprintf(w, "error")
		return
	}
}

func (h *handler) Error(w http.ResponseWriter, req *http.Request) {
	txn, _ := h.APM.StartWebTransaction("/error", w, req)
	defer h.APM.EndTransaction(txn, nil)
	h.APM.NoticeError(txn, errors.New("error"))
	fmt.Fprintf(w, "error")
}

func main() {
	monitor, err := monitoringsystem.New(monitoringsystem.Elastic, true, monitoringsystem.Option{
		ElasticServiceName:  "apm-test",
		ElasticServerAMPUrl: "http://localhost:8200",
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	apm := handler{
		APM: monitor,
	}
	http.HandleFunc("/user", apm.User)
	http.HandleFunc("/error", apm.Error)

	log.Printf("Server started listing on port 8000\n")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
