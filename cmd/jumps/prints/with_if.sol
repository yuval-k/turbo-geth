pragma solidity ^0.5.0;
contract with_if {
    constructor() public {
        if (1 > 2) {
            create(5);
        }
    }

    function create(uint newBalance) public {
    }
}