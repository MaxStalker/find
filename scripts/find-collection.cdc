import Profile from "../contracts/Profile.cdc"
import NonFungibleToken from "../contracts/standard/NonFungibleToken.cdc"

pub fun main(address: Address) : [String] {

	let account=getAccount(address)

	let profileCap= getAccount(address).getCapability<&{Profile.Public}>(Profile.publicPath)

	if !profileCap.check() {
		return ["Unknown, no profile created"]
	}

	let collections= profileCap.borrow()!.getCollections()

	var names:  [String] = []
	for col in collections {
		if col.type == Type<&{NonFungibleToken.CollectionPublic}>() {
			names.append(col.name)
		}
	}
	return names

}

