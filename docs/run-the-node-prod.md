# Getting started with ZKEVM: Spinning up an RPC Node 

Modules that will be spun up:

- Synchronizer
- RPC
- Databases
- Prover/Merkletree

First, build the corresponding image for ZKEVM Node:

`make build-docker`

Check there is no previous ZKEVM Node instance running:

`docker ps`

Then run the Node as RPC:

`make run-rpc` will run an RPC Node for ZKEVM that connects to testnet engine from Polygon.

Essentially this uses the L1/Sequencer/Broadcast from Polygon ZKEVM Testnet while using the Prover/Sync/RPC from your local docker instance.

### Configuration:

By default the L1 network will be `https://geth_goerli.hermez.io` but change as you see fit.

Broadcast, SequencerNode URIs should point to `internal.zkevm-test.net` with the correct ports.

Prover config resides `config/prover.config.local.json` and is already written for use in conjunction with the RPC-mode Node.