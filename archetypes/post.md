---
title: "{{ replace .Name "-" " " | title }}"
date: {{ .Date }}
draft: true # Auf false setzen, damit der Post auch veröffentlicht wird
behoben: "{{ .Date }}" # Kann weggelassen werden wenn der Missstand noch nicht behoben ist
versprechen: "gebrochen" # oder "gehalten" oder leer. Leer bedeutet, dass das Versprechen noch ausstehend ist.
infra_problem: "problem" # Kann statt verprechen genutzt. infra_problem und verpsrechen sollten nicht zusammen im frontmatter stehen
location:    # Ist Location gesetzt, taucht dieser Post auf der Karte auf. Weglassen um nicht auf der Karte aufutauchen.
  lat:       # Latitude der Google Maps Koordinaten
  lon:       # Longitude der Google Maps Koordinaten
---