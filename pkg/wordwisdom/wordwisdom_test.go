package wordwisdom

import (
	"alexander.rassanov/pow-tcp-server/pkg/cache"
	"alexander.rassanov/pow-tcp-server/pkg/pow"
	"alexander.rassanov/pow-tcp-server/pkg/protocol"
	cache2 "github.com/patrickmn/go-cache"
	"io"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestStream_ProcessMessage(t *testing.T) {
	type fields struct {
		cache     cache.Cache
		stream    io.ReadWriter
		clientID  string
		zeroCount int
	}
	type args struct {
		m protocol.Message
	}
	type testStruct struct {
		name    string
		fields  fields
		args    args
		want    protocol.Message
		wantErr bool
	}
	localCache := cache2.New(cache2.NoExpiration, time.Hour)
	tests := []func() testStruct{
		func() testStruct {
			rand.Seed(1)
			return testStruct{
				"Send request challenge and expect to receive challenge",
				fields{
					localCache,
					nil,
					"localhost",
					1,
				},
				args{protocol.Message{
					Header:  pow.RequestChallenge,
					Payload: "",
				}},
				protocol.Message{
					Header:  pow.ResponseChallenge,
					Payload: pow.NewHashCashDataChallenge("localhost", 1, rand.Intn(RandomForHashCash)),
				},
				false,
			}
		}, func() testStruct {
			rand.Seed(1)
			resolvedHash, _ := pow.NewHashCashDataChallenge("localhost", 1, rand.Intn(RandomForHashCash)).Resolve()
			return testStruct{
				"Send Request challenge with solved challenge and expect to receive a random quote",
				fields{
					localCache,
					nil,
					"localhost",
					1,
				},
				args{protocol.Message{
					Header:  pow.RequestService,
					Payload: resolvedHash,
				}},
				protocol.Message{
					Header:  pow.ResponseService,
					Payload: "Others see in the word of wisdom a teaching function.",
				},
				false,
			}
		},
	}
	for _, ttFunc := range tests {
		tt := ttFunc()
		t.Run(tt.name, func(t *testing.T) {
			rand.Seed(1)
			ww := StreamWithHashCash{
				cache:     tt.fields.cache,
				stream:    tt.fields.stream,
				clientID:  tt.fields.clientID,
				zeroCount: tt.fields.zeroCount,
			}
			got, err := ww.ProcessMessage(tt.args.m)

			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}
