

import Profile from "../contracts/Profile.cdc"
import NonFungibleToken from "../contracts/standard/NonFungibleToken.cdc"

pub fun main(address: Address, path: String, id:UInt64) : [String] {

	let account=getAccount(address)

	let collections= getAccount(address)
	.getCapability(Profile.publicPath)
	.borrow<&{Profile.Public}>()!
	.getCollections()

	for col in collections {
		if col.name == path && col.type == Type<&{NonFungibleToken.CollectionPublic}>() {
			let cap = col.collection.borrow<&{NonFungibleToken.CollectionPublic}>()! as &{NonFungibleToken.CollectionPublic}
			let nft=cap.borrowNFT(id: id)
			let views=nft.getViews()
			var viewIdentifiers : [String] = []
			for v in views {
				viewIdentifiers.append(v.identifier)
			}
			return viewIdentifiers
		}
	}
	return []
}

