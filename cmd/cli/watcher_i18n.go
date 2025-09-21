package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	i18nwatcher "github.com/kubex-ecosystem/grompt/internal/services/watcherz/i18n"
)

func Watcheri18nCmd() {
	if len(os.Args) < 2 {
		log.Fatalf("uso: %s <caminho-do-projeto>", filepath.Base(os.Args[0]))
	}
	root := os.Args[1]

	vault, err := i18nwatcher.NewJSONVault(filepath.Join(root, "i18n"))
	if err != nil {
		log.Fatal(err)
	}

	w, err := i18nwatcher.NewWatcher(root, func(file string, usages []i18nwatcher.Usage) {
		fmt.Printf("n[%s]n", file)
		if usages == nil {
			fmt.Println("  (removido)")
			return
		}
		for _, u := range usages {
			key := i18nwatcher.GenKey(u)
			item, _ := vault.UpsertDraft(u, key)
			fmt.Printf("  %-14s key=%q line=%d comp=%sn",
				u.CallType, item.Key, u.Line, u.Component)
		}
		if err := vault.Save(); err != nil {
			log.Printf("vault save: %v", err)
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	w.Start()
	defer w.Stop()
	select {} // ^C para sair

}
