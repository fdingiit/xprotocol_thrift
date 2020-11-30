PWD := $(shell pwd)
GOGO_PROTOBUF_PATH := "~/go/src/github.com/gogo/protobuf/protobuf"

.PHONY: all clean

all: clean
	@ protoc --gogofast_out=. \
		--proto_path=./pb/ \
		pb/*.proto
	@ echo "Generate Code Successfully"

clean:
	@ rm -f ./*.pb.go
	@ echo "Cleanup Finished"
