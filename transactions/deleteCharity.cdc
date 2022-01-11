import NonFungibleToken from "../contracts/standard/NonFungibleToken.cdc"
import CharityNFT from "../contracts/CharityNFT.cdc"

//This transaction will prepare the art collection
transaction(id: UInt64) {
	prepare(account: AuthAccount) {

		let collection=account.borrow<&NonFungibleToken.Collection>(from: CharityNFT.CollectionStoragePath)!

		destroy collection.withdraw(withdrawID: id)
	}
}
