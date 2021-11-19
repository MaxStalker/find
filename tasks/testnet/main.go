package main

import (
	"github.com/bjartek/go-with-the-flow/v2/gwtf"
)

func main() {

	//g := gwtf.NewGoWithTheFlowEmulator().InitializeContracts().CreateAccounts("emulator-account")
	g := gwtf.NewGoWithTheFlowDevNet()

	//first step create the adminClient as the fin user
	g.TransactionFromFile("setup_fin_1_create_client").
		SignProposeAndPayAs("find-admin").
		RunPrintEventsFull()

	//link in the server in the versus client
	g.TransactionFromFile("setup_fin_2_register_client").
		SignProposeAndPayAs("find").
		AccountArgument("find-admin").
		RunPrintEventsFull()

	//set up fin network as the fin user
	g.TransactionFromFile("setup_fin_3_create_network").
		SignProposeAndPayAs("find-admin").
		RunPrintEventsFull()

	g.TransactionFromFile("createProfile").
		SignProposeAndPayAs("find").
		StringArgument("find").
		RunPrintEventsFull()

	g.TransactionFromFile("createProfile").
		SignProposeAndPayAs("find-admin").
		StringArgument("find-admin").
		RunPrintEventsFull()

}
