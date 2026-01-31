// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract ArtifactRegistry {
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
        bytes32 contentHash
    );

    function publish(string calldata uri, bytes32 contentHash) external returns (uint256) {
        uint256 artifactId = ++totalArtifacts;
        artifacts[artifactId] = Artifact({
            creator: msg.sender,
            uri: uri,
            contentHash: contentHash,
            createdAt: block.timestamp
        });
        emit ArtifactPublished(artifactId, msg.sender, uri, contentHash);
        return artifactId;
    }
}
