// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IERC20 {
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
}

contract TipJar {
    event Tipped(
        uint256 indexed artifactId,
        address indexed from,
        address indexed to,
        address token,
        uint256 amount,
        string memo
    );

    error InvalidArtifactId();
    error InvalidRecipient();
    error InvalidAmount();
    error TransferFailed();

    function tipNative(uint256 artifactId, address to, string calldata memo) external payable {
        if (artifactId == 0) revert InvalidArtifactId();
        if (to == address(0)) revert InvalidRecipient();
        if (msg.value == 0) revert InvalidAmount();

        (bool ok, ) = to.call{value: msg.value}("");
        if (!ok) revert TransferFailed();

        emit Tipped(artifactId, msg.sender, to, address(0), msg.value, memo);
    }

    function tipToken(
        uint256 artifactId,
        address token,
        address to,
        uint256 amount,
        string calldata memo
    ) external {
        if (artifactId == 0) revert InvalidArtifactId();
        if (to == address(0)) revert InvalidRecipient();
        if (amount == 0) revert InvalidAmount();

        bool ok = IERC20(token).transferFrom(msg.sender, to, amount);
        if (!ok) revert TransferFailed();

        emit Tipped(artifactId, msg.sender, to, token, amount, memo);
    }
}
