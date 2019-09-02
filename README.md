# Task Tracker

Inspired by [@campoy](https://github.com/campoy)'s JustForFunc series #31 video: https://www.youtube.com/watch?v=uolTUtioIrc

## Usage

0. Clone repo

`git clone https://github.com/rmccorm4/gotodo-grpc`

1. Start server

```
# Start server listening on port 8888
go run cmd/server/main.go
```

2. Start adding tasks with client

```bash
go run cmd/client/main.go add Go to the gym
go run cmd/client/main.go add Write more code
go run cmd/client/main.go add Sleep
```

3. List your tasks

```bash
> go run cmd/client/main.go list

❌ Go to the gym
❌ Write more code
❌ Sleep
```

## Future Work
* [ ] Ability to check off tasks
* [ ] Ability to remove tasks

## Development Setup (Optional)

1. Download protoc binary and add it to your PATH

See [https://github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases) for
the binary download. Then just add the `bin/protoc` binary to your PATH.

2. Install protoc-gen-go plugin
  * `go get -u github.com/golang/protobuf/protoc-gen-go`
  * Add this binary or `${GOBIN}` to your path so `protoc` can find the plugin

3. Generate GRPC code stubs

```
cd todo/
protoc -I . todo.proto --go_out=plugins=grpc:.
```
## Debugging

If you have `protoc` installed, you can view your protobuf data at 
anytime with the following command:

`cat db.pb | protoc --decode_raw`
