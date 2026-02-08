const { expect } = require("chai");
const { ethers } = require("hardhat");
const { anyValue } = require("@nomicfoundation/hardhat-chai-matchers/withArgs");

describe("MVP provenance + attribution", function () {
  async function deployRegistryWithArtifact() {
    const registry = await ethers.deployContract("ArtifactRegistry");
    await registry.waitForDeployment();

    const uri = "ipfs://bafybeigdyrztv2xj6v4zv6m4n4yqvfp6mk7if4f4w7h2t5g2q4d6p3f2ey";
    const contentHash = ethers.keccak256(ethers.toUtf8Bytes("artifact-v1"));
    await registry.publish(uri, contentHash);

    return { registry, uri, contentHash };
  }

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
    expect(await registry.exists(1n)).to.equal(true);
    expect(await registry.exists(2n)).to.equal(false);
  });

  it("rejects invalid artifact input", async function () {
    const registry = await ethers.deployContract("ArtifactRegistry");
    await registry.waitForDeployment();

    const uri = "ipfs://bafybeigdyrztv2xj6v4zv6m4n4yqvfp6mk7if4f4w7h2t5g2q4d6p3f2ey";
    const contentHash = ethers.keccak256(ethers.toUtf8Bytes("artifact-v1"));

    await expect(registry.publish("", contentHash)).to.be.revertedWithCustomError(
      registry,
      "InvalidURI"
    );
    await expect(registry.publish(uri, ethers.ZeroHash)).to.be.revertedWithCustomError(
      registry,
      "InvalidContentHash"
    );
  });

  it("rejects invalid registry references in TipJar", async function () {
    const [owner] = await ethers.getSigners();
    const tipJarFactory = await ethers.getContractFactory("TipJar");
    const { registry } = await deployRegistryWithArtifact();
    const tipJar = await ethers.deployContract("TipJar", [await registry.getAddress()]);
    await tipJar.waitForDeployment();

    await expect(tipJarFactory.deploy(ethers.ZeroAddress)).to.be.revertedWithCustomError(
      tipJar,
      "InvalidRegistry"
    );
    await expect(tipJarFactory.deploy(owner.address)).to.be.revertedWithCustomError(
      tipJar,
      "InvalidRegistry"
    );
  });

  it("emits Tipped for native tips with artifact attribution", async function () {
    const [, recipient] = await ethers.getSigners();
    const { registry } = await deployRegistryWithArtifact();
    const tipJar = await ethers.deployContract("TipJar", [await registry.getAddress()]);
    await tipJar.waitForDeployment();

    const amount = ethers.parseEther("1");
    const memo = "thanks for your artifact";

    await expect(tipJar.tipNative(1n, recipient.address, memo, { value: amount }))
      .to.emit(tipJar, "Tipped")
      .withArgs(1n, anyValue, recipient.address, ethers.ZeroAddress, amount, memo);
  });

  it("reverts native tips for unknown artifact", async function () {
    const [, recipient] = await ethers.getSigners();
    const { registry } = await deployRegistryWithArtifact();
    const tipJar = await ethers.deployContract("TipJar", [await registry.getAddress()]);
    await tipJar.waitForDeployment();

    await expect(
      tipJar.tipNative(999n, recipient.address, "unknown", { value: ethers.parseEther("0.1") })
    ).to.be.revertedWithCustomError(tipJar, "ArtifactNotFound");
  });

  it("emits Tipped for ERC20 tips with artifact attribution", async function () {
    const [tipper, recipient] = await ethers.getSigners();
    const { registry } = await deployRegistryWithArtifact();
    const tipJar = await ethers.deployContract("TipJar", [await registry.getAddress()]);
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

  it("rejects invalid token addresses", async function () {
    const [tipper, recipient] = await ethers.getSigners();
    const { registry } = await deployRegistryWithArtifact();
    const tipJar = await ethers.deployContract("TipJar", [await registry.getAddress()]);
    await tipJar.waitForDeployment();

    await expect(
      tipJar.tipToken(1n, ethers.ZeroAddress, recipient.address, 1n, "invalid")
    ).to.be.revertedWithCustomError(tipJar, "InvalidToken");

    await expect(
      tipJar.tipToken(1n, tipper.address, recipient.address, 1n, "invalid")
    ).to.be.revertedWithCustomError(tipJar, "InvalidToken");
  });

  it("blocks reentrancy attempts during native tips", async function () {
    const { registry } = await deployRegistryWithArtifact();
    const tipJar = await ethers.deployContract("TipJar", [await registry.getAddress()]);
    await tipJar.waitForDeployment();

    const attacker = await ethers.deployContract("ReentrantRecipient", [await tipJar.getAddress(), 1n]);
    await attacker.waitForDeployment();

    await attacker.attack({ value: ethers.parseEther("0.1") });

    expect(await attacker.reentryBlocked()).to.equal(true);
  });
});
