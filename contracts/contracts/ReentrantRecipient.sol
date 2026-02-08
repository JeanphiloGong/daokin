// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface ITipJarForReentry {
    function tipNative(uint256 artifactId, address to, string calldata memo) external payable;
}

contract ReentrantRecipient {
    ITipJarForReentry public immutable tipJar;
    uint256 public immutable artifactId;
    bool public reentryBlocked;
    bool private attempted;

    constructor(address tipJarAddress, uint256 artifactId_) {
        tipJar = ITipJarForReentry(tipJarAddress);
        artifactId = artifactId_;
    }

    receive() external payable {
        if (attempted) {
            return;
        }

        attempted = true;
        // Best-effort reentry attempt: guarded TipJar should reject this second call.
        (bool ok, ) = address(tipJar).call{value: 1}(
            abi.encodeWithSignature("tipNative(uint256,address,string)", artifactId, address(this), "reenter")
        );
        if (!ok) {
            reentryBlocked = true;
        }
    }

    function attack() external payable {
        tipJar.tipNative{value: msg.value}(artifactId, address(this), "attack");
    }
}
