package config

import (
	"reflect"
	"testing"
)

func TestGetDirsCfg(t *testing.T) {
	tests := []struct {
		name string
		want *Dirs
	}{
		{
			name: "Should return the default dirs config",
			want: &Dirs{
				TerraDaggerDir:       terraDaggerDir,
				TerraDaggerCacheDir:  terraDaggerCacheDir,
				TerraDaggerExportDir: terraDaggerExportDir,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDirsCfg(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDirs() = %v, want %v", got, tt.want)
			}
		})
	}
}
