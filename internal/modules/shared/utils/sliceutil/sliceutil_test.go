package sliceutil

import (
	"reflect"
	"testing"
)

func TestSliceUtils(t *testing.T) {
	t.Run("Every", func(t *testing.T) {
		tests := []struct {
			name      string
			slice     []int
			predicate func(value int) bool
			want      bool
		}{
			{
				name:  "all elements satisfy predicate",
				slice: []int{2, 4, 6, 8},
				predicate: func(value int) bool {
					return value%2 == 0
				},
				want: true,
			},
			{
				name:  "not all elements satisfy predicate",
				slice: []int{2, 3, 4, 6},
				predicate: func(value int) bool {
					return value%2 == 0
				},
				want: false,
			},
			{
				name:  "no elements satisfy predicate",
				slice: []int{1, 3, 5},
				predicate: func(value int) bool {
					return value%2 == 0
				},
				want: false,
			},
			{
				name:  "empty slice",
				slice: []int{},
				predicate: func(value int) bool {
					return value > 0
				},
				want: true,
			},
			{
				name:  "nil slice",
				slice: nil,
				predicate: func(value int) bool {
					return value > 0
				},
				want: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := Every(tt.slice, tt.predicate)

				if got != tt.want {
					t.Errorf("Every() = %t, want %t", got, tt.want)
				}
			})
		}
	})

	t.Run("Filter", func(t *testing.T) {
		tests := []struct {
			name      string
			slice     []int
			predicate func(value int) bool
			want      []int
		}{
			{
				name:  "filter matching elements",
				slice: []int{1, 2, 3, 4, 5, 6},
				predicate: func(value int) bool {
					return value%2 == 0
				},
				want: []int{2, 4, 6},
			},
			{
				name:  "all elements match",
				slice: []int{2, 4, 6},
				predicate: func(value int) bool {
					return value%2 == 0
				},
				want: []int{2, 4, 6},
			},
			{
				name:  "no elements match",
				slice: []int{1, 3, 5},
				predicate: func(value int) bool {
					return value%2 == 0
				},
				want: []int{},
			},
			{
				name:  "empty slice",
				slice: []int{},
				predicate: func(value int) bool {
					return value > 0
				},
				want: []int{},
			},
			{
				name:  "nil slice",
				slice: nil,
				predicate: func(value int) bool {
					return value > 0
				},
				want: []int{},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := Filter(tt.slice, tt.predicate)

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Filter() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}
