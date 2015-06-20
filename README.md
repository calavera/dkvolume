# Docker volume extension api.

Go handler to create external volume extensions for Docker.

## Usage

This library is designed to be integrated in your program.

1. Implement the `dkv.Driver` interface.
2. Initialize a `dkv.Hander` with your implementation.
3. Call either `ServeTCP` or `ServeUnix` from the `dkv.Handler`.

### Example using TCP sockets:

```go
  d := MyVolumeDriver{}
  h := dkv.NewHandler(d)
  h.ServeTCP("test_volume", ":8080")
```

### Example using Unix sockets:

```go
  d := MyVolumeDriver{}
  h := dkv.NewHandler(d)
  h.ServeUnix("root", "/usr/share/docker/plugins/test_volume.sock")
```

See a full example in https://github.com/calavera/docker-volume-glusterfs

## License

MIT
