# 🔐 EVM Wallet Manager

A secure command-line tool for managing Ethereum wallets with encryption support. Create, store, and manage multiple Ethereum wallets with names for easy identification.

## ✨ Features

- 🛡️ Password-protected encryption for wallet storage
- 📝 Optional naming for wallets
- 🔍 Search wallets by name or address
- 📋 List all stored wallets
- 🗑️ Delete wallets by name or address
- 💾 Secure local storage in encrypted format

## 🚀 Installation

```bash
go install github.com/yourusername/goevmwallet@latest
```

## 🔧 Usage

### Creating a Wallet
Create a new wallet with an optional name:
```bash
# Create wallet with name
goevmwallet create -p yourpassword123 mywallet

# Create wallet without name
goevmwallet create -p yourpassword123
```

### Reading Wallets
View all wallets or search for specific ones:
```bash
# List all wallets
goevmwallet readAll -p yourpassword123

# Find wallet by name
goevmwallet read -p yourpassword123 mywallet

# Find wallet by address
goevmwallet read -p yourpassword123 0x123...
```

### Deleting Wallets
Remove wallets by name or address:
```bash
# Delete by name
goevmwallet delete -p yourpassword123 mywallet

# Delete by address
goevmwallet delete -p yourpassword123 0x123...
```

## 📝 Command Reference

| Command | Description | Example |
|---------|-------------|---------|
| `create [name]` | Create new wallet with optional name | `goevmwallet create -p pass123 mywallet` |
| `readAll` | Show all stored wallets | `goevmwallet readAll -p pass123` |
| `read <nameOrAddr>` | Find wallet by name or address | `goevmwallet read -p pass123 mywallet` |
| `delete <nameOrAddr>` | Delete wallet by name or address | `goevmwallet delete -p pass123 mywallet` |

## 🔒 Security Features

- 🔑 Mandatory password protection
- 🔐 All wallet data is encrypted on disk
- 🛡️ Private keys are never stored in plain text
- 📜 No default/fallback passwords allowed

## ⚠️ Important Notes

1. **Password Requirements**:
   - Password is mandatory for all operations
   - Use the same password that was used to create the wallet file
   - Store your password securely - if lost, wallet file cannot be decrypted

2. **Backup Recommendations**:
   - Keep secure backups of your wallet file
   - Never share your private keys or password

## 💡 Tips

- Use meaningful names for your wallets to easily identify them
- Regularly backup your wallet file (`wallets.dat`)
- Use a strong password with a mix of letters, numbers, and symbols
- Keep your password safe - without it, you cannot access your wallets

## 🛠️ Technical Details

- Written in Go
- Uses industry-standard encryption
- Stores wallets in local encrypted file
- Compatible with all EVM-based networks

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the GPLV3 License - see the LICENSE file for details.

## ⭐ Star History

If you find this project useful, please consider giving it a star on GitHub! Your support helps us keep maintaining and improving this tool.

## 🐛 Found a Bug?

Please open an issue with:
- Command you were trying to run
- Expected behavior
- Actual behavior
- Steps to reproduce

---
Made with ❤️ for the Ethereum community
