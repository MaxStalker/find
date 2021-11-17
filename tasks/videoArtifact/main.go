package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bjartek/go-with-the-flow/v2/gwtf"
)

func main() {

	//	g := gwtf.NewGoWithTheFlowInMemoryEmulator()

	clear()
	g := gwtf.NewGoWithTheFlow([]string{"flow.json"}, "emulator", true, 0).InitializeContracts().CreateAccounts("emulator-account")

	//first step create the adminClient as the find user
	g.TransactionFromFile("setup_fin_1_create_client").
		SignProposeAndPayAs("find").
		Run()

	//link in the server in the versus client
	g.TransactionFromFile("setup_fin_2_register_client").
		SignProposeAndPayAsService().
		AccountArgument("find").
		Run()

	//set up fin network as the fin user
	g.TransactionFromFile("setup_fin_3_create_network").
		SignProposeAndPayAs("find").
		Run()

	g.TransactionFromFile("createProfile").
		SignProposeAndPayAsService().
		StringArgument("Find").
		RunPrintEventsFull()

	g.TransactionFromFile("createProfile").
		SignProposeAndPayAs("find").
		StringArgument("Find").
		RunPrintEventsFull()

	g.TransactionFromFile("createProfile").
		SignProposeAndPayAs("user1").
		StringArgument("bjartek").
		Run()

	g.TransactionFromFile("mintFusd").
		SignProposeAndPayAsService().
		AccountArgument("user1").
		UFix64Argument("100.0").
		Run()

	g.TransactionFromFile("mintFlow").
		SignProposeAndPayAsService().
		AccountArgument("user1").
		UFix64Argument("100.0").
		Run()

	fmt.Println("We register a name in find to have human readable anchor for  content adressability")
	fmt.Println("-----------------------------------------------------------------------------------")
	g.TransactionFromFile("register").
		SignProposeAndPayAs("user1").
		StringArgument("bjartek").
		UFix64Argument("5.0").
		RunPrintEventsFull()

	fmt.Scanln()
	clear()
	fmt.Println("We buy a 'artifact' addon to be able to mint things and then we mint some artifacts, in this case some example Neo Motorcycles")
	fmt.Println("-----------------------------------------------------------------------------------")
	g.TransactionFromFile("buyAddon").SignProposeAndPayAs("user1").StringArgument("bjartek").StringArgument("artifact").UFix64Argument("50.0").RunPrintEventsFull()

	g.TransactionFromFile("mintArtifact").SignProposeAndPayAs("user1").StringArgument("bjartek").RunPrintEventsFull()

	fmt.Scanln()
	clear()

	fmt.Println("Time to find out what NFTs bjartek has, so we lookup his profile")
	fmt.Println("-----------------------------------------------------------------------------------")
	fmt.Println("find.xyz/bjartek")
	fmt.Scanln()
	result := g.ScriptFromFile("find-collection").AccountArgument("user1").RunFailOnError()
	fmt.Println(gwtf.CadenceValueToJsonString(result))
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()

	fmt.Println("bjartek has a collection of Artifact NFT, so lets look at what is in there")
	fmt.Println("-----------------------------------------------------------------------------------")
	fmt.Scanln()
	fmt.Println("find.xyz/bjartek/A.f8d6e0586b0a20c7.Artifact.Collection")
	result = g.ScriptFromFile("find-ids-profile").AccountArgument("user1").StringArgument("A.f8d6e0586b0a20c7.Artifact.Collection").RunFailOnError()
	fmt.Println(gwtf.CadenceValueToJsonString(result))
	fmt.Println("find.xyz/bjartek/A.f8d6e0586b0a20c7.Artifact.Collection/1")
	result = g.ScriptFromFile("find-schemes").AccountArgument("user1").StringArgument("A.f8d6e0586b0a20c7.Artifact.Collection").UInt64Argument(1).RunFailOnError()
	fmt.Println(gwtf.CadenceValueToJsonString(result))

	fmt.Println("There are a three NFT that has some basic types")
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("-----------------------------------------------------------------------------------")
	fmt.Scanln()
	resolveView(g, "String")
	fmt.Scanln()
	resolveView(g, "A.f8d6e0586b0a20c7.TypedMetadata.Display")

	fmt.Scanln()
	resolveView(g, "A.f8d6e0586b0a20c7.TypedMetadata.Royalties")
	resolveView(g, "A.f8d6e0586b0a20c7.Artifact.Minter")

}

func resolveView(g *gwtf.GoWithTheFlow, view string) {
	fmt.Println()
	fmt.Println("-------------------------------------------------------")
	fmt.Printf("find.xyz/bjartek/A.f8d6e0586b0a20c7.Artifact.Collection/0/%s\n", view)
	result := g.ScriptFromFile("find").AccountArgument("user1").StringArgument("A.f8d6e0586b0a20c7.Artifact.Collection").UInt64Argument(1).StringArgument(view).RunFailOnError()
	fmt.Println(gwtf.CadenceValueToJsonString(result))
	fmt.Println()

}

func clear() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}
