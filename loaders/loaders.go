package loaders

import (
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	ethkeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/eth"
	"github.com/iden3/go-iden3-core/identity/issuer"
	babykeystore "github.com/iden3/go-iden3-core/keystore"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3-servers/config"
	log "github.com/sirupsen/logrus"
)

var (
	dbMerkletreePrefix     = []byte{0}
	dbCounterfactualPrefix = []byte{1}
)

const (
	passwdPrefix = "passwd:"
	filePrefix   = "file:"
)

// func LoadKeyStore(cfgKeyStore *config.ConfigKeyStore, cfgEthKeys *config.ConfigEthKeys) (*ethkeystore.KeyStore, accounts.Account) {
// 	var err error
// 	var passwd string
//
// 	// Load keystore
// 	ks := ethkeystore.NewKeyStore(cfgKeyStore.Path, ethkeystore.StandardScryptN, ethkeystore.StandardScryptP)
//
// 	// Password can be prefixed by two options
// 	//   file: <path to file containing the password>
// 	//   passwd: raw password
// 	// if is not prefixed by any of those, file: is used
// 	// TODO: Handle the literal password / file password with a configuration type
// 	if strings.HasPrefix(cfgKeyStore.Password, passwdPrefix) {
// 		passwd = cfgKeyStore.Password[len(passwdPrefix):]
// 	} else {
// 		filename := cfgKeyStore.Password
// 		if strings.HasPrefix(filename, filePrefix) {
// 			filename = cfgKeyStore.Password[len(filePrefix):]
// 		}
// 		passwdbytes, err := ioutil.ReadFile(filename)
// 		Assert("Cannot read password", err)
// 		passwd = string(passwdbytes)
// 	}
//
// 	acc, err := ks.Find(accounts.Account{
// 		Address: cfgEthKeys.KUpdateRoot,
// 	})
// 	Assert("Cannot find keystore account", err)
// 	// KDis and KReen not used yet, but need to check if they exist
// 	_, err = ks.Find(accounts.Account{
// 		Address: cfgEthKeys.KDis,
// 	})
// 	Assert("Cannot find keystore account", err)
// 	_, err = ks.Find(accounts.Account{
// 		Address: cfgEthKeys.KReen,
// 	})
// 	Assert("Cannot find keystore account", err)
//
// 	Assert("Cannot unlock account", ks.Unlock(acc, string(passwd)))
// 	log.WithField("acc", acc.Address.Hex()).Info("Keystore and account unlocked successfully")
//
// 	return ks, acc
// }

func LoadKeyStoreBabyJub(cfgKeyStore *config.ConfigKeyStore, kOp *babyjub.PublicKey) (*babykeystore.KeyStore, *babyjub.PublicKeyComp) {
	storage := babykeystore.NewFileStorage(cfgKeyStore.Path)
	ks, err := babykeystore.NewKeyStore(storage, babykeystore.StandardKeyStoreParams)
	if err != nil {
		panic(err)
	}
	if kOp != nil {
		kOpComp := kOp.Compress()
		if err := ks.UnlockKey(&kOpComp, []byte(cfgKeyStore.Password)); err != nil {
			panic(err)
		}
		return ks, &kOpComp
	} else {
		return ks, nil
	}
}

func LoadEthClient2(ks *ethkeystore.KeyStore, acc *accounts.Account, web3Url string) (*eth.Client2, error) {
	// TODO: Handle the hidden: thing with a custon configuration type
	hidden := strings.HasPrefix(web3Url, "hidden:")
	if hidden {
		web3Url = web3Url[len("hidden:"):]
	}
	client, err := ethclient.Dial(web3Url)
	if err != nil {
		return nil, fmt.Errorf("Error dialing with ethclient: %w", err)
	}
	if hidden {
		log.WithField("url", "(hidden)").Info("Connection to web3 server opened")
	} else {
		log.WithField("url", web3Url).Info("Connection to web3 server opened")
	}
	return eth.NewClient2(client, acc, ks), nil
}

func LoadIdenPubOnChain(client *eth.Client2, idenStatesAddress common.Address) idenpubonchain.IdenPubOnChainer {

	addresses := idenpubonchain.ContractAddresses{
		IdenStates: idenStatesAddress,
	}
	return idenpubonchain.New(client, addresses)
}

