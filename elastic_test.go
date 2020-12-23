package monitoringsystem

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	newrelic "github.com/newrelic/go-agent"
)

func Test_NewElasticMonitor(t *testing.T) {
	type args struct {
		appName string
		url     string
	}
	tests := []struct {
		name    string
		args    args
		want    *elastic
		wantErr bool
	}{
		{"Expecting test case pass as invalid url provided", args{appName: "iam", url: "postgres://user:abc{DEf1=ghi@example.com:5432/db?sslmode=require"}, nil, true},
		{"Expecting test case pass as valid url provided", args{appName: "iam", url: "http://cms-in.bigtree.org:8200"}, nil, false},
		{"Expecting test case pass", args{appName: "iam"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ELASTIC_APM_SERVER_URL", tt.args.url)
			_, err := newElastic(tt.args.appName, "0984793213650b3adfc8e4d89a9bb9245fb5934e", tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("newElastic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewElastic_StartTransaction(t *testing.T) {

	tests := []struct {
		name            string
		elastic         *elastic
		transactionName string
		wantErr         bool
	}{
		{
			"Excepting error as application not initialized",
			&elastic{},
			"iam.manage.user.r",
			true,
		},
		{
			"Excepting transaction object",
			func() *elastic {
				n, _ := newElastic("iam", "", "")
				return n
			}(),
			"iam.manage.user.r",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.elastic.StartTransaction(tt.transactionName)
			if (err != nil) != tt.wantErr {
				t.Errorf("elastic.StartTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewElastic_StartWebTransaction(t *testing.T) {

	tests := []struct {
		name            string
		elastic         *elastic
		transactionName string
		wantErr         bool
	}{
		{
			"Excepting error as application not initialized",
			&elastic{nil},
			"iam.manage.user.r",
			true,
		},
		{
			"Excepting transaction object",
			func() *elastic {
				n, _ := newElastic("iam", "", "")
				return n
			}(),
			"iam.manage.user.r",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.elastic.StartWebTransaction(tt.transactionName, httptest.NewRecorder(), func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/user", nil)
				return req
			}())
			if (err != nil) != tt.wantErr {
				t.Errorf("elastic.StartWebTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewElastic_EndTransaction(t *testing.T) {
	tests := []struct {
		name        string
		elastic     *elastic
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&elastic{nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *elastic {
				n, _ := newElastic("iam", "", "")
				return n
			}(),
			func() interface{} {
				n, _ := newElastic("iam", "", "")
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.elastic.EndTransaction(tt.transaction); (err != nil) != tt.wantErr {
				t.Errorf("elastic.EndTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewElastic_StartSegment(t *testing.T) {
	tests := []struct {
		name        string
		elastic     *elastic
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&elastic{nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			func() interface{} {
				n, _ := newElastic("iam", "")
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.elastic.StartSegment("otp", tt.transaction); (err != nil) != tt.wantErr {
				t.Errorf("elastic.StartSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewElastic_EndSegment(t *testing.T) {
	tests := []struct {
		name        string
		elastic     *elastic
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&elastic{nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			func() interface{} {
				n, _ := newElastic("iam", "")
				trans, _ := n.StartTransaction("/user")
				segment, _ := n.StartSegment("otp_service", trans)
				return segment
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.elastic.EndSegment(tt.transaction); (err != nil) != tt.wantErr {
				t.Errorf("elastic.EndTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewElastic_StartDataStoreSegment(t *testing.T) {
	tests := []struct {
		name        string
		elastic     *elastic
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&elastic{nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			func() interface{} {
				n, _ := newElastic("iam", "")
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.elastic.StartDataStoreSegment("MongoDB", tt.transaction, "Find", "tblUsers", Operation{}); (err != nil) != tt.wantErr {
				t.Errorf("elastic.StartDataStoreSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewElastic_EndDataStoreSegment(t *testing.T) {
	tests := []struct {
		name        string
		elastic     *elastic
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&elastic{nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			func() interface{} {
				n, _ := newElastic("iam", "")
				trans, _ := n.StartTransaction("/user")
				segment, _ := n.StartDataStoreSegment("MongoDB", trans, "Find", "tblUsers")
				return segment
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.elastic.EndDataStoreSegment(tt.transaction); (err != nil) != tt.wantErr {
				t.Errorf("elastic.EndDataStoreSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewElastic_StartExternalSegment(t *testing.T) {
	tests := []struct {
		name        string
		elastic     *elastic
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&elastic{nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			func() interface{} {
				n, _ := newElastic("iam", "")
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.elastic.StartExternalSegment(tt.transaction, "http://iam.bookmyshow.com"); (err != nil) != tt.wantErr {
				t.Errorf("elastic.StartExternalSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewElastic_StartExternalWebSegment(t *testing.T) {
	tests := []struct {
		name        string
		elastic     *elastic
		transaction interface{}
		req         *http.Request
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&elastic{nil},
			nil,
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			func() interface{} {
				n, _ := newElastic("iam", "")
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			httptest.NewRequest(http.MethodGet, "/user", nil),
			false,
		},
		{
			"Excepting error as invalid request object provided",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			func() interface{} {
				n, _ := newElastic("iam", "")
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			new(http.Request),
			true,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			new(http.Request),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			"invalid",
			new(http.Request),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.elastic.StartExternalWebSegment(tt.transaction, tt.req); (err != nil) != tt.wantErr {
				t.Errorf("elastic.StartExternalWebSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewElastic_EndExternalSegment(t *testing.T) {
	tests := []struct {
		name        string
		elastic     *elastic
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&elastic{nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			func() interface{} {
				n, _ := newElastic("iam", "")
				trans, _ := n.StartTransaction("/user")
				externalSegment, _ := n.StartExternalSegment(trans, "http://iam.bookmyshow.com")
				return externalSegment
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction type provided",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			"invalid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.elastic.EndExternalSegment(tt.transaction); (err != nil) != tt.wantErr {
				t.Errorf("elastic.EndExternalSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewElastic_NoticeError(t *testing.T) {
	tests := []struct {
		name        string
		elastic     *elastic
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&elastic{nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			func() interface{} {
				n, _ := newElastic("iam", "")
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction object provided",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			"iam.manage.user.r",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.elastic.NoticeError(tt.transaction, errors.New("error")); (err != nil) != tt.wantErr {
				t.Errorf("elastic.NoticeError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewElastic_AddAttribute(t *testing.T) {
	tests := []struct {
		name        string
		elastic     *elastic
		transaction interface{}
		wantErr     bool
	}{
		{
			"Excepting error as application not initialized",
			&elastic{nil},
			nil,
			true,
		},
		{
			"Excepting transaction to end successfully",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			func() interface{} {
				n, _ := newElastic("iam", "")
				trans, _ := n.StartTransaction("/user")
				return trans
			}(),
			false,
		},
		{
			"Excepting error as transaction input value is nil",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			newrelic.NewWebRequest(nil),
			true,
		},
		{
			"Excepting error as invalid transaction object provided",
			func() *elastic {
				n, _ := newElastic("iam", "")
				return n
			}(),
			"iam.manage.user.r",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.elastic.AddAttribute(tt.transaction, "feature", "iam.manage.user.u"); (err != nil) != tt.wantErr {
				t.Errorf("elastic.AddAttribute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
