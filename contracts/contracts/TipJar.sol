// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IERC20 {
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
}

contract TipJar {
    event Tipped(
        address indexed from,
        address indexed to,
        address indexed token,
        uint256 amount,
        string memo
    );

    error InvalidRecipient();
    error InvalidAmount();
    error TransferFailed();

    function tipNative(address to, string calldata memo) external payable {
        if (to == address(0)) revert InvalidRecipient();
        if (msg.value == 0) revert InvalidAmount();

        (bool ok, ) = to.call{value: msg.value}("");
        if (!ok) revert TransferFailed();

        emit Tipped(msg.sender, to, address(0), msg.value, memo);
    }

    function tipToken(address token, address to, uint256 amount, string calldata memo) external {
        if (to == address(0)) revert InvalidRecipient();
        if (amount == 0) revert InvalidAmount();

        bool ok = IERC20(token).transferFrom(msg.sender, to, amount);
        if (!ok) revert TransferFailed();

        emit Tipped(msg.sender, to, token, amount, memo);
    }
}
