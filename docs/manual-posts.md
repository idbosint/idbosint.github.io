# Manual: crear posts, imágenes y datos

## 1. Un post nuevo

Crea `content/posts/mi-post.md`:

```markdown
---
title: "Título del post"
date: 2025-01-15T10:00:00+03:00
draft: false
summary: "Una línea que aparece en la lista de posts y en el RSS."
tags: ["osint", "spb"]
---

Primer párrafo. Se usa como resumen automático en la lista.

## Subtítulo

Texto normal, **negrita**, `código`, listas…
```

- `draft: true` → el post **no** se publica (`hugo server --buildDrafts` sí lo
  muestra para revisarlo).
- `tags` crea automáticamente `/tags/<tag>/`.
- El `summary` es lo que se ve en la portada; si no lo pones, Hugo corta el
  primer párrafo.

Comprobar:

```bash
hugo server --buildDrafts   # http://localhost:1313
```

---

## 2. Imágenes

**No enlaces imágenes externas** (en I2P no cargan y rompe el «cero CDN»).
Copia el archivo dentro del proyecto. Dos formas:

### A) Page bundle (recomendado)

```
content/posts/mi-post/
├── index.md
├── plaza.jpg
└── mapa-detalle.png
```

En el markdown, ruta **relativa**:

```markdown
![Plaza de la estación](plaza.jpg)
```

Hugo publica la imagen en `/posts/mi-post/plaza.jpg`, así que funciona igual en
la raíz del dominio y en `github.io/repo/`.

### B) Carpeta compartida

```
static/images/plaza.jpg
```

```markdown
![Plaza de la estación](images/plaza.jpg)
```

Hugo reescribe la ruta a relativa al generar el HTML.

### Con pie de foto

```markdown
{{< figure src="plaza.jpg" caption="Nevsky Prospekt, 2024" >}}
```

### Tamaño

I2P y las conexiones móviles son lentas. Recomendado:

- JPEG/WebP de **200–400 KB** como máximo.
- Ancho 1000–1400 px.

```bash
# redimensionar y recomprimir (ImageMagick)
mogrify -resize '1400x>' -quality 82 plaza.jpg
# o pasar a webp
cwebp -q 80 plaza.jpg -o plaza.webp
```

---

## 3. El mapa en un post

```markdown
{{< mapa
    id="spb-2025"
    tile="tiles/spb.pmtiles"
    incidentes="data/incidentes-spb.geojson"
    lat="59.9343"
    lon="30.3351"
    zoom="12"
    flavor="dark"
>}}
```

| Parámetro | Por defecto | Notas |
|---|---|---|
| `id` | `mapa-spb` | **único** por página si pones varios mapas |
| `tile` | `tiles/spb.pmtiles` | ruta dentro de `static/` |
| `incidentes` | `data/incidentes-spb.geojson` | ruta dentro de `static/` |
| `lat` / `lon` / `zoom` | `59.9343` / `30.3351` / `12` | vista inicial |
| `flavor` | `dark` | `dark`, `light`, `white`, `grayscale`, `black` |

El contador y el `<noscript>` se calculan en build-time a partir del GeoJSON.

---

## 4. Datos de incidentes

`static/data/incidentes-spb.geojson`:

```json
{
  "type": "Feature",
  "geometry": { "type": "Point", "coordinates": [30.3351, 59.9311] },
  "properties": {
    "titulo": "Stalking on Nevsky Prospekt",
    "fecha": "14.03.2024",
    "severidad": "alta",
    "descripcion": "Texto libre."
  }
}
```

- `coordinates`: `[longitud, latitud]` (ojo al orden).
- `severidad`: `alta` | `media` | `baja` → la UI los muestra como
  *high / medium / low* (rojo / naranja / verde).
- `fecha`: `DD.MM.YYYY` (se muestra tal cual).
- El esquema en español es intencional (es el contrato del proyecto).

---

## 5. Otro mapa (otra ciudad o recorte)

1. Añade una entrada a `mapas.json`:

```json
{
  "nombre": "moscu",
  "bbox": [37.30, 55.55, 37.90, 55.95],
  "salida": "static/tiles/moscu.pmtiles"
}
```

2. Genera el archivo:

```bash
go run fetch_maps.go
```

3. Úsalo en el post: `tile="tiles/moscu.pmtiles"`.

Recuerda que el `.pmtiles` pesa: a z15 una ciudad grande puede pasar de 100 MB
(límite de GitHub Pages). Baja `maxzoom` en `mapas.json` si hace falta.

---

## 6. Post bilingüe (como el de ejemplo)

```markdown
{{< idiomas defecto="en" >}}

{{< idioma en "English" >}}
Texto en inglés…
{{< /idioma >}}

{{< idioma ru "Русский" >}}
Текст по-русски…
{{< /idioma >}}

{{< /idiomas >}}
```

Sin JavaScript se ven los dos bloques seguidos; con JS se convierten en
pestañas.

---

## 7. Publicar

```bash
hugo server --buildDrafts         # revisar en local
git add -A && git commit -m "nuevo post"
git push                          # GitHub Actions despliega solo
```
