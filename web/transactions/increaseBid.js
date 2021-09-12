/** pragma type transaction **/

import {
  getEnvironment,
  replaceImportAddresses,
  reportMissingImports,
  reportMissing,
  sendTransaction
} from 'flow-cadut'

export const CODE = `
  import FIND from "../contracts/FIND.cdc"
import FungibleToken from "../contracts/standard/FungibleToken.cdc"
import FUSD from "../contracts/standard/FUSD.cdc"

transaction(name: String, amount: UFix64) {
	prepare(account: AuthAccount) {

		let vaultRef = account.borrow<&FUSD.Vault>(from: /storage/fusdVault) ?? panic("Could not borrow reference to the owner's Vault!")
		let seller=FIND.lookup(name)!.owner
		
		let fusdReceiver = account.getCapability<&{FungibleToken.Receiver}>(/public/fusdReceiver)

		let leaseCollection = account.getCapability<&{FIND.LeaseCollectionPublic}>(FIND.LeasePublicPath)
		if !leaseCollection.check() {
			account.save(<- FIND.createEmptyLeaseCollection(), to: FIND.LeaseStoragePath)
			account.link<&{FIND.LeaseCollectionPublic}>( FIND.LeasePublicPath, target: FIND.LeaseStoragePath)
		}

		let bidCollection = account.getCapability<&{FIND.BidCollectionPublic}>(FIND.BidPublicPath)
		if !bidCollection.check() {
			account.save(<- FIND.createEmptyBidCollection(receiver: fusdReceiver, leases: leaseCollection), to: FIND.BidStoragePath)
			account.link<&{FIND.BidCollectionPublic}>( FIND.BidPublicPath, target: FIND.BidStoragePath)
		}


		let vault <- vaultRef.withdraw(amount: amount) as! @FUSD.Vault
		let bids = account.borrow<&FIND.BidCollection>(from: FIND.BidStoragePath)!
		bids.increaseBid(name: name, vault: <- vault)

	}
}

`;

/**
* Method to generate cadence code for increaseBid transaction
* @param {Object.<string, string>} addressMap - contract name as a key and address where it's deployed as value
*/
export const increaseBidTemplate = async (addressMap = {}) => {
  const envMap = await getEnvironment();
  const fullMap = {
  ...envMap,
  ...addressMap,
  };

  // If there are any missing imports in fullMap it will be reported via console
  reportMissingImports(CODE, fullMap, `increaseBid =>`)

  return replaceImportAddresses(CODE, fullMap);
};


/**
* Sends increaseBid transaction to the network
* @param {Object.<string, string>} props.addressMap - contract name as a key and address where it's deployed as value
* @param Array<*> props.args - list of arguments
* @param Array<*> props.signers - list of signers
*/
export const increaseBid = async (props) => {
  const { addressMap, args = [], signers = [] } = props;
  const code = await increaseBidTemplate(addressMap);

  reportMissing("arguments", args.length, 2, `increaseBid =>`);
  reportMissing("signers", signers.length, 1, `increaseBid =>`);

  return sendTransaction({code, ...props})
}