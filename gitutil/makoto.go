package gitutil

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/afeldman/Makoto/kpc"
	"github.com/afeldman/Makoto/makoto"
)

func RegisterInMakoto(target string) error {
	// Suche nach Manifest
	manifest := filepath.Join(target, "nishimura.kpc")
	if _, err := os.Stat(manifest); err != nil {
		// Fallback: evtl. <name>.kpc
		matches, _ := filepath.Glob(filepath.Join(target, "*.kpc"))
		if len(matches) == 0 {
			return nil // nichts zu registrieren
		}
		manifest = matches[0]
	}

	k, err := kpc.ReadKPCFile(manifest)
	if err != nil {
		return fmt.Errorf("failed to read kpc: %w", err)
	}

	// DB initialisieren (falls nicht schon offen)
	m := makoto.InitMakoto()
	m.DBInit()

	if err := makoto.Append(k); err != nil {
		return fmt.Errorf("failed to append to DB: %w", err)
	}

	fmt.Printf("📦 registered %s@%s in Makoto DB\n", k.Name, k.Version)
	return nil
}
