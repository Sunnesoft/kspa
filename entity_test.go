package kspa

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGenerateRandomEntitiesJson(t *testing.T) {
	type args struct {
		path      string
		count     int
		removeOld bool
		c         RandomEntitySeqInfo
	}

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "generate_data_for_benchmark",
			args: args{
				path:      "./examples/v5000_e20000",
				count:     10,
				removeOld: true,
				c: RandomEntitySeqInfo{
					VertexCount:     5000,
					VertexStdFactor: 50,
					EdgesCount:      20000,
					RelationMin:     0.0,
					RelationMax:     100000.0,
					NoiseMean:       0.0,
					NoiseStdDev:     0.0001,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenerateRandomEntitiesJson(tt.args.path, tt.args.count, tt.args.removeOld, tt.args.c)

			files, _ := ioutil.ReadDir(tt.args.path)
			filesCount := len(files)

			if tt.args.count != filesCount {
				t.Errorf("Dfs.GenerateRandomEntitiesJson() create %v files, want %v", filesCount, tt.args.count)
			}

			os.RemoveAll(tt.args.path)
		})
	}
}
