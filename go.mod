module github.com/iden3/go-iden3-servers

go 1.12

// replace github.com/iden3/go-iden3-servers => ./

// replace github.com/iden3/go-iden3-core => ../go-iden3-core

require (
	github.com/appleboy/gin-jwt/v2 v2.6.2
	github.com/ethereum/go-ethereum v1.9.3
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.4.0
	github.com/iden3/go-iden3-core v0.0.7-0.20200128100725-56ab4417a3c7
	github.com/iden3/go-iden3-crypto v0.0.3-0.20190831180703-c95c95b7b161
	github.com/ipfs/go-ipfs-api v0.0.2
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.3.0
	github.com/urfave/cli v1.20.0
)
