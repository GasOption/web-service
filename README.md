# web-service

go get github.com/GasOption/web-service

Install dep by following: https://github.com/golang/dep#installation

dep ensure

go run server.go

# misc

Should you see the error "'libsecp256k1/include/secp256k1.h' file not found", use the following
hack:

go get github.com/ethereum/go-ethereum

cp -r \
  "${GOPATH}/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1" \
  "vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/"
