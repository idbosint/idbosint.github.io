# Manual: publicar el sitio en I2P (eepsite)

I2P sirve HTTP normal, así que publicar aquí es: **servir `public/` + abrir un
túnel de servidor en i2pd**. El sitio ya cumple las dos condiciones que importan:

- **Rutas relativas** → funciona en la raíz del `.i2p` sin tocar `baseURL`.
- **Cero recursos externos** → no llama a CDNs, fuentes ni teselas de la
  clearnet (solo hay enlaces de atribución, que son texto).

---

## 1. Servir el sitio

Genera el sitio y sírvelo con un servidor estático que soporte **HTTP Range**
(nginx o caddy; `python3 -m http.server` no sirve):

```bash
hugo --minify --gc
```

Ejemplo con nginx, sirviendo en `127.0.0.1:8080`:

```nginx
server {
    listen 8080;
    root /ruta/al/proyecto/public;
    location / { try_files $uri $uri/ $uri/index.html =404; }
    location ~* \.pmtiles$ {
        default_type application/octet-stream;
        add_header Accept-Ranges bytes;
    }
}
```

> Usa nginx, **no** `hugo server`: el servidor de desarrollo inyecta livereload y
> websockets que en I2P dan problemas. Y **no** uses `python3 -m http.server`:
> no soporta HTTP Range y el `.pmtiles` no se podría leer.

Comprobación rápida:

```bash
curl -I http://127.0.0.1:8080/
curl -r 0-1023 -o /dev/null -w '%{http_code}\n' http://127.0.0.1:8080/tiles/spb.pmtiles
# esperado: 206
```

---

## 2. Instalar i2pd (Arch Linux)

```bash
sudo pacman -S i2pd
sudo systemctl enable --now i2pd
```

| Qué | Dónde |
|---|---|
| Configuración | `/etc/i2pd/i2pd.conf` |
| Túneles | `/etc/i2pd/tunnels.conf` |
| Datos y claves | `/var/lib/i2pd/` |
| Consola web | `http://127.0.0.1:7070` |
| Proxy HTTP | `127.0.0.1:4444` |

*(Si prefieres i2pd en Docker, el datadir es el volumen del contenedor,
`tunnels.conf` se lee desde ahí y la consola sigue en `127.0.0.1:7070`.)*

---

## 3. Crear el túnel de servidor

Añade al final de `/etc/i2pd/tunnels.conf`:

```ini
[proyecto-osint]
type = http
host = 127.0.0.1
port = 8080
keys = proyecto-osint.dat
```

- `keys` es la identidad del eepsite. Se genera sola la primera vez en
  `/var/lib/i2pd/proyecto-osint.dat`. **Consérvala**: si se pierde, cambia tu
  dirección.
- Si i2pd corre en **otro contenedor Docker** (no en el host), apunta al gateway
  del bridge de Docker en lugar de `127.0.0.1` (normalmente `172.17.0.1`):

```ini
host = 172.17.0.1
port = 8080
```

(comprobado: desde un contenedor, `http://172.17.0.1:8080/` responde 200).

Reinicia:

```bash
sudo systemctl restart i2pd
sudo journalctl -u i2pd -n 30 --no-pager
```

---

## 4. Obtener la dirección

Abre `http://127.0.0.1:7070` → pestaña **I2P tunnels** → tu túnel
`proyecto-osint`. Verás:

- **Destination** (base64, 516 caracteres) — es la identidad real.
- **B32 address**: `xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.b32.i2p`
  (52 caracteres).

Esa b32 es tu sitio. Cópiala; no cambia mientras conserves `proyecto-osint.dat`.

---

## 5. Nombre legible (opcional)

**Solo para ti** (address book privada): añade a `/var/lib/i2pd/hosts.txt`
(o `~/.i2pd/hosts.txt` si i2pd corre como tu usuario):

```
proyecto-osint.i2p=<DESTINATION_BASE64>
```

**Público**: registra el nombre en un address book de I2P (por ejemplo
`http://stats.i2p/i2p/addkey.html`): subes el nombre y el *destination* y, tras
un tiempo, otros routers podrán resolver `proyecto-osint.i2p`.

---

## 6. Probar

Con el proxy HTTP de i2pd:

```bash
curl -x http://127.0.0.1:4444 -sI http://<B32>.b32.i2p/ | head
curl -x http://127.0.0.1:4444 -s http://<B32>.b32.i2p/ | grep -o '<title>[^<]*'
```

En el navegador: usa el **I2P Browser** o Firefox con proxy manual
`HTTP 127.0.0.1:4444` y *no usar proxy* para `localhost, 127.0.0.1`.

Comprueba que el mapa carga y que en la consola de red las peticiones al
`.pmtiles` son `206`.

---

## 7. Notas de rendimiento y seguridad

- **Tamaño**: 38 MB de `.pmtiles` por túneles I2P se nota. Si el eepsite va
  lento, baja `maxzoom` a 14 en `mapas.json`, regenera con `go run fetch_maps.go`
  y vuelve a desplegar.
- **Caché**: añade `Cache-Control: max-age=300` para el `.pmtiles`;
  i2pd también cachea en los peers mientras la tunnel está activa.
- **Solo local**: el túnel se conecta a `127.0.0.1:8080`; no hace falta abrir
  puertos en el router ni exponer el sitio a Internet.
- **Auditoría rápida** de que no hay recursos externos:

  ```bash
  grep -rIoE '(src|href)="?https?://[^" >]+' public/ | grep -v 'openstreetmap.org/copyright' | head
  ```

  Lo único que debería salir son los enlaces de atribución del mapa.
- **Actualizar el eepsite**: `git pull` (o edita), `hugo --minify --gc` y listo; no
  hay que tocar i2pd.
