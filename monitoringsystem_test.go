package monitoringsystem

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	newrelic "github.com/newrelic/go-agent"
)

func Example_newRelic() {
	apm, err := New(NewRelic, true)
	if err != nil {
		//Handler error
	}

	txn, _ := apm.StartWebTransaction("/users", nil, httptest.NewRequest(http.MethodGet, "/", nil))
	apm.AddAttribute(txn, "IAMFeature", "iam.manage.user.r")
	defer apm.EndTransaction(txn, nil)
	dataSegment, _ := apm.StartDataStoreSegment(txn, "Mongo", "find", "tblUsers")
	apm.EndSegment(dataSegment)
	segment, _ := apm.StartSegment(txn, "opt-service")
	apm.EndSegment(segment)
	externalSegment, _ := apm.StartExternalSegment(txn, "http://iam.bookmyhsow.com")
	apm.EndExternalSegment(externalSegment)
}

func Example_elastic() {
	apm, err := New(Elastic, true, Option{
		ElasticServiceName:  "apm-test",
		ElasticServerAMPUrl: "http://localhost:8200",
	})
	if err != nil {
		//Handler error
	}
	txn, _ := apm.StartWebTransaction("/users", nil, httptest.NewRequest(http.MethodGet, "/", nil))
	apm.AddAttribute(txn, "IAMFeature", "iam.manage.user.r")
	defer apm.EndTransaction(txn, nil)
	dataSegment, _ := apm.StartDataStoreSegment(txn, "Mongo", "find", "tblUsers")
	apm.EndSegment(dataSegment)
	segment, _ := apm.StartSegment(txn, "opt-service")
	apm.EndSegment(segment)
	externalSegment, _ := apm.StartExternalSegment(txn, "http://iam.bookmyhsow.com")
	apm.EndExternalSegment(externalSegment)
}

func Test_New(t *testing.T) {
	tests := []struct {
		name             string
		monitoringType   agentType
		appName          string
		licence          string
		enableMonitoring bool
		wantErr          bool
	}{
		{
			name:             "Expecting nil error as new relic required configuration not provided",
			monitoringType:   NewRelic,
			enableMonitoring: true,
			wantErr:          true,
		},
		{
			name:             "Expecting error as new relic required configuration not provided",
			monitoringType:   NewRelic,
			appName:          "iam",
			licence:          "",
			enableMonitoring: true,
			wantErr:          true,
		},
		{
			name:             "Expecting valid APM as monitoring is disabled",
			monitoringType:   NewRelic,
			enableMonitoring: false,
			wantErr:          false,
		},
		{
			name:             "Expecting invalid support error",
			monitoringType:   "",
			enableMonitoring: true,
			wantErr:          true,
		},
		{
			name:             "Expecting error invalid new relic licence key provided",
			monitoringType:   NewRelic,
			appName:          "iam",
			licence:          "adfc8e4d89a9bb9245fb5934e",
			enableMonitoring: true,
			wantErr:          true,
		},
		{
			name:             "Expecting NewRelic APM Handler",
			monitoringType:   NewRelic,
			appName:          "iam",
			licence:          "0984793213650b3adfc8e4d89a9bb9245fb5934e",
			enableMonitoring: true,
			wantErr:          false,
		},
		{
			name:             "Expecting Elastic APM Handler",
			monitoringType:   Elastic,
			appName:          "iam",
			licence:          "0984793213650b3adfc8e4d89a9bb9245fb5934e",
			enableMonitoring: true,
			wantErr:          false,
		},
		{
			name:             "Expecting Elastic APM Handler",
			monitoringType:   Elastic,
			appName:          "",
			licence:          "0984793213650b3adfc8e4d89a9bb9245fb5934e",
			enableMonitoring: true,
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("NEWRELIC_APP", tt.appName)
			os.Setenv("NEWRELIC_KEY", tt.licence)
			os.Setenv("ELASTIC_APM_SERVICE_NAME", tt.appName)
			_, err := New(tt.monitoringType, tt.enableMonitoring, Option{
				ElasticServiceName: os.Getenv("ELASTIC_APM_SERVICE_NAME"),
			})
			if (err != nil) != tt.wantErr {
				if err != nil {
					t.Errorf("Expected %v got %v", err.Error(), tt.wantErr)
					return
				}
			}
		})
	}
}

