import FIND from "../contracts/FIND.cdc"
import NonFungibleToken from "../contracts/standard/NonFungibleToken.cdc"
import FungibleToken from "../contracts/standard/FungibleToken.cdc"
import FUSD from "../contracts/standard/FUSD.cdc"
import FlowToken from "../contracts/standard/FlowToken.cdc"
import TypedMetadata from "../contracts/TypedMetadata.cdc"
import Dandy from "../contracts/Dandy.cdc"
import Profile from "../contracts/Profile.cdc"

transaction(name: String) {
	prepare(account: AuthAccount) {

		let finLeases= account.borrow<&FIND.LeaseCollection>(from:FIND.LeaseStoragePath)!

		let sharedContentCap =account.getCapability<&{NonFungibleToken.CollectionPublic}>(/private/sharedContent)
		if !sharedContentCap.check() {
			account.save<@NonFungibleToken.Collection>(<- Dandy.createEmptyCollection(), to: /storage/sharedContent)
			account.link<&{NonFungibleToken.CollectionPublic}>(/private/sharedContent, target: /storage/sharedContent)
		}

		//this will panic if you cannot borrow it
		finLeases.borrow(name) 

		let creativeWork=
		TypedMetadata.CreativeWork(artist:"Neo Motorcycles", name:"Neo Bike ", description: "Bringing the motorcycle world into the 21st century with cutting edge EV technology and advanced performance in a great classic British style, all here in the UK", type:"image")
		let media=TypedMetadata.Media(data:"https://neomotorcycles.co.uk/assets/img/neo_motorcycle_side.webp" , contentType: "image/webp", protocol: "http")


		let receiver=account.getCapability<&{FungibleToken.Receiver}>(/public/flowTokenReceiver)
		let sharedSchemas : [AnyStruct] = [
			media,
			creativeWork,
			Dandy.Royalties(royalty: {"artist" : Dandy.RoyaltyItem(receiver: receiver, cut: 0.05)})
		]

		let sharedNFT <- finLeases.mintDandy(name: name, nftName: "NeoBike", schemas:sharedSchemas)
		let sharedPointer= Dandy.Pointer(collection: sharedContentCap, id: sharedNFT.id, views: [Type<TypedMetadata.Media>(), Type<TypedMetadata.CreativeWork>(), Type<{TypedMetadata.Royalty}>()])
		sharedContentCap.borrow()!.deposit(token: <- sharedNFT)
	
		let cap = account.getCapability<&{NonFungibleToken.CollectionPublic}>(Dandy.DandyPublicPath)

		let collection=cap.borrow()!
		var i:UInt64=1
		let maxEdition:UInt64=3
		while i <= maxEdition {

			let editioned= TypedMetadata.Editioned(edition:i, maxEdition:maxEdition)
			let description=creativeWork.description.concat( " edition ").concat(i.toString()).concat( " of ").concat(maxEdition.toString())
			//TODO: do not send in Display but calculate it, send in thumbnail url if you do not have explicit media
			let display= TypedMetadata.Display(name: "Neo Motorcycle".concat(i.toString()).concat(" of ").concat(maxEdition.toString()), thumbnail: media.data, description: description, source: "dandy")
			let schemas: [AnyStruct] = [ editioned, display]
			let token <- finLeases.mintNFTWithSharedData(name: name, nftName: "Neo Motorcycle ".concat(i.toString()).concat(" of ").concat(maxEdition.toString()), schemas: schemas, sharedPointer: sharedPointer)

			collection.deposit(token: <- token)
			i=i+1
		}

	}
}
