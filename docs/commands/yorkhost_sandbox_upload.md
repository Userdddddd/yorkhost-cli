## yorkhost sandbox upload

Envoyer un fichier (1 Mio au plus, défaut /home/agent/<nom>)

```
yorkhost sandbox upload <id> <fichier local> [chemin distant] [options]
```

### Exemples

```
  yorkhost sandbox upload sbxabcdefghijkl ./setup.sh --mode 755
  yorkhost sandbox upload sbxabcdefghijkl .env app/.env
```

### Options

```
  -h, --help          aide de yorkhost sandbox upload
      --mode string   droits en octal (ex. 644, 755) ; défaut : 755 si exécutable localement
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

