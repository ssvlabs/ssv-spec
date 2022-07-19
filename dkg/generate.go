//go:generate protoc --proto_path=. --go_out=. --go_opt=paths=source_relative  ./types/messages.proto ./keygen/messages.proto
package dkg
