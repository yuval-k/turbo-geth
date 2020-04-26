pragma solidity ^0.5.0;
contract with_if_else {
    constructor() public {
        if (1 > 2) {
            create(5);
            return;
        } else {
            update(6);
        }
    }

    function create(uint newBalance) public {
    }

    function update(uint newBalance) public {
    }
}