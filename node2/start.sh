#!/bin/bash

geth --networkid 34683 --datadir "./data" --bootnodes enode://b25d7df2f3ec0d17c8cd4a7b3a18511463b37812793d3b093da45b10d39f9a193d6b3933d8dd991fb9194aab5cf40bb2c2de4110c945bf59dfe51ee6f2864382@127.0.0.1:0?discport=30301 --port 30305 --ipcdisable --unlock "0xEe072662B53dC708E6E4D5f2e47e6CB407A4035e" --password password.txt console
