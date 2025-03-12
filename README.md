# ssh-bip39gen

`ssh-bip39gen` is a command-line tool that generates deterministic Ed25519 SSH key pairs from a 24-word BIP-39 mnemonic phrase. Unlike traditional SSH key generation, this tool allows you to regenerate the same key pair using only the mnemonic—no need to back up key files. It’s ideal for scenarios where you want to recover your SSH keys without storing them, provided you keep the mnemonic secure.

## Features

- Generates Ed25519 SSH keys with 256-bit entropy (24-word mnemonic).
- Deterministic: Same mnemonic produces the same key pair.
- Cross-platform: Binaries for Linux, Windows, and macOS (amd64 and 386/arm64).
- Simple usage with a familiar `ssh-keygen`-like interface.

## Installation

### Pre-built Binaries

Download the latest release from the [Releases page](https://github.com/heqnx/ssh-bip39gen/releases) for your platform:
- `ssh-bip39gen-linux-amd64`
- `ssh-bip39gen-linux-386`
- `ssh-bip39gen-windows-amd64.exe`
- `ssh-bip39gen-windows-386.exe`
- `ssh-bip39gen-darwin-amd64`
- `ssh-bip39gen-darwin-arm64`

Move the binary to a directory in your PATH (e.g., `/usr/local/bin` on Unix-like systems).

### Build from Source

Requires Go 1.21+ and `make`.

1. Clone the repository:

```
$ git clone https://github.com/heqnx/ssh-bip39gen.git
$ cd ssh-bip39gen
```

2. Build all binaries:

```
$ make
```

- Output is in the build/ directory.

3. (Optional) Build for a specific platform:
```
$ make linux
$ make windows
$ make darwin
```

4. Clean up:

```
$ make clean
```

## Usage

### Generate a New Key Pair

```
$ ssh-bip39gen
```

- Creates id_ed25519 (private) and id_ed25519.pub (public).
- Outputs a 24-word mnemonic (e.g., "abandon ability able about ... actress").
- Save the mnemonic securely - it’s your only way to regenerate the keys!

### Generate with Custom Output File

```
$ ssh-bip39gen -f testkey
```

- Creates testkey (private) and testkey.pub (public).

### Regenerate from a Mnemonic

```
$ ssh-bip39gen -f testkey -mnemonic "abandon ability able about above absent absorb abstract absurd abuse access accident account accuse achieve acid acoustic acquire across act action actor actress"
```

- Regenerates the same key pair using the provided 24-word mnemonic.

### Help

```
$ ssh-bip39gen -h
```

- Displays usage instructions and flags.

## License

- GNU GENERAL PUBLIC LICENSE Version 3

## Security Notes

- Mnemonic Security: The mnemonic is your private key. Treat it like a secret - write it down on paper, store it in a safe, or use a hardware wallet. Do not store it digitally unless encrypted.
- Entropy: 256-bit entropy.
- Determinism: If the mnemonic is compromised, an attacker can regenerate your keys. Use a unique, randomly generated mnemonic for each key pair.

Contributing
Feel free to submit issues or pull requests on GitHub. Ensure any changes maintain the 24-word mnemonic requirement and Ed25519 key type.
License
MIT License (LICENSE) - Free to use, modify, and distribute.

## Acknowledgments

- Built with Go, go-bip39, and x/crypto.
- Inspired by SSH key management needs and BIP-39's mnemonic standard.
