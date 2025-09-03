package cli_test

import (
    "os"
    "testing"

    cli "github.com/rafa-mori/grompt/cmd/cli"
)

func TestGetDescriptions_BannerAndDescription(t *testing.T) {
    orig := os.Args
    os.Args = []string{"grompt", "-h"}
    t.Cleanup(func(){ os.Args = orig })

    m := cli.GetDescriptions([]string{"long description", "short description"}, true)
    if m == nil {
        t.Fatalf("expected map, got nil")
    }
    if m["description"] != "long description" {
        t.Fatalf("want long description with -h, got %q", m["description"])
    }
    if m["banner"] == "" {
        t.Fatalf("banner should not be empty")
    }
}

func TestAICmdList_ContainsExpectedCommands(t *testing.T) {
    cmds := cli.AICmdList()
    names := map[string]bool{}
    for _, c := range cmds { names[c.Use] = true }
    for _, want := range []string{"ask", "generate", "chat"} {
        if !names[want] { t.Fatalf("command %q not found", want) }
    }
}

func TestServerCmdList_StartPresent(t *testing.T) {
    cmds := cli.ServerCmdList()
    found := false
    for _, c := range cmds { if c.Use == "start" { found = true; break } }
    if !found { t.Fatalf("start server command not found") }
}

