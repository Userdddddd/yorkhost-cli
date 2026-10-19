## yorkhost sandbox create

Créer une sandbox (gratuite pendant sa fenêtre de fonctionnement)

```
yorkhost sandbox create [options]
```

### Exemples

```
  yorkhost sandbox create
  yorkhost sandbox create --kind desktop --wait
  yorkhost sandbox create --ssh-key ~/.ssh/work.pub --no-secrets --ssh
```

### Options

```
  -h, --help             aide de yorkhost sandbox create
      --kind string      type : linux, desktop, windows (default "linux")
      --no-secrets       ne pas exporter les secrets du compte dans l'essai
      --ssh              se connecter en SSH une fois démarrée
      --ssh-key string   clé publique SSH (défaut : ~/.ssh/id_ed25519.pub, id_ecdsa.pub ou id_rsa.pub)
      --wait             attendre que la sandbox soit démarrée
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

