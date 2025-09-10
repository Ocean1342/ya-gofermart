package orderprocessor

import "testing"

func TestAccrualResponse_GetAccrual(t *testing.T) {
	type fields struct {
		Order   string
		Status  string
		Accrual float64
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			"test",
			fields{
				Order:   "123",
				Status:  "NEW",
				Accrual: 12.128937123897,
			},
			1213,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccrualResponse{
				Order:   tt.fields.Order,
				Status:  tt.fields.Status,
				Accrual: tt.fields.Accrual,
			}
			if got := a.GetAccrual(); got != tt.want {
				t.Errorf("GetAccrual() = %v, want %v", got, tt.want)
			}
		})
	}
}
