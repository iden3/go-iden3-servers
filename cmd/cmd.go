package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/identity/issuer"
	babykeystore "github.com/iden3/go-iden3-core/keystore"
	zkutils "github.com/iden3/go-iden3-core/utils/zk"
	"github.com/iden3/go-iden3-servers/config"
	"github.com/iden3/go-iden3-servers/loaders"
	shell "github.com/ipfs/go-ipfs-api"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func WithCfg(cmd func(c *cli.Context, cfg *config.Config) error) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		var cfg config.Config
		if err := config.LoadFromCliFlag(c, &cfg); err != nil {
			return err
		}
		return cmd(c, &cfg)
	}
}

// Claim
// func CmdAddClaim(c *cli.Context, cfg *config.Config) error {
// 	indexData := c.Args().Get(0)
// 	outData := c.Args().Get(1)
//
// 	iden, err := loaders.LoadIdentity(cfg)
// 	if err != nil {
// 		return err
// 	}
//
// 	var indexSlot [400 / 8]byte
// 	var dataSlot [496 / 8]byte
// 	if len(indexData) != len(indexSlot) || len(outData) != len(dataSlot) {
// 		return fmt.Errorf(
// 			"Length of indexSlot and dataSlot must be %v and %v respectively",
// 			len(indexSlot), len(dataSlot))
// 	}
// 	copy(indexSlot[:], indexData)
// 	copy(dataSlot[:], outData)
// 	claim := claims.NewClaimBasic(indexSlot, dataSlot)
// 	fmt.Println("clam: " + common3.HexEncode(claim.Entry().Bytes()))
//
// 	err = iden.Manager.AddClaim(claim)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Print("root updated: " + iden.Mt.RootKey().Hex())
//
// 	mp, err := iden.Mt.GenerateProof(claim.Entry().HIndex(), nil)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Print("merkleproof: " + common3.HexEncode(mp.Bytes()))
//
// 	return nil
// }

// func CmdAddClaimsFromFile(c *cli.Context, cfg *config.Config) error {
// 	filepath := c.Args().Get(0)
//
// 	iden, err := loaders.LoadIdentity(cfg)
// 	if err != nil {
// 		return err
// 	}
//
// 	fmt.Print("\n---\nimporting claims\n---\n\n")
// 	// csv file will have the following structure: indexData, noindexData
// 	csvFile, _ := os.Open(filepath)
// 	reader := csv.NewReader(bufio.NewReader(csvFile))
// 	for {
// 		line, error := reader.Read()
// 		if error == io.EOF {
// 			break
// 		} else if error != nil {
// 			log.Fatal(error)
// 		}
//
// 		fmt.Println("importing claim with index: " + line[0] + ", outside index: " + line[1])
//
// 		var indexSlot [400 / 8]byte
// 		var dataSlot [496 / 8]byte
// 		if len(line[0]) != len(indexSlot) || len(line[1]) != len(dataSlot) {
// 			return fmt.Errorf(
// 				"Length of indexSlot and dataSlot must be %v and %v respectively",
// 				len(indexSlot), len(dataSlot))
// 		}
// 		copy(indexSlot[:], line[0])
// 		copy(dataSlot[:], line[1])
// 		claim := claims.NewClaimBasic(indexSlot, dataSlot)
// 		fmt.Println("clam: " + common3.HexEncode(claim.Entry().Bytes()) + "\n")
//
// 		// add claim to merkletree, without updating the root, that will be done on the end of the loop (csv file)
// 		err = iden.Mt.AddClaim(claim)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	fmt.Print("\n---\ngenerating proofs\n---\n\n")
// 	// now, let's generate the proofs
// 	csvFile, _ = os.Open(filepath)
// 	reader = csv.NewReader(bufio.NewReader(csvFile))
// 	for {
// 		line, error := reader.Read()
// 		if error == io.EOF {
// 			break
// 		} else if error != nil {
// 			log.Fatal(error)
// 		}
//
// 		fmt.Println("generating merkleproof of claim with index: " + line[0] + ", outside index: " + line[1])
//
// 		var indexSlot [400 / 8]byte
// 		var dataSlot [496 / 8]byte
// 		if len(line[0]) != len(indexSlot) || len(line[1]) != len(dataSlot) {
// 			return fmt.Errorf(
// 				"Length of indexSlot and dataSlot must be %v and %v respectively",
// 				len(indexSlot), len(dataSlot))
// 		}
// 		copy(indexSlot[:], line[0])
// 		copy(dataSlot[:], line[1])
// 		claim := claims.NewClaimBasic(indexSlot, dataSlot)
// 		fmt.Println("clam: " + common3.HexEncode(claim.Entry().Bytes()))
//
// 		// the proofs better generate them once all claims are added
// 		mp, err := iden.Mt.GenerateProof(claim.Entry().HIndex(), nil)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Println("merkleproof: " + common3.HexEncode(mp.Bytes()) + "\n")
// 	}
// 	// update the root in the smart contract
// 	iden.StateWriter.SetRoot(*iden.Mt.RootKey())
// 	fmt.Println("merkletree root: " + iden.Mt.RootKey().Hex())
//
// 	return nil
// }

