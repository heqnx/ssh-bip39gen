package main

import (
    "crypto/sha256"
    "encoding/pem"
    "flag"
    "fmt"
    "os"
    "strings"

    "github.com/tyler-smith/go-bip39"
    "golang.org/x/crypto/ed25519"
    "golang.org/x/crypto/ssh"
)

const helpMessage = `ssh-bip39gen: generate dd25519 ssh keys from a bip-39 mnemonic

this tool creates deterministic ed25519 ssh key pairs using a 24-word bip-39 mnemonic phrase.
if no mnemonic is provided, it generates a new one with 256-bit entropy. the mnemonic is your
key to regenerate the same SSH key pair later - keep it safe!

usage:
  ssh-bip39gen [-f output_file] [-mnemonic "24-word phrase"]

examples:
  generate a new key pair (default: id_ed25519):
    ssh-bip39gen

  generate with custom output file:
    ssh-bip39gen -f test

  regenerate from a mnemonic:
    ssh-bip39gen -f test -mnemonic "abandon ability able about ... actress"

flags:
`

type KeyType string

const (
    ED25519 KeyType = "ed25519"
)

type seededRand struct {
    seed []byte
    pos  int
}

func (r *seededRand) Read(p []byte) (n int, err error) {
    for n < len(p) {
        counter := uint32(r.pos / len(r.seed))
        hash := sha256.Sum256(append(r.seed, byte(counter), byte(counter>>8), byte(counter>>16), byte(counter>>24)))
        n += copy(p[n:], hash[:])
        r.pos += len(hash)
    }
    return len(p), nil
}

func GenerateEd25519Key(seed []byte) (ed25519.PublicKey, ed25519.PrivateKey, error) {
    if len(seed) < 32 {
        return nil, nil, fmt.Errorf("seed too short for ed25519; need 32 bytes, got %d", len(seed))
    }
    publicKey, privateKey, err := ed25519.GenerateKey(&seededRand{seed: seed[:32]})
    return publicKey, privateKey, err
}

func SaveKeys(privateKey ed25519.PrivateKey, privFile, pubFile string) error {
    block, err := ssh.MarshalPrivateKey(privateKey, "")
    if err != nil {
        return fmt.Errorf("failed to marshal ed25519 private key: %v", err)
    }
    privBytes := pem.EncodeToMemory(block)

    pub, err := ssh.NewPublicKey(privateKey.Public())
    if err != nil {
        return fmt.Errorf("failed to generate ed25519 public key: %v", err)
    }
    pubBytes := ssh.MarshalAuthorizedKey(pub)

    if err := os.WriteFile(privFile, privBytes, 0600); err != nil {
        return fmt.Errorf("failed to write private key to %s: %v", privFile, err)
    }
    if err := os.WriteFile(pubFile, pubBytes, 0644); err != nil {
        return fmt.Errorf("failed to write public key to %s: %v", pubFile, err)
    }
    return nil
}

func main() {
    flag.Usage = func() {
        fmt.Fprint(os.Stderr, helpMessage)
        flag.PrintDefaults()
    }
    
    mnemonic := flag.String("mnemonic", "", "bip-39 mnemonic phrase (leave empty to generate a new 24-word mnemonic)")
    outputFile := flag.String("f", "bip39-id_ed25519", "output file for private key (public key will be <file>.pub)")
    flag.Parse()

    privFile := *outputFile
    pubFile := *outputFile + ".pub"

    var seed []byte
    if *mnemonic == "" {
        entropy, err := bip39.NewEntropy(256)
        if err != nil {
            fmt.Printf("error generating entropy: %v\n", err)
            os.Exit(1)
        }
        *mnemonic, err = bip39.NewMnemonic(entropy)
        if err != nil {
            fmt.Printf("error generating mnemonic: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("new mnemonic generated - keep it secure:\n%s\n\n", *mnemonic)
        seed = bip39.NewSeed(*mnemonic, "")
    } else {
        wordCount := len(strings.Fields(*mnemonic))
        if wordCount != 24 {
            fmt.Printf("error: mnemonic must contain exactly 24 words, got %d\n", wordCount)
            os.Exit(1)
        }
        if !bip39.IsMnemonicValid(*mnemonic) {
            fmt.Println("error: invalid mnemonic phrase")
            os.Exit(1)
        }
        seed = bip39.NewSeed(*mnemonic, "")
    }

    _, priv, err := GenerateEd25519Key(seed)
    if err != nil {
        fmt.Printf("error generating ed25519 key: %v\n", err)
        os.Exit(1)
    }

    err = SaveKeys(priv, privFile, pubFile)
    if err != nil {
        fmt.Printf("error saving keys: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("ed25519 generated keys:\n  - %s (private)\n  - %s (public)\n", privFile, pubFile)
}
