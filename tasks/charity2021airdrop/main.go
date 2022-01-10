package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bjartek/go-with-the-flow/v2/gwtf"
)

func main() {
	g := gwtf.NewGoWithTheFlowMainNet()
	a := readCsvFile("charity_addresses.csv")

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })

	/*
		for _, addr := range addresses {
			value, err := g.ScriptFromFile("hasCharity").RawAccountArgument(addr).RunReturns()

			if err != nil {
				log.Fatal(err)
			}

			if value.String() == "true" {
				fmt.Printf("%s=%v", addr, value)
			}
		}

		Bronze Tier Cheque Neo Christmas Community Charity
		Silver Tier Cheque Neo Christmas Community Charity
		Gold Tier Cheque Neo Christmas Community Charity

		Neo Charity Airdrop 2021 Bronze
		Neo Charity Airdrop 2021 Silver
		Neo Charity Airdrop 2021 Gold

		A Bronze tier 3D Cheque to show participation in the Neo Collectibles x Flowverse Charity Auction in 2021.
		A Silver tier 3D Cheque to show participation in the Neo Collectibles x Flowverse Charity Auction in 2021.
		A Gold tier 3D Cheque to show participation in the Neo Collectibles x Flowverse Charity Auction in 2021.

	*/

	g.TransactionFromFile("mintCharity").
		SignProposeAndPayAs("find-admin").
		StringArgument("Christmas Tree 2021").
		StringArgument("ipfs://QmYGZXq39Ugazm9dwHz71fWCgxCf1Yub82y1kz3zkzQMyE").
		StringArgument("ipfs://QmXs9pejWe1opmDpRdS5cY6Uh7XTb1XApQQ3Dmt61ZwpKx").
		StringArgument("https://find.xyz/neo-x-flowverse-community-charity-tree").
		StringArgument(`This NFT is from the Neo x FlowVerse Charity Fundraiser 2021.
		It is a 1/1 NFT that was auctioned off with all of the proceeds going to “Women for Afghan Women”.
		The owner of this NFT is a legend for helping to make the world a better place!`).
		RawAccountArgument("0x7a1d854cbd4f84b9").
		Run()

}

func readCsvFile(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()

	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	results := []string{}
	for _, row := range records {
		results = append(results, row[0])
	}

	return results
}
