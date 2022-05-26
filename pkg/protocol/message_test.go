package protocol

import (
	"alexander.rassanov/pow-tcp-server/pkg/pow"
	"reflect"
	"testing"
)

func TestNewMessage(t *testing.T) {
	type args struct {
		h MessageHeader
		p interface{}
	}

	type MyPayload struct {
		Message int
	}
	m := Message{
		pow.RequestService, MyPayload{256},
	}
	tests := []struct {
		name string
		args args
		want Message
	}{
		{
			"Get fullfiled message structure",
			args{
				pow.RequestService, MyPayload{256},
			},
			m,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMessage(tt.args.h, tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_Encode(t *testing.T) {
	type fields struct {
		Header  MessageHeader
		Payload interface{}
	}
	type MyPayload struct{}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			"when Message is fullfiled",
			fields{
				pow.RequestService, MyPayload{},
			},
			[]byte{44, 255, 129, 3, 1, 1, 7, 77, 101, 115, 115, 97, 103, 101, 1, 255, 130, 0, 1, 2, 1, 6, 72, 101, 97, 100, 101, 114, 1, 4, 0, 1, 7, 80, 97, 121, 108, 111, 97, 100, 1, 16, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Message{
				Header:  tt.fields.Header,
				Payload: tt.fields.Payload,
			}
			if got, _ := m.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Message.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseMessage(t *testing.T) {
	type args struct {
		m []byte
	}
	encodedProperMessage, _ := NewMessage(pow.RequestService, "my awesome service").Encode()
	tests := []struct {
		name    string
		args    args
		want    Message
		wantErr bool
	}{
		{
			"trying to decode encoded message",
			args{
				encodedProperMessage,
			},
			NewMessage(pow.RequestService, "my awesome service"),
			false,
		}, {
			"trying to decode bad encoded message",
			args{
				[]byte("Hello worlds"),
			},
			NewMessage(0, nil),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMessage(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
