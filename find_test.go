package test_main

import (
	"testing"

	"github.com/bjartek/go-with-the-flow/v2/gwtf"
	"github.com/stretchr/testify/assert"
)

/*
Tests must be in the same folder as flow.json with contracts and transactions/scripts in subdirectories in order for the path resolver to work correctly
*/
func TestFIND(t *testing.T) {

	t.Run("Should be able to register a name", func(t *testing.T) {
		NewGWTFTest(t).
			setupFIND().
			createUser("100.0", "user1").
			registerUser("user1")
	})

	t.Run("Should get error if you try to register a name and dont have enough money", func(t *testing.T) {

		gt := NewGWTFTest(t).
			setupFIND().
			createUser("30.0", "user1")

		gt.GWTF.TransactionFromFile("register").
			SignProposeAndPayAs("user1").
			StringArgument("usr").
			UFix64Argument("500.0").
			Test(t).
			AssertFailure("Amount withdrawn must be less than or equal than the balance of the Vault")

	})

	t.Run("Should get error if you try to register a name that is too short", func(t *testing.T) {

		gt := NewGWTFTest(t).
			setupFIND().
			createUser("30.0", "user1")

		gt.GWTF.TransactionFromFile("register").
			SignProposeAndPayAs("user1").
			StringArgument("ur").
			UFix64Argument("5.0").
			Test(t).
			AssertFailure("A FIND name has to be lower-cased alphanumeric or dashes and between 3 and 16 characters")

	})

	t.Run("Should get error if you try to register a name that is already claimed", func(t *testing.T) {

		gt := NewGWTFTest(t).
			setupFIND().
			createUser("30.0", "user1").
			registerUser("user1")

		gt.GWTF.TransactionFromFile("register").
			SignProposeAndPayAs("user1").
			StringArgument("user1").
			UFix64Argument("5.0").
			Test(t).
			AssertFailure("Name already registered")

	})

	t.Run("Should allow registering a lease after it is freed", func(t *testing.T) {

		gt := NewGWTFTest(t).
			setupFIND().
			createUser("100.0", "user1").
			registerUser("user1")

		gt.expireLease().tickClock("2.0")

		gt.GWTF.Transaction(`
import FIND from "../contracts/FIND.cdc"

transaction(name: String) {

    prepare(account: AuthAccount) {
        let status=FIND.status(name)
				if status.status == FIND.LeaseStatus.LOCKED {
					panic("locked")
				}
				if status.status == FIND.LeaseStatus.FREE {
					panic("free")
				}
    }
}
`).
			SignProposeAndPayAs("user1").
			StringArgument("user1").
			Test(t).AssertFailure("locked")

		gt.expireLease()
		gt.registerUser("user1")
	})

	t.Run("Should be able to lookup address", func(t *testing.T) {

		NewGWTFTest(t).
			setupFIND().
			createUser("100.0", "user1").
			registerUser("user1").
			assertLookupAddress("user1", "0x179b6b1cb6755e31")
	})

	t.Run("Should not be able to lookup lease after expired", func(t *testing.T) {

		gt := NewGWTFTest(t).
			setupFIND().
			createUser("100.0", "user1").
			registerUser("user1").
			expireLease().
			tickClock("2.0")

		value := gt.GWTF.ScriptFromFile("status").StringArgument("user1").RunReturnsInterface()
		assert.Equal(t, "", value)

	})

	t.Run("Should be able to send ft to another name", func(t *testing.T) {

		gt := NewGWTFTest(t).
			setupFIND().
			createUser("100.0", "user1").
			createUser("100.0", "user2").
			registerUser("user1")

		gt.GWTF.TransactionFromFile("sendFT").
			SignProposeAndPayAs("user2").
			StringArgument("user1").
			UFix64Argument("5.0").
			StringArgument("fusd").
			Test(t).AssertSuccess().
			AssertEmitEvent(gwtf.NewTestEvent("A.f8d6e0586b0a20c7.FUSD.TokensDeposited", map[string]interface{}{
				"amount": "5.00000000",
				"to":     "0x179b6b1cb6755e31",
			})).
			AssertEmitEvent(gwtf.NewTestEvent("A.f8d6e0586b0a20c7.FUSD.TokensWithdrawn", map[string]interface{}{
				"amount": "5.00000000",
				"from":   "0xf3fcd2c1a78f5eee",
			}))

	})

	t.Run("Admin should be able to register without paying FUSD", func(t *testing.T) {

		gt := NewGWTFTest(t).
			setupFIND().
			createUser("10.0", "find")

		gt.GWTF.TransactionFromFile("registerAdmin").
			SignProposeAndPayAs("find").
			StringArrayArgument("find-admin").
			AccountArgument("find").
			Test(gt.T).
			AssertSuccess().
			AssertPartialEvent(gwtf.NewTestEvent("A.f8d6e0586b0a20c7.FIND.Register", map[string]interface{}{
				"name": "find-admin",
			}))

	})

	t.Run("Should be able to send lease to another name", func(t *testing.T) {

		gt := NewGWTFTest(t).
			setupFIND().
			createUser("100.0", "user1").
			createUser("100.0", "user2").
			registerUser("user1").
			registerUser("user2")

		gt.GWTF.TransactionFromFile("moveNameToName").
			SignProposeAndPayAs("user1").
			StringArgument("user1").
			StringArgument("user2").
			Test(t).AssertSuccess().
			AssertPartialEvent(gwtf.NewTestEvent("A.f8d6e0586b0a20c7.FIND.Register", map[string]interface{}{
				"name": "user1",
			})).
			AssertPartialEvent(gwtf.NewTestEvent("A.f8d6e0586b0a20c7.FIND.Moved", map[string]interface{}{
				"name": "user1",
			}))
	})

	t.Run("Should be able to send fusd with message", func(t *testing.T) {
		gt := NewGWTFTest(t).
			setupFIND().
			createUser("100.0", "user1").
			createUser("100.0", "user2").
			registerUser("user1").
			registerUser("user2")

		gt.GWTF.TransactionFromFile("sendFusdWithMessage").
			SignProposeAndPayAs("user1").
			StringArgument("user2").
			UFix64Argument("5.0").
			StringArgument("Happy to help").
			Test(t).AssertSuccess().
			AssertPartialEvent(gwtf.NewTestEvent("A.f8d6e0586b0a20c7.Profile.Verification", map[string]interface{}{
				"message": "user1 sent 5.00 FUSD with message:Happy to help",
			}))

	})

	t.Run("Should be able to send fusd to another name with tag and message", func(t *testing.T) {

		gt := NewGWTFTest(t).
			setupFIND().
			createUser("100.0", "user1").
			createUser("100.0", "user2").
			registerUser("user1").
			registerUser("user2")

		gt.GWTF.TransactionFromFile("sendFusdWithTagAndMessage").
			SignProposeAndPayAs("user2").
			StringArgument("user1").
			UFix64Argument("5.0").
			StringArgument("This is a test").
			StringArgument("test").
			Test(t).AssertSuccess().
			AssertPartialEvent(gwtf.NewTestEvent("A.f8d6e0586b0a20c7.FIND.FungibleTokenSent", map[string]interface{}{
				"from":      "0xf3fcd2c1a78f5eee",
				"fromName":  "user2",
				"toAddress": "0x179b6b1cb6755e31",
				"amount":    "5.00000000",
				"name":      "user1",
				"tag":       "test",
				"message":   "This is a test",
			})).
			AssertEmitEvent(gwtf.NewTestEvent("A.f8d6e0586b0a20c7.FUSD.TokensDeposited", map[string]interface{}{
				"amount": "5.00000000",
				"to":     "0x179b6b1cb6755e31",
			})).
			AssertEmitEvent(gwtf.NewTestEvent("A.f8d6e0586b0a20c7.FUSD.TokensWithdrawn", map[string]interface{}{
				"amount": "5.00000000",
				"from":   "0xf3fcd2c1a78f5eee",
			}))

	})

	t.Run("Should be able to send flow to another name with tag and message", func(t *testing.T) {

		gt := NewGWTFTest(t).
			setupFIND().
			createUser("100.0", "user1").
			createUser("100.0", "user2").
			registerUser("user1").
			registerUser("user2")

		gt.GWTF.TransactionFromFile("sendFlowWithTagAndMessage").
			SignProposeAndPayAs("user2").
			StringArgument("user1").
			UFix64Argument("5.0").
			StringArgument("This is a test").
			StringArgument("test").
			Test(t).AssertSuccess().
			AssertPartialEvent(gwtf.NewTestEvent("A.f8d6e0586b0a20c7.FIND.FungibleTokenSent", map[string]interface{}{
				"from":      "0xf3fcd2c1a78f5eee",
				"fromName":  "user2",
				"toAddress": "0x179b6b1cb6755e31",
				"amount":    "5.00000000",
				"name":      "user1",
				"tag":       "test",
				"message":   "This is a test",
			})).
			AssertEmitEvent(gwtf.NewTestEvent("A.0ae53cb6e3f42a79.FlowToken.TokensDeposited", map[string]interface{}{
				"amount": "5.00000000",
				"to":     "0x179b6b1cb6755e31",
			})).
			AssertEmitEvent(gwtf.NewTestEvent("A.0ae53cb6e3f42a79.FlowToken.TokensWithdrawn", map[string]interface{}{
				"amount": "5.00000000",
				"from":   "0xf3fcd2c1a78f5eee",
			}))

	})

	t.Run("Should be able to register related account and remove it", func(t *testing.T) {

		gt := NewGWTFTest(t).
			setupFIND().
			createUser("100.0", "user1").
			registerUser("user1")

		gt.GWTF.TransactionFromFile("setRelatedAccount").
			SignProposeAndPayAs("user1").
			StringArgument("dapper").
			AccountArgument("user2").
			Test(t).AssertSuccess().
			AssertPartialEvent(gwtf.NewTestEvent("A.f8d6e0586b0a20c7.RelatedAccounts.RelatedFlowAccountAdded", map[string]interface{}{
				"name":    "dapper",
				"address": "0x179b6b1cb6755e31",
				"related": "0xf3fcd2c1a78f5eee",
			}))

		value := gt.GWTF.ScriptFromFile("address_status").AccountArgument("user1").RunReturnsJsonString()
		assert.Contains(t, value, `"dapper": "0xf3fcd2c1a78f5eee"`)

		gt.GWTF.TransactionFromFile("removeRelatedAccount").
			SignProposeAndPayAs("user1").
			StringArgument("dapper").
			Test(t).AssertSuccess().
			AssertPartialEvent(gwtf.NewTestEvent("A.f8d6e0586b0a20c7.RelatedAccounts.RelatedFlowAccountRemoved", map[string]interface{}{
				"name":    "dapper",
				"address": "0x179b6b1cb6755e31",
				"related": "0xf3fcd2c1a78f5eee",
			}))

		value = gt.GWTF.ScriptFromFile("address_status").AccountArgument("user1").RunReturnsJsonString()
		assert.NotContains(t, value, `"dapper": "0xf3fcd2c1a78f5eee"`)

	})

	t.Run("Should be able to set private mode", func(t *testing.T) {

		gt := NewGWTFTest(t).
			setupFIND().
			createUser("100.0", "user1").
			registerUser("user1")

		gt.GWTF.TransactionFromFile("setPrivateMode").
			SignProposeAndPayAs("user1").
			BooleanArgument(true).
			Test(t).AssertSuccess()

		value := gt.GWTF.ScriptFromFile("address_status").AccountArgument("user1").RunReturnsJsonString()
		assert.Contains(t, value, `"privateMode": "true"`)

		gt.GWTF.TransactionFromFile("setPrivateMode").
			SignProposeAndPayAs("user1").
			BooleanArgument(false).
			Test(t).AssertSuccess()

		value = gt.GWTF.ScriptFromFile("address_status").AccountArgument("user1").RunReturnsJsonString()
		assert.Contains(t, value, `"privateMode": "false"`)

	})

}
