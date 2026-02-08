// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IERC20 {
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
}

interface IArtifactRegistry {
    function exists(uint256 artifactId) external view returns (bool);
}

contract TipJar {
    uint256 private constant _NOT_ENTERED = 1;
    uint256 private constant _ENTERED = 2;

    IArtifactRegistry public immutable artifactRegistry;
    uint256 private _status;

    event Tipped(
        uint256 indexed artifactId,
        address indexed from,
        address indexed to,
        address token,
        uint256 amount,
        string memo
    );

    error InvalidArtifactId();
    error ArtifactNotFound();
    error InvalidRegistry();
    error InvalidRecipient();
    error InvalidAmount();
    error InvalidToken();
    error ReentrancyBlocked();
    error TransferFailed();

    constructor(address artifactRegistryAddress) {
        if (artifactRegistryAddress == address(0) || artifactRegistryAddress.code.length == 0) {
            revert InvalidRegistry();
        }
        artifactRegistry = IArtifactRegistry(artifactRegistryAddress);
        _status = _NOT_ENTERED;
    }

    modifier nonReentrant() {
        if (_status == _ENTERED) revert ReentrancyBlocked();
        _status = _ENTERED;
        _;
        _status = _NOT_ENTERED;
    }

    function tipNative(uint256 artifactId, address to, string calldata memo) external payable nonReentrant {
        _validateArtifact(artifactId);
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
    ) external nonReentrant {
        _validateArtifact(artifactId);
        if (to == address(0)) revert InvalidRecipient();
        if (amount == 0) revert InvalidAmount();
        if (token == address(0) || token.code.length == 0) revert InvalidToken();

        bool ok = IERC20(token).transferFrom(msg.sender, to, amount);
        if (!ok) revert TransferFailed();

        emit Tipped(artifactId, msg.sender, to, token, amount, memo);
    }

    function _validateArtifact(uint256 artifactId) private view {
        if (artifactId == 0) revert InvalidArtifactId();
        if (!artifactRegistry.exists(artifactId)) revert ArtifactNotFound();
    }
}
