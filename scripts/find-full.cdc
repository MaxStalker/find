import NonFungibleToken from "../contracts/standard/NonFungibleToken.cdc"
import Profile from "../contracts/Profile.cdc"


pub fun main(address: Address) : {String : { UInt64 : { String : AnyStruct }}} {

	let results : {String : { UInt64 : { String : AnyStruct }}}={}

	let collections= getAccount(address).getCapability(Profile.publicPath)
	.borrow<&{Profile.Public}>()!
	.getCollections()

	for col in collections {
		if col.type ==Type<&{NonFungibleToken.CollectionPublic}>() {
			let name=col.name
			let collection : { UInt64 : { String : AnyStruct }}={}
			let vrc= col.collection.borrow<&{NonFungibleToken.CollectionPublic}>()!
			for id in vrc.getIDs() {
				let views : { String : AnyStruct }={}
				let nft=vrc.borrowNFT(id: id)
				for view in nft.getViews() {
					let resolved=nft.resolveView(view)
					views[view.identifier] = resolved
				}
				collection[id]=views
			}
			results[name]=collection
		}
	}
	return results
}





