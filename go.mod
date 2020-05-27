module github.com/iden3/go-iden3-servers

go 1.13

// replace github.com/iden3/go-iden3-core => ../go-iden3-core

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/ethereum/go-ethereum v1.9.13
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.5.0
	github.com/go-playground/validator/v10 v10.1.0
	github.com/iden3/go-iden3-core v0.0.8-0.20200527125702-3ace820b1db5
	github.com/iden3/go-iden3-crypto v0.0.5-0.20200525100545-2c471ab54594
	github.com/iden3/go-public-key-encryption v0.0.0-20200129111956-c21e08c0ca6d
	github.com/ipfs/go-ipfs-api v0.0.3
	github.com/robertkrimen/otto v0.0.0-20170205013659-6a77b7cbc37d // indirect
	github.com/sirupsen/logrus v1.5.0
	github.com/stretchr/testify v1.5.1
	github.com/urfave/cli v1.22.2
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
)
