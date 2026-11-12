---
title: "About"
---

Citizen log of harassment incidents in **Saint Petersburg, Russia**.
The data is published as GeoJSON and rendered over a self-hosted vector map:
Leaflet 1.9 + protomaps-leaflet and a `.pmtiles` archive clipped to the city.
No CDNs and no third-party tiles at runtime.

## Methodology

Each GeoJSON record includes `titulo`, `fecha` (DD.MM.YYYY), `severidad`
and `descripcion`. High-severity incidents involve stalking or physical
contact; medium severity means intimidation or blocking someone's path;
low severity means hostile comments without assault.

> The three incidents included here are **fictional** and serve as a template.
