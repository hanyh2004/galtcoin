package cli

import (
	"fmt"

	gcli "github.com/urfave/cli"

	"github.com/samoslab/galtcoin/src/api/webrpc"
	"github.com/samoslab/galtcoin/src/cipher"
	"github.com/samoslab/galtcoin/src/wallet"
)

func walletOutputsCmd(cfg Config) gcli.Command {
	name := "walletOutputs"
	return gcli.Command{
		Name:      name,
		Usage:     "Display outputs of specific wallet",
		ArgsUsage: "[wallet file]",
		Description: fmt.Sprintf(`Display outputs of specific wallet, the default
		wallet (%s) will be
		used if no wallet was specified, use ENV 'WALLET_NAME'
		to update default wallet file name, and 'WALLET_DIR' to update
		the default wallet directory`, cfg.FullWalletPath()),
		OnUsageError: onCommandUsageError(name),
		Action:       getWalletOutputsCmd,
	}
}

func addressOutputsCmd() gcli.Command {
	name := "addressOutputs"
	return gcli.Command{
		Name:      name,
		Usage:     "Display outputs of specific addresses",
		ArgsUsage: "[address list]",
		Description: `Display outputs of specific addresses, join multiple addresses with space,
        example: addressOutputs $addr1 $addr2 $addr3`,
		OnUsageError: onCommandUsageError(name),
		Action:       getAddressOutputsCmd,
	}

}

func getWalletOutputsCmd(c *gcli.Context) error {
	cfg := ConfigFromContext(c)
	rpcClient := RPCClientFromContext(c)

	w := ""
	if c.NArg() > 0 {
		w = c.Args().First()
	}

	var err error
	w, err = resolveWalletPath(cfg, w)
	if err != nil {
		return err
	}

	outputs, err := GetWalletOutputsFromFile(rpcClient, w)
	if err != nil {
		return err
	}

	return printJSON(outputs)
}

func getAddressOutputsCmd(c *gcli.Context) error {
	rpcClient := RPCClientFromContext(c)

	addrs := make([]string, c.NArg())
	var err error
	for i := 0; i < c.NArg(); i++ {
		addrs[i] = c.Args().Get(i)
		if _, err = cipher.DecodeBase58Address(addrs[i]); err != nil {
			return fmt.Errorf("invalid address: %v, err: %v", addrs[i], err)
		}
	}

	outputs, err := rpcClient.GetUnspentOutputs(addrs)
	if err != nil {
		return err
	}

	return printJSON(outputs)
}

// PUBLIC

// GetWalletOutputsFromFile returns unspent outputs associated with all addresses in a wallet file
func GetWalletOutputsFromFile(c *webrpc.Client, walletFile string) (*webrpc.OutputsResult, error) {
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return nil, err
	}

	return GetWalletOutputs(c, wlt)
}

// GetWalletOutputs returns unspent outputs associated with all addresses in a wallet.Wallet
func GetWalletOutputs(c *webrpc.Client, wlt *wallet.Wallet) (*webrpc.OutputsResult, error) {
	cipherAddrs := wlt.GetAddresses()
	addrs := make([]string, len(cipherAddrs))
	for i := range cipherAddrs {
		addrs[i] = cipherAddrs[i].String()
	}

	return c.GetUnspentOutputs(addrs)
}
