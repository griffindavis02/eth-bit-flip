# GETH Bit Flipping

The eth-bit-flip library allows users to simulate the effects of soft errors
in the Ethereum Virtual Machine with the use of the [go-ethereum library](https://github.com/ethereum/go-ethereum).

## Re-building the source

It would be great if we wrote a script that would take the few alterations we
make to the go-ethereum source code (i.e. where we flip bits) and it can insert
them in their proper places given the path to the root of the source code.

This way, users can most likely always use the latest version of the source code
instead of being stuck with ours. The method for them to insert our changes would
be something like:

Clone the [go-ethereum source code](https://github.com/ethereum/go-ethereum) to
your local machine, then run the following command with the <u>absolute</u> path
to the source code on your system <b>in single quotes</b>:

```shell
curl https://raw.githubusercontent.com/griffindavis02/eth-bit-flip/utils-test/patch.py -s \
python - '<path-to-the-go-ethereum-source-code>'
```

This should add the required functions to your source code and include the
configuration CLI for your use. You can then run the `make` commands seen on the
go-ethereum GitHub page to build the tools.

[//]: # "Add a section for example injections into the go-ethereum source code so people can implement their own calls to it"

## License

The eth-bit-flip library (all code within this repository) is licensed under
the [GNU General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html),
also included within the `COPYING` file.
