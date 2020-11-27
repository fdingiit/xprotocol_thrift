go build -mod=vendor -buildmode=plugin -o thrift.so  \
  ./thrift/buffer.go \
  ./thrift/command.go \
  ./thrift/decoder.go \
  ./thrift/encoder.go \
  ./thrift/mapping.go \
  ./thrift/matcher.go \
  ./thrift/protocol.go \
  ./thrift/types.go