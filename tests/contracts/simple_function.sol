pragma solidity ^0.5.0;
contract simple_function {
    mapping(address => uint) public balances;

    constructor() public {
        balances[msg.sender] = 100;
    }

    function create(uint newBalance) public {
        balances[msg.sender] = newBalance;
    }
}