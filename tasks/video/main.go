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
		SignProposeAndPayAs("user1").
		StringArgument("bjartek").
		Run()

	g.TransactionFromFile("mintFusd").
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

	g.TransactionFromFile("createProfile").
		SignProposeAndPayAs("find").
		StringArgument("Find").
		Run()

	fmt.Scanln()
	clear()
	fmt.Println("We mint a piece of versus art. and put into 'bjartek' Art collection")
	fmt.Println("-----------------------------------------------------------------------------------")
	g.TransactionFromFile("mintArt").
		SignProposeAndPayAs("find").
		AccountArgument("user1").
		StringArgument("Versus").
		StringArgument("Versus").
		StringArgument("Versus promo nft to early testers").
		AccountArgument("user1").
		StringArgument("image/jpeg").
		UFix64Argument("0.05").
		UFix64Argument("0.025").
		StringArgument("https://res.cloudinary.com/dxra4agvf/image/upload/w_600/v1629285775/maincache5.jpg").
		RunPrintEventsFull()

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

	fmt.Println("bjartek has a collection of versus Art nft, so lets look at what is in there")
	fmt.Println("-----------------------------------------------------------------------------------")
	fmt.Scanln()
	fmt.Println("find.xyz/bjartek/A.f8d6e0586b0a20c7.Art.Collection")
	result = g.ScriptFromFile("find-ids-profile").AccountArgument("user1").StringArgument("A.f8d6e0586b0a20c7.Art.Collection").RunFailOnError()
	fmt.Println(gwtf.CadenceValueToJsonString(result))
	fmt.Println("find.xyz/bjartek/A.f8d6e0586b0a20c7.Art.Collection/0")
	result = g.ScriptFromFile("find-schemes").AccountArgument("user1").StringArgument("A.f8d6e0586b0a20c7.Art.Collection").UInt64Argument(0).RunFailOnError()
	fmt.Println(gwtf.CadenceValueToJsonString(result))

	fmt.Println("There are a single NFT that has some basic types")
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

	fmt.Println("Another example using finds build in Dandy NFT with first class support for the metadata proposal")
}

func resolveView(g *gwtf.GoWithTheFlow, view string) {
	fmt.Println()
	fmt.Println("-------------------------------------------------------")
	fmt.Printf("find.xyz/bjartek/A.f8d6e0586b0a20c7.Art.Collection/0/%s\n", view)
	result := g.ScriptFromFile("find").AccountArgument("user1").StringArgument("A.f8d6e0586b0a20c7.Art.Collection").UInt64Argument(0).StringArgument(view).RunFailOnError()
	fmt.Println(gwtf.CadenceValueToJsonString(result))
	fmt.Println()

}

func clear() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}
