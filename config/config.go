package config

import (

	// common3 "github.com/iden3/go-iden3-core/common"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/urfave/cli"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(data []byte) error {
	duration, err := time.ParseDuration(string(data))
	if err != nil {
		return err
	}
	d.Duration = duration
	return nil
}

type Contract struct {
	JsonABI string         `validate:"required"`
	Address common.Address `validate:"required"`
}

type Server struct {
	ServiceApi string `validate:"required"`
	AdminApi   string `validate:"required"`
}

type Password struct {
	Value string  // private content
	Path  *string // path of the file with the password
}

func (p *Password) String() string {
	if p.Path == nil {
		return fmt.Sprintf("%v***", prefixPassword)
	}
	return fmt.Sprintf("%v%v", prefixFile, *p.Path)
}

const (
	prefixPassword = "password://"
	prefixFile     = "file://"
)

// UnmarshalText unmarshals the Password using the following rules
// Password can be prefixed by two options
//   'file://': <path to file containing the password>
//   'password//': raw password
func (p *Password) UnmarshalText(data []byte) error {
	var passwd string
	input := string(data)
	if strings.HasPrefix(input, prefixPassword) {
		passwd = input[len(prefixPassword):]
	} else if strings.HasPrefix(input, prefixFile) {
		filename := input[len(prefixFile):]
		p.Path = &filename
		passwdbytes, err := ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("Cannot read password: %w", err)
		}
		passwd = string(passwdbytes)
	} else {
		return fmt.Errorf("Prefix is missing. Use 'password://' or 'file://'")
	}
	p.Value = passwd
	return nil
}

type KeyStore struct {
	Path     string   `validate:"required"`
	Password Password `validate:"required"`
}

type EthKeys struct {
	KDis        common.Address `validate:"required"`
	KReen       common.Address `validate:"required"`
	KUpdateRoot common.Address `validate:"required"`
}

type Contracts struct {
	IdenStates Contract `validate:"required"`
	// Iden3Impl     Contract `validate:"required"`
	// Iden3Deployer Contract `validate:"required"`
	// Iden3Proxy    Contract `validate:"required"`
}

type KeysBabyJub struct {
	KOp babyjub.PublicKey `validate:"required"`
}

type Identity struct {
	Id   core.ID `validate:"required"`
	Keys struct {
		// Ethereum EthKeys `validate:"required"`
		BabyJub KeysBabyJub `validate:"required"`
	} `validate:"required"`
}

type Web3 struct {
	Url string `validate:"required"`
}

type IdenPubOffChain struct {
	Http struct {
		Url string `validate:"required"`
	} `validate:"required"`
}

type Config struct {
	Identity Identity `validate:"required"`
	// Domain    string       `validate:"required"`
	// Namespace string       `validate:"required"`
	Server       Server    `validate:"required"`
	Web3         Web3      `validate:"required"`
	KeyStore     KeyStore  `validate:"required"`
	KeyStoreBaby KeyStore  `validate:"required"`
	Contracts    Contracts `validate:"required"`
	Account      struct {
		Address common.Address `validate:"required"`
	} `validate:"required"`
	Storage struct {
		Path string
	} `validate:"required"`
	Issuer struct {
		PublishStatePeriod        Duration `validate:"required"`
		SyncIdenStatePublicPeriod Duration `validate:"required"`
	}
	IdenPubOffChain IdenPubOffChain `validate:"required"`
	// Names struct {
	// 	Path string `validate:"required"`
	// } `validate:"required"`
	// Entitites struct {
	// 	Path string `validate:"required"`
	// } `validate:"required"`
}

func LoadFromCliFlag(c *cli.Context, cfg interface{}) error {
	cfgFilePath := c.GlobalString("config")
	if cfgFilePath == "" {
		return fmt.Errorf("No config file path specified")
	}
	bs, err := ioutil.ReadFile(cfgFilePath)
	if err != nil {
		return err
	}
	if err := Load(string(bs), cfg); err != nil {
		return fmt.Errorf("Error loading configuration from cli flag: %w", err)
	}
	return nil
}

func Load(cfgToml string, cfg interface{}) error {
	if _, err := toml.Decode(cfgToml, cfg); err != nil {
		return err
	}
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("Error validating configuration file: %w", err)
	}
	return nil
}
