package module_test

import (
	"testing"

	mod "github.com/kubex-ecosystem/grompt/internal/module"
	"github.com/spf13/cobra"
)

func findCmd(root *cobra.Command, name string) *cobra.Command {
	if root.Use == name {
		return root
	}
	for _, c := range root.Commands() {
		if got := findCmd(c, name); got != nil {
			return got
		}
	}
	return nil
}

func TestFlagsFromRoot_AskGenerateChatStart(t *testing.T) {
	m := mod.RegX()
	root := m.Command()
	if root == nil {
		t.Fatalf("root command is nil")
	}

	// ask
	ask := findCmd(root, "ask")
	if ask == nil {
		t.Fatalf("ask command not found from root")
	}
	if f := ask.Flags().Lookup("max-tokens"); f == nil || f.DefValue != "1000" {
		t.Fatalf("ask max-tokens default = %v, want 1000", f)
	}
	if f := ask.Flags().Lookup("ollama-endpoint"); f == nil || f.DefValue != "http://localhost:11434" {
		t.Fatalf("ask ollama-endpoint default mismatch: %#v", f)
	}

	// generate
	gen := findCmd(root, "generate")
	if gen == nil {
		t.Fatalf("generate command not found from root")
	}
	if f := gen.Flags().Lookup("purpose-type"); f == nil || f.DefValue != "code" {
		t.Fatalf("generate purpose-type default = %v, want code", f)
	}
	if f := gen.Flags().Lookup("lang"); f == nil || f.DefValue != "english" {
		t.Fatalf("generate lang default = %v, want english", f)
	}
	if f := gen.Flags().Lookup("max-tokens"); f == nil || f.DefValue != "2048" {
		t.Fatalf("generate max-tokens default = %v, want 2048", f)
	}

	// chat
	chat := findCmd(root, "chat")
	if chat == nil {
		t.Fatalf("chat command not found from root")
	}
	if f := chat.Flags().Lookup("max-tokens"); f == nil || f.DefValue != "1000" {
		t.Fatalf("chat max-tokens default = %v, want 1000", f)
	}

	// server start
	start := findCmd(root, "start")
	if start == nil {
		t.Fatalf("server start command not found from root")
	}
	if f := start.Flags().Lookup("bind"); f == nil || f.DefValue != "localhost" {
		t.Fatalf("server bind default = %v, want localhost", f)
	}
	if f := start.Flags().Lookup("port"); f == nil || f.DefValue != "8080" {
		t.Fatalf("server port default = %v, want 8080", f)
	}
}
