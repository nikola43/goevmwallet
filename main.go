package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nikola43/gocypher/cypher"
	"golang.org/x/crypto/sha3"
)

type Wallet struct {
	Name       string `json:"Name"`
	Address    string `json:"Address"`
	PrivateKey string `json:"PrivateKey"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	password := getPassword(os.Args)
	if password == "" {
		fmt.Println("Error: Password is required. Use -p or --password to specify the password.")
		printUsage()
		return
	}

	c := cypher.NewCypher(password)
	walletFile := "wallets.dat"

	switch command {
	case "create":
		args := filterPasswordArgs(os.Args[2:])
		var name string
		if len(args) > 0 {
			name = args[0]
		}
		handleCreate(walletFile, c, name)
	case "readAll":
		handleReadAll(walletFile, c)
	case "read":
		args := filterPasswordArgs(os.Args[2:])
		if len(args) < 1 {
			fmt.Println("Error: Please provide a name or address to read")
			return
		}
		handleRead(walletFile, c, args[0])
	case "delete":
		args := filterPasswordArgs(os.Args[2:])
		if len(args) < 1 {
			fmt.Println("Error: Please provide a name or address to delete")
			return
		}
		handleDelete(walletFile, c, args[0])
	default:
		fmt.Printf("Error: Unknown command '%s'\n", command)
		printUsage()
	}
}

func getPassword(args []string) string {
	for i, arg := range args {
		if (arg == "-p" || arg == "--password") && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func filterPasswordArgs(args []string) []string {
	filtered := make([]string, 0)
	for i := 0; i < len(args); i++ {
		if args[i] == "-p" || args[i] == "--password" {
			i++ // Skip the next argument (password value)
			continue
		}
		filtered = append(filtered, args[i])
	}
	return filtered
}

func printUsage() {
	fmt.Println("\nUsage:")
	fmt.Println("  goevmwallet [command] -p <password> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  create [name]     - Create new wallet with optional name")
	fmt.Println("  readAll           - Show all wallets")
	fmt.Println("  read <nameOrAddr> - Show specific wallet by name or address")
	fmt.Println("  delete <nameOrAddr> - Delete wallet by name or address")
	fmt.Println("\nRequired:")
	fmt.Println("  -p, --password <password>   - Password for encryption/decryption")
	fmt.Println("\nExamples:")
	fmt.Println("  goevmwallet create -p mypass123 mywallet")
	fmt.Println("  goevmwallet readAll -p mypass123")
	fmt.Println("  goevmwallet read -p mypass123 0x123...")
	fmt.Println("  goevmwallet delete -p mypass123 mywallet")
	fmt.Println("")
}

func handleCreate(walletFile string, c *cypher.Cypher, name string) {
	wallets, err := loadWallets(walletFile, c)
	if err != nil {
		if strings.Contains(err.Error(), "cipher: message authentication failed") {
			fmt.Println("Error: Invalid password or corrupted wallet file")
		} else {
			fmt.Println("Error loading wallets:", err)
		}
		return
	}

	wallet, err := generateWallet()
	if err != nil {
		fmt.Println("Error generating wallet:", err)
		return
	}

	wallet.Name = name
	wallets = append(wallets, wallet)

	err = saveWallets(wallets, walletFile, c)
	if err != nil {
		fmt.Println("Error saving wallets:", err)
		return
	}

	fmt.Println("\nWallet created successfully:")
	printWallet(wallet)
}

func handleReadAll(walletFile string, c *cypher.Cypher) {
	wallets, err := loadWallets(walletFile, c)
	if err != nil {
		if strings.Contains(err.Error(), "cipher: message authentication failed") {
			fmt.Println("Error: Invalid password or corrupted wallet file")
		} else {
			fmt.Println("Error loading wallets:", err)
		}
		return
	}

	if len(wallets) == 0 {
		fmt.Println("No wallets found")
		return
	}

	fmt.Printf("\nTotal wallets: %d\n\n", len(wallets))
	for i, wallet := range wallets {
		fmt.Printf("Wallet #%d:\n", i+1)
		printWallet(wallet)
	}
}

func handleRead(walletFile string, c *cypher.Cypher, query string) {
	wallets, err := loadWallets(walletFile, c)
	if err != nil {
		if strings.Contains(err.Error(), "cipher: message authentication failed") {
			fmt.Println("Error: Invalid password or corrupted wallet file")
		} else {
			fmt.Println("Error loading wallets:", err)
		}
		return
	}

	found := false
	for _, wallet := range wallets {
		if strings.EqualFold(wallet.Name, query) || strings.EqualFold(wallet.Address, query) {
			printWallet(wallet)
			found = true
		}
	}

	if !found {
		fmt.Println("No wallet found with that name or address")
	}
}

func handleDelete(walletFile string, c *cypher.Cypher, query string) {
	wallets, err := loadWallets(walletFile, c)
	if err != nil {
		if strings.Contains(err.Error(), "cipher: message authentication failed") {
			fmt.Println("Error: Invalid password or corrupted wallet file")
		} else {
			fmt.Println("Error loading wallets:", err)
		}
		return
	}

	var updatedWallets []*Wallet
	deleted := false

	for _, wallet := range wallets {
		if !strings.EqualFold(wallet.Name, query) && !strings.EqualFold(wallet.Address, query) {
			updatedWallets = append(updatedWallets, wallet)
		} else {
			deleted = true
		}
	}

	if !deleted {
		fmt.Println("No wallet found with that name or address")
		return
	}

	err = saveWallets(updatedWallets, walletFile, c)
	if err != nil {
		fmt.Println("Error saving wallets:", err)
		return
	}

	fmt.Println("Wallet deleted successfully")
}

func printWallet(wallet *Wallet) {
	if wallet.Name != "" {
		fmt.Printf("Name: %s\n", wallet.Name)
	}
	fmt.Printf("Address: %s\n", wallet.Address)
	fmt.Printf("PrivateKey: %s\n\n", wallet.PrivateKey)
}

func loadWallets(walletFile string, c *cypher.Cypher) ([]*Wallet, error) {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return make([]*Wallet, 0), nil
	}

	data, err := loadFile(walletFile)
	if err != nil {
		return nil, err
	}

	decrypted, err := c.Decrypt(data)
	if err != nil {
		return nil, err
	}

	var wallets []*Wallet
	err = json.Unmarshal(decrypted, &wallets)
	if err != nil {
		return nil, err
	}

	return wallets, nil
}

func loadFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func saveWallets(wallets []*Wallet, walletFile string, c *cypher.Cypher) error {
	walletsJSON, err := json.Marshal(wallets)
	if err != nil {
		return err
	}

	encryptedWallets, err := c.Encrypt(walletsJSON)
	if err != nil {
		return err
	}

	err = saveFile(walletFile, encryptedWallets)
	if err != nil {
		return err
	}

	return nil
}

func saveFile(fileName string, data []byte) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func generateWallet() (*Wallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])

	wallet := &Wallet{
		Address:    address,
		PrivateKey: hexutil.Encode(privateKeyBytes)[2:],
	}

	return wallet, nil
}