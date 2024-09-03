package tokenizer

import (
	"github.com/ahmadrezamusthafa/search-engine/internal/structs"
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name    string
		content structs.Content
		want    []string
	}{
		{
			name: "remove unused characters",
			content: structs.Content{
				String: "FT24245L5RRD TRF Dari - 451 - KHOERIYAH APENDI",
			},
			want: []string{"ft24245l5rrd", "451", "khoeriyah", "apendi"},
		},
		{
			name: "remove unused characters 2",
			content: structs.Content{
				String: "Journal no: 946874 TRANSFER DARI DANDUNG CATUR SUGENG PURWADI",
			},
			want: []string{"946874", "dandung", "catur", "sugeng", "purwadi"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Tokenize(tt.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}
