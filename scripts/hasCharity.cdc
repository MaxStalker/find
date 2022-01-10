import CharityNFT from "../contracts/CharityNFT.cdc"
import NonFungibleToken from "../contracts/standard/NonFungibleToken.cdc"


//Check the status of a fin user
pub fun main(user: Address) : Bool {
	let account=getAccount(user)
	let charityCap=account.getCapability<&{NonFungibleToken.CollectionPublic}>(CharityNFT.CollectionPublicPath)
	return charityCap.check()
}
