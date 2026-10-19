## yorkhost auth

Sessions et identifiants (status, token, login, logout)

### Options

```
  -h, --help   aide de yorkhost auth
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
* [yorkhost auth login](yorkhost_auth_login.md)	 - Se connecter (navigateur + A2F, ou clé API)
* [yorkhost auth logout](yorkhost_auth_logout.md)	 - Se déconnecter (révoque la session OAuth, oublie les identifiants)
* [yorkhost auth status](yorkhost_auth_status.md)	 - Identifiants de chaque profil (sans les secrets)
* [yorkhost auth token](yorkhost_auth_token.md)	 - Afficher le jeton d'accès courant (pour curl, scripts)