func Test_Enable(t *testing.T) {

	tests := []struct {
		name   string
		APM    *Agent
		enable bool
	}{
		{
			"Excepting monitoring should be disabled",
			&Agent{"iam", true, nil, &sync.Mutex{}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.APM.Enable(tt.enable)
		})
	}
}

func Test_StartTransaction(t *testing.T) {

	tests := []struct {
		name            string
		APM             *Agent
		transactionName string
		wantErr         bool
	}{
		{
			"Excepting nil error as application not initialized",
			&Agent{"iam", true, nil, &sync.Mutex{}},
			"iam.manage.user.r",
			false,
		},
		{
			"Excepting transaction object",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			"iam.manage.user.r",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.APM.StartTransaction(tt.transactionName)
			if (err != nil) != tt.wantErr {
				t.Errorf("New.StartTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_StartWebTransaction(t *testing.T) {

	tests := []struct {
		name            string
		APM             *Agent
		transactionName string
		wantErr         bool
	}{
		{
			"Excepting nil error as application not initialized",
			&Agent{"iam", true, nil, &sync.Mutex{}},
			"iam.manage.user.r",
			false,
		},
		{
			"Excepting transaction object",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			"iam.manage.user.r",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.APM.StartWebTransaction(tt.transactionName, httptest.NewRecorder(), new(http.Request))
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRelic.StartWebTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_EndTransaction(t *testing.T) {
	tests := []struct {
		name        string
		APM         *Agent
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting nil error as application not initialized",
			&Agent{"iam", true, nil, &sync.Mutex{}},
			nil,
			false,
		},
		{
			"Excepting transaction to end successfully",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			func() interface{} {
				n, _ := New(NewRelic, true)
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.APM.EndTransaction(tt.transaction, errors.New("invalid access code provided")); (err != nil) != tt.wantErr {
				t.Errorf("New.EndTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew_StartSegment(t *testing.T) {
	tests := []struct {
		name        string
		APM         *Agent
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting nil error as application not initialized",
			&Agent{"iam", true, nil, &sync.Mutex{}},
			nil,
			false,
		},
		{
			"Excepting transaction to end successfully",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			func() interface{} {
				n, _ := New(NewRelic, true)
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting nil error as transaction input value is nil",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting nil error as invalid transaction type provided",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.APM.StartSegment(tt.transaction, "otp_service"); (err != nil) != tt.wantErr {
				t.Errorf("New.StartSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew_EndSegment(t *testing.T) {
	tests := []struct {
		name        string
		APM         *Agent
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting nil error as application not initialized",
			&Agent{"iam", true, nil, &sync.Mutex{}},
			nil,
			false,
		},
		{
			"Excepting nil transaction to end successfully",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			func() interface{} {
				n, _ := New(NewRelic, true)
				trans, _ := n.StartTransaction("/user")
				segment, _ := n.StartSegment(trans, "otp_service")
				return segment
			}(),
			false,
		},
		{
			"Excepting nil error as transaction input value is nil",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting nil error as invalid transaction type provided",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.APM.EndSegment(tt.transaction); (err != nil) != tt.wantErr {
				t.Errorf("New.EndTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew_StartDataStoreSegment(t *testing.T) {
	tests := []struct {
		name        string
		APM         *Agent
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting nil error as application not initialized",
			&Agent{"iam", true, nil, &sync.Mutex{}},
			nil,
			false,
		},
		{
			"Excepting transaction to end successfully",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			func() interface{} {
				n, _ := New(NewRelic, true)
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.APM.StartDataStoreSegment(tt.transaction, "MongoDB", "Find", "tblUsers"); (err != nil) != tt.wantErr {
				t.Errorf("New.StartDataStoreSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew_EndDataStoreSegment(t *testing.T) {
	tests := []struct {
		name        string
		APM         *Agent
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting nil error as application not initialized",
			&Agent{"iam", true, nil, &sync.Mutex{}},
			nil,
			false,
		},
		{
			"Excepting transaction to end successfully",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			func() interface{} {
				n, _ := New(NewRelic, true)
				trans, _ := n.StartTransaction("/user")
				segment, _ := n.StartDataStoreSegment(trans, "MongoDB", "Find", "tblUsers")
				return segment
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.APM.EndDataStoreSegment(tt.transaction); (err != nil) != tt.wantErr {
				t.Errorf("New.EndDataStoreSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew_StartExternalSegment(t *testing.T) {
	tests := []struct {
		name        string
		APM         *Agent
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting nil error as application not initialized",
			&Agent{"iam", true, nil, &sync.Mutex{}},
			nil,
			false,
		},
		{
			"Excepting transaction to end successfully",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			func() interface{} {
				n, _ := New(NewRelic, true)
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.APM.StartExternalSegment(tt.transaction, "http://iam.bookmyshow.com"); (err != nil) != tt.wantErr {
				t.Errorf("New.StartExternalSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew_StartExternalWebSegment(t *testing.T) {
	tests := []struct {
		name        string
		APM         *Agent
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&Agent{"iam", true, nil, &sync.Mutex{}},
			nil,
			false,
		},
		{
			"Excepting transaction to end successfully",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			func() interface{} {
				n, _ := New(NewRelic, true)
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.APM.StartExternalWebSegment(tt.transaction, new(http.Request)); (err != nil) != tt.wantErr {
				t.Errorf("New.StartExternalWebSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew_EndExternalSegment(t *testing.T) {
	tests := []struct {
		name        string
		APM         *Agent
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting nil error as application not initialized",
			&Agent{"iam", true, nil, &sync.Mutex{}},
			nil,
			false,
		},
		{
			"Excepting transaction to end successfully",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			func() interface{} {
				n, _ := New(NewRelic, true)
				trans, _ := n.StartTransaction("/user")
				externalSegment, _ := n.StartExternalSegment(trans, "http://iam.bookmyshow.com")
				return externalSegment
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.APM.EndExternalSegment(tt.transaction); (err != nil) != tt.wantErr {
				t.Errorf("New.EndExternalSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew_NoticeError(t *testing.T) {
	tests := []struct {
		name        string
		APM         *Agent
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting nil error as application not initialized",
			&Agent{"iam", true, nil, &sync.Mutex{}},
			nil,
			false,
		},
		{
			"Excepting transaction to end successfully",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			func() interface{} {
				n, _ := New(NewRelic, true)
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction object provided",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			"iam.manage.user.r",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.APM.NoticeError(tt.transaction, errors.New("error")); (err != nil) != tt.wantErr {
				t.Errorf("New.NoticeError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew_AddAttribute(t *testing.T) {
	tests := []struct {
		name        string
		APM         *Agent
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting nil error as application not initialized",
			&Agent{"iam", true, nil, &sync.Mutex{}},
			nil,
			false,
		},
		{
			"Excepting transaction to end successfully",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			func() interface{} {
				n, _ := New(NewRelic, true)
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction object provided",
			func() *Agent {
				os.Setenv("NEWRELIC_APP", "iam")
				os.Setenv("NEWRELIC_KEY", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				n, _ := New(NewRelic, true)
				return n
			}(),
			"iam.manage.user.r",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.APM.AddAttribute(tt.transaction, "feature", "iam.manage.user.u"); (err != nil) != tt.wantErr {
				t.Errorf("New.AddAttribute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