func LoadWeb3(ks *ethkeystore.KeyStore, acc *accounts.Account, web3Url string) (*eth.Web3Client, error) {
	// Create geth client
	hidden := strings.HasPrefix(web3Url, "hidden:")
	if hidden {
		web3Url = web3Url[len("hidden:"):]
	}
	web3cli, err := eth.NewWeb3Client(web3Url, ks, acc)
	if err != nil {
		return nil, fmt.Errorf("Error creating web3 client: %w", err)
	}
	if hidden {
		log.WithField("url", "(hidden)").Info("Connection to web3 server opened")
	} else {
		log.WithField("url", web3Url).Info("Connection to web3 server opened")
	}
	return web3cli, nil
}

func LoadStorage(storagePath string) (db.Storage, error) {
	// Open database
	storage, err := db.NewLevelDbStorage(storagePath, false)
	if err != nil {
		return nil, fmt.Errorf("Error opening leveldb storage: %w", err)
	}
	log.WithField("path", storagePath).Info("Storage opened")
	return storage, nil
}

// func LoadMerkele(storage db.Storage) (*merkletree.MerkleTree, error) {
// 	mtstorage := storage.WithPrefix(dbMerkletreePrefix)
// 	mt, err := merkletree.NewMerkleTree(mtstorage, 140)
// 	if err != nil {
// 		return nil, fmt.Errorf("Error opening merkle tree: %w", err)
// 	}
// 	log.WithField("hash", mt.RootKey().Hex()).Info("Current root")
//
// 	return mt, nil
// }

func LoadContract(client eth.Client, jsonabifile string, address *common.Address) (*eth.Contract, error) {
	abiFile, err := os.Open(jsonabifile)
	if err != nil {
		return nil, fmt.Errorf("Error reading contract %s: %w", jsonabifile, err)
	}

	abi, code, err := eth.UnmarshallSolcAbiJson(abiFile)
	if err != nil {
		return nil, fmt.Errorf("Error parsing contract %s: %w", jsonabifile, err)
	}

	return eth.NewContract(client, abi, code, address), nil
}

// func LoadIdenManager(mt *merkletree.MerkleTree, rootservice idenstatewriter.IdenStateWriter, ks *babykeystore.KeyStore, pk *babyjub.PublicKey, id *core.ID) *idenmanager.IdenManager {
// 	log.WithField("id", id.String()).Info("Running claim service")
// 	signer := idensigner.New(ks, *pk)
// 	return idenmanager.New(id, mt, rootservice, *signer)
// }
//
// func LoadIdenAdminUtils(mt *merkletree.MerkleTree, rootservice idenstatewriter.IdenStateWriter, claimservice *idenmanager.IdenManager) *idenadminutils.IdenAdminUtils {
// 	return idenadminutils.New(mt, rootservice, claimservice)
// }

func LoadIssuer(id *core.ID, storage db.Storage, keyStore *babykeystore.KeyStore,
	idenPubOnChain idenpubonchain.IdenPubOnChainer,
	idenPubOffChainWrite idenpuboffchain.IdenPubOffChainWriter) (*issuer.Issuer, error) {

	idenStorage := storage.WithPrefix([]byte(fmt.Sprintf("%v:", id)))
	is, err := issuer.Load(idenStorage, keyStore, idenPubOnChain, idenPubOffChainWrite)
	if err != nil {
		return nil, err
	}
	return is, nil
}

