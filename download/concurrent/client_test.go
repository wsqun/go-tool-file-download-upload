package concurrent

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestClient_Download(t *testing.T) {
	type fields struct {
		httpReq http.Client
	}
	type args struct {
		ctx   context.Context
		param DownloadParam
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "",
			fields:  fields{},
			args:    args{
				ctx:   context.TODO(),
				param: DownloadParam{
					Url: "http://localhost:8080/file-server",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient()
			if err := c.Download(tt.args.ctx, tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}