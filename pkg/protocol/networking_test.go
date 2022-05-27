package protocol

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNetworkPackage(t *testing.T) {
	type args struct {
		message Message
	}
	tests := []struct {
		name    string
		args    args
		want    Message
		wantErr bool
	}{
		{
			"Usual package",
			args{
				NewMessage(1, "Hello"),
			},
			NewMessage(1, "Hello"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := &bytes.Buffer{}
			err := SendPackage(tt.args.message, conn)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendPackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := ReadPackage(conn)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendPackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SendPackage() gotConn = %v, want %v", got, tt.want)
			}
		})
	}
}
