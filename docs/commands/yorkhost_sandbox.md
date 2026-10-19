## yorkhost sandbox

Sandboxes agents : VM d'essai jetables (SSH via yorkhost.sh)

### Description

VM d'essai pour agents IA (Linux, bureau Linux, bureau Windows), sans IP publique :
SSH par la passerelle (ssh <id>@yorkhost.sh), applications web sur leur URL de preview.
Gratuites pendant leur fenêtre de fonctionnement, puis arrêtées ; supprimées avec leurs
fichiers au bout de 24 h, sauf passage en VPS (yorkhost sandbox keep).

### Exemples

```
  yorkhost sandbox create --ssh            # crée un essai Linux et s'y connecte
  yorkhost sandbox exec sbxabcdefghijkl -- uname -a
  yorkhost sandbox upload sbxabcdefghijkl ./app.py
  yorkhost sandbox keep sbxabcdefghijkl    # le garder : passage en VPS-2
```

### Options

```
  -h, --help   aide de yorkhost sandbox
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
* [yorkhost sandbox cat](yorkhost_sandbox_cat.md)	 - Afficher un fichier (1 Mio au plus)
* [yorkhost sandbox claim](yorkhost_sandbox_claim.md)	 - Rattacher un essai anonyme (ssh yorkhost.sh) et le passer en VPS
* [yorkhost sandbox create](yorkhost_sandbox_create.md)	 - Créer une sandbox (gratuite pendant sa fenêtre de fonctionnement)
* [yorkhost sandbox delete](yorkhost_sandbox_delete.md)	 - Supprimer une sandbox et ses fichiers
* [yorkhost sandbox download](yorkhost_sandbox_download.md)	 - Télécharger un fichier (1 Mio au plus)
* [yorkhost sandbox exec](yorkhost_sandbox_exec.md)	 - Exécuter une commande (bash -lc, utilisateur agent) sans SSH
* [yorkhost sandbox get](yorkhost_sandbox_get.md)	 - Détail d'une sandbox
* [yorkhost sandbox keep](yorkhost_sandbox_keep.md)	 - Garder la sandbox : commander un VPS-2 et y copier /home/agent
* [yorkhost sandbox key](yorkhost_sandbox_key.md)	 - Clés SSH du compte acceptées sur toutes vos sandboxes
* [yorkhost sandbox links](yorkhost_sandbox_links.md)	 - Nouveau lien bureau (l'ancien cesse de fonctionner)
* [yorkhost sandbox list](yorkhost_sandbox_list.md)	 - Lister vos sandboxes
* [yorkhost sandbox ls](yorkhost_sandbox_ls.md)	 - Lister un dossier (défaut /home/agent)
* [yorkhost sandbox offers](yorkhost_sandbox_offers.md)	 - Types de sandbox proposés et leurs limites
* [yorkhost sandbox secret](yorkhost_sandbox_secret.md)	 - Secrets exportés dans vos sandboxes Linux (ANTHROPIC_API_KEY…)
* [yorkhost sandbox ssh](yorkhost_sandbox_ssh.md)	 - Se connecter en SSH (passerelle yorkhost.sh)
* [yorkhost sandbox upload](yorkhost_sandbox_upload.md)	 - Envoyer un fichier (1 Mio au plus, défaut /home/agent/<nom>)

