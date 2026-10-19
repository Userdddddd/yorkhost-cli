## yorkhost vps

VPS : état, démarrage, arrêt, métriques, SSH

### Description

Gère vos VPS. L'ID est celui du service (yorkhost vps list) ; un nom de
service est aussi accepté, et l'ID peut être omis avec « yorkhost config set default_vps <id> ».

### Exemples

```
  yorkhost vps list
  yorkhost vps get 1234
  yorkhost vps restart 1234 --wait
  yorkhost vps ssh 1234
```

### Options

```
  -h, --help   aide de yorkhost vps
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
* [yorkhost vps get](yorkhost_vps_get.md)	 - État et ressources en direct
* [yorkhost vps list](yorkhost_vps_list.md)	 - Lister vos VPS
* [yorkhost vps metrics](yorkhost_vps_metrics.md)	 - CPU et RAM de la dernière heure
* [yorkhost vps restart](yorkhost_vps_restart.md)	 - Redémarrer un VPS
* [yorkhost vps ssh](yorkhost_vps_ssh.md)	 - Se connecter en SSH à l'IP principale du VPS
* [yorkhost vps start](yorkhost_vps_start.md)	 - Démarrer un VPS
* [yorkhost vps stop](yorkhost_vps_stop.md)	 - Arrêter un VPS

