

import FUSD from "../contracts/standard/FUSD.cdc"
import FIN from "../contracts/FIN.cdc"
import FungibleToken from "../contracts/standard/FungibleToken.cdc"

transaction(leasePeriod: UFix64) {

	prepare(account: AuthAccount) {
		let wallet=account.getCapability<&{FungibleToken.Receiver}>(/public/fusdReceiver)
		if !wallet.check() {
			let fusd <- FUSD.createEmptyVault()
			account.save(<- fusd, to: /storage/fusdVault)
			account.link<&FUSD.Vault{FungibleToken.Receiver}>( /public/fusdReceiver, target: /storage/fusdVault)
			account.link<&FUSD.Vault{FungibleToken.Balance}>( /public/fusdBalance, target: /storage/fusdVault)
		}

		let adminClient=account.borrow<&FIN.Admin>(from: FIN.AdminClientStoragePath)!

		adminClient.createNetwork(
			admin: account, 
			leasePeriod: leasePeriod,
			lockPeriod: leasePeriod / 2.0,
			secondaryCut: 0.025,
			defaultPrice: 5.0,
			lengthPrices: { 3: 500.0, 4:100.0},
			wallet: wallet)
		}
	}

