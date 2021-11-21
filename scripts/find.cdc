import NonFungibleToken from "../contracts/standard/NonFungibleToken.cdc"
import Profile from "../contracts/Profile.cdc"

pub fun main(address: Address, path: String, id: UInt64, identifier: String) : AnyStruct? {


	/*
	let paths =getAccount(address).getPublicPaths<&{NonFungibleToken.Collection}>()
	for path in paths {
		let col=path.getCapability().borrow()
		return col.borrowNFT(id).resolveView(identifier)
	}
	*/

	let collections= getAccount(address).getCapability(Profile.publicPath)
	.borrow<&{Profile.Public}>()!
	.getCollections()

	for col in collections {
		if col.name == path && col.type == Type<&{NonFungibleToken.CollectionPublic}>() {
			let cap = col.collection.borrow<&{NonFungibleToken.CollectionPublic}>()! as &{NonFungibleToken.CollectionPublic}

			let nft=cap.borrowNFT(id: id)
			for v in nft.getViews() {
				if v.identifier== identifier {
					return nft.resolveView(v)
				}
			}
		}
	}
	return nil
}

