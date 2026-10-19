# Sécurité

## Signaler une faille

**N'ouvrez pas d'issue publique.** Utilisez le
[signalement privé GitHub](https://github.com/YorkHost-fr/yorkhost-cli/security/advisories/new),
ou écrivez à **security@yorkhost.fr**. Nous accusons réception sous 72 h.

Merci d'indiquer la version (`yorkhost version`), le système, et les étapes pour reproduire.
Ne joignez jamais une clé API ou un jeton réel : révoquez-le d'abord dans l'espace client.

Sont concernés ce CLI et ses scripts d'installation. Une faille de l'API (`api.yorkhost.fr`) ou
de l'espace client se signale à la même adresse.

## Ce que fait le CLI pour protéger vos accès

- Aucun secret embarqué ; tous les contrôles d'accès sont faits par l'API.
- Connexion OAuth avec PKCE S256 et `state`, callback sur 127.0.0.1 uniquement.
- Identifiants stockés dans un fichier `0600`, écrit de façon atomique ; refresh tokens à usage
  unique, rafraîchis sous verrou.
- `--debug` n'affiche jamais l'en-tête `Authorization` ; `whoami` ne montre que l'identifiant
  public d'une clé.
- Aucune action modifiante n'est rejouée automatiquement ; les actions sensibles demandent une
  confirmation.
- Archives de release accompagnées de sommes SHA-256, vérifiées par les scripts d'installation.

## Versions suivies

Les correctifs de sécurité sont publiés sur la dernière version mineure.
