package common

import "testing"

func TestMoneyFloatToInt(t *testing.T) {
	type args struct {
		in float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"case 1",
			args{in: 100.50},
			10050,
		},
		{
			"case 2",
			args{in: 100.51},
			10051,
		},
		{
			"case 3",
			args{in: 0.47},
			47,
		},
		{
			"case 4",
			args{in: 100},
			10000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MoneyFloatToInt(tt.args.in); got != tt.want {
				t.Errorf("MoneyFloatToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoneyIntToFloat(t *testing.T) {
	type args struct {
		in int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"case 1 729.28",
			args{in: 72928},
			729.28,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MoneyIntToFloat(tt.args.in); got != tt.want {
				t.Errorf("MoneyIntToFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}
