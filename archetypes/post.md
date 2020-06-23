---
title: "{{ replace .Name "-" " " | title }}"
date: {{ .Date }}
draft: true
behoben: "{{ .Date }}"
versprechen: "gebrochen" # oder "gehalten" oder leer.
infra_problem: "problem" # Kann statt verprechen genutzt
location:
  lat:
  lon:
---