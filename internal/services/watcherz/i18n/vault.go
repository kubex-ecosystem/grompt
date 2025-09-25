package wi18nast

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type vaultState struct {
	Items map[string]VaultItem `json:"items"`
}

type Vault interface {
	UpsertDraft(u Usage, suggestedKey string) (VaultItem, error)
	Approve(key string) error
	Rename(oldKey, newKey string) error
	List(status Status) ([]VaultItem, error)
	Save() error
	Stats() (total int, drafts int, approved int)
}

type JSONVault struct {
	path  string
	state vaultState
	mu    sync.Mutex
}

func NewJSONVault(dir string) (*JSONVault, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	v := &JSONVault{
		path: filepath.Join(dir, "i18n.vault.json"),
		state: vaultState{
			Items: map[string]VaultItem{},
		},
	}
	_ = v.load()
	return v, nil
}

func (v *JSONVault) UpsertDraft(u Usage, suggestedKey string) (VaultItem, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	k := normalizeKey(suggestedKey)
	now := time.Now()
	item, ok := v.state.Items[k]
	if !ok {
		item = VaultItem{
			Key:       k,
			Text:      extractGuessText(u),
			File:      u.FilePath,
			Line:      u.Line,
			Component: u.Component,
			Element:   inferElementFromJSX(u.JSXCtx),
			Status:    StatusDraft,
			FirstSeen: now,
			LastSeen:  now,
		}
	} else {
		item.LastSeen = now
		if item.Text == "" {
			item.Text = extractGuessText(u)
		}
	}
	if len(item.Contexts) < 5 && u.JSXCtx != "" {
		item.Contexts = append(item.Contexts, trimMax(u.JSXCtx, 120))
	}
	v.state.Items[k] = item
	return item, nil

}

func (v *JSONVault) Approve(key string) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	k := normalizeKey(key)
	it, ok := v.state.Items[k]
	if !ok {
		return errors.New("key não encontrada")
	}
	it.Status = StatusApproved
	v.state.Items[k] = it
	return nil
}

func (v *JSONVault) Rename(oldKey, newKey string) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	o := normalizeKey(oldKey)
	n := normalizeKey(newKey)
	it, ok := v.state.Items[o]
	if !ok {
		return errors.New("key antiga não encontrada")
	}
	delete(v.state.Items, o)
	it.Key = n
	v.state.Items[n] = it
	return nil
}

func (v *JSONVault) List(status Status) ([]VaultItem, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	var out []VaultItem
	for _, it := range v.state.Items {
		if status == "" || it.Status == status {
			out = append(out, it)
		}
	}
	return out, nil
}

func (v *JSONVault) Stats() (total int, drafts int, approved int) {
	v.mu.Lock()
	defer v.mu.Unlock()
	for _, it := range v.state.Items {
		total++
		switch it.Status {
		case StatusDraft:
			drafts++
		case StatusApproved:
			approved++
		}
	}
	return
}

func (v *JSONVault) Save() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	tmp := v.path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v.state); err != nil {
		_ = f.Close()
		return err
	}
	_ = f.Close()
	return os.Rename(tmp, v.path)
}

func (v *JSONVault) load() error {
	b, err := os.ReadFile(v.path)
	if err != nil {
		return nil // first run
	}
	var st vaultState
	if err := json.Unmarshal(b, &st); err != nil {
		return err
	}
	if st.Items == nil {
		st.Items = map[string]VaultItem{}
	}
	v.state = st
	return nil
}

func normalizeKey(k string) string {
	k = strings.TrimSpace(k)
	k = strings.Trim(k, ".")
	k = strings.ReplaceAll(k, "..", ".")
	return k
}

func trimMax(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
