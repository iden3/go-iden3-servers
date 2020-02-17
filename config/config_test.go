package config

import (
	"testing"

	"github.com/iden3/go-iden3-core/core"
	"github.com/stretchr/testify/require"
)

var cfgTomlGood = `
[Identity]
Id = "113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeUf"

[Server]
ServiceApi = "0.0.0.0:6000"
AdminApi = "0.0.0.0:6001"

[Web3]
Url = "http://127.0.0.1:8545"

[KeyStore]
Path = "/var/config/keystore"
Password = "/var/config/keystore.password"

[Contracts]

[Contracts.RootCommits]
JsonABI = "/compiled_contracts/rootcmt.json"
Address = "0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe"

[Storage]
Path = "/var/data/merkletree.db"
`

var cfgTomlBad1 = `
Id = "113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeef"

[Server]
ServiceApi = "0.0.0.0:6000"
AdminApi = "0.0.0.0:6001"
`

var cfgTomlBad2 = `
Id = "113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeUf"

[Server]
ServiceApi = "0.0.0.0:6000"
`

func TestLoad(t *testing.T) {
	var cfg0 struct {
		Identity ConfigIdentity
		Server   ConfigServer
		Web3     struct {
			Url string
		}
		Contracts struct {
			RootCommits ConfigContract
		}
		Storage struct {
			Path string
		}
		KeyStore ConfigKeyStore
	}
	err := Load(cfgTomlGood, &cfg0)
	require.Nil(t, err)

	var cfg1 struct {
		Id     core.ID
		Server ConfigServer
	}
	err = Load(cfgTomlBad1, &cfg1)
	require.NotNil(t, err)

	var cfg2 struct {
		Id     core.ID
		Server ConfigServer
	}
	err = Load(cfgTomlBad2, &cfg2)
	require.NotNil(t, err)
}
