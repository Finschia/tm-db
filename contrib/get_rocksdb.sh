set -e

version="6.20.3"
rocksdb="rocksdb"
archive="v${version}.tar.gz"

rm -rf ${rocksdb} ${rocksdb}-${archive}
wget -q -O ${rocksdb}-${archive} https://github.com/facebook/rocksdb/archive/${archive}
tar -zxvf ${rocksdb}-${archive}
mv ${rocksdb}-${version} ${rocksdb}