// LoadGenesis will calculate the genesis id from the keys in the configuration
// file and check it against the id in the configuration.  It will populate the
// merkle tree with the genesis claims if it's empty or check that the claims
// exist in the merkle tree otherwise.  It returns the ProofClaims of the
// genesis claims.
// func LoadGenesis(mt *merkletree.MerkleTree, id *core.ID, kOp *babyjub.PublicKey, cfgEthKeys *config.ConfigEthKeys) *genesis.GenesisProofClaims {
// 	kDis := cfgEthKeys.KDis
// 	kReen := cfgEthKeys.KReen
// 	kUpdateRoot := cfgEthKeys.KUpdateRoot
// 	id0, proofClaims, err := genesis.CalculateIdGenesisFrom4Keys(kOp, kDis, kReen, kUpdateRoot)
// 	Assert("CalculateIdGenesis failed", err)
//
// 	if *id0 != *id {
// 		Assert("Error", fmt.Errorf("Calculated genesis id (%v) "+
// 			"doesn't match configuration id (%v)", id0.String(), id0.String()))
// 	}
//
// 	proofClaimsList := []proof.ProofClaim{proofClaims.KOp, proofClaims.KDis,
// 		proofClaims.KReen, proofClaims.KUpdateRoot}
// 	root := mt.RootKey()
// 	if bytes.Equal(root[:], merkletree.HashZero[:]) {
// 		// Merklee tree DB is empty
// 		// Add genesis claims to merkle tree
// 		log.WithField("root", root.Hex()).Info("Merkle tree is empty")
// 		for _, proofClaim := range proofClaimsList {
// 			if err := mt.AddEntry(proofClaim.Claim); err != nil {
// 				Assert("Error adding claim to merkle tree", err)
// 			}
// 		}
// 	} else {
// 		// MerkleTree DB has already been initialized
// 		// Check that the geneiss claims are in the merkle tree
// 		log.WithField("root", root.Hex()).Info("Merkle tree already initialized")
// 		for _, proofClaim := range proofClaimsList {
// 			entry := proofClaim.Claim
// 			data, err := mt.GetDataByIndex(entry.HIndex())
// 			if err != nil {
// 				Assert("Error getting claim from the merkle tree", err)
// 			}
// 			if !entry.Data.Equal(data) {
// 				Assert("Error", fmt.Errorf("Claim from the merkle tree (%v) "+
// 					"doesn't match the expected claim (%v)",
// 					data.String(), entry.Data.String()))
// 			}
// 		}
// 	}
//
// 	return proofClaims
// }

type Server struct {
	Id             core.ID
	Mt             *merkletree.MerkleTree
	Issuer         *issuer.Issuer
	IdenPubOnChain idenpubonchain.IdenPubOnChainer
	KeyStore       *ethkeystore.KeyStore
	KeyStoreBaby   *babykeystore.KeyStore
	Web3           *eth.Web3Client
	EthClient2     *eth.Client2
	KOp            *babyjub.PublicKey
}

func LoadServer(cfg *config.Config) (*Server, error) {
	// ks, acc := LoadKeyStore(&cfg.KeyStore, &cfg.Keys.Ethereum)
	ksBaby, kOp := LoadKeyStoreBabyJub(&cfg.KeyStoreBaby, &cfg.Identity.Keys.BabyJub.KOp)
	pk, err := kOp.Decompress()
	if err != nil {
		return nil, err
	}
	client, err := LoadWeb3(nil, nil, cfg.Web3.Url)
	if err != nil {
		return nil, err
	}
	// client2 := LoadEthClient2(ks, &acc, cfg.Web3.Url)
	client2, err := LoadEthClient2(nil, nil, cfg.Web3.Url)
	if err != nil {
		return nil, err
	}
	idenPubOnChain := LoadIdenPubOnChain(client2, cfg.Contracts.IdenStates.Address)
	storage, err := LoadStorage(cfg.Storage.Path)
	if err != nil {
		return nil, err
	}
	// mt, err := LoadMerkele(storage)
	// if err != nil {
	// 	return nil, err
	// }

	println("AAA")
	is, err := LoadIssuer(&cfg.Identity.Id, storage, ksBaby, idenPubOnChain, nil)
	if err != nil {
		return nil, err
	}
	println("BBB")

	// proofClaims := LoadGenesis(mt, &cfg.Id, &cfg.Keys.BabyJub.KOp, &cfg.Keys.Ethereum)
	// kUpdateMtp := proofClaims.KUpdateRoot.Proof.Mtp0.Bytes()

	return &Server{
		Issuer: is,
		// Mt:             mt,
		IdenPubOnChain: idenPubOnChain,
		// KeyStore:       ks,
		KeyStore:     nil,
		KeyStoreBaby: ksBaby,
		Web3:         client,
		EthClient2:   client2,
		KOp:          pk,
	}, nil
}
