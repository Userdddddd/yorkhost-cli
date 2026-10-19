# Configuration et profils

## Fichier

| Système | Emplacement |
|---|---|
| Linux, macOS, BSD | `$XDG_CONFIG_HOME/yorkhost/config.yaml`, sinon `~/.config/yorkhost/config.yaml` |
| Windows | `%AppData%\yorkhost\config.yaml` |
| Partout | `$YORKHOST_CONFIG_FILE` s'il est défini |

`yorkhost config path` affiche le chemin utilisé. Le fichier est écrit de façon atomique avec
les droits `0600` : il contient vos jetons.

```yaml
active_profile: default
profiles:
  default:
    output: table
    default_vps: "1234"
    auth:
      type: oauth
      access_token: yhk_client_v1_…
      refresh_token: yhr_client_v1_…
      expires_at: 2026-10-19T20:00:00Z
      client_id: yhc_…
  pro:
    api_url: https://api.yorkhost.fr
    auth:
      type: api_key
      api_key: yhk_client_v1_…
```

Préférez les commandes `yorkhost config …` et `yorkhost login` à l'édition manuelle.

## Réglages

| Réglage | Rôle | Défaut |
|---|---|---|
| `api_url` | URL de l'API | `https://api.yorkhost.fr` |
| `output` | format par défaut : `table`, `wide`, `json`, `yaml` | `table` |
| `default_vps` | VPS visé quand l'ID est omis (`yorkhost vps get`) | — |
| `default_game` | serveur de jeu visé quand l'ID est omis | — |
| `ssh_public_key` | clé publique pour les sandboxes | `~/.ssh/id_ed25519.pub`, `id_ecdsa.pub`, `id_rsa.pub` |

```sh
yorkhost config list
yorkhost config set output json
yorkhost config get output
yorkhost config unset output
```

## Profils

Un profil = un compte (ou un environnement) avec ses identifiants et ses réglages.

```sh
yorkhost --profile pro login          # crée le profil "pro" et s'y connecte
yorkhost config profile list
yorkhost config profile use pro       # profil actif par défaut
yorkhost -p perso invoice list        # ponctuellement
yorkhost config profile delete pro    # révoque la session et supprime le profil
```

## Priorités

Pour chaque valeur : **option de ligne de commande > variable d'environnement > profil > défaut**.

| Valeur | Option | Variable |
|---|---|---|
| Profil | `--profile`, `-p` | `YORKHOST_PROFILE` |
| URL de l'API | `--api-url` | `YORKHOST_API_URL` |
| Format de sortie | `--output`, `-o` | `YORKHOST_OUTPUT` |
| Identifiants | `login` | `YORKHOST_API_KEY` (ou `YORKHOST_TOKEN`) |
| Confirmations | `--yes`, `-y` | `YORKHOST_YES=1` |
| Couleurs | `--no-color` | `NO_COLOR` |
| Fichier de configuration | — | `YORKHOST_CONFIG_FILE` |
| Navigateur | `login --no-browser` | `YORKHOST_NO_BROWSER=1` |