// DB
func CmdDbRawDump(c *cli.Context, storagePath string) error {
	storage, err := loaders.LoadStorage(storagePath)
	if err != nil {
		return err
	}
	ldb := (storage.(*db.LevelDbStorage)).LevelDB()
	iter := ldb.NewIterator(nil, nil)
	for iter.Next() {
		fmt.Println(common3.HexEncode(iter.Key()) + ", " + common3.HexEncode(iter.Value()))
	}
	iter.Release()
	return nil
}

func CmdDbRawImport(c *cli.Context, storagePath string) error {
	path := c.Args().Get(0)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println("importing raw dump from file " + path)

	count := 0

	storage, err := loaders.LoadStorage(storagePath)
	if err != nil {
		return err
	}
	tx, err := storage.NewTx()
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Close()
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ", ")
		if len(line) < 2 {
			fmt.Println("error in line ", strconv.Itoa(count))
			return err
		}

		var kBytes []byte
		kBytes, err = common3.HexDecode(line[0])
		if err != nil {
			return err
		}
		var vBytes []byte
		vBytes, err = common3.HexDecode(line[1])
		if err != nil {
			return err
		}
		tx.Put(kBytes, vBytes)
		count++
	}
	fmt.Println("imported " + strconv.Itoa(count) + " lines")
	return nil
}

func CmdDbIPFSexport(c *cli.Context, storagePath string) error {
	storage, err := loaders.LoadStorage(storagePath)
	if err != nil {
		return err
	}
	ldb := (storage.(*db.LevelDbStorage)).LevelDB()
	iter := ldb.NewIterator(nil, nil)
	for iter.Next() {
		sh := shell.NewShell("localhost:5001") // ipfs daemon IP:Port
		cid, err := sh.Add(bytes.NewReader(iter.Value()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s", err)
			os.Exit(1)
		}
		fmt.Println("value of key "+common3.HexEncode(iter.Key())+" added, ipfs hash: ", cid)
	}
	iter.Release()
	return nil
}

func NewIssuer(storagePath, keyStoreBabyPath, keyStoreBabyPassword string, confirmBlocks uint64) error {
	// Open babyjub keystore
	params := babykeystore.StandardKeyStoreParams
	keyStoreStorage := babykeystore.NewFileStorage(keyStoreBabyPath)
	keyStore, err := babykeystore.NewKeyStore(keyStoreStorage, params)
	if err != nil {
		return err
	}
	defer keyStore.Close()

	// Create babyjub keys
	kOpComp, err := keyStore.NewKey([]byte(keyStoreBabyPassword))
	if err = keyStore.UnlockKey(kOpComp, []byte(keyStoreBabyPassword)); err != nil {
		return err
	}

	// Create the Issuer in a memory db and later transfer it to the storage under the identity prefix
	memStorage := db.NewMemoryStorage()
	cfg := issuer.ConfigDefault
	cfg.ConfirmBlocks = confirmBlocks
	id, err := issuer.Create(cfg, kOpComp, nil, memStorage, keyStore)
	if err != nil {
		return err
	}
	storage, err := loaders.LoadStorage(storagePath)
	if err != nil {
		return err
	}
	idenStorage := storage.WithPrefix([]byte(fmt.Sprintf("%v:", id)))
	tx, err := idenStorage.NewTx()
	if err != nil {
		return err
	}
	memStorage.Iterate(func(k []byte, v []byte) (bool, error) {
		tx.Put(k, v)
		return true, nil
	})
	if err := tx.Commit(); err != nil {
		return err
	}

	var config struct {
		Identity config.Identity
	}
	config.Identity.Id = *id
	kOp, err := kOpComp.Decompress()
	if err != nil {
		return err
	}
	config.Identity.Keys.BabyJub.KOp = *kOp

	var configTOML bytes.Buffer
	if err := toml.NewEncoder(&configTOML).Encode(&config); err != nil {
		return nil
	}

	fmt.Fprintf(os.Stderr, "Keys and identity created successfully."+
		" Copy & paste the lines between '---' into the config file:\n---\n")
	fmt.Print(configTOML.String())
	fmt.Fprintf(os.Stderr, "---\n")
	return nil
}

func CmdNewIssuer(c *cli.Context) error {
	var cfg struct {
		KeyStore     config.KeyStore  `validate:"required"`
		KeyStoreBaby config.KeyStore  `validate:"required"`
		Contracts    config.Contracts `validate:"required"`
		Storage      struct {
			Path string
		} `validate:"required"`
		Issuer struct {
			ConfirmBlocks uint64 `validate:"required"`
		}
	}
	if err := config.LoadFromCliFlag(c, &cfg); err != nil {
		return err
	}
	return NewIssuer(cfg.Storage.Path, cfg.KeyStoreBaby.Path, cfg.KeyStoreBaby.Password.Value,
		cfg.Issuer.ConfirmBlocks)
}

func CmdStop(c *cli.Context, cfg *config.Config) error {
	if err := PostAdminApi(&cfg.Server, "stop", nil); err != nil {
		return err
	}
	return nil
}

// func CmdInfo(c *cli.Context, cfg *config.Config) error {
// 	output, err := PostAdminApi(cfg.Server.AdminApi, "info")
// 	if err == nil {
// 		log.Info("Server response: ", output)
// 	}
// 	return err
// }

func CmdSync(c *cli.Context, cfg *config.Config) error {
	if err := PostAdminApi(&cfg.Server, "issuer/syncidenstatepublic", nil); err != nil {
		return err
	}
	return nil
}

func CmdStart(c *cli.Context, cfg *config.Config, endpointServe func(cfg *config.Config, srv *loaders.Server)) error {
	srv, err := loaders.LoadServer(cfg)
	if err != nil {
		return err
	}

	// Check for funds
	balance, err := srv.EthClient.BalanceAt(srv.EthClient.Account().Address)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"balance": balance.String(),
		"address": srv.EthClient.Account().Address.Hex(),
	}).Info("Account balance retrieved")
	if balance.Cmp(new(big.Int).SetUint64(3000000)) == -1 {
		return fmt.Errorf("Not enough funds in the ethereum address")
	}

	srv.Start()

	endpointServe(cfg, srv)

	srv.StopAndJoin()

	return nil
}

