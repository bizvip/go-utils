# gozlog

`gozlog` is a zerolog-based logging engine for Go services.

It provides:
- explicit manager construction
- concurrency-safe scoped logger caching
- request/context logger propagation
- file/stdout/stderr outputs with lumberjack rotation
- service/module child loggers for dependency injection

Example:

```go
cfg := gozlog.DefaultConfig()
cfg.Output = "file"
cfg.LogDir = "runtime/logs"

mgr, err := gozlog.NewManager(cfg)
if err != nil {
    panic(err)
}

log := mgr.Service("api")
log.Info().Msg("server starting")
```
