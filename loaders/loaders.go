package loaders

import (
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	ethkeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	idenpuboffchainwriterhttp "github.com/iden3/go-iden3-core/components/idenpuboffchain/writerhttp"
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

func LoadKeyStore(cfgKeyStore *config.KeyStore, accountAddr *common.Address) (*ethkeystore.KeyStore, *accounts.Account, error) {
	// Load keystore
	ks := ethkeystore.NewKeyStore(cfgKeyStore.Path, ethkeystore.StandardScryptN, ethkeystore.StandardScryptP)

	acc, err := ks.Find(accounts.Account{
		Address: *accountAddr,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("Cannot find keystore account: %w", err)
	}
	// KDis and KReen not used yet, but need to check if they exist
	// _, err = ks.Find(accounts.Account{
	// 	Address: cfgEthKeys.KDis,
	// })
	// Assert("Cannot find keystore account", err)
	// _, err = ks.Find(accounts.Account{
	// 	Address: cfgEthKeys.KReen,
	// })
	// Assert("Cannot find keystore account", err)

	if err := ks.Unlock(acc, string(cfgKeyStore.Password.Value)); err != nil {
		return nil, nil, fmt.Errorf("Cannot unlock account: %w", err)
	}
	log.WithField("acc", acc.Address.Hex()).Info("Keystore and account unlocked successfully")

	return ks, &acc, nil
}

func LoadKeyStoreBabyJub(cfgKeyStore *config.KeyStore, kOp *babyjub.PublicKey) (*babykeystore.KeyStore, error) {
	storage := babykeystore.NewFileStorage(cfgKeyStore.Path)
	ks, err := babykeystore.NewKeyStore(storage, babykeystore.StandardKeyStoreParams)
	if err != nil {
		return nil, fmt.Errorf("Error creating/opening babyjub keystore: %w", err)
	}
	if kOp != nil {
		kOpComp := kOp.Compress()
		if err := ks.UnlockKey(&kOpComp, []byte(cfgKeyStore.Password.Value)); err != nil {
			return nil, fmt.Errorf("Error unlocking babyjub key from keystore: %w", err)
		}
		log.WithField("kOp", kOpComp.String()).Info("Babyjub Keystore and key unlocked successfully")
		return ks, nil
	} else {
		return ks, nil
	}
}

func LoadEthClient(ks *ethkeystore.KeyStore, acc *accounts.Account, web3Url string) (*eth.Client, error) {
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
	return eth.NewClient(client, acc, ks), nil
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

func LoadIdenPubOffChainWriteHttp(storage db.Storage, id *core.ID,
	url string) (*idenpuboffchainwriterhttp.IdenPubOffChainWriteHttp, error) {
	return idenpuboffchainwriterhttp.NewIdenPubOffChainWriteHttp(
		idenpuboffchainwriterhttp.NewConfigDefault(url),
		storage.WithPrefix([]byte(fmt.Sprintf("%v:writerhttp:", id))),
	)
}

func LoadIssuer(id *core.ID, storage db.Storage, keyStore *babykeystore.KeyStore,
	idenPubOnChain idenpubonchain.IdenPubOnChainer,
	idenStateZkProofConf *issuer.IdenStateZkProofConf,
	idenPubOffChainWrite idenpuboffchain.IdenPubOffChainWriter) (*issuer.Issuer, error) {

	idenStorage := storage.WithPrefix([]byte(fmt.Sprintf("%v:", id)))
	is, err := issuer.Load(idenStorage, keyStore, idenPubOnChain, idenStateZkProofConf, idenPubOffChainWrite)
	if err != nil {
		return nil, fmt.Errorf("Error loading issuer: %w", err)
	}
	log.WithField("id", id).Info("Issuer loaded successfully")
	return is, nil
}

type Server struct {
	Cfg                      *config.Config
	stopchPublish            chan (interface{})
	stoppedchPublish         chan (interface{})
	stopchSync               chan (interface{})
	stoppedchSync            chan (interface{})
	Id                       core.ID
	Mt                       *merkletree.MerkleTree
	Issuer                   *issuer.Issuer
	IdenPubOnChain           idenpubonchain.IdenPubOnChainer
	IdenPubOffChainWriteHttp *idenpuboffchainwriterhttp.IdenPubOffChainWriteHttp
	KeyStore                 *ethkeystore.KeyStore
	KeyStoreBaby             *babykeystore.KeyStore
	EthClient                *eth.Client
	KOp                      *babyjub.PublicKey
}

func (s *Server) Start() error {
	log.Info("Starting Issuer Server")
	go func() {
		log.Info("Starting periodic Issuer PublishState")
		for {
			select {
			case <-s.stopchPublish:
				log.Info("Issuer server finalized")
				s.stoppedchPublish <- nil
				return
			case <-time.After(s.Cfg.Issuer.PublishStatePeriod.Duration):
				log.Debug("Issuer.PublishState()...")
				if err := s.Issuer.PublishState(); err != nil {
					if err != issuer.ErrIdenStatePendingNotNil {
						log.WithField("err", err).Error("Issuer.PublishState")
					}
				}
				state, _ := s.Issuer.State()
				onchain := s.Issuer.IdenStateOnChain()
				pending := s.Issuer.IdenStatePending()
				log.WithField("state", state).WithField("onchain", onchain).
					WithField("pending", pending).Debug("Issuer.PublishState()")
			}
		}
	}()
	go func() {
		log.Info("Starting periodic Issuer SyncIdenStatePublic")
		for {
			select {
			case <-s.stopchSync:
				log.Info("Issuer server finalized")
				s.stoppedchSync <- nil
				return
			case <-time.After(s.Cfg.Issuer.SyncIdenStatePublicPeriod.Duration):
				log.Debug("Issuer.SyncIdenStatePublic()...")
				if err := s.Issuer.SyncIdenStatePublic(); err != nil {
					log.WithField("err", err).Error("Issuer.SyncIdenStatePublicPeriod")
				}
				state, _ := s.Issuer.State()
				pending := s.Issuer.IdenStatePending()
				onchain := s.Issuer.IdenStateOnChain()
				log.WithField("state", state).WithField("onchain", onchain).
					WithField("pending", pending).Debug("Issuer.SyncIdenStatePublic()")
			}
		}
	}()
	return nil
}

func (s *Server) StopAndJoin() {
	go func() {
		s.stopchPublish <- nil
		s.stopchSync <- nil
	}()
	<-s.stoppedchPublish
	<-s.stoppedchSync
}

func LoadServer(cfg *config.Config) (*Server, error) {
	ks, acc, err := LoadKeyStore(&cfg.KeyStore, &cfg.Account.Address)
	if err != nil {
		return nil, err
	}
	kOp := &cfg.Identity.Keys.BabyJub.KOp
	ksBaby, err := LoadKeyStoreBabyJub(&cfg.KeyStoreBaby, kOp)
	if err != nil {
		return nil, err
	}
	ethClient, err := LoadEthClient(ks, acc, cfg.Web3.Url)
	if err != nil {
		return nil, err
	}

	idenPubOnChain := idenpubonchain.New(ethClient, idenpubonchain.ContractAddresses{
		IdenStates: cfg.Contracts.IdenStates.Address,
	})
	storage, err := LoadStorage(cfg.Storage.Path)
	if err != nil {
		return nil, err
	}

	idenPubOffChainWriteHttp, err := LoadIdenPubOffChainWriteHttp(storage, &cfg.Identity.Id,
		cfg.IdenPubOffChain.Http.Url)
	if err != nil {
		return nil, err
	}

	zkFilesIdenState := cfg.IdenStateZKProof.Files.Value()
	if err := zkFilesIdenState.LoadAll(); err != nil {
		return nil, err
	}
	idenStateZkProofConf := issuer.IdenStateZkProofConf{Levels: cfg.IdenStateZKProof.Levels, Files: *zkFilesIdenState}

	is, err := LoadIssuer(&cfg.Identity.Id, storage, ksBaby, idenPubOnChain, &idenStateZkProofConf,
		idenPubOffChainWriteHttp)
	if err != nil {
		return nil, err
	}

	// proofClaims := LoadGenesis(mt, &cfg.Id, &cfg.Keys.BabyJub.KOp, &cfg.Keys.Ethereum)
	// kUpdateMtp := proofClaims.KUpdateRoot.Proof.Mtp0.Bytes()

	return &Server{
		Cfg:                      cfg,
		stopchPublish:            make(chan (interface{})),
		stoppedchPublish:         make(chan (interface{})),
		stopchSync:               make(chan (interface{})),
		stoppedchSync:            make(chan (interface{})),
		Issuer:                   is,
		IdenPubOnChain:           idenPubOnChain,
		IdenPubOffChainWriteHttp: idenPubOffChainWriteHttp,
		// KeyStore:       ks,
		KeyStore:     nil,
		KeyStoreBaby: ksBaby,
		EthClient:    ethClient,
		KOp:          kOp,
	}, nil
}
