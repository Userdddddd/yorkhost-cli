## yorkhost auth logout

Se déconnecter (révoque la session OAuth, oublie les identifiants)

### Description

Supprime les identifiants du profil. Une session OAuth est aussi révoquée
côté YorkHost. Une clé API n'est pas révoquée : faites-le dans l'espace client
(Clés API) si elle ne doit plus servir.

```
yorkhost auth logout [options]
```

### Options

```
      --all    tous les profils
  -h, --help   aide de yorkhost auth logout
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

