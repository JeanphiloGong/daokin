// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract ArtifactRegistry {
    error InvalidURI();
    error InvalidContentHash();

    struct Artifact {
        address creator;
        string uri;
        bytes32 contentHash;
        uint256 createdAt;
    }

    uint256 public totalArtifacts;
    mapping(uint256 => Artifact) public artifacts;

    event ArtifactPublished(
        uint256 indexed artifactId,
        address indexed creator,
        string uri,
        bytes32 contentHash,
        uint256 createdAt
    );

    function publish(string calldata uri, bytes32 contentHash) external returns (uint256) {
        if (bytes(uri).length == 0) revert InvalidURI();
        if (contentHash == bytes32(0)) revert InvalidContentHash();

        uint256 artifactId = ++totalArtifacts;
        uint256 createdAt = block.timestamp;
        artifacts[artifactId] = Artifact({
            creator: msg.sender,
            uri: uri,
            contentHash: contentHash,
            createdAt: createdAt
        });
        emit ArtifactPublished(artifactId, msg.sender, uri, contentHash, createdAt);
        return artifactId;
    }

    function exists(uint256 artifactId) external view returns (bool) {
        return artifactId > 0 && artifactId <= totalArtifacts;
    }
}
