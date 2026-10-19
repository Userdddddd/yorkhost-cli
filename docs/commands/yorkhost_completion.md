## yorkhost completion

Script d'autocomplétion pour bash, zsh, fish ou PowerShell

### Description

Génère le script d'autocomplétion du shell.

  bash :       source <(yorkhost completion bash)
  zsh :        yorkhost completion zsh > "${fpath[1]}/_yorkhost"
  fish :       yorkhost completion fish > ~/.config/fish/completions/yorkhost.fish
  PowerShell : yorkhost completion powershell | Out-String | Invoke-Expression

### Options

```
  -h, --help   aide de yorkhost completion
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
* [yorkhost completion bash](yorkhost_completion_bash.md)	 - Script d'autocomplétion pour bash
* [yorkhost completion fish](yorkhost_completion_fish.md)	 - Script d'autocomplétion pour fish
* [yorkhost completion powershell](yorkhost_completion_powershell.md)	 - Script d'autocomplétion pour powershell
* [yorkhost completion zsh](yorkhost_completion_zsh.md)	 - Script d'autocomplétion pour zsh

