//go:generate protoc --proto_path=. --go_out=. --go_opt=paths=source_relative  ./base/messages.proto ./keygen/messages.proto
package dkg
