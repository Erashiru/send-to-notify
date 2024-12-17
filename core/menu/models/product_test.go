package models

import (
	"reflect"
	"testing"
)

func TestProducts_Unique(t *testing.T) {
	tests := []struct {
		name string
		p    Products
		want Products
	}{
		{
			name: "test",
			p: Products{
				{
					ExtID: "123",
				},
				{
					ExtID: "23",
				},
				{
					ExtID: "123",
				},
				{
					ExtID: "256",
				},
			},
			want: Products{
				{
					ExtID: "123",
				},
				{
					ExtID: "23",
				},
				{
					ExtID: "256",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Unique(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}
