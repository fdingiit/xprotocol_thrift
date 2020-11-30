PWD := $(shell pwd)
GOGO_PROTOBUF_PATH := "~/go/src/github.com/golang/protobuf/protobuf"

.PHONY: all clean

all: clean
	@ protoc --gofast_out=. \
		--proto_path=./pb/ \
		pb/*.proto
	@ echo "Generate Code Successfully"

clean:
	@ rm -f ./*.pb.go
	@ echo "Cleanup Finished"