func CmdImportEthAccount(c *cli.Context) error {
	keyfile := c.Args().First()
	if len(keyfile) == 0 {
		return fmt.Errorf("keyfile must be given as argument")
	}
	key, err := crypto.LoadECDSA(keyfile)
	if err != nil {
		return err
	}

	var cfg struct {
		KeyStore config.KeyStore `validate:"required"`
	}
	if err := config.LoadFromCliFlag(c, &cfg); err != nil {
		return err
	}

	ks := keystore.NewKeyStore(cfg.KeyStore.Path, keystore.StandardScryptN, keystore.StandardScryptP)

	account, err := ks.ImportECDSA(key, cfg.KeyStore.Password.Value)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"keyfile": keyfile,
		"path":    cfg.KeyStore.Path,
		"address": account.Address.Hex(),
	}).Info("Imported Eth Account from private key")

	fmt.Fprintf(os.Stderr, "Ethereum Account imported successfully."+
		" Copy & paste the lines between '---' into the config file:\n")
	printAccountToml(account)

	return nil
}

func CmdNewEthAccount(c *cli.Context) error {
	var cfg struct {
		KeyStore config.KeyStore `validate:"required"`
	}
	if err := config.LoadFromCliFlag(c, &cfg); err != nil {
		return err
	}

	ks := keystore.NewKeyStore(cfg.KeyStore.Path, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(cfg.KeyStore.Password.Value)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"path":    cfg.KeyStore.Path,
		"address": account.Address.Hex(),
	}).Info("Created new Eth Account")

	fmt.Fprintf(os.Stderr, "Ethereum Account created successfully."+
		" Copy & paste the lines between '---' into the config file:\n")
	printAccountToml(account)

	return nil
}

