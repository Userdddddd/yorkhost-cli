## yorkhost game

Serveurs de jeu : état, démarrage, arrêt, métriques

### Description

Gère vos serveurs de jeu (FiveM, Minecraft, Rust, bots…). L'ID est celui du
service ; un nom est accepté, et l'ID peut être omis avec « yorkhost config set default_game <id> ».

### Exemples

```
  yorkhost game list
  yorkhost game get mon-serveur
  yorkhost game restart 5678 --wait
  yorkhost game open 5678          # ouvre le panel
```

### Options

```
  -h, --help   aide de yorkhost game
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

* [yorkhost](yorkhost.md)	 - CLI officiel YorkHost : VPS, serveurs de jeu, facturation, support, sandboxes
* [yorkhost game get](yorkhost_game_get.md)	 - État et ressources en direct
* [yorkhost game kill](yorkhost_game_kill.md)	 - Tuer un serveur de jeu bloqué (arrêt forcé)
* [yorkhost game list](yorkhost_game_list.md)	 - Lister vos serveurs de jeu
* [yorkhost game metrics](yorkhost_game_metrics.md)	 - CPU et RAM de la dernière heure
* [yorkhost game open](yorkhost_game_open.md)	 - Ouvrir le panel du serveur dans le navigateur
* [yorkhost game restart](yorkhost_game_restart.md)	 - Redémarrer un serveur de jeu
* [yorkhost game start](yorkhost_game_start.md)	 - Démarrer un serveur de jeu
* [yorkhost game stop](yorkhost_game_stop.md)	 - Arrêter un serveur de jeu

