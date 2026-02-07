const { expect } = require("chai");
const { ethers } = require("hardhat");
const { anyValue } = require("@nomicfoundation/hardhat-chai-matchers/withArgs");

describe("MVP provenance + attribution", function () {
  it("emits ArtifactPublished with provenance fields", async function () {
    const [creator] = await ethers.getSigners();
    const registry = await ethers.deployContract("ArtifactRegistry");
    await registry.waitForDeployment();

    const uri = "ipfs://bafybeigdyrztv2xj6v4zv6m4n4yqvfp6mk7if4f4w7h2t5g2q4d6p3f2ey";
    const contentHash = ethers.keccak256(ethers.toUtf8Bytes("artifact-v1"));

    await expect(registry.publish(uri, contentHash))
      .to.emit(registry, "ArtifactPublished")
      .withArgs(1n, creator.address, uri, contentHash, anyValue);

    const artifact = await registry.artifacts(1n);
    expect(artifact.creator).to.equal(creator.address);
    expect(artifact.uri).to.equal(uri);
    expect(artifact.contentHash).to.equal(contentHash);
    expect(artifact.createdAt).to.be.gt(0n);
  });

  it("emits Tipped for native tips with artifact attribution", async function () {
    const [tipper, recipient] = await ethers.getSigners();
    const tipJar = await ethers.deployContract("TipJar");
    await tipJar.waitForDeployment();

    const amount = ethers.parseEther("1");
    const memo = "thanks for your artifact";

    await expect(
      tipJar.connect(tipper).tipNative(1n, recipient.address, memo, { value: amount })
    )
      .to.emit(tipJar, "Tipped")
      .withArgs(1n, tipper.address, recipient.address, ethers.ZeroAddress, amount, memo);
  });

  it("emits Tipped for ERC20 tips with artifact attribution", async function () {
    const [tipper, recipient] = await ethers.getSigners();
    const tipJar = await ethers.deployContract("TipJar");
    const token = await ethers.deployContract("MockERC20", [ethers.parseEther("100")]);
    await tipJar.waitForDeployment();
    await token.waitForDeployment();

    const amount = ethers.parseEther("5");
    const tokenAddress = await token.getAddress();
    const memo = "token tip";

    await token.approve(await tipJar.getAddress(), amount);

    await expect(tipJar.tipToken(1n, tokenAddress, recipient.address, amount, memo))
      .to.emit(tipJar, "Tipped")
      .withArgs(1n, tipper.address, recipient.address, tokenAddress, amount, memo);

    expect(await token.balanceOf(recipient.address)).to.equal(amount);
  });
});
