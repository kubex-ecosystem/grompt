package bitflags

import (
	"fmt"
)

// Example of composing middleware decisions from flags.

func ExampleMwDecision() {
	var sec SecFlag
	var reg FlagReg32[SecFlag]
	reg.Store(0)
	reg.Set(SecAuth | SecSanitize)
	sec = reg.Load()

	chain := []string{"trace", "logging"}
	if sec&SecAuth != 0 {
		chain = append(chain, "authentication")
	}
	if sec&SecSanitize != 0 {
		chain = append(chain, "sanitize")
	}
	fmt.Println(chain)
	// Output: [trace logging authentication sanitize]
}
