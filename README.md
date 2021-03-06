# Keep Random Beacon Dappnode Package
This dappnode package will allow you to participate in the keep network as a random beacon operator.

see:
- https://keep.network/
- https://discordapp.com/invite/wYezN7v (see 🥩Staking channels)
- https://keep-network.gitbook.io/staking-documentation/

# Requirements
This package requires a public, static ip address.

Additionally, at the start of the stakedrop, 90,000 KEEP will be required to operate a random beacon.
This amount will decrease over time until it reaches 10,000 KEEP ([source](https://keep-network.gitbook.io/staking-documentation/help/faq#is-there-a-minimum-staking-amount)).

# Install
Access this link using your dappnode wifi/vpn:
http://my.dappnode/#/installer/%2Fipfs%2FQmd2CjLMNhg9QZWr8o3fEQ3CSatuvdmxUzJkYGeFCSBPJR

current ipfs hash `Qmd2CjLMNhg9QZWr8o3fEQ3CSatuvdmxUzJkYGeFCSBPJR`

## Quick Start
1. Set `KEEP_ETHEREUM_PASSWORD` and `ANNOUNCED_ADDRESSES` environment variables in the Config section of the package.


ANNOUNCED_ADDRESSES should use the [libp2p multiaddr format](https://docs.libp2p.io/concepts/addressing/)
```
/ip4/255.255.255.255/tcp/3919
```

(Always use port 3919)


If announcing multiple addresses, the address should be listed as a comma-delimited string with no spaces or quotation marks, eg:
```
/ip4/255.255.255.255/tcp/3919,/dns4/my.dappnode.tech/tcp/3919
```

2. Copy operator address from package logs
3. GOTO https://dashboard.keep.network/
5. Delegate your keep tokens to the operator address (see [delegation](#delegation))
6. Authorize the Random Beacon contract
7. Send some eth to the operator address
  - ~0.5 eth is fine to start with, but be sure to monitor the balance!

### Delegation
You will need to set up the delegation to this operator account in the keep dashboard at https://dashboard.keep.network

It is recommended that you only use the account generated by this DNP package as the Operator Address.
Use a more secure account for the Authorizer and Beneficiary addresses.

Remember to Authorize the Keep Random Beacon Operator Contract after delegating!


# Risks
The primary risks in running a random beacon are downtime and loss of persistent data.

Educate yourself on the risks before attempting to operate a random beacon client!

resources:
- https://hackmd.io/@protocollayer/BkUBl7zIw#Random-beacon
- https://keep-network.gitbook.io/staking-documentation/about-staking/slashing

## Backups

### Extracting your operator account
The operator account is generated automatically for you when this package is initialized.
The operator account (ie private key) is stored in the data volume for this package,
so if you delete the dnp completely, including data, then you will lose your operator account.

It is a good idea to back up the account file / json someplace safe like a password manager.
The (encrypted) account is written to `/mnt/keystore/keep_wallet.json`.

To save this file, simply browse to the `File Manager` section of the DNP package and enter
this path into the `DOWNLOAD FILE FROM PACKAGE` input.
