/*
- collection
 - type
 - dictionary id ->
  - name
  - imageurl
  - hash
	*/

import TypedMetadata from "../contracts/TypedMetadata.cdc"
import Profile from "../contracts/Profile.cdc"


pub struct MetadataCollection{
	pub let type: String
	pub let items: [MetadataCollectionItem]

	init(type:String, items: [MetadataCollectionItem]) {
		self.type=type
		self.items=items
	}
}

pub struct MetadataCollectionItem {
	pub let id:UInt64
	pub let name: String
	pub let url: String


	init(id:UInt64, name:String, url:String) {
		self.id=id
		self.name=name
		self.url=url
	}

}


pub fun main(address: Address) : {String : MetadataCollection} {

	let results : {String :  MetadataCollection}={}

	let collections= getAccount(address).getCapability(Profile.publicPath).borrow<&{Profile.Public}>()!.getCollections()

	for col in collections {
		if col.type ==Type<&{TypedMetadata.ViewResolverCollection}>() {
			let name=col.name
			let collection : { UInt64 : { String : AnyStruct }}={}
			let vrc= col.collection.borrow<&{TypedMetadata.ViewResolverCollection}>()!

			let items: [MetadataCollectionItem]=[]
			for id in vrc.getIDs() {
				let nft=vrc.borrowViewResolver(id: id)
				//just assume everything supports it
				let display = nft.resolveView(Type<TypedMetadata.Display>()) as! TypedMetadata.Display
				items.append(MetadataCollectionItem(id:id, name:display.name,  url:display.thumbnail))
			}
			results[name]= MetadataCollection(type: vrc.getType().identifier, items: items)
		}
	}
	return results
}





