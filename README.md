# EthFlip

This branch solely serves to outline the purpose of all of the other
branches. It lives here so as not to create bloat.

I recommend making a directory to store all of the necessary folders that
interact with one another. i.e.

```shell
mkdir capstone-project
cd capstone-project
```

And continue on with the steps on which branches to include and how from within
that directory.

## Branches

### v0.1.0

This is the format by which [pkg.go.dev](pkg.go.dev) manages their packages.
This branch should operate as the traditional `main` or `master` branch and be
stepped up to higher versions, and eventuall v1.0.0 when it has all the features
we want to begin with.

Before everything else, clone this into its own folder, and check out `v0.1.0`.

```shell
git clone git@github.com:griffindavis02/eth-bit-flip.git eth-bit-flip
cd eth-bit-flip
git checkout v0.1.0
```

### geth-source

This is the altered source code which allows you to use the EthFlip library on
the hashing function of transactions. You will want to pull this into its own
folder to build the `geth` executables for use on a blockchain.

```bash
git clone git@github.com:griffindavis02/eth-bit-flip.git geth-source
cd geth-source
git checkout geth-source
```

Open up `go.mod` and make sure that the 'replace' command for the eth-bit-flip
package points to your local folder for eth-bit-flip. This can likely be avoided
once the main branch is beyond v1.0.0 and go.pkg.dev updates the source code
more frequently.

If you folloewd the steps for branch `v0.1.0`:

```shell
replace github.com/griffindavis02/eth-bit-flip v0.1.0 => ..\eth-bit-flip
```

Also update the `go.mod` file of the eth-bit-flip source code to point here:

```shell
cd ..\geth-source
[open go.mod]
replace github.com/ethereum/go-ethereum => ..\geth-source
```

You will also need to update the line that requires eth-bit-flip to the latest
version when updated.

It would be a good idea to maintain a record of what changes are being made
now that the project is a maintainable state. This will allow you to keep a
record of where flips are occurring and will be invaluable if a patch script is
written.

### blockchaintest

This branch has a fully configured blockchain with local nodes. If you fetch it
and checkout in one folder with the `geth-source` branch checked out in another,
you can use the altered source code to launch this blockchain and run tests.

```bash
git clone git@github.com:griffindavis02/eth-bit-flip.git blockchain
cd blockchain
git checkout blockchaintest
```

## Building the Source Code

Re-enter the directory that has `geth-source` checked out. Run `make all` to
build the executables. I trust you to research installing make yourself if you
don't have it installed.

```shell
make all
```

If everything is configured correctly, it should build geth, puppeth, flipconfig,
everything you need for your tests. Try running `./build/bin/flipconfig.exe` to
be sure.

You'll likely want to add `.../capstone-project/geth-source/build/bin` or
whatever directory tree you chose to `PATH` so you can just run `geth` and
`flipconfig` as regular tests.

## Running the Blockchain

:)
