package test_main

import (
	"testing"

	"github.com/bjartek/overflow/overflow"
	"github.com/stretchr/testify/assert"
)

func TestFTRegistry(t *testing.T) {
	t.Run("Should be able to registry flow token", func(t *testing.T) {
		NewOverflowTest(t).
			setupFIND().
			registerFTInFtRegistry("flow", "A.f8d6e0586b0a20c7.FTRegistry.FTInfoRegistered", map[string]interface{}{
				"alias":          "Flow",
				"typeIdentifier": "A.0ae53cb6e3f42a79.FlowToken.Vault",
			})

	})

	t.Run("Should be able to registry flow token and get it", func(t *testing.T) {
		otu := NewOverflowTest(t).
			setupFIND().
			registerFTInFtRegistry("flow", "A.f8d6e0586b0a20c7.FTRegistry.FTInfoRegistered", map[string]interface{}{
				"alias":          "Flow",
				"typeIdentifier": "A.0ae53cb6e3f42a79.FlowToken.Vault",
			})

		o := otu.O
		result := o.ScriptFromFile("getFTInfoByTypeIdentifier").
			Args(o.Arguments().String("A.0ae53cb6e3f42a79.FlowToken.Vault")).
			RunReturnsInterface()

		expected := map[string]interface{}{
			"alias":                  "Flow",
			"balancePath":            "/public/flowTokenBalance",
			"balancePathIdentifier":  "flowTokenBalance",
			"icon":                   "https://static.flowscan.org/mainnet/icons/A.1654653399040a61.FlowToken.png",
			"receiverPath":           "/public/flowTokenReceiver",
			"receiverPathIdentifier": "flowTokenReceiver",
			"tag":                    []interface{}{"utility coin"},
			"type":                   "Type<A.0ae53cb6e3f42a79.FlowToken.Vault>()",
			"typeIdentifier":         "A.0ae53cb6e3f42a79.FlowToken.Vault",
			"vaultPath":              "/storage/flowTokenVault",
			"vaultPathIdentifier":    "flowTokenVault",
		}
		assert.Equal(t, expected, result)

	})

	t.Run("Should be able to registry flow token and get it by alias", func(t *testing.T) {
		otu := NewOverflowTest(t).
			setupFIND().
			registerFTInFtRegistry("flow", "A.f8d6e0586b0a20c7.FTRegistry.FTInfoRegistered", map[string]interface{}{
				"alias":          "Flow",
				"typeIdentifier": "A.0ae53cb6e3f42a79.FlowToken.Vault",
			})

		o := otu.O
		result := o.ScriptFromFile("getFTInfoByAlias").
			Args(o.Arguments().String("Flow")).
			RunReturnsInterface()

		expected := map[string]interface{}{
			"alias":                  "Flow",
			"balancePath":            "/public/flowTokenBalance",
			"balancePathIdentifier":  "flowTokenBalance",
			"icon":                   "https://static.flowscan.org/mainnet/icons/A.1654653399040a61.FlowToken.png",
			"receiverPath":           "/public/flowTokenReceiver",
			"receiverPathIdentifier": "flowTokenReceiver",
			"tag":                    []interface{}{"utility coin"},
			"type":                   "Type<A.0ae53cb6e3f42a79.FlowToken.Vault>()",
			"typeIdentifier":         "A.0ae53cb6e3f42a79.FlowToken.Vault",
			"vaultPath":              "/storage/flowTokenVault",
			"vaultPathIdentifier":    "flowTokenVault",
		}
		assert.Equal(t, expected, result)

	})

	t.Run("Should be able to registry flow token, fusd token and get list from it", func(t *testing.T) {
		expected := `
			{
			    "A.0ae53cb6e3f42a79.FlowToken.Vault": {
			        "alias": "Flow",
			        "balancePath": "/public/flowTokenBalance",
					"balancePathIdentifier":"flowTokenBalance",
			        "icon": "https://static.flowscan.org/mainnet/icons/A.1654653399040a61.FlowToken.png",
					"tag" : ["utility coin"],
			        "receiverPath": "/public/flowTokenReceiver",
					"receiverPathIdentifier":"flowTokenReceiver",
			        "type": "Type\u003cA.0ae53cb6e3f42a79.FlowToken.Vault\u003e()",
			        "typeIdentifier": "A.0ae53cb6e3f42a79.FlowToken.Vault",
			        "vaultPath": "/storage/flowTokenVault",
					"vaultPathIdentifier":"flowTokenVault"
			    },
				"A.f8d6e0586b0a20c7.FUSD.Vault": {
			        "alias": "FUSD",
			        "balancePath": "/public/fusdBalance",
					"balancePathIdentifier":"fusdBalance",
			        "icon": "https://static.flowscan.org/mainnet/icons/A.3c5959b568896393.FUSD.png",
					"tag" : ["stablecoin"],
			        "receiverPath": "/public/fusdReceiver",
					"receiverPathIdentifier":"fusdReceiver",
			        "type": "Type\u003cA.f8d6e0586b0a20c7.FUSD.Vault\u003e()",
			        "typeIdentifier": "A.f8d6e0586b0a20c7.FUSD.Vault",
			        "vaultPath": "/storage/fusdVault",
					"vaultPathIdentifier":"fusdVault"
			    }
			}
			`
		NewOverflowTest(t).
			setupFIND().
			registerFTInFtRegistry("flow", "A.f8d6e0586b0a20c7.FTRegistry.FTInfoRegistered", map[string]interface{}{
				"alias":          "Flow",
				"typeIdentifier": "A.0ae53cb6e3f42a79.FlowToken.Vault",
			}).
			registerFTInFtRegistry("fusd", "A.f8d6e0586b0a20c7.FTRegistry.FTInfoRegistered", map[string]interface{}{
				"alias":          "FUSD",
				"typeIdentifier": "A.f8d6e0586b0a20c7.FUSD.Vault",
			}).
			scriptEqualToJson("getFTInfoAll", expected)

	})

	t.Run("Should not be able to overrride a ft without removing it first", func(t *testing.T) {
		otu := NewOverflowTest(t).
			setupFIND().
			registerFTInFtRegistry("flow", "A.f8d6e0586b0a20c7.FTRegistry.FTInfoRegistered", map[string]interface{}{
				"alias":          "Flow",
				"typeIdentifier": "A.0ae53cb6e3f42a79.FlowToken.Vault",
			})

		o := otu.O
		o.TransactionFromFile("setFTInfo_flow").
			SignProposeAndPayAs("find").
			Args(o.Arguments()).
			Test(t).
			AssertFailure("This FungibleToken Register already exist")
	})

	t.Run("Should be able to registry and remove flow token by Alias, as well as return nil on scripts", func(t *testing.T) {
		otu := NewOverflowTest(t).
			setupFIND().
			registerFTInFtRegistry("flow", "A.f8d6e0586b0a20c7.FTRegistry.FTInfoRegistered", map[string]interface{}{
				"alias":          "Flow",
				"typeIdentifier": "A.0ae53cb6e3f42a79.FlowToken.Vault",
			}).
			removeFTInFtRegistry("removeFTInfoByAlias", "Flow",
				"A.f8d6e0586b0a20c7.FTRegistry.FTInfoRemoved", map[string]interface{}{
					"alias":          "Flow",
					"typeIdentifier": "A.0ae53cb6e3f42a79.FlowToken.Vault",
				})

		o := otu.O
		aliasResult := o.ScriptFromFile("getFTInfoByAlias").
			Args(o.Arguments().String("Flow")).
			RunReturnsInterface()
		assert.Equal(t, "", aliasResult)

		infoResult := o.ScriptFromFile("getFTInfoByAlias").
			Args(o.Arguments().String("Flow")).
			RunReturnsInterface()
		assert.Equal(t, "", infoResult)

	})

	t.Run("Should be able to registry and remove flow token by Type Identifier, as well as return nil on scripts", func(t *testing.T) {
		otu := NewOverflowTest(t).
			setupFIND().
			registerFTInFtRegistry("flow", "A.f8d6e0586b0a20c7.FTRegistry.FTInfoRegistered", map[string]interface{}{
				"alias":          "Flow",
				"typeIdentifier": "A.0ae53cb6e3f42a79.FlowToken.Vault",
			}).
			removeFTInFtRegistry("removeFTInfoByTypeIdentifier", "A.0ae53cb6e3f42a79.FlowToken.Vault",
				"A.f8d6e0586b0a20c7.FTRegistry.FTInfoRemoved", map[string]interface{}{
					"alias":          "Flow",
					"typeIdentifier": "A.0ae53cb6e3f42a79.FlowToken.Vault",
				})

		o := otu.O
		aliasResult := o.ScriptFromFile("getFTInfoByTypeIdentifier").
			Args(o.Arguments().String("A.0ae53cb6e3f42a79.FlowToken.Vault")).
			RunReturnsInterface()
		assert.Equal(t, "", aliasResult)

		infoResult := o.ScriptFromFile("getFTInfoByAlias").
			Args(o.Arguments().String("Flow")).
			RunReturnsInterface()
		assert.Equal(t, "", infoResult)

	})

	t.Run("Should be able to registry usdc token and get it", func(t *testing.T) {
		otu := NewOverflowTest(t).
			setupFIND().
			registerFTInFtRegistry("usdc", "A.f8d6e0586b0a20c7.FTRegistry.FTInfoRegistered", map[string]interface{}{
				"alias":          "USDC",
				"typeIdentifier": "A.f8d6e0586b0a20c7.FiatToken.Vault",
			})

		o := otu.O
		result := o.ScriptFromFile("getFTInfoByTypeIdentifier").
			Args(o.Arguments().String("A.f8d6e0586b0a20c7.FiatToken.Vault")).
			RunReturnsInterface()

		expected := map[string]interface{}{
			"alias":                  "USDC",
			"balancePath":            "/public/USDCVaultBalance",
			"balancePathIdentifier":  "USDCVaultBalance",
			"icon":                   "https://static.flowscan.org/mainnet/icons/A.b19436aae4d94622.FiatToken.png",
			"receiverPath":           "/public/USDCVaultReceiver",
			"receiverPathIdentifier": "USDCVaultReceiver",
			"tag":                    []interface{}{"stablecoin"},
			"type":                   "Type<A.f8d6e0586b0a20c7.FiatToken.Vault>()",
			"typeIdentifier":         "A.f8d6e0586b0a20c7.FiatToken.Vault",
			"vaultPath":              "/storage/USDCVault",
			"vaultPathIdentifier":    "USDCVault",
		}

		assert.Equal(t, expected, result)

	})

	t.Run("Should be able to send usdc to another name", func(t *testing.T) {

		otu := NewOverflowTest(t).
			setupFIND().
			createUser(100.0, "user1").
			createUser(100.0, "user2").
			registerFtInRegistry().
			registerUser("user1").
			registerUser("user2")

		otu.O.TransactionFromFile("sendFT").
			SignProposeAndPayAs("user2").
			Args(otu.O.Arguments().
				String("user1").
				UFix64(5.0).
				String("USDC").
				String("test").
				String("This is a message")).
			Test(t).AssertSuccess().
			AssertEmitEvent(overflow.NewTestEvent("A.f8d6e0586b0a20c7.FiatToken.TokensDeposited", map[string]interface{}{
				"amount": "5.00000000",
				"to":     "0x179b6b1cb6755e31",
			})).
			AssertEmitEvent(overflow.NewTestEvent("A.f8d6e0586b0a20c7.FiatToken.TokensWithdrawn", map[string]interface{}{
				"amount": "5.00000000",
				"from":   "0xf3fcd2c1a78f5eee",
			})).
			AssertPartialEvent(overflow.NewTestEvent("A.f8d6e0586b0a20c7.FIND.FungibleTokenSent", map[string]interface{}{
				"from":      "0xf3fcd2c1a78f5eee",
				"fromName":  "user2",
				"toAddress": "0x179b6b1cb6755e31",
				"amount":    "5.00000000",
				"name":      "user1",
				"tag":       "test",
				"message":   "This is a message",
			}))
	})

	t.Run("Should be able to send fusd to another name", func(t *testing.T) {

		otu := NewOverflowTest(t).
			setupFIND().
			createUser(100.0, "user1").
			createUser(100.0, "user2").
			registerFtInRegistry().
			registerUser("user1").
			registerUser("user2")

		otu.O.TransactionFromFile("sendFT").
			SignProposeAndPayAs("user2").
			Args(otu.O.Arguments().
				String("user1").
				UFix64(5.0).
				String("FUSD").
				String("test").
				String("This is a message")).
			Test(t).AssertSuccess().
			AssertEmitEvent(overflow.NewTestEvent("A.f8d6e0586b0a20c7.FUSD.TokensDeposited", map[string]interface{}{
				"amount": "5.00000000",
				"to":     "0x179b6b1cb6755e31",
			})).
			AssertEmitEvent(overflow.NewTestEvent("A.f8d6e0586b0a20c7.FUSD.TokensWithdrawn", map[string]interface{}{
				"amount": "5.00000000",
				"from":   "0xf3fcd2c1a78f5eee",
			})).
			AssertPartialEvent(overflow.NewTestEvent("A.f8d6e0586b0a20c7.FIND.FungibleTokenSent", map[string]interface{}{
				"from":      "0xf3fcd2c1a78f5eee",
				"fromName":  "user2",
				"toAddress": "0x179b6b1cb6755e31",
				"amount":    "5.00000000",
				"name":      "user1",
				"tag":       "test",
				"message":   "This is a message",
			}))
	})

	t.Run("Should be able to send flow to another name", func(t *testing.T) {

		otu := NewOverflowTest(t).
			setupFIND().
			createUser(100.0, "user1").
			createUser(100.0, "user2").
			registerFtInRegistry().
			registerUser("user1").
			registerUser("user2")

		otu.O.TransactionFromFile("sendFT").
			SignProposeAndPayAs("user2").
			Args(otu.O.Arguments().
				String("user1").
				UFix64(5.0).
				String("Flow").
				String("test").
				String("This is a message")).
			Test(t).AssertSuccess().
			AssertEmitEvent(overflow.NewTestEvent("A.0ae53cb6e3f42a79.FlowToken.TokensDeposited", map[string]interface{}{
				"amount": "5.00000000",
				"to":     "0x179b6b1cb6755e31",
			})).
			AssertEmitEvent(overflow.NewTestEvent("A.0ae53cb6e3f42a79.FlowToken.TokensWithdrawn", map[string]interface{}{
				"amount": "5.00000000",
				"from":   "0xf3fcd2c1a78f5eee",
			})).
			AssertPartialEvent(overflow.NewTestEvent("A.f8d6e0586b0a20c7.FIND.FungibleTokenSent", map[string]interface{}{
				"from":      "0xf3fcd2c1a78f5eee",
				"fromName":  "user2",
				"toAddress": "0x179b6b1cb6755e31",
				"amount":    "5.00000000",
				"name":      "user1",
				"tag":       "test",
				"message":   "This is a message",
			}))
	})
}
