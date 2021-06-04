# Changelog

## [Unreleased]

### Features
* (cleveldb/rocksdb) [\#3](https://github.com/line/tm-db/pull/3) Make path for cleveldb, rocksdb
* (prefix) [\#10](https://github.com/line/tm-db/pull/10) Prefix iterator (#10)
* (api) [\#15](https://github.com/line/tm-db/pull/15) Add AvailableDBBackends function (#15)

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
