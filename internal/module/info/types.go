package info

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
