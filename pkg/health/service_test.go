package health

import (
	"context"
	"reflect"
	"testing"
)

func Test_service_GetLiveness(t *testing.T) {
	type fields struct {
		services []canBeReady
	}
	type args struct {
		ctx context.Context
		req *GetLivenessRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *GetLivenessRespone
		wantErr bool
	}{
		{
			"get",
			fields{},
			args{
				context.Background(),
				&GetLivenessRequest{},
			},
			&GetLivenessRespone{Status: "ok"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				services: tt.fields.services,
			}
			got, err := s.GetLiveness(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetLiveness() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetLiveness() = %v, want %v", got, tt.want)
			}
		})
	}
}
