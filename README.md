# SPb Incidents — OSINT

Blog estático para documentar incidentes sobre un mapa vectorial servido desde
el propio sitio: **Leaflet 1.9 + protomaps-leaflet + un archivo `.pmtiles`**
recortado a la zona de interés (el de ejemplo cubre San Petersburgo). Sin CDNs
ni teselas externas en tiempo de ejecución.

Los posts pueden publicarse en **inglés y ruso** con pestañas; el resto del
sitio está en inglés.

> El post de ejemplo es una **plantilla**: sustituye su texto y el GeoJSON por
> los datos reales antes de publicar.

---

## Stack

| Pieza | Versión / nota |
|---|---|
| Hugo | 0.169.0 **extended** (el tema compila SCSS) |
| Tema | [Solite](https://git.sr.ht/~mlthlschr/solite) — incluido en `themes/` (MIT) |
| Node | 22 (postcss + autoprefixer para el SCSS del tema) |
| Mapa | Leaflet 1.9.4 + protomaps-leaflet 5.1.0 (vendorizados en `static/js/`) |
| Teselas | `static/tiles/spb.pmtiles` (38 MB, z0–z15) |
| Datos | `static/data/incidentes-spb.geojson` |

Todo el JS/CSS es local. Rutas relativas (`baseURL = "/"` + `relativeURLs`), así
que el sitio funciona en la raíz de un dominio, en `usuario.github.io/repo/` y
en un eepsite `.i2p`.

---

## Estructura

```
.
├── content/
│   ├── _index.md               # portada = lista de posts
│   ├── about/_index.md
│   └── posts/
│       ├── _index.md
│       └── report-ch.md        # post de ejemplo
├── layouts/
│   ├── _default/baseof.html    # <title>, bloque "scripts", locale
│   ├── _default/summary.html   # tarjeta de post con .RelPermalink
│   ├── index.html              # lista de posts paginada
│   ├── partials/header.html    # enlace home relativo
│   ├── partials/footer.html    # barra inferior con About
│   ├── rss.xml
│   └── shortcodes/
│       ├── mapa.html           # Leaflet + PMTiles + GeoJSON
│       ├── idiomas.html        # contenedor de pestañas
│       └── idioma.html         # bloque de idioma
├── static/
│   ├── data/incidentes-spb.geojson
│   ├── js/                     # leaflet, protomaps-leaflet, css propios
│   ├── tiles/spb.pmtiles
│   └── favicon.svg
├── themes/solite/              # tema vendorizado (sin .git)
├── docs/                       # manuales
├── hugo.toml
├── mapas.json                  # fuente de teselas + bbox
├── fetch_maps.go               # genera el .pmtiles con pmtiles extract
├── package.json / package-lock.json
└── postcss.config.js
```

---

## Uso local

Requisitos: **Hugo extended 0.169+** y **Node 22**.

```bash
npm ci
hugo server --buildDrafts      # http://localhost:1313
```

Build de producción:

```bash
hugo --minify --gc             # genera public/
```

---

## Contenido

- **Crear un post, imágenes y shortcodes**: [`docs/manual-posts.md`](docs/manual-posts.md)
- El mapa se inserta con:

  ```markdown
  {{< mapa id="spb-2025" tile="tiles/spb.pmtiles"
           incidentes="data/incidentes-spb.geojson"
           lat="59.9343" lon="30.3351" zoom="12" flavor="dark" >}}
  ```

- **Datos**: `static/data/incidentes-spb.geojson` (`titulo`, `fecha`
  `DD.MM.YYYY`, `severidad` `alta|media|baja`, `descripcion`).

---

## Teselas (`.pmtiles`)

El archivo ya está en el repo. Para regenerarlo o cambiar de ciudad:

```bash
# requiere el binario `pmtiles` en el PATH
# https://github.com/protomaps/go-pmtiles/releases
go run fetch_maps.go
```

`mapas.json` define la fuente, el bbox y `maxzoom`. Límite de GitHub Pages:
**100 MB por archivo** (el de SPb son 38 MB).

---

## Despliegue

- **GitHub Pages**: [`docs/manual-github-pages.md`](docs/manual-github-pages.md)
  El workflow `.github/workflows/deploy.yml` construye y publica en cada push a
  `main`. Basta con activar *Settings → Pages → Source: GitHub Actions*.
- **I2P (eepsite)**: [`docs/manual-i2p.md`](docs/manual-i2p.md)

---

## Créditos y licencias

- Tema **Solite**, MIT — © mlthlschr.
- Teselas vectoriales: © [Protomaps](https://protomaps.com) / datos de
  [OpenStreetMap](https://openstreetmap.org/copyright) (ODbL).
- Los datos de incidentes son ficticios, creados como plantilla.
