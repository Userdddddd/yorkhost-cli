## yorkhost game metrics

CPU et RAM de la dernière heure

### Description

CPU et RAM de la dernière heure. En -o wide, chaque point est listé ; en -o json, la série brute.

```
yorkhost game metrics [id|nom] [options]
```

### Options

```
  -h, --help   aide de yorkhost game metrics
```

### Options globales

```
      --api-url string   URL de l'API (env YORKHOST_API_URL)
      --debug            journal des requêtes HTTP sur stderr (sans secrets)
      --no-color         désactiver les couleurs (env NO_COLOR)
  -o, --output string    format de sortie : table, wide, json, yaml (env YORKHOST_OUTPUT)
  -p, --profile string   profil de configuration (env YORKHOST_PROFILE)
  -q, --quiet            n'afficher que le résultat
  -y, --yes              répondre oui aux confirmations (env YORKHOST_YES=1)
```

### Voir aussi

* [yorkhost game](yorkhost_game.md)	 - Serveurs de jeu : état, démarrage, arrêt, métriques

