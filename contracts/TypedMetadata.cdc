import FungibleToken from "../contracts/standard/FungibleToken.cdc"
import NonFungibleToken from "../contracts/standard/NonFungibleToken.cdc"
import FlowToken from "../contracts/standard/FlowToken.cdc"
import FUSD from "../contracts/standard/FUSD.cdc"

pub contract TypedMetadata {

	pub resource interface ViewResolver {
		pub fun getViews() : [Type]
		pub fun resolveView(_ view:Type): AnyStruct
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


	pub struct Royalties{
		pub let royalty: { String : Royalty}
		init(royalty: {String : Royalty}) {
			self.royalty=royalty
		}
	}

	pub struct Royalty{
		pub let wallets: { String : Capability<&{FungibleToken.Receiver}>  }
		pub let cut: UFix64

		pub let percentage: Bool
		pub let owner: Address

		init(wallets:{ String: Capability<&{FungibleToken.Receiver}>}, cut: UFix64, percentage: Bool, owner: Address ){
			self.wallets=wallets
			self.cut=cut
			self.percentage=percentage
			self.owner=owner
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

	pub fun createPercentageRoyalty(user:Address, cut: UFix64) : Royalty {
		let userAccount=getAccount(user)
		let fusdReceiver = userAccount.getCapability<&{FungibleToken.Receiver}>(/public/fusdReceiver)
		let flowReceiver = userAccount.getCapability<&{FungibleToken.Receiver}>(/public/flowTokenReceiver)
		let walletDicts :{ String : Capability<&{FungibleToken.Receiver}> }= {}
		walletDicts[Type<@FUSD.Vault>().identifier]=fusdReceiver
		walletDicts[Type<@FlowToken.Vault>().identifier]=flowReceiver
		let userRoyalty = TypedMetadata.Royalty(wallets: walletDicts, cut: cut, percentage:true, owner:user)

		return userRoyalty
	}

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
