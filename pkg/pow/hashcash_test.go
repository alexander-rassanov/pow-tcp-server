package pow

import (
	"reflect"
	"testing"
)

func TestHashCashData_ToString(t *testing.T) {
	type fields struct {
		Ver      int
		Bits     int
		Date     string
		Resource string
		Rand     int
		Counter  int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"standard hashcah",
			fields{
				1,
				1,
				"060408",
				"test@localhost",
				100,
				100,
			},
			"1:1:060408:test@localhost::MTAw:MTAw",	
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HashCashData{
				Ver:      tt.fields.Ver,
				Bits:     tt.fields.Bits,
				Date:     tt.fields.Date,
				Resource: tt.fields.Resource,
				Rand:     tt.fields.Rand,
				Counter:  tt.fields.Counter,
			}
			if got := h.ToString(); got != tt.want {
				t.Errorf("HashCashData.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashCashData_Sha1Hash(t *testing.T) {
	type fields struct {
		Ver      int
		Bits     int
		Date     string
		Resource string
		Rand     int
		Counter  int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"standard hashcah",
			fields{
				1,
				1,
				"060408",
				"test@localhost",
				100,
				100,
			},
			"7396a3708521232845c1ee56b44ea245a0240675",	
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HashCashData{
				Ver:      tt.fields.Ver,
				Bits:     tt.fields.Bits,
				Date:     tt.fields.Date,
				Resource: tt.fields.Resource,
				Rand:     tt.fields.Rand,
				Counter:  tt.fields.Counter,
			}
			if got := h.Sha1Hash(); got != tt.want {
				t.Errorf("HashCashData.Sha1Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashCashData_Resolve(t *testing.T) {
	type fields struct {
		Ver      int
		Bits     int
		Date     string
		Resource string
		Rand     int
		Counter  int
	}
	tests := []struct {
		name    string
		fields  fields
		want    HashCashData
		wantErr bool
	}{
		{
			"Want 1 zero",
			fields{
				1,
				1,
				"060408",
				"test@localhost",
				100,
				100,
			},
			HashCashData{
				1,
				1,
				"060408",
				"test@localhost",
				100,
				126,
			},
			false,
		},{
			"Want 2 zeroes",
			fields{
				1,
				2,
				"060408",
				"test@localhost",
				100,
				3500,
			},
			HashCashData{
				1,
				2,
				"060408",
				"test@localhost",
				100,
				3674,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HashCashData{
				Ver:      tt.fields.Ver,
				Bits:     tt.fields.Bits,
				Date:     tt.fields.Date,
				Resource: tt.fields.Resource,
				Rand:     tt.fields.Rand,
				Counter:  tt.fields.Counter,
			}
			got, err := h.Resolve()
			if (err != nil) != tt.wantErr {
				t.Errorf("HashCashData.Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HashCashData.Resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}
