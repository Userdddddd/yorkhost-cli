## yorkhost ticket create

Ouvrir un ticket

```
yorkhost ticket create [options]
```

### Exemples

```
  yorkhost ticket create --subject "VPS injoignable" --body "Depuis 10h…" --service 1234
  yorkhost ticket create --subject "Question facture" --type BILLING --body @message.txt
  journalctl -n 50 | yorkhost ticket create --subject "Erreurs au démarrage" --body -
```

### Options

```
  -b, --body string      message (texte, @fichier ou - pour stdin)
  -h, --help             aide de yorkhost ticket create
      --service string   service concerné (ID ou nom)
  -s, --subject string   sujet
  -t, --type string      type de demande : HELP, BILLING, CANCELLATION, OTHER (default "HELP")
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

* [yorkhost ticket](yorkhost_ticket.md)	 - Support : tickets, réponses

