package server

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestCar_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		car     *Car
		args    args
		wantErr bool
	}{
		{"ValidCar", &Car{}, args{[]byte("{ \"id\": 4, \"seats\": 5 }")}, false},
		{"CarWithtId0", &Car{}, args{[]byte("{ \"id\": 0, \"seats\": 5 }")}, true},
		{"CarWithoutId", &Car{}, args{[]byte("{ \"seats\": 5 }")}, true},
		{"CarWithoutSeats", &Car{}, args{[]byte("{ \"id\": 3 }")}, true},
		{"CarWithLessThan1Seat", &Car{}, args{[]byte("{ \"id\": 4, \"seats\": 0 }")}, true},
		{"CarWithMoreThan6Seat", &Car{}, args{[]byte("{ \"id\": 4, \"seats\": 7 }")}, true},
		{"CarInvalidJson", &Car{}, args{[]byte("{ \"id\": 4 \"seats\": 5 ")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.car.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Car.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			} else if tt.wantErr {
				t.Logf(err.Error())
			}
		})
	}
}

func TestGroup_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		group   *Group
		args    args
		wantErr bool
	}{
		{"ValidGroup", &Group{}, args{[]byte("{ \"id\": 4, \"people\": 3 }")}, false},
		{"GroupWithoutId", &Group{}, args{[]byte("{ \"people\": 3 }")}, true},
		{"GroupWithoutPeople", &Group{}, args{[]byte("{ \"id\": 3 }")}, true},
		{"GroupWithLessThan1People", &Group{}, args{[]byte("{ \"id\": 4, \"people\": 0 }")}, true},
		{"GroupInvalidJson", &Group{}, args{[]byte("{ \"id\": 4, \"people\": 3 ")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.group.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Group.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			} else if tt.wantErr {
				t.Logf(err.Error())
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want *http.Server
	}{
		// TODO: Add test cases.
		{"NewServerOK", args{":9091"}, &http.Server{Addr: ":9091"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			got := New(tt.args.addr)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
			got.Shutdown(ctx)
		})
	}
}
