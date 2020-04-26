pragma solidity ^0.5.0;
contract with_if_continue {
    constructor() public {
        if (1 > 2) {
            create(5);
            return;
        }
        update(6);
    }

    function create(uint newBalance) public {
    }

    function update(uint newBalance) public {
    }
}