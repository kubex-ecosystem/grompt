package module_test

import (
    "strings"
    "testing"

    mod "github.com/kubex-ecosystem/grompt/internal/module"
)

func TestModule_CommandAnnotationsAndExamples(t *testing.T) {
    m := mod.RegX()
    m.SetParentCmdName("grompt")
    cmd := m.Command()
    if cmd == nil { t.Fatalf("Command() returned nil") }

    if cmd.Annotations == nil || cmd.Annotations["description"] == "" {
        t.Fatalf("command annotations missing description")
    }
    if !strings.Contains(cmd.Example, "grompt ") {
        t.Fatalf("examples should be prefixed with parent cmd name; got %q", cmd.Example)
    }
}

func TestModule_Basics(t *testing.T) {
    m := mod.RegX()
    if !m.Active() { t.Fatalf("module should be active") }
    if m.Module() != "grompt" { t.Fatalf("Module() = %q, want 'grompt'", m.Module()) }
    if m.ShortDescription() == "" || m.LongDescription() == "" { t.Fatalf("descriptions must not be empty") }
    if m.Usage() == "" { t.Fatalf("Usage must not be empty") }
}

