docker run --rm -v $(pwd):/defs -v $GOPATH/pkg:/defs/pkg namely/protoc-all \
    -d api/proto \
    -i scripts \
    -i pkg/mod \
    -o cmd/service/pb \
    -l go
