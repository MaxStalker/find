import FungibleToken from "./standard/FungibleToken.cdc"
import FlowToken from "./standard/FlowToken.cdc"
import NonFungibleToken from "./standard/NonFungibleToken.cdc"
import MetadataViews from "./standard/MetadataViews.cdc"
import FindViews from "../contracts/FindViews.cdc"
import Profile from "./Profile.cdc"
import Clock from "./Clock.cdc"
import Debug from "./Debug.cdc"
import FIND from "./FIND.cdc"
import FindMarket from "./FindMarket.cdc"
import FindMarketTenant from "../contracts/FindMarketTenant.cdc"
import NFTRegistry from "../contracts/NFTRegistry.cdc"
import FTRegistry from "../contracts/FTRegistry.cdc"

/*

A Find Market for direct sales
*/
pub contract FindMarketSale {

	pub event ForSale(tenant: String, id: UInt64, seller: Address, sellerName: String?, amount: UFix64, status: String, vaultType:String, nft: FindMarket.NFTInfo, buyer:Address?, buyerName:String?)

	//A sale item for a direct sale
	pub resource SaleItem : FindMarket.SaleItem{

		//this is set when bought so that pay will work
		access(self) var buyer: Address?

		access(contract) let vaultType: Type //The type of vault to use for this sale Item
		access(contract) var pointer: FindViews.AuthNFTPointer

		//this field is set if this is a saleItem
		access(contract) var salePrice: UFix64

		//TODO: add valid until?
		init(pointer: FindViews.AuthNFTPointer, vaultType: Type, price:UFix64) {
			self.vaultType=vaultType
			self.pointer=pointer
			self.salePrice=price
			self.buyer=nil
		}

		pub fun getSaleType() : String {
			return "directSale"
		}

		pub fun getListingTypeIdentifier(): String {
			return Type<@SaleItem>().identifier
		}

		pub fun setBuyer(_ address:Address) {
			self.buyer=address
		}

		pub fun getBuyer(): Address? {
			return self.buyer
		}

		pub fun getBuyerName() : String? {
			if let address = self.buyer {
				return FIND.reverseLookup(address)
			}
			return nil
		}

		pub fun getId() : UInt64{
			return self.pointer.getUUID()
		}

		pub fun getItemID() : UInt64 {
			return self.pointer.id
		}

		pub fun getItemType() : Type {
			return self.pointer.getItemType()
		}

		pub fun getItemCollectionAlias() : String {
			return NFTRegistry.getNFTInfoByTypeIdentifier(self.getItemType().identifier)!.alias
		}

		pub fun getRoyalty() : MetadataViews.Royalties? {
			if self.pointer.getViews().contains(Type<MetadataViews.Royalties>()) {
				return self.pointer.resolveView(Type<MetadataViews.Royalties>())! as! MetadataViews.Royalties
			}

			return  nil
		}

		pub fun getSeller() : Address {
			return self.pointer.owner()
		}

		pub fun getSellerName() : String? {
			let address = self.pointer.owner()
			return FIND.reverseLookup(address)
		}

		pub fun getBalance() : UFix64 {
			return self.salePrice
		}

		pub fun getAuction(): FindMarket.AuctionItem? {
			return nil
		}

		pub fun getFtType() : Type  {
			return self.vaultType
		}

		pub fun getFtAlias() : String {
			return FTRegistry.getFTInfoByTypeIdentifier(self.getFtType().identifier)!.alias
		}

		pub fun getValidUntil() : UFix64? {
			return nil 
		}

		pub fun toNFTInfo() : FindMarket.NFTInfo{
			return FindMarket.NFTInfo(self.pointer.getViewResolver(), id: self.pointer.id)
		}

	}

	pub resource interface SaleItemCollectionPublic {
		//fetch all the tokens in the collection
		pub fun getIds(): [UInt64]
		pub fun getSaleInformation(_ id:UInt64) : FindMarket.SaleItemInformation?
		pub fun getSaleItemReport() : FindMarket.SaleItemCollectionReport

		pub fun buy(id: UInt64, vault: @FungibleToken.Vault, nftCap: Capability<&{NonFungibleToken.Receiver}>) 
	}

	pub resource SaleItemCollection: SaleItemCollectionPublic, FindMarket.SaleItemCollectionPublic {
		//is this the best approach now or just put the NFT inside the saleItem?
		access(contract) var items: @{UInt64: SaleItem}

		access(contract) let tenantCapability: Capability<&FindMarketTenant.Tenant{FindMarketTenant.TenantPublic}>

		init (_ tenantCapability: Capability<&FindMarketTenant.Tenant{FindMarketTenant.TenantPublic}>) {
			self.items <- {}
			self.tenantCapability=tenantCapability
		}

		access(self) fun getTenant() : &FindMarketTenant.Tenant{FindMarketTenant.TenantPublic} {
			pre{
				self.tenantCapability.check() : "Tenant client is not linked anymore"
			}
			return self.tenantCapability.borrow()!
		}

		pub fun getSaleInformation(_ id:UInt64) : FindMarket.SaleItemInformation? {
			pre {
				self.items.containsKey(id) : "Invalid id=".concat(id.toString())
			}
			let item=self.borrow(id)
			let tenant=self.getTenant()
			let info = self.checkSaleInformation(tenant: tenant, ids: [id], getGhost: false)
			if info.items.length > 0 {
				return info.items[0]
			}
			return nil
		}

		// todo: do we need this here?
		pub fun getSaleItemReport() : FindMarket.SaleItemCollectionReport {
			let tenant=self.getTenant()
			return self.checkSaleInformation(tenant: tenant, ids: self.getIds(), getGhost: true)
		}

		access(contract) fun checkSaleInformation(tenant: &FindMarketTenant.Tenant{FindMarketTenant.TenantPublic}, ids: [UInt64], getGhost:Bool) : FindMarket.SaleItemCollectionReport {
			let ghost: [FindMarket.GhostListing] =[]
			let info: [FindMarket.SaleItemInformation] =[]
			let listingType = self.getListingType()
			for id in ids {
				let item=self.borrow(id)
				if !item.pointer.valid() {
					if getGhost {
						ghost.append(FindMarket.GhostListing(listingType: listingType, id:id))
					}
					continue
				} 
				//TODO: do we need to be smarter about this?
				let stopped=tenant.allowedAction(listingType: listingType, nftType: item.getItemType(), ftType: item.getFtType(), action: FindMarketTenant.MarketAction(listing:false, "delist item for sale"))
				var status="active"
				if !stopped.allowed {
					status="stopped"
				}
				let deprecated=tenant.allowedAction(listingType: listingType, nftType: item.getItemType(), ftType: item.getFtType(), action: FindMarketTenant.MarketAction(listing:true, "delist item for sale"))

				if !deprecated.allowed {
					status="deprecated"
				}

				if let validTime = item.getValidUntil() {
					if validTime >= getCurrentBlock().timestamp{
						status="ended"
					}
				}
				info.append(FindMarket.SaleItemInformation(item, status))
			}

			return FindMarket.SaleItemCollectionReport(items: info, ghosts: ghost)
		}

		pub fun getListingType() : Type {
			return Type<@SaleItem>()
		}

		pub fun buy(id: UInt64, vault: @FungibleToken.Vault, nftCap: Capability<&{NonFungibleToken.Receiver}>) {
			pre {
				self.items.containsKey(id) : "Invalid id=".concat(id.toString())
			}

			//TODO: check valid until
			let saleItem=self.borrow(id)

			if saleItem.salePrice != vault.balance {
				panic("Incorrect balance sent in vault. Expected ".concat(saleItem.salePrice.toString()).concat(" got ").concat(vault.balance.toString()))
			}

			if saleItem.vaultType != vault.getType() {
				panic("This item can be baught using ".concat(saleItem.vaultType.identifier).concat(" you have sent in ").concat(vault.getType().identifier))
			}

			let actionResult=self.getTenant().allowedAction(listingType: Type<@FindMarketSale.SaleItem>(), nftType: saleItem.getItemType(), ftType: saleItem.getFtType(), action: FindMarketTenant.MarketAction(listing:false, "buy item for sale"))

			if !actionResult.allowed {
				panic(actionResult.message)
			}

			let cuts= self.getTenant().getTeantCut(name: actionResult.name, listingType: Type<@FindMarketSale.SaleItem>(), nftType: saleItem.getItemType(), ftType: saleItem.getFtType())

			let ftType=saleItem.vaultType
			let owner=saleItem.getSeller()
			let nftInfo= saleItem.toNFTInfo()

			let royalty=saleItem.getRoyalty()
			let soldFor=saleItem.getBalance()
			saleItem.setBuyer(nftCap.address)
			let buyer=nftCap.address

			emit ForSale(tenant:self.getTenant().name, id: id, seller:owner, sellerName: FIND.reverseLookup(owner), amount: soldFor, status:"sold", vaultType: ftType.identifier, nft:nftInfo, buyer: buyer, buyerName: FIND.reverseLookup(buyer))

			FindMarket.pay(tenant:self.getTenant().name, id:id, saleItem: saleItem, vault: <- vault, royalty:royalty, nftInfo:nftInfo, cuts:cuts)
			nftCap.borrow()!.deposit(token: <- saleItem.pointer.withdraw())

			destroy <- self.items.remove(key: id)
		}

		pub fun listForSale(pointer: FindViews.AuthNFTPointer, vaultType: Type, directSellPrice:UFix64) {

			// What happends if we relist  
			let saleItem <- create SaleItem(pointer: pointer, vaultType:vaultType, price: directSellPrice)

			let actionResult=self.getTenant().allowedAction(listingType: Type<@FindMarketSale.SaleItem>(), nftType: saleItem.getItemType(), ftType: saleItem.getFtType(), action: FindMarketTenant.MarketAction(listing:true, "list item for sale"))

			if !actionResult.allowed {
				panic(actionResult.message)
			}

			let owner=self.owner!.address
			emit ForSale(tenant: self.getTenant().name, id: pointer.getUUID(), seller:owner, sellerName: FIND.reverseLookup(owner), amount: saleItem.salePrice, status: "listed", vaultType: vaultType.identifier, nft:FindMarket.NFTInfo(pointer.getViewResolver(), id: pointer.id), buyer: nil, buyerName:nil)
			let old <- self.items[pointer.getUUID()] <- saleItem
			destroy old

		}

		pub fun delist(_ id: UInt64) {
			pre {
				self.items.containsKey(id) : "Unknown item with id=".concat(id.toString())
			}

			let saleItem <- self.items.remove(key: id)!

			let actionResult=self.getTenant().allowedAction(listingType: Type<@FindMarketSale.SaleItem>(), nftType: saleItem.getItemType(), ftType: saleItem.getFtType(), action: FindMarketTenant.MarketAction(listing:false, "delist item for sale"))

			if !actionResult.allowed {
				panic(actionResult.message)
			}

			let owner=self.owner!.address
			emit ForSale(tenant:self.getTenant().name, id: id, seller:owner, sellerName:FIND.reverseLookup(owner), amount: saleItem.salePrice, status: "cancelled", vaultType: saleItem.vaultType.identifier,nft: FindMarket.NFTInfo(saleItem.pointer.getViewResolver(), id:saleItem.pointer.id), buyer:nil, buyerName:nil)
			destroy saleItem
		}

		pub fun getIds(): [UInt64] {
			return self.items.keys
		}

		pub fun borrow(_ id: UInt64): &SaleItem {
			return &self.items[id] as &SaleItem
		}

		destroy() {
			destroy self.items
		}
	}


	//Create an empty lease collection that store your leases to a name
	pub fun createEmptySaleItemCollection(_ tenantCapability: Capability<&FindMarketTenant.Tenant{FindMarketTenant.TenantPublic}>): @SaleItemCollection {
		let wallet=FindMarketSale.account.getCapability<&{FungibleToken.Receiver}>(Profile.publicReceiverPath)
		return <- create SaleItemCollection(tenantCapability)
	}

	//BAM: will this work?
	//pub fun getFindSaleItemCapability(_ user: Address) : Capability<&FindMarketSale.SaleItemCollection{FindMarketSale.SaleItemCollectionPublic, FindMarket.SaleItemCollection}>? {
	pub fun getFindSaleItemCapability(_ user: Address) : Capability<&FindMarketSale.SaleItemCollection{FindMarketSale.SaleItemCollectionPublic, FindMarket.SaleItemCollectionPublic}>? {
		return FindMarketSale.getSaleItemCapability(marketplace: FindMarketSale.account.address, user:user) 
	}

	pub fun getSaleItemCapability(marketplace:Address, user:Address) : Capability<&FindMarketSale.SaleItemCollection{FindMarketSale.SaleItemCollectionPublic, FindMarket.SaleItemCollectionPublic}>? {
		pre{
			FindMarketTenant.getTenantCapability(marketplace) != nil : "Invalid tenant"
		}
		if let tenant=FindMarketTenant.getTenantCapability(marketplace)!.borrow() {
			return getAccount(user).getCapability<&FindMarketSale.SaleItemCollection{FindMarketSale.SaleItemCollectionPublic, FindMarket.SaleItemCollectionPublic}>(tenant.getPublicPath(Type<@FindMarketSale.SaleItemCollection>()))
		}
		return nil
	}
}
