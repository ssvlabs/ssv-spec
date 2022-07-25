//go:generate protoc --proto_path=. --go_out=. --go_opt=paths=source_relative  ./dkg/types/messages.proto ./gg20/types/messages.proto
package ssv_spec
