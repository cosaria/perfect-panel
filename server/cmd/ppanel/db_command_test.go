package ppanel

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestDBCommandTree(t *testing.T) {
	if !hasCommand(rootCmd, "db") {
		t.Fatalf("expected root command to include db")
	}
	for _, name := range []string{"bootstrap", "seed", "reset", "revisions"} {
		if !hasCommand(dbCmd, name) {
			t.Fatalf("expected db command to include %s", name)
		}
	}
}

func TestDBBootstrapRejectsUnknownRevisionSource(t *testing.T) {
	t.Cleanup(func() {
		_ = dbCmd.PersistentFlags().Set("source", "embedded")
	})
	if err := dbCmd.PersistentFlags().Set("source", "bogus"); err != nil {
		t.Fatalf("set source flag: %v", err)
	}

	err := dbBootstrapCmd.RunE(dbBootstrapCmd, nil)
	if err == nil {
		t.Fatal("expected bootstrap to reject unknown revision source")
	}
	if got := err.Error(); !strings.Contains(got, "unknown revision source") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func hasCommand(cmd interface{ Commands() []*cobra.Command }, name string) bool {
	for _, child := range cmd.Commands() {
		if child.Name() == name {
			return true
		}
	}
	return false
}
