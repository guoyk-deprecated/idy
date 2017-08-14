# Idy

[![Travis Build Status](https://travis-ci.org/yanke-guo/idy.svg)](https://travis-ci.org/yanke-guo/idy)

Idy is a unique random id generation service, it starts with 0 and iterates over all uint64 numbers.

Idy supports clustering by pre-defined sharding pattern.

Idy has a database format easy to read and does not grow.

# Install / Upgrade

```base
go get -u github.com/yanke-guo/idy/cmd/idyd
```

# Usage

## Start

```
Usage of ./idyd:
-b string
  HTTP service bind address (default "127.0.0.1:8865")
-d string
  location of data directory (default "data")
-p string
  comma seperated pool names, only ^[0-9a-zA-Z._-]+$ is allowed (default "user")
-s string
  sharding partten of this instance, see README.md (default "1:1:2048:4096")
```

For example

```bash
idyd -b 0.0.0.0:8866 -d /data/idyd-data -p users,orders,comments -s 1:1:4096:9999
```

Idy will create database if not exists.

## HTTP API

Get a new id in hex format

```
curl http://127.0.0.1:8865/user/_hex

1e
```

Get a new id in decimal format

```
curl http://127.0.0.1:8865/user/_dec

24
```

# How It Works

Idy can provide multiple pools, each pool is a independent unique random number service.

You can use one pool for user id, another for order id.

For each pool, Idy creates a slice of continious UInt64 numbers, shuffles it with Fisher-Yates algorithm and fetch numbers from it.

If current slice drained, Idy move to next one, until max possible UInt64 number.

Idy database file is basically a json file, Idy will update database file automatically.

```javascript
{
  "version":  1,                // version of database file
  "shard":    "1:3:2048:4096"   // sharding pattern
  "seed":     "1234567890998",  // seed of current slice, used for shuffling, string as int64
  "start":    "0",              // start of current slice (inclusive), string as uint64
  "index":    0,                // lates used index in slice, int
}
```

* `version`

`version` should be and only could be 1.

* `shard`

`shard` is sharding pattern, it describes how slices are created.

For single instance, first two numbers should be `1:1`

For 3 instances cluster, first two numbers should be set to `1:3`, `2:3` and `3:3`.

For last two number, `2048:4096` means each slice have 4096 numbers and 2048 of them will be used.

This will allow Idy skip some numbers, making your id unpredicatable.

* `seed`

`seed` is the random source of shuffling in current slice, regenerated and saved each time Idy moves to a new slice.

* `start`

`start` will be the start of this slice across the cluster.

That means, if `start` is `0`, `node-1` will use `[0,4096)`, `node-2` will use `[4096,8192)`, `node-3` will use `[8192,12288)`.

And when current slice drained, the `start` of next slice will be set to `12288`.

* `index`

the position of last used id in current slice.

Idy uses a local copy of golang `math/rand` package, once `seed` is determined (persisted in database file), a sequence of pseudo-random number is determined.

Thats why Idy can resume a whole suffled slice from `seed` and `start`, this will make database file small and no grow.

# Credits

Yanke Guo <ryan@islandzero.net>, see LICENSE file for MIT License.
