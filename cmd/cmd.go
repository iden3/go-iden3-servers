package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/identity/issuer"
	babykeystore "github.com/iden3/go-iden3-core/keystore"
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

func NewIssuer(storagePath, keyStoreBabyPath, keyStoreBabyPassword string) error {
	// Open babyjub keystore
	params := babykeystore.StandardKeyStoreParams
	keyStoreStorage := babykeystore.NewFileStorage(keyStoreBabyPath)
	keyStore, err := babykeystore.NewKeyStore(keyStoreStorage, params)
	if err != nil {
		return err
	}
	defer keyStore.Close()

	// Create babyjub keys
	kOp, err := keyStore.NewKey([]byte(keyStoreBabyPassword))
	if err = keyStore.UnlockKey(kOp, []byte(keyStoreBabyPassword)); err != nil {
		return err
	}

	// Create the Issuer in a memory db and later transfer it to the storage under the identity prefix
	memStorage := db.NewMemoryStorage()
	cfg := issuer.ConfigDefault
	is, err := issuer.New(cfg, kOp, nil, memStorage, keyStore, nil, nil)
	if err != nil {
		return err
	}
	id := is.ID()
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

	// Verify that the issuer can be loaded successfully
	_, err = loaders.LoadIssuer(id, storage, keyStore, nil, nil)
	if err != nil {
		return err
	}

	var m bytes.Buffer
	fmt.Fprintf(&m, "Id = %v\n", id.String())
	fmt.Fprintf(&m, "KOp = %v\n", kOp.String())

	fmt.Fprintf(os.Stderr, "Keys and identity created successfully."+
		" Copy & paste the lines between '---' into the config file:\n---\n")
	fmt.Print(m.String())
	fmt.Fprintf(os.Stderr, "---\n")
	return nil
}

func CmdNewIdentity(c *cli.Context, cfg *config.Config) error {
	return NewIssuer(cfg.Storage.Path, cfg.KeyStoreBaby.Path, cfg.KeyStoreBaby.Password)
}

func CmdStop(c *cli.Context, cfg *config.Config) error {
	output, err := PostAdminApi(cfg.Server.AdminApi, "stop")
	if err == nil {
		log.Info("Server response: ", output)
	}
	return err
}

func CmdInfo(c *cli.Context, cfg *config.Config) error {
	output, err := PostAdminApi(cfg.Server.AdminApi, "info")
	if err == nil {
		log.Info("Server response: ", output)
	}
	return err
}

// func CmdStart(c *cli.Context, cfg *config.Config, endpointServe func(cfg *config.Config, iden *loaders.Identity)) error {
// 	iden, err := loaders.LoadIdentity(cfg)
// 	if err != nil {
// 		return err
// 	}
//
// 	// Check for funds
// 	balance, err := iden.Web3.BalanceAt(iden.Web3.Account().Address)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	log.WithFields(log.Fields{
// 		"balance": balance.String(),
// 		"address": iden.Web3.Account().Address.Hex(),
// 	}).Info("Account balance retrieved")
// 	if balance.Int64() < 3000000 {
// 		log.Panic("Not enough funds in the relay address")
// 	}
//
// 	endpointServe(cfg, iden)
//
// 	iden.StateWriter.StopAndJoin()
//
// 	return nil
// }
