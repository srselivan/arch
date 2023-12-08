.PHONY: migrate_create, proto

# Migration name
MIGRATE_NAME=create_permisions
# Migration files path
MIGRATE_PATH=.\migrations\

# destination to generated proto files
PROTO_DST_DIR := proto/pb
# proto files home
PROTO_SRC_DIR := proto/
# proto file name
PROTO_FILE_NAME := entity.proto

migrate_create:
	migrate create -ext sql -dir ${MIGRATE_PATH} -seq ${MIGRATE_NAME}

proto:
	protoc --proto_path=${PROTO_SRC_DIR} --go_opt=paths=source_relative \
	-I=${PROTO_SRC_DIR} --go_out=${PROTO_DST_DIR} ${PROTO_SRC_DIR}/${PROTO_FILE_NAME}