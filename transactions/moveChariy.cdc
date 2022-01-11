import NonFungibleToken from "../contracts/standard/NonFungibleToken.cdc"
import CharityNFT from "../contracts/CharityNFT.cdc"

//This transaction will prepare the art collection
transaction(to: Address, id: UInt64) {
	prepare(account: AuthAccount) {

		let collection=account.borrow<&NonFungibleToken.Collection>(from: CharityNFT.CollectionStoragePath)!
		let charityCap=getAccount(to).getCapability<&{NonFungibleToken.CollectionPublic}>(CharityNFT.CollectionPublicPath)

		charityCap.borrow()!.deposit(token: <- collection.withdraw(withdrawID: id))
	}
}
