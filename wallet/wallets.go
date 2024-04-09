// 第一个文件

package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

const walletFile = "./tmp/wallets.json"

type Wallets struct {
	Wallets map[string]*Wallet `json:"wallets"`
}

func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFile()
	return &wallets, err
}

func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())
	ws.Wallets[address] = wallet
	return address
}

func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

func (ws *Wallets) GetWallet(address string) *Wallet {
	return ws.Wallets[address]
}

func (ws *Wallets) LoadFile() error {
	fileContent, err := os.ReadFile(walletFile)
	if err != nil {
		if os.IsNotExist(err) {
			ws.Wallets = make(map[string]*Wallet)
			return nil
		}
		return err
	}

	lines := bytes.Split(fileContent, []byte("\n"))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		if err := json.Unmarshal(line, ws); err != nil {
			return err
		}
	}

	return nil
}

func (ws *Wallets) SaveFile() error {
	file, err := os.OpenFile(walletFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := json.Marshal(ws)
	if err != nil {
		return err
	}

	// 添加换行符以区分不同的JSON对象
	content = append(content, '\n')

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}
