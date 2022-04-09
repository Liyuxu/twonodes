#!/bin/bash

set -ex

# Bring the test network down
pushd ../test-network
./network.sh down
popd

# clean out any old identites in the wallets
rm -rf client0/sdkgo/wallet/*
rm -rf client0/sdkgo/keystore
rm -rf client0/ledger/*
rm -rf client0/results/*

rm -rf client1/sdkgo/wallet/*
rm -rf client1/sdkgo/keystore
rm -rf client1/ledger/*
rm -rf client1/results/*
