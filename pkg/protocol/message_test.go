package protocol

import (
	"reflect"
	"testing"
)

const testHeader = iota

func TestNewMessage(t *testing.T) {
	type args struct {
		h MessageHeader
		p interface{}
	}

	type MyPayload struct {
		Message int
	}
	m := Message{
		testHeader, MyPayload{256},
	}
	tests := []struct {
		name string
		args args
		want Message
	}{
		{
			"Get fullfiled message structure",
			args{
				testHeader, MyPayload{256},
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
	RegisterType(MyPayload{})
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			"when Message is fullfiled",
			fields{
				testHeader, MyPayload{},
			},
			[]byte{44, 255, 129, 3, 1, 1, 7, 77, 101, 115, 115, 97, 103, 101, 1, 255, 130, 0, 1, 2, 1, 6, 72, 101, 97, 100, 101, 114, 1, 4, 0, 1, 7, 80, 97, 121, 108, 111, 97, 100, 1, 16, 0, 0, 0, 81, 255, 130, 2, 56, 97, 108, 101, 120, 97, 110, 100, 101, 114, 46, 114, 97, 115, 115, 97, 110, 111, 118, 47, 112, 111, 119, 45, 116, 99, 112, 45, 115, 101, 114, 118, 101, 114, 47, 112, 107, 103, 47, 112, 114, 111, 116, 111, 99, 111, 108, 46, 77, 121, 80, 97, 121, 108, 111, 97, 100, 255, 131, 3, 1, 1, 9, 77, 121, 80, 97, 121, 108, 111, 97, 100, 1, 255, 132, 0, 0, 0, 5, 255, 132, 1, 0, 0},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Message{
				Header:  tt.fields.Header,
				Payload: tt.fields.Payload,
			}
			if got, err := m.Encode(); (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Message.Encode() = %v\n want %v", got, tt.want)
			}
		})
	}
}

func TestParseMessage(t *testing.T) {
	type args struct {
		m []byte
	}
	encodedProperMessage, _ := NewMessage(testHeader, "my awesome service").Encode()
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
			NewMessage(testHeader, "my awesome service"),
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
