package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bjartek/go-with-the-flow/v2/gwtf"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
)

func main() {
	tier := "Gold"
	g := gwtf.NewGoWithTheFlowMainNet()
	//	a := readCsvFile("charity_addresses.csv")

	a := []string{"0x886f3aeaf848c535"}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	var addresses []cadence.Value
	for _, key := range a {
		address := cadence.BytesToAddress(flow.HexToAddress(key).Bytes())
		addresses = append(addresses, address)
	}
	cadenceArray := cadence.NewArray(addresses)

	thumbnails := map[string]string{
		"Bronze": "QmcxXHLADpcw5R7xi6WmPjnKAEayK3eiEh85gzjgdzfwN6",
		"Silver": "QmeNsnmaPJsquCQZwGMsRvHyn5mKDXach5TZsdenEQ6Tsg",
		"Gold":   "QmbP9KKXEBVLWN66CcN12hNmrqhoE1Nd11MW2yvXr14PbZ",
	}

	images := map[string]string{
		"Bronze": "QmZNCA3hpierS95qiR7p5hS4XRY4KA8zhLTtW7BHbR6Sbp",
		"Silver": "QmVPrb43RuYopTdB4PAPamscg5rDRrqFt7uAU6vb5zp8FT",
		"Gold":   "QmT1ZrLgFFiC4Fjg5Sda195B8EkYMbimuacpXn6mTPxTC1",
	}

	title := "Neo Charity Airdrop 2021 %s "
	description := "A %s tier 3D Cheque to show participation in the Neo Collectibles x Flowverse Charity Auction in 2021.This NFT is from the Neo x FlowVerse Charity Fundraiser 2021."

	g.TransactionFromFile("mintCharity").
		SignProposeAndPayAs("find-admin").
		StringArgument(fmt.Sprintf(title, tier)).
		StringArgument("ipfs://" + images[tier]).     //image
		StringArgument("ipfs://" + thumbnails[tier]). //thumbnail
		StringArgument("https://find.xyz/neo-x-flowverse-community-charity-tree").
		StringArgument(tier).
		StringArgument(fmt.Sprintf(description, tier)).
		Argument(cadenceArray).
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
