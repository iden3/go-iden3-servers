module github.com/iden3/go-iden3-servers

go 1.13

// replace github.com/iden3/go-iden3-servers => ./

replace github.com/iden3/go-iden3-core => ../go-iden3-core

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/ethereum/go-ethereum v1.9.10
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.5.0
	github.com/go-playground/validator/v10 v10.1.0
	github.com/iden3/go-iden3-core v0.0.7-0.20200213165305-161fa5bfa30f
	github.com/iden3/go-iden3-crypto v0.0.3-0.20190831180703-c95c95b7b161
	github.com/iden3/go-public-key-encryption v0.0.0-20200129111956-c21e08c0ca6d
	github.com/ipfs/go-ipfs-api v0.0.3
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	github.com/urfave/cli v1.22.2
)
