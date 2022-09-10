#!/bin/bash

geth --networkid 34683 --datadir "./data" --bootnodes enode://b25d7df2f3ec0d17c8cd4a7b3a18511463b37812793d3b093da45b10d39f9a193d6b3933d8dd991fb9194aab5cf40bb2c2de4110c945bf59dfe51ee6f2864382@127.0.0.1:0?discport=30301 --port 30303 --unlock "0x85C58Cf29d98cF731520224dA7954d527Cb78cf0" --ipcdisable --password password.txt --mine console
