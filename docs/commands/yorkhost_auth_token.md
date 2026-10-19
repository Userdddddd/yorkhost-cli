## yorkhost auth token

Afficher le jeton d'accès courant (pour curl, scripts)

```
yorkhost auth token [options]
```

### Exemples

```
  curl -H "Authorization: Bearer $(yorkhost auth token)" https://api.yorkhost.fr/client/v1/services
```

### Options

```
  -h, --help   aide de yorkhost auth token
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

* [yorkhost auth](yorkhost_auth.md)	 - Sessions et identifiants (status, token, login, logout)

