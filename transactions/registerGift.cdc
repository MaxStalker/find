import FungibleToken from "../contracts/standard/FungibleToken.cdc"
import FUSD from "../contracts/standard/FUSD.cdc"
import FlowToken from "../contracts/standard/FlowToken.cdc"
import Profile from "../contracts/Profile.cdc"
import FIND from "../contracts/FIND.cdc"

transaction(name: String, amount: UFix64, recipient: Address) {
	prepare(acct: AuthAccount) {

		let price=FIND.calculateCost(name)
		if price != amount {
			panic("Calculated cost does not match expected cost")
		}
		log("The cost for registering this name is ".concat(price.toString()))

		let vaultRef = acct.borrow<&FUSD.Vault>(from: /storage/fusdVault) ?? panic("Could not borrow reference to the fusdVault!")
		let payVault <- vaultRef.withdraw(amount: price) as! @FUSD.Vault

		let leases=acct.borrow<&FIND.LeaseCollection>(from: FIND.LeaseStoragePath)!
		leases.register(name: name, vault: <- payVault)

		let receiver = getAccount(recipient)
		let receiverLease = receiver.getCapability<&FIND.LeaseCollection{FIND.LeaseCollectionPublic}>(FIND.LeasePublicPath)
		let receiverProfile = receiver.getCapability<&{Profile.Public}>(Profile.publicPath)
		if !receiverLease.check() {
			panic("Receiver is not a find user")
		}
		leases.move(name: name, profile: receiverProfile, to: receiverLease)
	}
}
