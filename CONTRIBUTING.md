# Contribuer

Merci de votre intérêt ! Issues et pull requests sont les bienvenues.

## Mettre en place

```sh
git clone https://github.com/YorkHost-fr/yorkhost-cli.git
cd yorkhost-cli
make build          # → bin/yorkhost
make test           # tests unitaires + bout en bout contre un faux serveur API
```

Prérequis : Go 1.24+. Optionnel : [golangci-lint](https://golangci-lint.run/) et
[GoReleaser](https://goreleaser.com/) (`make snapshot`).

Essayer contre l'API réelle sans toucher à votre configuration :

```sh
YORKHOST_CONFIG_FILE=/tmp/yh.yaml ./bin/yorkhost login
```

## Organisation

Voir [docs/architecture.md](docs/architecture.md#organisation-du-code). En bref : une ressource
= un fichier dans `internal/cli/`, le transport dans `internal/api/`, l'authentification dans
`internal/auth/`.

## Ajouter une commande

1. La route doit exister dans l'API client (`/client/v1`, dépôt `yorkhost-api`) : le CLI
   n'appelle jamais d'autre API.
2. Suivez la forme `yorkhost <ressource> <verbe>` et les conventions existantes :
   - résultat sur stdout, messages sur stderr (`a.Printer.Info/Success/Warn`) ;
   - `-o json` renvoie la réponse de l'API telle quelle (`a.Printer.List/Object/Raw`) ;
   - confirmation (`a.Prompt.Confirm(…, a.Yes)`) pour toute action destructive ou payante ;
   - erreurs d'usage via `usageErr(…)` (code de sortie 2) ;
   - textes en français, aide avec `Example`.
3. Ajoutez un test dans `internal/cli/cli_test.go` (le faux serveur reprend les formats réels).
4. Mettez à jour `docs/api.md`, `CHANGELOG.md` et régénérez la référence : `make docs`.

## Avant d'ouvrir une pull request

```sh
make lint test docs-check
```

Messages de commit au format [Conventional Commits](https://www.conventionalcommits.org/fr/)
(`feat:`, `fix:`, `docs:`…) : ils alimentent le changelog des releases.

Ne committez jamais de clé, de jeton ni de fichier de configuration réel.
