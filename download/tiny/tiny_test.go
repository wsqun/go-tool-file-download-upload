package tiny

import (
	"os"
	"testing"
)

func TestClient_Download(t *testing.T) {
	type args struct {
		param FileParam
	}
	dir,_ := os.Getwd()
	println(dir)
	fileName := dir + "/demo.jpg"
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{param: FileParam{
				Url:      "http://cdn-ali-img-shstaticbz.shanhutech.cn/bizhi/staticwp/202103/d84979ac1cf2a88225077e71416f4b74--70509261.jpg",
				DestFile: fileName,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Client{}
			s.Download(tt.args.param)
		})
	}
}