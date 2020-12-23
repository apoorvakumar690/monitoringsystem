package monitoringsystem

import "testing"

func Test_PanicRecovery(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Expected panic error",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer panicRecover()
			panic("error")
		})
	}
}
