package pow

import (
	"alexander.rassanov/pow-tcp-server/pkg/protocol"
	"bytes"
	"io"
	"math/rand"
	"reflect"
	"testing"
)

func TestGetServiceByStream(t *testing.T) {
	type args struct {
		stream io.ReadWriter
	}
	tests := []struct {
		name             string
		args             args
		responseHashCash protocol.Message
		responseService  protocol.Message
		want             interface{}
		wantErr          bool
	}{
		{
			"Try to get a basic service",
			args{&bytes.Buffer{}},
			protocol.NewMessage(ResponseChallenge, NewHashCashDataChallenge("", 1, rand.Int())),
			protocol.NewMessage(ResponseService, interface{}("Hello World")),
			interface{}("Hello World"),
			false,
		},
		{
			"Try to get a basic service but server returns bad header for hash cash challenge",
			args{&bytes.Buffer{}},
			protocol.NewMessage(-1, NewHashCashDataChallenge("", 1, rand.Int())),
			protocol.NewMessage(ResponseService, interface{}("Hello World")),
			interface{}("Hello World"),
			true,
		}, {
			"Try to get a basic service but server returns bad header for response service",
			args{&bytes.Buffer{}},
			protocol.NewMessage(ResponseChallenge, NewHashCashDataChallenge("", 1, rand.Int())),
			protocol.NewMessage(-1, interface{}("Hello World")),
			interface{}("Hello World"),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			protocol.RegisterType(HashCashData{})
			if err := protocol.SendPackage(tt.responseHashCash, tt.args.stream); err != nil {
				t.Errorf("unexpected err in SendPackage: %v", err)
				return
			}
			go func() {
				_, _ = protocol.ReadPackage(tt.args.stream)
				if err := protocol.SendPackage(tt.responseService, tt.args.stream); err != nil {
					t.Errorf("unexpected err in SendPackage: %v", err)
					return
				}
			}()
			got, err := GetServiceByStream(tt.args.stream)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetServiceByStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetServiceByStream() got = %v, want %v", got, tt.want)
			}
		})
	}
}
