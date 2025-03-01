# .find

.find is a solution that aims to make it easier to .find people and their things on the flow blockchain. It is live on mainnet since Dec 13th 2021 at https://find.xyz


## Developing .find

See [developing](developing.md) for how to develop on .find

## Integrate with .find
The following solutions are integrated with .find
 - https://versus.auction
 - https://schwap.io
 - https://flovatar.com
 - https://flowscan.org

see [integrating](integrating.md) for how to integrate with .find

## Using .find
.find contracts could be found deployed on both testnet and mainnet networks:
- testnet - **0xa16ab1d0abde3625** ([Flow View Source - Testnet - FIND Contract](https://flow-view-source.com/testnet/account/0xa16ab1d0abde3625/contract/FIND))
- mainnet - **0x097bafa4e0b48eef** ([Flow View Source - Mainnet - FIND Contract](https://flow-view-source.com/mainnet/account/0x097bafa4e0b48eef/contract/FIND))

You can use [resolve.cdc](https://github.com/MaxStalker/find/blob/main/scripts/resolve.cdc) and [name.cdc](https://github.com/MaxStalker/find/blob/main/scripts/name.cdc) scripts for direct and reverse lookup. Simply replace add in import statement.

## Testing
  
 `gotestsum -f testname --watch`

