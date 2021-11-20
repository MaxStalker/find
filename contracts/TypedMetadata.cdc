import FungibleToken from "../contracts/standard/FungibleToken.cdc"
import NonFungibleToken from "../contracts/standard/NonFungibleToken.cdc"
import FlowToken from "../contracts/standard/FlowToken.cdc"
import FUSD from "../contracts/standard/FUSD.cdc"

pub contract TypedMetadata {

	pub resource interface ViewResolver {
		pub fun getViews() : [Type]
		pub fun resolveView(_ view:Type): AnyStruct
	}

	pub struct interface Royalty {
					
		/// if nil cannot pay this type
		/// if not nill withdraw that from main vault and put it into distributeRoyalty 

		pub fun calculateRoyalty(type:Type, amount:UFix64) : UFix64?

		/// call this with a vault containing the amount given in calculate royalty and it will be distributed accordingly
		pub fun distributeRoyalty(vault: @FungibleToken.Vault) 


		/// generate a string that represents all the royalties this NFT has for display purposes
		pub fun displayRoyalty() : String?  

	}


	pub struct Display{
		pub let name: String
		pub let thumbnail: String
		pub let description: String
		pub let source: String

		init(name:String, thumbnail: String, description: String, source:String) {
			self.source=source
			self.name=name
			self.thumbnail=thumbnail
			self.description=description
		}
	}

	pub struct Royalties : Royalty {
		pub let royalty: { String : RoyaltyItem}
		init(royalty: {String : RoyaltyItem}) {
			self.royalty=royalty
		}

		pub fun calculateRoyalty(type:Type, amount:UFix64) : UFix64? {
			var sum:UFix64=0.0
			for key in self.royalty.keys {
				let item= self.royalty[key]!
				sum=sum+amount*item.cut
			}
			return sum
		}
	
		pub fun distributeRoyalty(vault: @FungibleToken.Vault) {
			let totalAmount=vault.balance
			var sumCuts:UFix64=0.0
			for key in self.royalty.keys {
				let item= self.royalty[key]!
				sumCuts=sumCuts+item.cut
			}

			let totalKeys=self.royalty.keys.length
			var currentKey=1
			var lastReceiver: Capability<&{FungibleToken.Receiver}>?=nil
			for key in self.royalty.keys {
				let item= self.royalty[key]!
				let relativeCut=item.cut / sumCuts

				if currentKey!=totalKeys {
					item.receiver.borrow()!.deposit(from: <-  vault.withdraw(amount: totalAmount*relativeCut))
				} else { 
					//we cannot calculate the last cut as it will have rounding errors
					lastReceiver=item.receiver
				}
				currentKey=currentKey+1
			}
			if let r=lastReceiver {
				r.borrow()!.deposit(from: <-  vault)
			}else {
				destroy vault
			}
		}

		pub fun displayRoyalty() : String?  {
			var text=""
			for key in self.royalty.keys {
				let item= self.royalty[key]!
				text.concat(key).concat(" ").concat((item.cut * 100.0).toString()).concat("%\n")
			}
			return text
		}
	}

	pub struct RoyaltyItem{
		pub let receiver: Capability<&{FungibleToken.Receiver}> 
		pub let cut: UFix64

		init(receiver: Capability<&{FungibleToken.Receiver}>, cut: UFix64) {
			self.cut=cut
			self.receiver=receiver
		}
	}

	pub struct Medias {
		pub let media : {String:  Media}

		init(_ items: {String: Media}) {
			self.media=items
		}
	}

	pub struct Media {
		pub let data: String
		pub let contentType: String
		pub let protocol: String

		init(data:String, contentType: String, protocol: String) {
			self.data=data
			self.protocol=protocol
			self.contentType=contentType
		}
	}

	pub struct CreativeWork {
		pub let artist: String
		pub let name: String
		pub let description: String
		pub let type: String

		init(artist: String, name: String, description: String, type: String) {
			self.artist=artist
			self.name=name
			self.description=description
			self.type=type
		}
	}

	pub struct Editioned {
		pub let edition: UInt64
		pub let maxEdition: UInt64

		init(edition:UInt64, maxEdition:UInt64){
			self.edition=edition
			self.maxEdition=maxEdition
		}
	}

	//end

	//This interface is here to get this to work before the standard is merged in Artifact
	pub resource interface ViewResolverCollection {
		pub fun borrowViewResolver(id: UInt64): &{ViewResolver}
		pub fun deposit(token: @NonFungibleToken.NFT)
		pub fun getIDs(): [UInt64]
		pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT
	}

	// A struct for Rarity
	// A struct for Rarity Data parts like on flovatar
	// A Display struct for showing the name/thumbnail of something


	pub resource interface TypeConverter {
		pub fun convert(to: Type, value:AnyStruct) : AnyStruct
		pub fun convertTo() : [Type]
		pub fun convertFrom() : Type
	}
}
