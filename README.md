# Keep Random Beacon Dappnode Package
This dappnode package will allow you to participate in the keep random beacon.

see:
- https://keep.network/
- https://discordapp.com/invite/wYezN7v (see 🥩Staking channels)

## Requirements
This package requires a public, static ip address.

## Risks
see: https://hackmd.io/@protocollayer/BkUBl7zIw#Random-beacon

## Install
Access this link using your dappnode wifi/vpn:
http://my.dappnode/#/installer/%2Fipfs%2FQmUwN3obuWiLZ1WbmE7PwuPg6ejpXXAGH7uM2e9c23uPLr

current ipfs hash `QmUwN3obuWiLZ1WbmE7PwuPg6ejpXXAGH7uM2e9c23uPLr`

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
5. Delegate your keep tokens to the operator address
6. Authorize the Random Beacon contract
7. Send some eth to the operator address
  - ~0.5 eth is fine to start with, but be sure to monitor the balance!

## Extracting your operator account
The operator account is generated automatically for you when this package is initialized.
The operator account (ie private key) is stored in the data volume for this package,
so if you delete the dnp completely, including data, then you will lose your operator account.

It is a good idea to back up the account file / json someplace safe like a password manager.
The (encrypted) account is written to `/mnt/keystore/keep_wallet.json`.

To save this file, simply browse to the `File Manager` section of the DNP package and enter
this path into the `DOWNLOAD FILE FROM PACKAGE` input.

## Delegation
You will need to set up the delegation to this operator account in the keep dashboard at https://dashboard.keep.network

It is recomended that you only use the account generated by this DNP package as the Operator Address.
Use a more secure account for the Authorizer and Beneficiary addresses.

Remember to Authorize the Keep Random Beacon Operator Contract after delegating!
