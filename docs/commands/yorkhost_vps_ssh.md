## yorkhost vps ssh

Se connecter en SSH à l'IP principale du VPS

```
yorkhost vps ssh [id|nom] [-- commande…] [options]
```

### Exemples

```
  yorkhost vps ssh 1234
  yorkhost vps ssh 1234 --user debian -- uptime
```

### Options

```
  -h, --help          aide de yorkhost vps ssh
      --port int      port SSH (default 22)
  -l, --user string   utilisateur (défaut : celui du VPS, sinon root)
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

* [yorkhost vps](yorkhost_vps.md)	 - VPS : état, démarrage, arrêt, métriques, SSH

