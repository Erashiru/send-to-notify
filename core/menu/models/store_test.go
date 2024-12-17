package models

import (
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"reflect"
	"testing"
	"time"
)

func TestStoreDSMenus_GetActiveMenu(t *testing.T) {

	t1, err := time.Parse(time.RFC3339, "2022-10-19T16:48:08.087+00:00")
	if err != nil {
		t.Error()
	}

	t2, err := time.Parse(time.RFC3339, "2022-10-31T16:00:33.293+00:00")
	if err != nil {
		t.Error()
	}

	t3, err := time.Parse(time.RFC3339, "2023-03-02T11:39:24.412+00:00")
	if err != nil {
		t.Error()
	}

	t4, _ := time.Parse(time.RFC3339, "2023-04-01T16:00:33.293+00:00")

	tests := []struct {
		name string
		s    storeModels.StoreDSMenus
		arg  storeModels.AggregatorName
		want storeModels.StoreDSMenu
	}{
		{
			name: "#get glovo",
			arg:  storeModels.AggregatorName("glovo"),
			s: storeModels.StoreDSMenus{
				{
					ID:        "634fd5e87c41c0c0e35c46d0",
					Delivery:  "glovo",
					UpdatedAt: coreModels.FromTime(t1).Time,
				},
				{
					ID:        "635f9cc176e03cfb9aeb756c",
					Delivery:  "glovo",
					UpdatedAt: coreModels.FromTime(t2).Time,
				},
				{
					ID:        "64008aec595886e63c6584e8",
					Delivery:  "wolt",
					UpdatedAt: coreModels.FromTime(t3).Time,
				},
				{
					ID:        "634fd5e87c41c0c0e35c46d0",
					Delivery:  "glovo",
					UpdatedAt: coreModels.FromTime(t4).Time,
				},
			},
			want: storeModels.StoreDSMenu{
				ID:        "634fd5e87c41c0c0e35c46d0",
				Delivery:  "glovo",
				UpdatedAt: coreModels.FromTime(t4).Time,
			},
		},
		{
			name: "#get wolt",
			arg:  storeModels.AggregatorName("wolt"),
			s: storeModels.StoreDSMenus{
				{
					ID:        "634fd5e87c41c0c0e35c46d0",
					Delivery:  "glovo",
					UpdatedAt: coreModels.FromTime(t1).Time,
				},
				{
					ID:        "635f9cc176e03cfb9aeb756c",
					Delivery:  "glovo",
					UpdatedAt: coreModels.FromTime(t2).Time,
				},
				{
					ID:        "64008aec595886e63c6584e8",
					Delivery:  "wolt",
					UpdatedAt: coreModels.FromTime(t3).Time,
				},
				{
					ID:        "634fd5e87c41c0c0e35c46d0",
					Delivery:  "glovo",
					UpdatedAt: coreModels.FromTime(t4).Time,
				},
			},
			want: storeModels.StoreDSMenu{
				ID:        "64008aec595886e63c6584e8",
				Delivery:  "wolt",
				UpdatedAt: coreModels.FromTime(t3).Time,
			},
		},
		{
			name: "#emty aggregator name",
			arg:  storeModels.AggregatorName(""),
			s: storeModels.StoreDSMenus{
				{
					ID:        "634fd5e87c41c0c0e35c46d0",
					Delivery:  "glovo",
					UpdatedAt: coreModels.FromTime(t1).Time,
				},
				{
					ID:        "635f9cc176e03cfb9aeb756c",
					Delivery:  "glovo",
					UpdatedAt: coreModels.FromTime(t2).Time,
				},
				{
					ID:        "64008aec595886e63c6584e8",
					Delivery:  "wolt",
					UpdatedAt: coreModels.FromTime(t3).Time,
				},
			},
			want: storeModels.StoreDSMenu{},
		},
		{
			name: "#3",
			arg:  storeModels.AggregatorName(""),
			s:    storeModels.StoreDSMenus{},
			want: storeModels.StoreDSMenu{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetActiveMenu(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetActiveMenu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreDSMenus_IsUniqueMenu(t *testing.T) {
	tests := []struct {
		name string
		s    storeModels.StoreDSMenus
		arg  storeModels.AggregatorName
		want storeModels.StoreDSMenu
	}{
		{
			name: "#get glovo",
			arg:  storeModels.AggregatorName("glovo"),
			s: storeModels.StoreDSMenus{
				{
					ID:       "634fd5e87c41c0c0e35c46d0",
					Delivery: "glovo",
				},
				{
					ID:       "635f9cc176e03cfb9aeb756c",
					Delivery: "glovo",
				},
				{
					ID:       "64008aec595886e63c6584e8",
					Delivery: "wolt",
				},
				{
					ID:       "634fd5e87c41c0c0e35c46d0",
					Delivery: "glovo",
				},
			},
			want: storeModels.StoreDSMenu{
				ID:       "634fd5e87c41c0c0e35c46d0",
				Delivery: "glovo",
			},
		},
		{
			name: "#get wolt",
			arg:  storeModels.AggregatorName("wolt"),
			s: storeModels.StoreDSMenus{
				{
					ID:       "634fd5e87c41c0c0e35c46d0",
					Delivery: "glovo",
				},
				{
					ID:       "635f9cc176e03cfb9aeb756c",
					Delivery: "glovo",
				},
				{
					ID:       "64008aec595886e63c6584e8",
					Delivery: "wolt",
				},
				{
					ID:       "634fd5e87c41c0c0e35c46d0",
					Delivery: "glovo",
				},
			},
			want: storeModels.StoreDSMenu{
				ID:       "64008aec595886e63c6584e8",
				Delivery: "wolt",
			},
		},
		{
			name: "#emty aggregator name",
			arg:  storeModels.AggregatorName(""),
			s: storeModels.StoreDSMenus{
				{
					ID:       "634fd5e87c41c0c0e35c46d0",
					Delivery: "glovo",
				},
				{
					ID:       "635f9cc176e03cfb9aeb756c",
					Delivery: "glovo",
				},
				{
					ID:       "64008aec595886e63c6584e8",
					Delivery: "wolt",
				},
			},
			want: storeModels.StoreDSMenu{},
		},
		{
			name: "#3",
			arg:  storeModels.AggregatorName(""),
			s:    storeModels.StoreDSMenus{},
			want: storeModels.StoreDSMenu{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetActiveMenu(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetActiveMenu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreDSMenus_SetActiveMenu(t *testing.T) {
	t1, err := time.Parse(time.RFC3339, "2022-10-19T16:48:08.087+00:00")
	if err != nil {
		t.Error()
	}

	t2, err := time.Parse(time.RFC3339, "2022-10-31T16:00:33.293+00:00")
	if err != nil {
		t.Error()
	}

	t3, err := time.Parse(time.RFC3339, "2023-03-02T11:39:24.412+00:00")
	if err != nil {
		t.Error()
	}

	t4, _ := time.Parse(time.RFC3339, "2023-04-01T16:00:33.293+00:00")

	tests := []struct {
		name string
		s    storeModels.StoreDSMenus
		arg  storeModels.AggregatorName
		want storeModels.StoreDSMenu
	}{
		{
			name: "#get glovo",
			arg:  storeModels.AggregatorName("glovo"),
			s: storeModels.StoreDSMenus{
				{
					ID:        "634fd5e87c41c0c0e35c46d0",
					Delivery:  "glovo",
					UpdatedAt: coreModels.FromTime(t1).Time,
				},
				{
					ID:        "635f9cc176e03cfb9aeb756c",
					Delivery:  "glovo",
					UpdatedAt: coreModels.FromTime(t2).Time,
				},
				{
					ID:        "64008aec595886e63c6584e8",
					Delivery:  "wolt",
					UpdatedAt: coreModels.FromTime(t3).Time,
				},
				{
					ID:        "634fd5e87c41c0c0e35c46d0",
					Delivery:  "glovo",
					UpdatedAt: coreModels.FromTime(t4).Time,
				},
			},
			want: storeModels.StoreDSMenu{
				ID:        "634fd5e87c41c0c0e35c46d0",
				Delivery:  "glovo",
				UpdatedAt: coreModels.FromTime(t4).Time,
			},
		},
		{
			name: "#get wolt",
			arg:  storeModels.AggregatorName("wolt"),
			s: storeModels.StoreDSMenus{
				{
					ID:        "634fd5e87c41c0c0e35c46d0",
					Delivery:  "glovo",
					UpdatedAt: coreModels.FromTime(t1).Time,
				},
				{
					ID:        "635f9cc176e03cfb9aeb756c",
					Delivery:  "glovo",
					UpdatedAt: coreModels.FromTime(t2).Time,
				},
				{
					ID:        "64008aec595886e63c6584e8",
					Delivery:  "wolt",
					UpdatedAt: coreModels.FromTime(t3).Time,
				},
				{
					ID:        "634fd5e87c41c0c0e35c46d0",
					Delivery:  "glovo",
					UpdatedAt: coreModels.FromTime(t4).Time,
				},
			},
			want: storeModels.StoreDSMenu{
				ID:        "64008aec595886e63c6584e8",
				Delivery:  "wolt",
				UpdatedAt: coreModels.FromTime(t3).Time,
			},
		},
		{
			name: "#emty aggregator name",
			arg:  storeModels.AggregatorName(""),
			s: storeModels.StoreDSMenus{
				{
					ID:        "634fd5e87c41c0c0e35c46d0",
					Delivery:  "glovo",
					UpdatedAt: coreModels.FromTime(t1).Time,
				},
				{
					ID:        "635f9cc176e03cfb9aeb756c",
					Delivery:  "glovo",
					UpdatedAt: coreModels.FromTime(t2).Time,
				},
				{
					ID:        "64008aec595886e63c6584e8",
					Delivery:  "wolt",
					UpdatedAt: coreModels.FromTime(t3).Time,
				},
			},
			want: storeModels.StoreDSMenu{},
		},
		{
			name: "#3",
			arg:  storeModels.AggregatorName(""),
			s:    storeModels.StoreDSMenus{},
			want: storeModels.StoreDSMenu{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.SetActiveMenu()
		})
	}
}
