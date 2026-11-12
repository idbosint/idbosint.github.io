// fetch_maps.go — descarga recortes de tiles vectoriales en formato PMTiles.
//
// Uso:
//   go run fetch_maps.go            # lee ./mapas.json y genera todos los mapas
//   go run fetch_maps.go --dry-run  # solo muestra los comandos
//
// Lee mapas.json, y por cada entrada ejecuta:
//
//	pmtiles extract <fuente> <salida> --bbox=<minLon,minLat,maxLon,maxLat> --maxzoom=<n>
//
// El archivo resultante se sirve tal cual desde static/tiles/ (sin CDN, sin
// servidor de tiles: el navegador hace peticiones Range al .pmtiles).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	Fuente  string `json:"fuente"`
	Maxzoom int    `json:"maxzoom"`
	Mapas   []Mapa `json:"mapas"`
}

type Mapa struct {
	Nombre string    `json:"nombre"`
	BBox   []float64 `json:"bbox"` // [minLon, minLat, maxLon, maxLat]
	Salida string    `json:"salida"`
}

func main() {
	configPath := flag.String("config", "mapas.json", "ruta al JSON de configuración")
	dryRun := flag.Bool("dry-run", false, "imprime los comandos sin ejecutarlos")
	flag.Parse()

	if _, err := exec.LookPath("pmtiles"); err != nil {
		fatal("pmtiles no está en el PATH: descarga el binario desde\n" +
			"  https://github.com/protomaps/go-pmtiles/releases")
	}

	raw, err := os.ReadFile(*configPath)
	if err != nil {
		fatal("no se puede leer %s: %v", *configPath, err)
	}

	var cfg Config
	if err := json.Unmarshal(raw, &cfg); err != nil {
		fatal("%s no es JSON válido: %v", *configPath, err)
	}
	if cfg.Fuente == "" {
		fatal("falta el campo \"fuente\" en %s", *configPath)
	}
	if cfg.Maxzoom == 0 {
		cfg.Maxzoom = 15
	}
	if len(cfg.Mapas) == 0 {
		fatal("no hay mapas definidos en %s", *configPath)
	}

	fmt.Printf("fuente : %s\n", cfg.Fuente)
	fmt.Printf("maxzoom: %d\n\n", cfg.Maxzoom)

	var fallos int
	for _, m := range cfg.Mapas {
		if err := m.validar(); err != nil {
			fmt.Fprintf(os.Stderr, "[%s] configuración inválida: %v\n", m.Nombre, err)
			fallos++
			continue
		}
		if err := extraer(cfg.Fuente, m, cfg.Maxzoom, *dryRun); err != nil {
			fmt.Fprintf(os.Stderr, "[%s] ERROR: %v\n", m.Nombre, err)
			fallos++
		}
	}

	if fallos > 0 {
		os.Exit(1)
	}
	fmt.Println("\nListo. Sirve los .pmtiles desde static/tiles/ (sin CDN).")
}

func (m Mapa) validar() error {
	if m.Nombre == "" {
		return fmt.Errorf("falta \"nombre\"")
	}
	if len(m.BBox) != 4 {
		return fmt.Errorf("bbox debe tener 4 valores [minLon,minLat,maxLon,maxLat]")
	}
	if m.BBox[0] >= m.BBox[2] || m.BBox[1] >= m.BBox[3] {
		return fmt.Errorf("bbox mal ordenado: %v", m.BBox)
	}
	if m.Salida == "" {
		return fmt.Errorf("falta \"salida\"")
	}
	return nil
}

func extraer(fuente string, m Mapa, maxzoom int, dryRun bool) error {
	bbox := fmt.Sprintf("%.2f,%.2f,%.2f,%.2f", m.BBox[0], m.BBox[1], m.BBox[2], m.BBox[3])

	args := []string{
		"extract", fuente, m.Salida,
		"--bbox=" + bbox,
		fmt.Sprintf("--maxzoom=%d", maxzoom),
	}

	fmt.Printf("[%s] pmtiles %s\n", m.Nombre, strings.Join(args, " "))
	if dryRun {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(m.Salida), 0o755); err != nil {
		return err
	}

	cmd := exec.Command("pmtiles", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pmtiles extract falló: %w", err)
	}

	info, err := os.Stat(m.Salida)
	if err != nil {
		return err
	}
	mb := float64(info.Size()) / (1024 * 1024)
	fmt.Printf("[%s] generado %s — %.1f MB\n", m.Nombre, m.Salida, mb)
	if mb >= 100 {
		fmt.Fprintf(os.Stderr, "[%s] AVISO: %.1f MB supera el límite de 100 MB\n"+
			"      por archivo de GitHub Pages. Baja --maxzoom o recorta el bbox.\n", m.Nombre, mb)
	}
	return nil
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
