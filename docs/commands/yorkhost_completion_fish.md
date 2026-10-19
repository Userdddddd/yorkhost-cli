## yorkhost completion fish

Script d'autocomplétion pour fish

### Description

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	yorkhost completion fish | source

To load completions for every new session, execute once:

	yorkhost completion fish > ~/.config/fish/completions/yorkhost.fish

You will need to start a new shell for this setup to take effect.


```
yorkhost completion fish [options]
```

### Options

```
  -h, --help              aide de yorkhost completion fish
      --no-descriptions   sans les descriptions des propositions
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

* [yorkhost completion](yorkhost_completion.md)	 - Script d'autocomplétion pour bash, zsh, fish ou PowerShell

