---
title: "{{ replace .File.ContentBaseName `-` ` ` | title }}"
date: "{{ .Date }}"
draft: true
# you can add multiple tracks
tracks:
  - id: track1
    file: 250308.gpx
    color: "#3388ff"
    weight: 5
    opacity: 0.7
# recommended colors (solarized):
# #b1cc00
# #2eb2a7
# #3093d9
# #5f62bf
# #c92c78
# #c02321
# #dc5318
# #f5b800
tags:
  - outdoor

layout: outdoor
---

## Lorem ipsum
