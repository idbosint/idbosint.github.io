# Manual: desplegar en GitHub Pages

## 0. Tres cosas que hay que mirar antes

1. **El tema ya está vendorizado**: `themes/solite/` va en el repo como archivos
   normales (sin `.git`), así que no hace falta `submodules: true` en el
   workflow.

2. **El `.pmtiles` va versionado** (38 MB): es el mapa que consume el sitio y
   Pages no puede generarlo. Ya **no** está en `.gitignore`; comprueba que entra
   en el commit:

   ```bash
   git check-ignore static/tiles/spb.pmtiles && echo "¡ojo, ignorado!"
   ```

   38 MB entra de sobra: el límite de GitHub es **100 MB por archivo** y **1 GB
   por sitio**.

3. **Rutas relativas.** El sitio usa `baseURL = "/"` + `relativeURLs = true`, así
   que funciona en `usuario.github.io/repo/` sin tocar nada.

   **No compiles con `--baseURL https://usuario.github.io/repo/`**: Hugo
   duplicaría el prefijo en los recursos (`./repo/js/site.css`) y el CSS
   daría 404. Se compila con `/` (el de `hugo.toml`) y el workflow absolutiza
   después solo el sitemap y el RSS, que sí necesitan URLs completas.

---

## 1. Crear el repositorio y subir

```bash
cd proyecto-osint
git init -b main
git add -A
git commit -m "Sitio inicial"

# crea el repo vacío en GitHub (sin README) y:
git remote add origin git@github.com:TUUSUARIO/TUREPO.git
git push -u origin main
```

Comprueba que el tile va en el repo:

```bash
git ls-files static/tiles
# static/tiles/spb.pmtiles
```

---

## 1.1 Antes de subir

Comprueba que el `.pmtiles` entra en el commit y que no se cuelan dependencias:

```bash
git check-ignore static/tiles/spb.pmtiles && echo "¡ojo, ignorado!"
git status --short | grep -E 'node_modules|public/|resources/' || echo "limpio"
```

## 2. Workflow de GitHub Actions

Crea `.github/workflows/hugo.yml`:

```yaml
name: Deploy Hugo site to Pages

on:
  push:
    branches: [main]
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: false

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v7

      - name: Instalar Hugo (extended)
        run: |
          curl -fsSL -o hugo.tgz \
            https://github.com/gohugoio/hugo/releases/download/v0.169.0/hugo_extended_0.169.0_linux-amd64.tar.gz
          sudo tar -xzf hugo.tgz -C /usr/local/bin hugo
          hugo version

      - uses: actions/setup-node@v7
        with:
          node-version: 22
          cache: npm

      - name: Dependencias (postcss del tema)
        run: npm ci

      - uses: actions/configure-pages@v6
        id: pages

      - name: Build
        run: hugo --minify --gc

      - name: Absolutizar sitemap y RSS
        run: |
          BASE="${{ steps.pages.outputs.base_url }}"
          sed -i "s|<loc>/|<loc>${BASE}/|g" public/sitemap.xml
          [ -f public/index.xml ] && sed -i "s|<link>/|<link>${BASE}/|g" public/index.xml

      - uses: actions/upload-pages-artifact@v5
        with:
          path: ./public

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - id: deployment
        uses: actions/deploy-pages@v5
```

> Versiones comprobadas en noviembre de 2026: `checkout@v7`,
> `setup-node@v7`, `configure-pages@v6`, `upload-pages-artifact@v5`,
> `deploy-pages@v5`.

**Por qué hace falta Node**: el tema Solite compila el SCSS vía `postCSS` y
necesita `postcss-cli` + `autoprefixer` en `node_modules`. Sin `npm ci` el build
falla con «POSTCSS: ...».

---

## 3. Activar Pages

1. Repo → **Settings → Pages**.
2. *Build and deployment* → **Source: GitHub Actions**.
3. Push a `main` (o *Run workflow*). En ~1-3 minutos:
   `https://TUUSUARIO.github.io/TUREPO/`.

---

## 4. Verificar el despliegue

- La portada lista los posts.
- Abre DevTools → **Network** → filtra `pmtiles`.
  Debe responder **206 Partial Content**; si responde 200 completo o 404, el
  mapa saldrá gris.
- Comprueba que el archivo está publicado:
  `https://TUUSUARIO.github.io/TUREPO/tiles/spb.pmtiles`.

Si el mapa no carga:

| Causa | Cómo se ve | Solución |
|---|---|---|
| El `.pmtiles` no se subió | 404 en `/tiles/spb.pmtiles` | `git add -f`, commit y push |
| Supera 100 MB | El deploy falla o el archivo no aparece | Baja `maxzoom` en `mapas.json` y regenera con `go run fetch_maps.go` |
| Ruta mal escrita | 404 con otro nombre | Revisa el parámetro `tile` del shortcode |

---

## 5. Actualizar el sitio

```bash
# contenido nuevo o cambios
git add -A && git commit -m "nuevo post" && git push
```

El workflow reconstruye y publica. El historial de despliegues está en
**Actions → Deploy Hugo site to Pages**.

---

## 6. Alternativa sin Actions (rama `gh-pages`)

```bash
hugo --minify --gc        # y luego prefija sitemap/RSS como en el workflow

# publicar public/ en la rama gh-pages
cd public && git init -b gh-pages && git add -A \
  && git commit -m "deploy" \
  && git remote add origin git@github.com:TUUSUARIO/TUREPO.git \
  && git push -f origin gh-pages
```

Luego en **Settings → Pages** elige *Deploy from a branch* → `gh-pages` / `/`.

---

## 7. Notas sobre límites

- **100 MB por archivo**, 1 GB por sitio.
- El `.pmtiles` de San Petersburgo (z15) pesa **38 MB**: entra con margen.
- Si más adelante añades otra ciudad, o bajas `maxzoom` a 14, o usas un
  repositorio aparte para los tiles (y cambias `tile` por una URL).
- Pages no ejecuta código: todo tiene que estar pre-generado en `public/`
  (por eso el `.pmtiles` se versiona).
