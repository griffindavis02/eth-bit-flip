## Transaction Sending

```go
web3.eth.sendTransaction({from: "0x85c58cf29d98cf731520224da7954d527cb78cf0",to: "0xee072662b53dc708e6e4d5f2e47e6cb407a4035e",value: 1000000000000000000})
```

[//]: # "drive home on ability to configure error laden
experiment is successful due to hash zeroing
	this shows that any transactions that are zeroed will be overwritten by future zeroed hashes which removes data from the blockchain"

## Zero Hash Override

<b>A zero-hash transaction was commited to block 130.</b>

> INFO [11-29|17:21:40.408] Submitted transaction hash=0x0000000000000000000000000000000000000000000000000000000000000000 from=0x85C58Cf29d98cF731520224dA7954d527Cb78cf0 nonce=1 recipient=0xEe072662B53dC708E6E4D5f2e47e6CB407A4035e
> value=1,000,000,000

> INFO [11-29|17:21:41.041] Commit new mining work number=130 sealhash=1537b7..6dec8a uncles=0 txs=1 gas=21000 fees=2.1e-05 elapsed=2.008ms

<b>But `eth.getTransaction("0x0000...0000")` returns block 146</b>

```javascript
{
  blockHash: "0x7564d196b575428109b0aec06142543c2dd4878368e78eaa34278381121f0fba",
  blockNumber: 146,
  from: "0x85c58cf29d98cf731520224da7954d527cb78cf0",
  gas: 21000,
  gasPrice: 1000000000,
  hash: "0x0000000000000000000000000000000000000000000000000000000000000000",
  ...
}
```

## Failed Transaction Search

> ERROR[11-29|17:37:14.070] Transaction not found number=139 hash=398f84..e6921b txhash=966e25..4161e8
> null

```javascript
{
  blockHash: "0x398f841efe53fba2e2201eb7cb9fb9ee5f9ffcd611fc94caada13a75d6e6921b",
  blockNumber: 139,
  from: "0x85c58cf29d98cf731520224da7954d527cb78cf0",
  gas: 21000,
  gasPrice: 1000000000,
  hash: "0x966e258937564886a3d82258a7242ec7058b36fdbf3d5018a490c060594161e8",
  input: "0x",
  nonce: 2,
  r: "0xc6240809e2b40bb14467a06c39e25761b2fe00fddff72c8c24a8710c298f67f8",
  s: "0x114b54a08b93c37d9026ba2125ae47d83e823dfac853416badc58a59a9104473",
  to: "0xee072662b53dc708e6e4d5f2e47e6cb407a4035e",
  transactionIndex: 0,
  type: "0x0",
  v: "0x10f19",
  value: 1000000000
}
```

## Failure to Send a Transaction

> WARN [12-02|09:24:45.379] Served eth_sendTransaction reqid=31 t=33.8966ms err="already known"
> Error: already known

        at web3.js:6357:37(47)
        at send (web3.js:5091:62(35))
        at <eval>:1:25(11)

## Asynchronous Failure

> ERROR[12-02|09:27:21.038] No transaction found to be deleted hash=000000..ad59c0
