# GoLockbox

GoLockbox is a secure, command-line password and file manager written in Go. It stores files inside an **encrypted vault** protected by a master password, providing a fully interactive **Terminal User Interface (TUI)** that requires no GUI or external dependencies beyond the Go toolchain.

---

## Table of Contents

- [Features](#features)
- [How It Works](#how-it-works)
  - [Vault Structure on Disk](#vault-structure-on-disk)
  - [Encryption Model](#encryption-model)
  - [Application Flow](#application-flow)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Building](#building)
- [Running](#running)
- [Keyboard Reference](#keyboard-reference)
- [Testing](#testing)
- [Developer Guide](#developer-guide)
  - [Adding a New Vault Operation](#adding-a-new-vault-operation)
  - [Adding a New TUI Screen](#adding-a-new-tui-screen)
  - [Modifying Cryptographic Parameters](#modifying-cryptographic-parameters)
- [Security Details](#security-details)
- [Dependencies](#dependencies)
- [License](#license)

---

## Features

- Create and open **encrypted vaults** protected by a master password
- **View, create, edit, rename, delete**, and **export** files inside a vault
- **Import** any file from disk directly into the vault
- All vault content (file data and metadata) is encrypted at rest with **ChaCha20-Poly1305**
- Password-based key derivation using **Argon2id**
- Master key is **zeroed from memory** when the vault closes
- Duplicate filename handling — imported files are automatically renamed to avoid collisions
- `~/` path expansion for convenient vault paths
- Works entirely in the terminal — no GUI required

---

## How It Works

### Vault Structure on Disk

```
<vault-directory>/
├── master.key      # Encrypted master key (version byte + salt + ciphertext)
├── entries.enc     # Encrypted JSON file index (list of file names)
└── files/
    ├── secret.txt  # Each file stored in its own encrypted blob
    ├── notes.md
    └── ...
```

`master.key` layout:

```
[1 byte: format version][16 bytes: Argon2id salt][ChaCha20-Poly1305 ciphertext of master key]
```

`entries.enc` is a JSON document (`{"version":1,"entries":[{"name":"secret.txt"},…]}`)
encrypted with the master key. The TUI reads it on every open to list available files.

### Encryption Model

1. **Password → key (Argon2id):** When creating or opening a vault the user's password is stretched into a 32-byte derived key with Argon2id (time=1, memory=64 MB, threads=4).
2. **Derived key → master key:** The derived key encrypts (or decrypts) a random 32-byte master key stored in `master.key`.
3. **Master key → file data:** Every individual file and the `entries.enc` metadata file are independently encrypted with the master key using ChaCha20-Poly1305 (an authenticated encryption scheme that provides both confidentiality and integrity).
4. **Memory safety:** The derived key and master key are zeroed (`zeroBytes`) as soon as they are no longer needed.

### Application Flow

```
┌──────────┐    choose      ┌─────────────┐   enter path  ┌──────────────┐
│   Menu   │ ─────────────► │  Vault Path │ ────────────► │   Password   │
└──────────┘                └─────────────┘               └──────┬───────┘
                                                                  │ unlock
                                                                  ▼
              ┌──────────────────────────────────────────────────────────────┐
              │                       File Browser                           │
              │   ↑/↓ navigate  ·  Enter view  ·  e edit  ·  n new          │
              │   a add  ·  d delete  ·  r rename  ·  b back  ·  Ctrl+C quit │
              └─────────────────────────────────────────────────────────────┘
                        │           │           │
                        ▼           ▼           ▼
                  File Viewer   File Editor   Rename / Add / New screens
```

---

## Project Structure

```
GoLockbox/
├── cmd/GoLockbox/
│   └── main.go                      # Entry point — starts the Bubble Tea program
├── internal/
│   ├── crypto/
│   │   ├── encryption.go            # EncryptBytes / DecryptBytes (ChaCha20-Poly1305)
│   │   ├── errors.go                # Crypto-specific error values
│   │   ├── key.go                   # GenerateRandomKey, GenerateSalt
│   │   ├── key_derivation.go        # DeriveKey (Argon2id)
│   │   └── password_gen.go          # GeneratePassword utility
│   ├── vault/
│   │   ├── constants.go             # File/format version constants
│   │   ├── errors.go                # Vault-specific error values (ErrInvalidPassword, …)
│   │   ├── files.go                 # AddFile, OpenFile, SaveFile, DeleteFile, RenameFile, ExportFile, CreateFile
│   │   ├── path.go                  # expandPath (tilde expansion, separator normalization)
│   │   ├── types.go                 # FileEntry, metadata structs
│   │   ├── vault.go                 # CreateVault, OpenVault, ListFiles, SaveVault, Close
│   │   └── tests/
│   │       ├── vault_test.go        # Core functionality tests
│   │       ├── savevault_test.go    # Persistence / save-vault tests
│   │       ├── vault_property_test.go # Property-based tests
│   │       ├── vault_fuzz_test.go   # Fuzz tests
│   │       └── more_tests.go        # Additional edge-case tests
│   └── tui/
│       ├── root.go                  # RootModel — orchestrates all screens
│       ├── types.go                 # ScreenID constants and inter-screen message types
│       ├── menu.go                  # Main menu (New Vault / Open Vault)
│       ├── vault_path.go            # Vault path input screen
│       ├── password.go              # Password input screen
│       ├── browser.go               # File list / browser screen
│       ├── file_view.go             # Read-only file viewer
│       ├── editor.go                # In-terminal file editor
│       ├── add_file.go              # Import existing file into vault
│       ├── rename.go                # Rename file screen
│       └── new_file_name.go         # New file name input screen
├── go.mod
├── go.sum
└── LICENSE
```

---

## Prerequisites

- **Go 1.21 or later** (the module declares `go 1.25`)

No other tools or external libraries need to be installed manually — `go build` downloads all dependencies automatically.

---

## Building

```bash
# Clone the repository (if you haven't already)
git clone https://github.com/emmyie/GoLockbox.git
cd GoLockbox

# Build a binary named GoLockbox in the current directory
go build -o GoLockbox ./cmd/GoLockbox

# Or just verify it compiles (no output binary)
go build ./...
```

The resulting `GoLockbox` binary has no runtime dependencies and can be copied to any directory in your `$PATH`.

---

## Running

```bash
./GoLockbox
```

The interactive TUI starts immediately. No flags or arguments are required.

**Typical first-run workflow:**

1. Select **New Vault** with ↑/↓ and press Enter.
2. Type a directory path for the vault (e.g. `~/my-vault`) and press Enter.
3. Enter a strong master password and press Enter.
4. You are now in the **file browser**. Use the keys listed below to manage files.

**Opening an existing vault:**

1. Select **Open Vault**.
2. Enter the path to the vault directory you created previously.
3. Enter the master password.

---

## Keyboard Reference

### Main Menu / Path / Password screens

| Key | Action |
|-----|--------|
| ↑ / ↓ | Navigate |
| Enter | Confirm selection / input |
| Ctrl+C | Quit |

### File Browser

| Key | Action |
|-----|--------|
| ↑ / ↓ | Move cursor |
| Enter | View selected file |
| `e` | Edit selected file |
| `n` | Create a new empty file |
| `a` | Add (import) an existing file from disk |
| `d` | Delete selected file |
| `r` | Rename selected file |
| `b` | Return to main menu |
| Ctrl+C | Save vault and quit |

### File Viewer

| Key | Action |
|-----|--------|
| Esc | Back to browser |
| Ctrl+C | Quit |

### File Editor

| Key | Action |
|-----|--------|
| Ctrl+S | Save changes |
| Esc | Discard and return to browser |
| Ctrl+C | Quit |

---

## Testing

The test suite lives in `internal/vault/tests/` and covers vault creation, file operations, metadata persistence, property-based invariants, and fuzz testing.

```bash
# Run all tests
go test ./...

# Verbose output
go test ./... -v

# Run only vault tests
go test ./internal/vault/tests/... -v

# Run the fuzz tests (requires Go 1.18+)
go test ./internal/vault/tests/ -fuzz=FuzzDecryptFile -fuzztime=30s
```

---

## Developer Guide

### Adding a New Vault Operation

1. **Implement the method** in `internal/vault/files.go` (file operations) or `internal/vault/vault.go` (lifecycle operations):

   ```go
   // ExampleOperation demonstrates how to add a vault operation.
   func (v *Vault) ExampleOperation(index int) error {
       if index < 0 || index >= len(v.entries) {
           return fmt.Errorf("index out of range")
       }
       // ... do work ...
       return v.saveMetadata() // persist any metadata changes
   }
   ```

2. **Write a test** in `internal/vault/tests/vault_test.go` following the existing patterns:

   ```go
   func TestExampleOperation(t *testing.T) {
       dir := tempDir(t)
       v, err := vault.CreateVault(dir, "password")
       if err != nil {
           t.Fatal(err)
       }
       // ... exercise the operation ...
   }
   ```

3. **Wire it into the TUI** if you want it to be accessible from the terminal interface (see next section).

### Adding a New TUI Screen

The TUI is built with [Bubble Tea](https://github.com/charmbracelet/bubbletea). Every screen is an independent `tea.Model`.

1. **Define a screen ID** in `internal/tui/types.go`:

   ```go
   const (
       // existing IDs …
       ScreenMyFeature  // add yours here
   )
   ```

2. **Create the model** in a new file, e.g. `internal/tui/my_feature.go`:

   ```go
   package tui

   import tea "github.com/charmbracelet/bubbletea"

   type MyFeatureModel struct{ /* state fields */ }

   func NewMyFeatureModel() MyFeatureModel { return MyFeatureModel{} }

   func (m MyFeatureModel) Init() tea.Cmd { return nil }

   func (m MyFeatureModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
       // handle key presses, return ScreenMsg to navigate away
       return m, nil
   }

   func (m MyFeatureModel) View() string {
       return "My feature screen"
   }
   ```

3. **Register the screen** in `internal/tui/root.go`:
   - Add a field to `RootModel`.
   - Handle `ScreenMyFeature` in the `View()` switch.
   - Handle the incoming `ScreenMsg` that transitions to this screen in `Update()`.

### Modifying Cryptographic Parameters

**Argon2id parameters** (key derivation cost) are in `internal/crypto/key_derivation.go`:

```go
// Increase time/memory for stronger protection at the cost of slower unlock.
argon2.IDKey([]byte(password), salt, 1 /*time*/, 64*1024 /*memory KiB*/, 4 /*threads*/, KeySize)
```

Changing these parameters will make existing vaults **unreadable** because the derived key will differ. Only change them before creating new vaults or implement a migration path.

**Encryption algorithm** is chosen in `internal/crypto/encryption.go`. The current scheme (ChaCha20-Poly1305) is a modern authenticated cipher and should not need to be changed.

---

## Security Details

| Property | Detail |
|----------|--------|
| Symmetric cipher | ChaCha20-Poly1305 (AEAD) |
| Key derivation | Argon2id — time=1, memory=64 MiB, parallelism=4 |
| Master key size | 256 bits (32 bytes) |
| Salt size | 128 bits (16 bytes), random per vault |
| Key isolation | Derived key is zeroed after use; master key is zeroed on `Close()` |
| File permissions | Vault directory and files created with mode `0700` / `0600` |
| Metadata integrity | `entries.enc` is authenticated with the same AEAD as file data |

Each encrypted blob includes an authentication tag — any tampering with vault files is detected and rejected.

---

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/charmbracelet/bubbletea` | TUI event loop and program model |
| `github.com/charmbracelet/bubbles` | Reusable TUI components (text input, list, viewport) |
| `golang.org/x/crypto` | Argon2id, ChaCha20-Poly1305 |

All other entries in `go.sum` are transitive dependencies pulled in by the above.

---

## License

See [LICENSE](LICENSE).
