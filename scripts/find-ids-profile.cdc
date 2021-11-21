import Profile from "../contracts/Profile.cdc"
import NonFungibleToken from "../contracts/standard/NonFungibleToken.cdc"

pub fun main(address: Address, path: String) : [UInt64] {

	let account=getAccount(address)

	let collections= getAccount(address)
	.getCapability(Profile.publicPath)
	.borrow<&{Profile.Public}>()!
	.getCollections()

	for col in collections {
		if col.name == path && col.type == Type<&{NonFungibleToken.CollectionPublic}>() {
			let cap = col.collection.borrow<&{NonFungibleToken.CollectionPublic}>()!
			return cap.getIDs()
		}
	}
	return []
}

