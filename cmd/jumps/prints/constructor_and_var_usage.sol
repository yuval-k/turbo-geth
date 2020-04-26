pragma solidity ^0.5.0;
contract constructor_and_var_usage {
    mapping(address => uint) public balances;

    constructor() public {
        balances[msg.sender] = 100;
    }
}