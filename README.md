# Forked Tendermint DB for LFB SDK and Ostracon 

Common database interface for various database backends. This is forked for applications built on [Ostracon](https://github.com/line/ostracon), such as the [LFB SDK](https://github.com/line/lfb-sdk) from [Tendermint tm-db](https://github.com/tendermint/tm-db), but can be used independently of these as well.

### Minimum Go Version

Go 1.18+

## Supported Database Backends

- **[GoLevelDB](https://github.com/syndtr/goleveldb) [stable]**: A pure Go implementation of [LevelDB](https://github.com/google/leveldb) (see below). Currently the default on-disk database used in the Cosmos SDK.

- **MemDB [stable]:** An in-memory database using [Google's B-tree package](https://github.com/google/btree). Has very high performance both for reads, writes, and range scans, but is not durable and will lose all data on process exit. Does not support transactions. Suitable for e.g. caches, working sets, and tests. Used for [IAVL](https://github.com/tendermint/iavl) working sets when the pruning strategy allows it.

- **[LevelDB](https://github.com/google/leveldb) [experimental]:** A [Go wrapper](https://github.com/jmhodges/levigo) around [LevelDB](https://github.com/google/leveldb). Uses LSM-trees for on-disk storage, which have good performance for write-heavy workloads, particularly on spinning disks, but requires periodic compaction to maintain decent read performance and reclaim disk space. Does not support transactions.

- **[BoltDB](https://github.com/etcd-io/bbolt) [experimental]:** A [fork](https://github.com/etcd-io/bbolt) of [BoltDB](https://github.com/boltdb/bolt). Uses B+trees for on-disk storage, which have good performance for read-heavy workloads and range scans. Supports serializable ACID transactions.

- **[RocksDB](https://github.com/line/gorocksdb) [experimental]:** A [Go wrapper](https://github.com/line/gorocksdb) around [RocksDB](https://rocksdb.org). Similarly to LevelDB (above) it uses LSM-trees for on-disk storage, but is optimized for fast storage media such as SSDs and memory. Supports atomic transactions, but not full ACID transactions.

- **[BadgerDB](https://github.com/dgraph-io/badger) [experimental]:** A key-value database written as a pure-Go alternative to e.g. LevelDB and RocksDB, with LSM-tree storage. Makes use of multiple goroutines for performance, and includes advanced features such as serializable ACID transactions, write batches, compression, and more.

## Meta-databases

- **PrefixDB [stable]:** A database which wraps another database and uses a static prefix for all keys. This allows multiple logical databases to be stored in a common underlying databases by using different namespaces. Used by the Cosmos SDK to give different modules their own namespaced database in a single application database.

- **RemoteDB [experimental]:** A database that connects to distributed Tendermint db instances via [gRPC](https://grpc.io/). This can help with detaching difficult deployments such as LevelDB, and can also ease dependency management for Tendermint developers.

## Tests

To test common databases, run `make test`. If all databases are available on the local machine, use `make test-all` to test them all.
