package commands

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/urfave/cli"
)

// func loadIdService() eth.Client {
// 	ks, acc := genericserver.LoadKeyStore()
// 	ksBaby, pkc := genericserver.LoadKeyStoreBabyJub()
// 	pk, err := pkc.Decompress()
// 	if err != nil {
// 		panic(err)
// 	}
// 	client := genericserver.LoadWeb3(ks, &acc)
// 	client2 := genericserver.LoadEthClient2(ks, &acc)
// 	ethsrv := genericserver.LoadEthService(client2)
// 	storage := genericserver.LoadStorage()
// 	mt := genericserver.LoadMerkele(storage)
// 	proofClaims := genericserver.LoadGenesis(mt)
// 	kUpdateMtp := proofClaims.KUpdateRoot.Proof.Mtp0.Bytes()
//
// 	rootService := genericserver.LoadRootsService(ethsrv, kUpdateMtp)
// 	claimService := genericserver.LoadClaimService(mt, rootService, ksBaby, pk)
// 	return client, genericserver.LoadIdentityService(claimService)
// }

var IdCommands = []cli.Command{{
	Name:        "id",
	Aliases:     []string{},
	Usage:       "operate with identities",
	Subcommands: []cli.Command{
		// {
		// 	Name:   "info",
		// 	Usage:  "show information about identity",
		// 	Action: cmdIdInfo,
		// },
		// {
		// 	Name:   "list",
		// 	Usage:  "list identities",
		// 	Action: cmdIdList,
		// },
		// {
		// 	Name:   "add",
		// 	Usage:  "add new identity to db",
		// 	Action: cmdIdAdd,
		// },
		// {
		// 	Name:   "deploy",
		// 	Usage:  "deploy new identity",
		// 	Action: cmdIdDeploy,
		// },
	},
}}

type idInfo struct {
	IdAddr common.Address
}
