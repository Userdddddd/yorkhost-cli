## yorkhost sandbox secret

Secrets exportés dans vos sandboxes Linux (ANTHROPIC_API_KEY…)

### Description

Variables d'environnement chiffrées (AES-256-GCM) exportées dans les sandboxes Linux
créées ensuite (/home/agent/.config/yorkhost/env). Les valeurs ne sont jamais relues.

### Options

```
  -h, --help   aide de yorkhost sandbox secret
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

* [yorkhost sandbox](yorkhost_sandbox.md)	 - Sandboxes agents : VM d'essai jetables (SSH via yorkhost.sh)
* [yorkhost sandbox secret delete](yorkhost_sandbox_secret_delete.md)	 - Supprimer un secret
* [yorkhost sandbox secret list](yorkhost_sandbox_secret_list.md)	 - Lister les secrets (noms seulement)
* [yorkhost sandbox secret set](yorkhost_sandbox_secret_set.md)	 - Créer ou remplacer un secret (saisie masquée si la valeur est omise)

