## yorkhost config profile

Gérer les profils

### Options

```
  -h, --help   aide de yorkhost config profile
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

* [yorkhost config](yorkhost_config.md)	 - Réglages et profils (plusieurs comptes, comme aws --profile)
* [yorkhost config profile delete](yorkhost_config_profile_delete.md)	 - Supprimer un profil (et ses identifiants)
* [yorkhost config profile list](yorkhost_config_profile_list.md)	 - Lister les profils
* [yorkhost config profile use](yorkhost_config_profile_use.md)	 - Choisir le profil actif (créé s'il n'existe pas)

