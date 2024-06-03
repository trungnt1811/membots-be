import {
  assert,
  describe,
  test,
  clearStore,
  beforeAll,
  afterAll
} from "matchstick-as/assembly/index"
import { Address, BigInt } from "@graphprotocol/graph-ts"
import { CollectFees } from "../src/types/schema"
import { CollectFees as CollectFeesEvent } from "../src/types/TruglyMemeception/TruglyMemeception"
import { handleCollectFees } from "../src/mapping/trugly-memeception"
import { createCollectFeesEvent } from "./trugly-memeception-utils"

// Tests structure (matchstick-as >=0.5.0)
// https://thegraph.com/docs/en/developer/matchstick/#tests-structure-0-5-0

describe("Describe entity assertions", () => {
  beforeAll(() => {
    let memeToken = Address.fromString(
      "0x0000000000000000000000000000000000000001"
    )
    let recipient = Address.fromString(
      "0x0000000000000000000000000000000000000001"
    )
    let amount0 = BigInt.fromI32(234)
    let amount1 = BigInt.fromI32(234)
    let fee0 = BigInt.fromI32(234)
    let fee1 = BigInt.fromI32(234)
    let newCollectFeesEvent = createCollectFeesEvent(
      memeToken,
      recipient,
      amount0,
      amount1,
      fee0,
      fee1
    )
    handleCollectFees(newCollectFeesEvent)
  })

  afterAll(() => {
    clearStore()
  })

  // For more test scenarios, see:
  // https://thegraph.com/docs/en/developer/matchstick/#write-a-unit-test

  test("CollectFees created and stored", () => {
    assert.entityCount("CollectFees", 1)

    // 0xa16081f360e3847006db660bae1c6d1b2e17ec2a is the default address used in newMockEvent() function
    assert.fieldEquals(
      "CollectFees",
      "0xa16081f360e3847006db660bae1c6d1b2e17ec2a-1",
      "memeToken",
      "0x0000000000000000000000000000000000000001"
    )
    assert.fieldEquals(
      "CollectFees",
      "0xa16081f360e3847006db660bae1c6d1b2e17ec2a-1",
      "recipient",
      "0x0000000000000000000000000000000000000001"
    )
    assert.fieldEquals(
      "CollectFees",
      "0xa16081f360e3847006db660bae1c6d1b2e17ec2a-1",
      "amount0",
      "234"
    )
    assert.fieldEquals(
      "CollectFees",
      "0xa16081f360e3847006db660bae1c6d1b2e17ec2a-1",
      "amount1",
      "234"
    )
    assert.fieldEquals(
      "CollectFees",
      "0xa16081f360e3847006db660bae1c6d1b2e17ec2a-1",
      "fee0",
      "234"
    )
    assert.fieldEquals(
      "CollectFees",
      "0xa16081f360e3847006db660bae1c6d1b2e17ec2a-1",
      "fee1",
      "234"
    )

    // More assert options:
    // https://thegraph.com/docs/en/developer/matchstick/#asserts
  })
})
