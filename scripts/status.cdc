import FiNS from "../contracts/FiNS.cdc"

//Check the status of a fin user
pub fun main(tag: String) : UInt8 {
    let status=FiNS.status(tag)
    return status.rawValue
}
