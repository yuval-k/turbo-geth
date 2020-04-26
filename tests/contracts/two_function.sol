pragma solidity ^0.5.0;
contract two_function {
    mapping(address => uint) public balances;

    constructor() public {
        balances[msg.sender] = 100;
    }

    function create(uint newBalance) public {
        balances[msg.sender] = newBalance;
    }

    function update(uint newBalance) public {
        balances[msg.sender] = newBalance;
    }
}