func CmdDeployState(c *cli.Context) error {
	var cfg struct {
		Web3     config.Web3     `validate:"required"`
		KeyStore config.KeyStore `validate:"required"`
		Account  struct {
			Address common.Address `validate:"required"`
		} `validate:"required"`
	}
	if err := config.LoadFromCliFlag(c, &cfg); err != nil {
		return err
	}
	ks, acc, err := loaders.LoadKeyStore(&cfg.KeyStore, &cfg.Account.Address)
	if err != nil {
		return err
	}
	ethClient, err := loaders.LoadEthClient(ks, acc, cfg.Web3.Url)
	if err != nil {
		return err
	}
	result, err := idenpubonchain.DeployState(ethClient, nil)
	if err != nil {
		return err
	}

	var cfgOut struct {
		Contracts config.Contracts `validate:"required"`
	}
	cfgOut.Contracts.IdenStates.Address = result.State.Address
	var cfgOutTOML bytes.Buffer
	if err := toml.NewEncoder(&cfgOutTOML).Encode(&cfgOut); err != nil {
		log.Error(err)
	}

	fmt.Fprintf(os.Stderr, "\n---\n")
	fmt.Print(cfgOutTOML.String())
	fmt.Fprintf(os.Stderr, "---\n")

	return nil
}

func printAccountToml(account accounts.Account) {
	var cfg struct {
		Account struct {
			Address common.Address
		}
	}
	cfg.Account.Address = account.Address

	var cfgTOML bytes.Buffer
	if err := toml.NewEncoder(&cfgTOML).Encode(&cfg); err != nil {
		log.Error(err)
	}

	fmt.Fprintf(os.Stderr, "\n---\n")
	fmt.Print(cfgTOML.String())
	fmt.Fprintf(os.Stderr, "---\n")
}

func ValidateZKFilesProvingKeyFormat(format string) (zkutils.ProvingKeyFormat, error) {
	match := false
	provingKeyFormat := zkutils.ProvingKeyFormat(format)
	for _, f := range []zkutils.ProvingKeyFormat{zkutils.ProvingKeyFormatJSON,
		zkutils.ProvingKeyFormatBin, zkutils.ProvingKeyFormatGoBin} {
		if provingKeyFormat == f {
			match = true
			break
		}
	}
	if !match {
		return "", fmt.Errorf("invalid ProvingKeyFormat %v", provingKeyFormat)
	}
	return provingKeyFormat, nil
}

func CmdDownloadZKFiles(c *cli.Context) error {
	url := c.GlobalString("url")
	if url == "" {
		return fmt.Errorf("No url specified")
	}
	path := c.GlobalString("path")
	if path == "" {
		return fmt.Errorf("No path specified")
	}
	format, err := ValidateZKFilesProvingKeyFormat(c.GlobalString("format"))
	if err != nil {
		return err
	}
	zkfiles := zkutils.NewZkFiles(url, path, format, zkutils.ZkFilesHashes{}, false)
	if err := zkfiles.InsecureDownloadAll(); err != nil {
		return err
	}
	return nil
}

func CmdHashZKFiles(c *cli.Context) error {
	path := c.GlobalString("path")
	if path == "" {
		return fmt.Errorf("No path specified")
	}
	format, err := ValidateZKFilesProvingKeyFormat(c.GlobalString("format"))
	if err != nil {
		return err
	}
	zkfiles := zkutils.NewZkFiles("", path, format, zkutils.ZkFilesHashes{}, false)
	zkFilesHashes, err := zkfiles.InsecureCalcHashes()
	if err != nil {
		return err
	}

	var cfg struct {
		Files config.ZkFiles
	}
	cfg.Files.Path = path
	cfg.Files.Hashes = config.ZkFilesHashes{
		ProvingKey:      zkFilesHashes.ProvingKey,
		VerificationKey: zkFilesHashes.VerificationKey,
		WitnessCalcWASM: zkFilesHashes.WitnessCalcWASM,
	}

	var cfgTOML bytes.Buffer
	if err := toml.NewEncoder(&cfgTOML).Encode(&cfg.Files.Hashes); err != nil {
		log.Error(err)
	}

	fmt.Fprintf(os.Stderr, "ZkFiles hashes calculated successfully."+
		" Copy & paste the lines between '---' into the config file:\n")
	fmt.Fprintf(os.Stderr, "\n---\n")
	fmt.Print(cfgTOML.String())
	fmt.Fprintf(os.Stderr, "---\n")
	return nil
}

func CmdCheckZKFiles(c *cli.Context) error {
	url := c.GlobalString("url")
	if url == "" {
		return fmt.Errorf("No url specified")
	}
	path := c.GlobalString("path")
	if path == "" {
		return fmt.Errorf("No path specified")
	}
	format, err := ValidateZKFilesProvingKeyFormat(c.GlobalString("format"))
	if err != nil {
		return err
	}
	zkfiles := zkutils.NewZkFiles(url, path, format, zkutils.ZkFilesHashes{}, false)
	if err := zkfiles.InsecureDownloadAll(); err != nil {
		return err
	}
	return nil
}
