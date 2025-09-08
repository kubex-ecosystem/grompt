package cli_test

import (
    "testing"

    cli "github.com/kubex-ecosystem/grompt/cmd/cli"
    "github.com/spf13/cobra"
)

func getCmdByUse(cmds []*cobra.Command, use string) *cobra.Command { // shim for readability
    for _, c := range cmds { if c.Use == use { return c } }
    return nil
}

func TestAskCommand_FlagsAndDefaults(t *testing.T) {
    var ask *cobra.Command
    for _, c := range cli.AICmdList() { if c.Use == "ask" { ask = c; break } }
    if ask == nil { t.Fatalf("ask command not found") }

    cases := []struct{ name, def string }{
        {"prompt", ""},
        {"provider", ""},
        {"model", ""},
        {"max-tokens", "1000"},
        {"config", ""},
        {"apikey", ""},
        {"ollama-endpoint", "http://localhost:11434"},
    }
    for _, c := range cases {
        f := ask.Flags().Lookup(c.name)
        if f == nil { t.Fatalf("flag %q not found", c.name) }
        if f.DefValue != c.def { t.Fatalf("flag %q default = %q, want %q", c.name, f.DefValue, c.def) }
    }
    // debug is boolean false by default
    if f := ask.Flags().Lookup("debug"); f == nil || f.DefValue != "false" {
        t.Fatalf("flag debug missing or default != false: %#v", f)
    }
}

func TestGenerateCommand_FlagsAndDefaults(t *testing.T) {
    var gen *cobra.Command
    for _, c := range cli.AICmdList() { if c.Use == "generate" { gen = c; break } }
    if gen == nil { t.Fatalf("generate command not found") }

    cases := []struct{ name, def string }{
        {"ideas", "[]"}, // pflag StringSlice default renders as []
        {"purpose", ""},
        {"purpose-type", "code"},
        {"lang", "english"},
        {"max-tokens", "2048"},
        {"provider", ""},
        {"model", ""},
        {"config", ""},
        {"output", ""},
        {"apikey", ""},
        {"ollama-endpoint", "http://localhost:11434"},
    }
    for _, c := range cases {
        f := gen.Flags().Lookup(c.name)
        if f == nil { t.Fatalf("flag %q not found", c.name) }
        if f.DefValue != c.def { t.Fatalf("flag %q default = %q, want %q", c.name, f.DefValue, c.def) }
    }
    if f := gen.Flags().Lookup("debug"); f == nil || f.DefValue != "false" {
        t.Fatalf("flag debug missing or default != false: %#v", f)
    }
}

func TestChatCommand_FlagsAndDefaults(t *testing.T) {
    var chat *cobra.Command
    for _, c := range cli.AICmdList() { if c.Use == "chat" { chat = c; break } }
    if chat == nil { t.Fatalf("chat command not found") }

    cases := []struct{ name, def string }{
        {"provider", ""},
        {"model", ""},
        {"max-tokens", "1000"},
        {"config", ""},
        {"apikey", ""},
        {"ollama-endpoint", "http://localhost:11434"},
    }
    for _, c := range cases {
        f := chat.Flags().Lookup(c.name)
        if f == nil { t.Fatalf("flag %q not found", c.name) }
        if f.DefValue != c.def { t.Fatalf("flag %q default = %q, want %q", c.name, f.DefValue, c.def) }
    }
    if f := chat.Flags().Lookup("debug"); f == nil || f.DefValue != "false" {
        t.Fatalf("flag debug missing or default != false: %#v", f)
    }
}

func TestStartServerCommand_FlagsAndDefaults(t *testing.T) {
    var start *cobra.Command
    for _, c := range cli.ServerCmdList() { if c.Use == "start" { start = c; break } }
    if start == nil { t.Fatalf("start server command not found") }

    cases := []struct{ name, def string }{
        {"bind", "localhost"},
        {"port", "8080"},
        {"config", ""},
        {"openai-key", ""},
        {"deepseek-key", ""},
        {"ollama-endpoint", "http://localhost:11434"},
        {"gemini-key", ""},
        {"chatgpt-key", ""},
        {"claude-key", ""},
    }
    for _, c := range cases {
        f := start.Flags().Lookup(c.name)
        if f == nil { t.Fatalf("flag %q not found", c.name) }
        if f.DefValue != c.def { t.Fatalf("flag %q default = %q, want %q", c.name, f.DefValue, c.def) }
    }
    // boolean flags
    if f := start.Flags().Lookup("debug"); f == nil || f.DefValue != "false" {
        t.Fatalf("flag debug missing or default != false: %#v", f)
    }
    if f := start.Flags().Lookup("background"); f == nil || f.DefValue != "false" {
        t.Fatalf("flag background missing or default != false: %#v", f)
    }
}
