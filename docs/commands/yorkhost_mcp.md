## yorkhost mcp

Brancher un assistant IA (Claude Code, Cursor, VS Code…) sur votre compte

### Description

Le serveur MCP YorkHost (https://api.yorkhost.fr/client/mcp) expose les mêmes
opérations que ce CLI comme outils pour assistants IA. Cette commande affiche la
configuration à copier pour votre client. Deux méthodes :

  - OAuth (recommandé) : le client ouvre lui-même la page de connexion YorkHost ;
  - clé API : créez une clé dédiée (espace client → Clés API), scopes minimum.

```
yorkhost mcp [options]
```

### Exemples

```
  yorkhost mcp
  yorkhost mcp --client cursor
```

### Options

```
      --client string   client MCP : claude, cursor, vscode, json (default "claude")
  -h, --help            aide de yorkhost mcp
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

