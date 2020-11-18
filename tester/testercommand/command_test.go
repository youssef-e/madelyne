package testercommand

import (
	"testing"
)

func TestRun(t *testing.T) {
	tests := []string{
		"ls",
		"ls ; sleep 0",
		"ls&",
		"ls& sleep 0;",
		"./../../global-coverage.sh",
	}

	for _, cmd := range tests {
		err := Run(cmd)
		if err != nil {
			t.Fatalf("failed cmd %s : %v", cmd, err)
		}
	}
}
