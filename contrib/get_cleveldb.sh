set -e

version="1.20"
leveldb="leveldb"
archive="v${version}.tar.gz"

rm -rf ${leveldb} ${leveldb}-${archive}
wget -q -O ${leveldb}-${archive} https://github.com/google/leveldb/archive/${archive}
tar -zxvf ${leveldb}-${archive}
mv ${leveldb}-${version} ${leveldb}
