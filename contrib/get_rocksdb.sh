set -e

version="6.20.3"
rocksdb="rocksdb"
rocksdb_dir="rocksdb.build"
archive="v${version}.tar.gz"

rm -rf ${rocksdb_dir} ${rocksdb}-${archive}
wget -O ${rocksdb}-${archive} https://github.com/facebook/rocksdb/archive/${archive}
tar -zxvf ${rocksdb}-${archive}
mv ${rocksdb}-${version} ${rocksdb_dir}
