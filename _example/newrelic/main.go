package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	h.APM.AddAttribute(txn, "IAMFeature", "iam.manage.user.r")
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
	// h.APM.NoticeError(txn, errors.New("error"))
	fmt.Fprintf(w, "ok")
}

func main() {
	// Required NEWRELIC_APM, NEWRELIC_KEY environemt variable need
	os.Setenv("NEWRELIC_APM", "iam")
	os.Setenv("NEWRELIC_KEY", "a9e8741821c4e8e35303e434a1f028c692dc041b")
	monitor, err := monitoringsystem.New(monitoringsystem.Elastic, true)
	if err != nil {
		log.Fatal(err)
		return
	}
	apm := handler{
		APM: monitor,
	}
	http.HandleFunc("/user", apm.User)

	log.Printf("Server started listing on port 8000\n")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
