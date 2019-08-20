module go-iden3-servers

go 1.12

replace github.com/iden3/go-iden3-servers => ./

replace github.com/iden3/go-iden3-core => ../go-iden3-core

require (
	github.com/appleboy/gin-jwt/v2 v2.6.2
	github.com/ethereum/go-ethereum v1.9.1
	github.com/gballet/go-libpcsclite v0.0.0-20190607065134-2772fd86a8ff // indirect
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.4.0
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/huin/goupnp v1.0.0 // indirect
	github.com/iden3/go-iden3-core v0.0.0-00010101000000-000000000000
	github.com/iden3/go-iden3-crypto v0.0.1
	github.com/iden3/go-iden3-servers v0.0.0-00010101000000-000000000000
	github.com/ipfs/go-ipfs-api v0.0.2
	github.com/jackpal/go-nat-pmp v1.0.1 // indirect
	github.com/karalabe/usb v0.0.0-20190819132248-550797b1cad8 // indirect
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/olekukonko/tablewriter v0.0.1 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.4.0
	github.com/status-im/keycard-go v0.0.0-20190424133014-d95853db0f48 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/tyler-smith/go-bip39 v1.0.2 // indirect
	github.com/urfave/cli v1.20.0
	github.com/wsddn/go-ecdh v0.0.0-20161211032359-48726bab9208 // indirect
)
