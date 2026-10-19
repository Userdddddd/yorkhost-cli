## yorkhost login

Se connecter (navigateur + A2F, ou clé API)

### Description

Connecte le profil courant à votre compte YorkHost.

Par défaut, le navigateur s'ouvre sur votre espace client : connectez-vous
(double authentification requise), choisissez les permissions et les services
accessibles au CLI, puis validez. Le CLI reçoit un jeton d'accès (1 h) et un
jeton de rafraîchissement (30 jours glissants), renouvelés automatiquement.
L'accès se retire dans l'espace client → Clés API → Applications connectées.

Avec --api-key, le CLI utilise une clé créée dans l'espace client
(Clés API → Nouvelle clé) : pratique pour un serveur ou une CI. La clé est
demandée en saisie masquée, ou lue sur stdin avec --with-token.

```
yorkhost login [options]
```

### Exemples

```
  yorkhost login
  yorkhost login --no-browser            # machine distante : ouvrez le lien vous-même
  yorkhost login --api-key
  echo "$YORKHOST_KEY" | yorkhost login --with-token --profile ci
  yorkhost login --scopes services:read,vps:read,vps:power
```

### Options

```
      --api-key          utiliser une clé API (saisie masquée)
  -h, --help             aide de yorkhost login
      --no-browser       ne pas ouvrir le navigateur, afficher le lien
      --scopes strings   permissions demandées (défaut : toutes, à affiner sur la page de consentement)
      --with-token       lire la clé API sur stdin
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

