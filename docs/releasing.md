# Publier une version

Les versions suivent le [versionnage sémantique](https://semver.org/lang/fr/) : `vMAJEUR.MINEUR.CORRECTIF`.
Une commande ou une option retirée, un format de sortie ou un code de sortie modifié = version majeure.

## Étapes

1. Mettre à jour `CHANGELOG.md` : déplacer « Non publié » sous la nouvelle version, avec la date.
2. Vérifier en local :

   ```sh
   make lint test docs-check
   make snapshot          # goreleaser à blanc → dist/
   ```

3. Taguer et pousser :

   ```sh
   git tag -a v1.0.0 -m "v1.0.0"
   git push origin v1.0.0
   ```

   Ou, sans tag local : *Actions → Release → Run workflow*, avec la version (`v1.0.0`).
   Le workflow crée le tag sur `main` puis publie.

Le workflow [`release.yml`](../.github/workflows/release.yml) lance les tests puis GoReleaser,
qui publie sur la release GitHub :

- les binaires Linux, macOS, Windows, FreeBSD (amd64, arm64, armv7) en `.tar.gz` / `.zip`,
  avec README, licence, autocomplétions et pages `man` ;
- les paquets `.deb`, `.rpm` et `.apk` ;
- `checksums.txt` (SHA-256), utilisé par les scripts d'installation ;
- un changelog groupé à partir des messages de commit (`feat:`, `fix:`…).

Un tag contenant un suffixe (`v1.1.0-rc.1`) est publié comme *pre-release* et n'est pas pris
par `install.sh` (qui suit la dernière release stable).

## Version dans le binaire

GoReleaser injecte version, commit et date (`-ldflags -X …/internal/version.*`).
`go install …@v1.0.0` sans ldflags lit la version du module. `yorkhost version` affiche le tout.
