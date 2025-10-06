package info

import (
	"math/rand"
	"os"
	"strings"
)

type IPC struct {
	Type   string `json:"type"`
	Socket string `json:"socket"`
	Mode   string `json:"mode,omitempty"`
}

type Bitreg struct {
	BrfPath string `json:"brf_path"`
	NSBits  int    `json:"ns_bits"`
	Policy  string `json:"policy,omitempty"`

	// CapMask is a hexadecimal string representing the capability mask.
	CapMask string `json:"cap_mask,omitempty"`

	// StateHex is a hexadecimal string representing the state.
	StateHex string `json:"state_hex,omitempty"`
}

type KV struct {
	DeclareHashes []KeyHash `json:"declare_hashes,omitempty"`
	Values        []KVValue `json:"values,omitempty"`
	Encoding      string    `json:"encoding,omitempty"`
}

type KeyHash struct {
	KeyHash string `json:"key_hash"`
}

type KVValue struct {
	KeyHash string `json:"key_hash"`
	U64Hex  string `json:"u64_hex,omitempty"`
}

var banners = []string{
	`
  ______                                            __
 /      \                                          |  \
|  ▓▓▓▓▓▓\ ______   ______  ______ ____   ______  _| ▓▓_
| ▓▓ __\▓▓/      \ /      \|      \    \ /      \|   ▓▓ \
| ▓▓|    \  ▓▓▓▓▓▓\  ▓▓▓▓▓▓\ ▓▓▓▓▓▓\▓▓▓▓\  ▓▓▓▓▓▓\\▓▓▓▓▓▓
| ▓▓ \▓▓▓▓ ▓▓   \▓▓ ▓▓  | ▓▓ ▓▓ | ▓▓ | ▓▓ ▓▓  | ▓▓ | ▓▓ __
| ▓▓__| ▓▓ ▓▓     | ▓▓__/ ▓▓ ▓▓ | ▓▓ | ▓▓ ▓▓__/ ▓▓ | ▓▓|  \
 \▓▓    ▓▓ ▓▓      \▓▓    ▓▓ ▓▓ | ▓▓ | ▓▓ ▓▓    ▓▓  \▓▓  ▓▓
  \▓▓▓▓▓▓ \▓▓       \▓▓▓▓▓▓ \▓▓  \▓▓  \▓▓ ▓▓▓▓▓▓▓    \▓▓▓▓
                                        | ▓▓
                                        | ▓▓
                                         \▓▓
`,
}

func GetDescriptions(descriptionArg []string, hideBanner bool) map[string]string {
	var description, banner string

	if descriptionArg != nil {
		if strings.Contains(strings.Join(os.Args[0:], ""), "-h") {
			description = descriptionArg[0]
		} else {
			description = descriptionArg[1]
		}
	} else {
		description = ""
	}

	bannerRandLen := len(banners)
	bannerRandIndex := rand.Intn(bannerRandLen)
	banner = banners[bannerRandIndex]

	if hideBanner {
		return map[string]string{"description": description}
	}

	return map[string]string{"banner": banner, "description": description}
}
