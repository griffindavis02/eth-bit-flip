# Geth Source Code Changes

> This document can serve to keep track of active changes in the source code
> for what's being flipped, as well as the third section containing the
> _absolutely necessary_ additions for `flipconfig` and `geth flip` to work.

## List of Changes

Here is a list of all files with bitflips currently in the `geth-source` branch.
Just `Ctrl + F` or `grep` each file for `injection.BitFlip`

-   core/genesis.go
-   core/types/transaction.go
-   core/types/transaction_signing.go
-   miner/stress/1559/main.go

## How to Use the Injection Library

Find something you think bit flipping is worthwhile on? Bitflipping is simple.
All you have to do is get the variable into a primitive data type (int, uint,
string, etc.), run the BitFlip function, and convert that back to the original
data type.

The conversions to a primitive type is unnecessary if the variable is already
primitive, but the conversion back is necessary regardless as the function
returns an interface.

You will also need to make sure to import the injection library at the top of
the file.

Given this line:

```go
return hash
```

You could return the bitflipped hash by converting the hash to a string, and
the bitflipped string back to a hash:

```go
string_hash = hash.Hex()                        // get primitive type
bitflipped = injection.BitFlip(string_hash)     // interface returned
bitflipped_string = bitflipped.(string)         // converted back to string
hash = common.HexToHash(bitflipped_string)      // converted back to hash
return hash

// All in one line
return common.HexToHash(injection.BitFlip(hash.Hex()).(string))
```

For making sense of the bit flip output, you can attach an optional message to
the bitflip:

```go
return common.HexToHash(injection.BitFlip(hash.Hex()).(string), "Flipping hash")
```

It would be a good idea to add common geth types like `common.Hash` to
injection.go to cut out some of this syntactic sugar.

## Necessary Changes

Must make a 'flipconfig' folder in the `./cmd/` folder. This folder contains the
files in the `config` folder of this repo (presently main.go, manage_config.go,
and new_config.go).

```bash
cp -r eth-bit-flip/config geth-source/cmd/flipconfig
```

In `geth-source/utils/flags,go` add the flags from `eth-bit-flip/cmd/flags.go`.

Copy the flip.go file:

```bash
cp  eth-bit-flip/cmd/flip.go geth-source/cmd/geth/flip.go
```

Make sure to change the package from `cmd` to `main` at the top of the file in
the source code.

In `geth-source/cmd/geth/main.go`, add the injection library to the import list
("github.com/griffindavis02/eth-bit-flip/injection")

In the global variables, add the flag array:

```go
var (
    ...
    flipFlags = []cli.Flag {
        utils.FlipStart,
        utils.FlipStop,
        utils.FlipRestart
    }
)
```

In the `init` function, include the flipCommand from flip.go

```go
func init() {
    ...
    app.Commands = []cli.Command {
        ...
        flipCommand
    }
}
```

Include the flags afterwards:

```go
app.Flags = append(app.Flags, flipFlags...)
```
