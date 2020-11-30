go build -gcflags "all=-N -l" -buildmode=plugin -o thrift.so  \
  ./thrift/api.go \
  ./thrift/buffer.go \
  ./thrift/command.go \
  ./thrift/decoder.go \
  ./thrift/encoder.go \
  ./thrift/mapping.go \
  ./thrift/matcher.go \
  ./thrift/protocol.go \
  ./thrift/types.go \
  ./thrift/logger.go