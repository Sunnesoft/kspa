package kspa

import "testing"

func TestEdgesCsvToJson(t *testing.T) {
	type args struct {
		csvFp  string
		jsonFp string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "csv-2-json",
			args: args{
				csvFp:  "./examples/data.txt",
				jsonFp: "./examples/data.json",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EdgesCsvToJson(tt.args.csvFp, tt.args.jsonFp); (err != nil) != tt.wantErr {
				t.Errorf("EdgesCsvToJson() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
