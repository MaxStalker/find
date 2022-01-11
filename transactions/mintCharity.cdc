import NonFungibleToken from "../contracts/standard/NonFungibleToken.cdc"
import CharityNFT from "../contracts/CharityNFT.cdc"
import Admin from "../contracts/Admin.cdc"

//mint an art and add it to a users collection
transaction(
	name: String,
	image: String,
	thumbnail: String,
	originUrl: String,
	description: String,
	tier: String, 
	recipients: [Address],
	fallback: Address
) {

	prepare(account: AuthAccount) {
		let  client= account.borrow<&Admin.AdminProxy>(from: Admin.AdminProxyStoragePath)!

		let maxEdition=recipients.length

		var i=1
		for recipient in recipients {
			let metadata = {"name" : name.concat(i.toString()).concat("/").concat(maxEdition.toString()), "image" : image, "thumbnail": thumbnail, "originUrl": originUrl, "description":description, "edition": i.toString(), "maxEdition" :  maxEdition.toString(), "rarity" : tier }

			var receiverCap= getAccount(recipient).getCapability<&{NonFungibleToken.CollectionPublic}>(CharityNFT.CollectionPublicPath)
			if !receiverCap.check() {
				receiverCap= getAccount(fallback).getCapability<&{NonFungibleToken.CollectionPublic}>(CharityNFT.CollectionPublicPath)
			}
			client.mintCharity(metadata: metadata, recipient: receiverCap)

			i=i+1
		}
	}
}

