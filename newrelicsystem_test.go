package monitoringsystem

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	newrelic "github.com/newrelic/go-agent"
)

func Test_NewRelicMonitor(t *testing.T) {
	type args struct {
		appName string
		license string
	}
	tests := []struct {
		name    string
		args    args
		want    *newRelicHandler
		wantErr bool
	}{
		{"Expecting test case to fail as appName not provided", args{}, nil, true},
		{"Expecting test case pass", args{"iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newRelicMonitor(tt.args.appName, tt.args.license)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRelicMonitor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewRelic_StartTransaction(t *testing.T) {

	tests := []struct {
		name            string
		newRelic        *newRelicHandler
		transactionName string
		wantErr         bool
	}{
		{
			"Excepting error as application not initialized",
			&newRelicHandler{"iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e", nil},
			"iam.manage.user.r",
			true,
		},
		{
			"Excepting transaction object",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			"iam.manage.user.r",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.newRelic.StartTransaction(tt.transactionName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRelic.StartTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewRelic_StartWebTransaction(t *testing.T) {

	tests := []struct {
		name            string
		newRelic        *newRelicHandler
		transactionName string
		wantErr         bool
	}{
		{
			"Excepting error as application not initialized",
			&newRelicHandler{"iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e", nil},
			"iam.manage.user.r",
			true,
		},
		{
			"Excepting transaction object",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			"iam.manage.user.r",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.newRelic.StartWebTransaction(tt.transactionName, httptest.NewRecorder(), new(http.Request))
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRelic.StartWebTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewRelic_EndTransaction(t *testing.T) {
	tests := []struct {
		name        string
		newRelic    *newRelicHandler
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&newRelicHandler{"iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e", nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			func() interface{} {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.newRelic.EndTransaction(tt.transaction); (err != nil) != tt.wantErr {
				t.Errorf("NewRelic.EndTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewRelic_StartSegment(t *testing.T) {
	tests := []struct {
		name        string
		newRelic    *newRelicHandler
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&newRelicHandler{"iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e", nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			func() interface{} {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.newRelic.StartSegment("otp", tt.transaction); (err != nil) != tt.wantErr {
				t.Errorf("NewRelic.StartSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewRelic_EndSegment(t *testing.T) {
	tests := []struct {
		name        string
		newRelic    *newRelicHandler
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&newRelicHandler{"iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e", nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			func() interface{} {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				trans, _ := n.StartTransaction("/user")
				segment, _ := n.StartSegment("otp_service", trans)
				return segment
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.newRelic.EndSegment(tt.transaction); (err != nil) != tt.wantErr {
				t.Errorf("NewRelic.EndTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewRelic_StartDataStoreSegment(t *testing.T) {
	tests := []struct {
		name        string
		newRelic    *newRelicHandler
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&newRelicHandler{"iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e", nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			func() interface{} {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.newRelic.StartDataStoreSegment("MongoDB", tt.transaction, "Find", "tblUsers"); (err != nil) != tt.wantErr {
				t.Errorf("NewRelic.StartDataStoreSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewRelic_EndDataStoreSegment(t *testing.T) {
	tests := []struct {
		name        string
		newRelic    *newRelicHandler
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&newRelicHandler{"iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e", nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			func() interface{} {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				trans, _ := n.StartTransaction("/user")
				segment, _ := n.StartDataStoreSegment("MongoDB", trans, "Find", "tblUsers")
				return segment
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.newRelic.EndDataStoreSegment(tt.transaction); (err != nil) != tt.wantErr {
				t.Errorf("NewRelic.EndDataStoreSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewRelic_StartExternalSegment(t *testing.T) {
	tests := []struct {
		name        string
		newRelic    *newRelicHandler
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&newRelicHandler{"iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e", nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			func() interface{} {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.newRelic.StartExternalSegment(tt.transaction, "http://iam.bookmyshow.com"); (err != nil) != tt.wantErr {
				t.Errorf("NewRelic.StartExternalSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewRelic_StartExternalWebSegment(t *testing.T) {
	tests := []struct {
		name        string
		newRelic    *newRelicHandler
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&newRelicHandler{"iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e", nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			func() interface{} {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.newRelic.StartExternalWebSegment(tt.transaction, new(http.Request)); (err != nil) != tt.wantErr {
				t.Errorf("NewRelic.StartExternalWebSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewRelic_EndExternalSegment(t *testing.T) {
	tests := []struct {
		name        string
		newRelic    *newRelicHandler
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&newRelicHandler{"iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e", nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			func() interface{} {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				trans, _ := n.StartTransaction("/user")
				externalSegment, _ := n.StartExternalSegment(trans, "http://iam.bookmyshow.com")
				return externalSegment
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.newRelic.EndExternalSegment(tt.transaction); (err != nil) != tt.wantErr {
				t.Errorf("NewRelic.EndExternalSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewRelic_NoticeError(t *testing.T) {
	tests := []struct {
		name        string
		newRelic    *newRelicHandler
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&newRelicHandler{"iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e", nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			func() interface{} {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction object provided",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			"iam.manage.user.r",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.newRelic.NoticeError(tt.transaction, errors.New("error")); (err != nil) != tt.wantErr {
				t.Errorf("NewRelic.NoticeError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewRelic_AddAttribute(t *testing.T) {
	tests := []struct {
		name        string
		newRelic    *newRelicHandler
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&newRelicHandler{"iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e", nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			func() interface{} {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction object provided",
			func() *newRelicHandler {
				n, _ := newRelicMonitor("iam", "0984793213650b3adfc8e4d89a9bb9245fb5934e")
				return n
			}(),
			"iam.manage.user.r",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.newRelic.AddAttribute(tt.transaction, "feature", "iam.manage.user.u"); (err != nil) != tt.wantErr {
				t.Errorf("NewRelic.AddAttribute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
