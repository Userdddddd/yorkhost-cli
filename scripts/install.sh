#!/bin/sh
# Installe le CLI YorkHost (dernière version, ou YORKHOST_VERSION=v1.2.3).
#
#   curl -fsSL https://raw.githubusercontent.com/YorkHost-fr/yorkhost-cli/main/scripts/install.sh | sh
#
# Variables :
#   YORKHOST_VERSION      version à installer (défaut : dernière release)
#   YORKHOST_INSTALL_DIR  dossier de destination (défaut : /usr/local/bin si
#                         accessible en écriture, sinon ~/.local/bin)
#
# La somme SHA-256 de l'archive est vérifiée contre checksums.txt de la release.
set -eu

REPO="YorkHost-fr/yorkhost-cli"
BIN="yorkhost"

say() { printf '%s\n' "$*" >&2; }
die() { say "Erreur : $*"; exit 1; }

need() { command -v "$1" >/dev/null 2>&1 || die "commande requise introuvable : $1"; }

fetch() {
	if command -v curl >/dev/null 2>&1; then
		curl -fsSL "$1" -o "$2"
	elif command -v wget >/dev/null 2>&1; then
		wget -qO "$2" "$1"
	else
		die "curl ou wget requis"
	fi
}

need uname
need tar

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
	linux | darwin | freebsd) ;;
	*) die "système non pris en charge : $os (Windows : voir docs/installation.md)" ;;
esac

arch=$(uname -m)
case "$arch" in
	x86_64 | amd64) arch=x86_64 ;;
	aarch64 | arm64) arch=arm64 ;;
	armv7* | armv8l) arch=armv7 ;;
	*) die "architecture non prise en charge : $arch" ;;
esac

version=${YORKHOST_VERSION:-}
if [ -z "$version" ]; then
	tmpv=$(mktemp)
	fetch "https://api.github.com/repos/$REPO/releases/latest" "$tmpv" || die "impossible de lire la dernière version"
	version=$(sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' "$tmpv" | head -n1)
	rm -f "$tmpv"
	[ -n "$version" ] || die "aucune release publiée"
fi
num=${version#v}

archive="${BIN}_${num}_${os}_${arch}.tar.gz"
base="https://github.com/$REPO/releases/download/$version"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT INT TERM

say "Téléchargement de $BIN $version ($os/$arch)…"
fetch "$base/$archive" "$tmp/$archive" || die "archive introuvable : $base/$archive"
fetch "$base/checksums.txt" "$tmp/checksums.txt" || die "checksums.txt introuvable"

expected=$(grep "  $archive\$" "$tmp/checksums.txt" | cut -d' ' -f1)
[ -n "$expected" ] || die "somme de contrôle absente pour $archive"
if command -v sha256sum >/dev/null 2>&1; then
	actual=$(sha256sum "$tmp/$archive" | cut -d' ' -f1)
elif command -v shasum >/dev/null 2>&1; then
	actual=$(shasum -a 256 "$tmp/$archive" | cut -d' ' -f1)
else
	die "sha256sum ou shasum requis pour vérifier l'archive"
fi
[ "$expected" = "$actual" ] || die "somme de contrôle invalide (attendu $expected, obtenu $actual)"

tar -xzf "$tmp/$archive" -C "$tmp" "$BIN"

dir=${YORKHOST_INSTALL_DIR:-}
if [ -z "$dir" ]; then
	if [ -w /usr/local/bin ]; then
		dir=/usr/local/bin
	else
		dir="$HOME/.local/bin"
	fi
fi
mkdir -p "$dir"
install -m 0755 "$tmp/$BIN" "$dir/$BIN" 2>/dev/null || { cp "$tmp/$BIN" "$dir/$BIN" && chmod 0755 "$dir/$BIN"; }

say "Installé : $dir/$BIN"
case ":$PATH:" in
	*":$dir:"*) ;;
	*) say "Ajoutez $dir à votre PATH, par exemple : export PATH=\"$dir:\$PATH\"" ;;
esac
say ""
say "Pour commencer : $BIN login"
