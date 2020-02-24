package main

import (
	"reflect"
	"sort"
	"testing"
)

func Test_concurrentSortedList_insert1Lock(t *testing.T) {
	type args struct {
		v []int
	}
	tests := []struct {
		name string
		list *concurrentSortedList
		args args
	}{
		{
			name: "1 2 3",
			list: &concurrentSortedList{head: nil},
			args: args{v: []int{1, 2, 3}},
		},
		{
			name: "3 2 1",
			list: &concurrentSortedList{head: nil},
			args: args{v: []int{3, 2, 1}},
		},
		{
			name: "9,0,2,3,2,5,1",
			list: &concurrentSortedList{head: nil},
			args: args{v: []int{9, 0, 2, 3, 2, 5, 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.args.v {
				tt.list.insert1Lock(v)
			}
			testOut := tt.list.toList()
			temp := tt.args.v
			sort.Ints(temp)
			if !reflect.DeepEqual(testOut, temp) {
				t.Errorf("got:%v, expect:%v, intput:%v", testOut, temp, tt.args.v)
			}
		})
	}
}
