import Admin from "../contracts/Admin.cdc"
import FIND from "../contracts/FIND.cdc"
import FUSD from "../contracts/standard/FUSD.cdc"
import FungibleToken from "../contracts/standard/FungibleToken.cdc"

//link together the administrator to the client, signed by the owner of the contract
transaction(ownerAddress: Address) {

	//versus account
	prepare(account: AuthAccount) {

		let owner= getAccount(ownerAddress)
		let wallet=account.getCapability<&{FungibleToken.Receiver}>(/public/fusdReceiver)
		if !wallet.check() {
			let fusd <- FUSD.createEmptyVault()
			account.save(<- fusd, to: /storage/fusdVault)
			account.link<&FUSD.Vault{FungibleToken.Receiver}>( /public/fusdReceiver, target: /storage/fusdVault)
			account.link<&FUSD.Vault{FungibleToken.Balance}>( /public/fusdBalance, target: /storage/fusdVault)
		}


		let client= owner.getCapability<&{Admin.AdminProxyClient}>(Admin.AdminProxyPublicPath)
		.borrow() ?? panic("Could not borrow admin client")

		let network=account.getCapability<&FIND.Network>(FIND.NetworkPrivatePath)
		client.addCapability(network)

	}
}

