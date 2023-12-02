package utils

import (
	"testing"
)

func Test_RunCommand(t *testing.T) {
	tests := []struct {
		name    string
		cmd     Command
		wantErr bool
	}{
		{
			name: "simple binary command",
			cmd: Command{
				Binary:   "echo",
				Commands: []string{"hello"},
			},
			wantErr: false,
		},
		{
			name: "binary command with multiple arguments",
			cmd: Command{
				Binary:   "ls",
				Commands: []string{"-l"},
			},
			wantErr: false,
		},
		{
			name: "unknown binary command",
			cmd: Command{
				Binary:   "unknown",
				Commands: []string{"-l"},
			},
			wantErr: true,
		},
		{
			name: "command with no binary specified",
			cmd: Command{
				Commands: []string{"ls", "-l"},
			},
			wantErr: false,
		},
		{
			name:    "empty command",
			cmd:     Command{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := RunCommand(tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunCommand() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("RunCommand() output = %v", out)
		})
	}
}
