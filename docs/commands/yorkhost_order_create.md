## yorkhost order create

Commander (devis, confirmation, facture)

### Description

Commande un produit du catalogue en deux temps, comme l'impose l'API :

  1. devis : prix, frais d'installation et montant dû maintenant (rien n'est commandé) ;
  2. confirmation explicite (ou --yes) : commande et facture créées.

La facture est réglée par le crédit du compte s'il suffit ; sinon le lien de
paiement sécurisé est affiché (et ouvert avec --pay).

```
yorkhost order create <product-id> [options]
```

### Exemples

```
  yorkhost catalog
  yorkhost order create 42 --cycle a
  yorkhost order create 42 --yes --pay
```

### Options

```
      --cycle string   cycle de facturation : m, q, s, a (default "m")
  -h, --help           aide de yorkhost order create
      --pay            ouvrir la page de paiement si la facture est impayée
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

* [yorkhost order](yorkhost_order.md)	 - Commander un produit

