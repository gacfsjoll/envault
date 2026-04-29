# envault

A local secrets manager for development environments that syncs `.env` files with encrypted storage.

---

## Installation

```bash
go install github.com/yourname/envault@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envault.git && cd envault && go build -o envault .
```

---

## Usage

**Initialize a vault in your project:**

```bash
envault init
```

**Add a secret:**

```bash
envault set DATABASE_URL "postgres://localhost:5432/mydb"
```

**Sync secrets to your `.env` file:**

```bash
envault sync
```

**Load secrets into a command's environment:**

```bash
envault run -- go run main.go
```

Secrets are encrypted at rest using AES-256 and stored in `~/.envault/`. Each project vault is identified by its directory path.

---

## Configuration

Envault looks for a `.envault.toml` file in your project root for project-specific settings such as vault name and sync targets.

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)