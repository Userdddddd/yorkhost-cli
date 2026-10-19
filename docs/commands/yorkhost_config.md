## yorkhost config

Réglages et profils (plusieurs comptes, comme aws --profile)

### Description

Chaque profil a sa propre URL d'API, ses identifiants et ses réglages.
Le profil actif est utilisé par défaut ; --profile ou YORKHOST_PROFILE en choisissent un autre.

Réglages :
  api_url         URL de l'API YorkHost (défaut https://api.yorkhost.fr)
  output          format de sortie par défaut : table, wide, json, yaml
  default_vps     VPS utilisé quand l'ID est omis (yorkhost vps …)
  default_game    serveur de jeu utilisé quand l'ID est omis (yorkhost game …)
  ssh_public_key  fichier de clé publique SSH pour les sandboxes (défaut ~/.ssh/id_ed25519.pub)


### Exemples

```
  yorkhost config set output json
  yorkhost config set default_vps 1234
  yorkhost --profile pro login
  yorkhost config profile use pro
```

### Options

```
  -h, --help   aide de yorkhost config
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
* [yorkhost config get](yorkhost_config_get.md)	 - Lire un réglage du profil
* [yorkhost config list](yorkhost_config_list.md)	 - Réglages du profil courant
* [yorkhost config path](yorkhost_config_path.md)	 - Chemin du fichier de configuration
* [yorkhost config profile](yorkhost_config_profile.md)	 - Gérer les profils
* [yorkhost config set](yorkhost_config_set.md)	 - Modifier un réglage du profil
* [yorkhost config unset](yorkhost_config_unset.md)	 - Revenir à la valeur par défaut

