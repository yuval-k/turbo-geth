pragma solidity ^0.5.0;
contract two_function_call_with_return {
    mapping(address => uint) public balances;

    constructor() public {
        balances[msg.sender] = 100;
    }

    function create(uint newBalance) public {
        balances[msg.sender] = update(newBalance);
    }

    function update(uint newBalance) public returns(uint) {
        return newBalance;
    }
}