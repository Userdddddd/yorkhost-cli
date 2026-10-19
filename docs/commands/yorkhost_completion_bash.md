## yorkhost completion bash

Script d'autocomplétion pour bash

### Description

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(yorkhost completion bash)

To load completions for every new session, execute once:

#### Linux:

	yorkhost completion bash > /etc/bash_completion.d/yorkhost

#### macOS:

	yorkhost completion bash > $(brew --prefix)/etc/bash_completion.d/yorkhost

You will need to start a new shell for this setup to take effect.


```
yorkhost completion bash
```

### Options

```
  -h, --help              aide de yorkhost completion bash
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

