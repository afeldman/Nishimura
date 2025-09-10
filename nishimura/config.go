package nishimura

import (
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
)

// NishimuraConfig beschreibt die globale Konfiguration / Verzeichnisse.
type NishimuraConfig struct {
	RootDir  string // Basisordner: ~/.nishimura
	ConfFile string // z.B. ~/.nishimura/config.toml
	CacheDir string // z.B. ~/.nishimura/src
}

// DefaultConfPath liefert den Default-Pfad zur Config-Datei (~/.nishimura/config.toml).
func DefaultConfPath() string {
	return filepath.Join(configdir.LocalConfig("nishimura"), "config.toml")
}

// InitNishimura initialisiert die Config-Struktur, erzeugt nötige Ordner falls nicht vorhanden.
func InitNishimura(confFile string) *NishimuraConfig {
	root := configdir.LocalConfig("nishimura")

	// Fallback: wenn kein ConfFile angegeben, nimm default
	if confFile == "" {
		confFile = filepath.Join(root, "config.toml")
	}

	cfg := &NishimuraConfig{
		RootDir:  root,
		ConfFile: confFile,
		CacheDir: filepath.Join(root, "src"),
	}

	// Ordner anlegen falls nötig
	_ = os.MkdirAll(cfg.RootDir, 0755)
	_ = os.MkdirAll(cfg.CacheDir, 0755)

	return cfg
}
