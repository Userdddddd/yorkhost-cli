## yorkhost sandbox exec

Exécuter une commande (bash -lc, utilisateur agent) sans SSH

### Description

Exécute une commande dans la sandbox (Linux et bureau Linux en marche) sous
l'utilisateur agent, via bash -lc, dans /home/agent ou --cwd. stdout et stderr sont
relayés ; le code de sortie du CLI est celui de la commande. 60 s par défaut, 300 s au plus.

```
yorkhost sandbox exec <id> -- <commande…> [options]
```

### Exemples

```
  yorkhost sandbox exec sbxabcdefghijkl -- 'cd app && npm test'
  echo 'print(1+1)' | yorkhost sandbox exec sbxabcdefghijkl -- python3
  yorkhost sandbox exec sbxabcdefghijkl --timeout 300 -- ./build.sh
```

### Options

```
      --cwd string          dossier de travail (sous /home/agent)
  -h, --help                aide de yorkhost sandbox exec
      --stdin-file string   fichier passé sur l'entrée standard (32 Kio au plus)
      --timeout int         délai en secondes (1 à 300, défaut 60)
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

