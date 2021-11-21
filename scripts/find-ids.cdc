import NonFungibleToken from "../contracts/standard/NonFungibleToken.cdc"

pub fun main(address: Address, path: PublicPath) : [UInt64] {

	let account=getAccount(address)
	return account.getCapability(path).borrow<&{NonFungibleToken.CollectionPublic}>()!.getIDs()

}

