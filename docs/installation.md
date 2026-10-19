# Installation

Le CLI est un binaire unique, sans dépendance. Choisissez la méthode qui vous convient.

## Script (Linux, macOS, FreeBSD)

```sh
curl -fsSL https://raw.githubusercontent.com/YorkHost-fr/yorkhost-cli/main/scripts/install.sh | sh
```

Le script détecte le système et l'architecture (x86_64, arm64, armv7), télécharge l'archive de
la dernière release, **vérifie sa somme SHA-256** contre `checksums.txt`, puis installe
`yorkhost` dans `/usr/local/bin` (ou `~/.local/bin` sans droits d'écriture).

```sh
# Version précise, dossier précis
curl -fsSL https://raw.githubusercontent.com/YorkHost-fr/yorkhost-cli/main/scripts/install.sh \
  | YORKHOST_VERSION=v1.0.0 YORKHOST_INSTALL_DIR="$HOME/bin" sh
```

Vous préférez lire avant d'exécuter ? Téléchargez [`scripts/install.sh`](../scripts/install.sh),
lisez-le, puis lancez `sh install.sh`.

## Windows

PowerShell :

```powershell
irm https://raw.githubusercontent.com/YorkHost-fr/yorkhost-cli/main/scripts/install.ps1 | iex
```

Installe `yorkhost.exe` dans `%LOCALAPPDATA%\Programs\yorkhost` et l'ajoute au `PATH` de
l'utilisateur. Ou téléchargez `yorkhost_<version>_windows_x86_64.zip` depuis les releases.

## Paquets

Sur la [page des releases](https://github.com/YorkHost-fr/yorkhost-cli/releases/latest) :

```sh
# Debian, Ubuntu
sudo apt install ./yorkhost_<version>_linux_amd64.deb
# Fedora, RHEL, Rocky, Alma
sudo dnf install ./yorkhost_<version>_linux_amd64.rpm
# Alpine
sudo apk add --allow-untrusted ./yorkhost_<version>_linux_amd64.apk
```

Les paquets installent aussi l'autocomplétion bash, zsh et fish, et les pages `man yorkhost`.

## Go

```sh
go install github.com/YorkHost-fr/yorkhost-cli/cmd/yorkhost@latest
```

Go 1.26 ou plus récent. Le binaire est placé dans `$(go env GOPATH)/bin`.

## Archive manuelle

```sh
VERSION=1.0.0
curl -fsSLO https://github.com/YorkHost-fr/yorkhost-cli/releases/download/v$VERSION/yorkhost_${VERSION}_linux_x86_64.tar.gz
curl -fsSLO https://github.com/YorkHost-fr/yorkhost-cli/releases/download/v$VERSION/checksums.txt
sha256sum --ignore-missing -c checksums.txt
tar -xzf yorkhost_${VERSION}_linux_x86_64.tar.gz yorkhost
sudo install -m 0755 yorkhost /usr/local/bin/
```

## Depuis les sources

```sh
git clone https://github.com/YorkHost-fr/yorkhost-cli.git
cd yorkhost-cli
make build        # → bin/yorkhost
```

## Autocomplétion

```sh
# bash (~/.bashrc)
source <(yorkhost completion bash)
# zsh
yorkhost completion zsh > "${fpath[1]}/_yorkhost"
# fish
yorkhost completion fish > ~/.config/fish/completions/yorkhost.fish
# PowerShell ($PROFILE)
yorkhost completion powershell | Out-String | Invoke-Expression
```

L'autocomplétion propose aussi vos **IDs de services et de sandboxes** (avec leur nom) quand
vous êtes connecté.

## Mise à jour et désinstallation

```sh
yorkhost version --check    # nouvelle version disponible ?
```

Pour mettre à jour, relancez la méthode d'installation. Pour désinstaller : supprimez le
binaire (ou le paquet), puis éventuellement `yorkhost logout --all` et le dossier de
configuration (`yorkhost config path`).
