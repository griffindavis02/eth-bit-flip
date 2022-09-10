# EthFlip

This branch solely serves to outline the purpose of all of the other
branches in this repo. It lives here so as not to create bloat.

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
stepped up to higher versions, and eventually v1.0.0 when it has all the features
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

If you followed the steps for branch `v0.1.0`:

```shell
replace github.com/griffindavis02/eth-bit-flip v0.1.0 => ..\eth-bit-flip
```

Also update the `go.mod` file of the eth-bit-flip source code to point here:

```shell
cd ..\eth-bit-flip
[open go.mod]
replace github.com/ethereum/go-ethereum => ..\geth-source
```

You will also need to update the line that requires eth-bit-flip to the latest
version when updated.

It would be a good idea to maintain a record of what changes are being made
now that the project is in a maintainable state. This will allow you to keep a
record of where flips are occurring and will be invaluable if a patch script is
written.

### blockchaintest

This branch has a fully configured blockchain with local nodes. If you fetch it
and checkout in one folder with the `geth-source` branch checked out in another,
you can use the altered source code to launch this blockchain and run tests.

The commands for operating an Ethereum blockchain will not be included here.

```bash
git clone git@github.com:griffindavis02/eth-bit-flip.git blockchain
cd blockchain
git checkout blockchaintest
```

## Building the Source Code

Re-enter the directory that has `geth-source` checked out. Run `make all` to
build the executables. I trust you to research setting up `make` yourself if
your system does not recognize it (I recommend looking into `chocolatey`).

```shell
make all
```

If everything is configured correctly, it should build geth, puppeth, flipconfig,
everything you need for your tests. Try running `./build/bin/flipconfig.exe` to
be sure.

You'll likely want to add `.../capstone-project/geth-source/build/bin` or
whatever directory tree you chose to `PATH` so you can just run `geth` and
`flipconfig` as regular commands.

i.e. `flipconfig` will launch the tool instead of needing
`.../capstone-project/geth-source/build/bin/flipconfig.exe`.

## Flip Commands Available to You

To configure your test environment, run `flipconfig` and create a new configuration.
This configuration is saved in `~/.flipconfig` (`C:\Users\<you>\.flipconfig` on
Windows). 

Make sure to enter `y` for using an API if intending to use the frontend
(documentation in its repository). If developing locally, use whatever port
it is hosted on for the hostname (likely http://localhost:5000) or the URL
to the frontend if using the "production" frontend.

To control whether bits are getting flipped in the blockchain, the `flip`
command is added to `geth` with possible arguments:

```shell
geth flip --flipstart # start bit flipping
geth flip --flipstop # stop bit flipping
geth flip --fliprestart # reset the time/variable/bit counter
# based on test type from flipconfig and start bit flipping
```

## Running the Blockchain

You will likely want to execute `geth flip --flipstop` so as not to run into
bitflips while unlocking accounts on each node.

Each node directory has its own local `start.sh` script for Unix based OS and
`start.bat` for Windows. You should, however, be able to run the `start_nodes.sh`
script in the parent directory to open each node in its own terminal. This
script assumes you have `bash` installed, which you should if you have git.

Then, in your original terminal, just execute `geth flip --flipstart` to begin
bit flips tests. i.e.:

```shell
geth flip --flipstop
./start_nodes.sh
geth flip --flipstart
```
If you make changes to any scripts in this directory, be wary of your adds and
commits back to the remote branch. It will congest your version history if you
add and commit the many many files that get changed and deleted while the
blockchain is running. This branch is mainly for quick testing on your local
machine.
