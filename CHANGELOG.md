# Changelog

## 0.6.6

**2022-02-23**

**Important node:** Line version was being developed on v0.6.5 line, that's
been reverted to flat directory structure in v0.6.6. 

### Breaking Changes
* (line/tm-db) [\#36](https://github.com/line/tm-db/pull/36) Merge v0.6.6

**2021-11-08**

**Important note:** Version v0.6.5 was accidentally tagged and should be
avoided.  This version is identical to v0.6.4 in package structure and API, but
has updated the version marker so that normal `go get` upgrades will not
require modifying existing use of v0.6.4.

### Version bumps (since v0.6.4)

- Bump grpc from to 1.42.0.
- Bump dgraph/badger to v2 2.2007.3.
- Bump go.etcd.io/bbolt to 1.3.6.

## 0.6.5

**2021-08-04**

**Important note**: This version was tagged by accident, and should not be
used. The tag now points to the [package-reorg
branch](https://github.com/tendermint/tm-db/tree/package-reorg) so that
any existing dependencies should not break.

## 0.6.4

### Features
* (cleveldb/rocksdb) [\#3](https://github.com/line/tm-db/pull/3) Make path for cleveldb, rocksdb
* (prefix) [\#10](https://github.com/line/tm-db/pull/10) Prefix iterator (#10)
* (api) [\#15](https://github.com/line/tm-db/pull/15) Add AvailableDBBackends function (#15)
* (rdb) [\#34](https://github.com/line/tm-db/pull/34) Name & WriteLowPri methods (#34)

### Improvements
* (global) [\#1](https://github.com/line/tm-db/pull/1) Revise module path
* (perf) [\#4](https://github.com/line/tm-db/pull/4) Optimize 2 `[]byte`s concatenation (#4)
* (prefixdb) [\#6](https://github.com/line/tm-db/pull/6) Package `prefixdb` (#6)
* (badgerdb) [\#8](https://github.com/line/tm-db/pull/8) Re-org badgerdb files to follow up convention (#8)
* (perf) [\#9](https://github.com/line/tm-db/pull/9) Pointer receiver for cLevelDBIterator (#9)
* (global) [\#11](https://github.com/line/tm-db/pull/11) Remove Iterator.Domain() (#11)
* (goleveldb) [\#12](https://github.com/line/tm-db/pull/12) Revise goleveldb iterator (#12)
* (prefixdb) [\#13](https://github.com/line/tm-db/pull/13) Revise `PrefixRange()` to return err as well (#13)
* (perf) [\#14](https://github.com/line/tm-db/pull/14) Optimize cleveldb iterator (#14)
* (rocksdb) [\#16](https://github.com/line/tm-db/pull/16) Revise rocksdb iterator (#16)
* (prefixdb) [\#17](https://github.com/line/tm-db/pull/17) Revise prefixdb iterator (#17)
* (test) [\#18](https://github.com/line/tm-db/pull/18) No writes on an iterator (#18)
* (prefixdb) [\#19](https://github.com/line/tm-db/pull/19) Simplify prefixdb iterator and add tests (#19)
* (perf) [\#20](https://github.com/line/tm-db/pull/20) Remote prefixdb.mtx (#20)

### Bug Fixes
* (test) [\#5](https://github.com/line/tm-db/pull/5) Fix test (#5)
* (test) [\#7](https://github.com/line/tm-db/pull/7) Revise test (#7)

### Breaking Changes

## [tendermint/tm-db v0.6.4] - 2021-03-15
Initial line/tm-db is based on the tendermint/tm-db v0.6.4

* (tendermint/tm-db) [v0.6.4](https://github.com/tendermint/tm-db/releases/tag/v0.6.4).

Please refer [CHANGELOG_OF_tendermint/tm-db_v0.6.4](https://github.com/tendermint/tm-db/blob/v0.6.4/CHANGELOG.md)
<!-- Release links -->
