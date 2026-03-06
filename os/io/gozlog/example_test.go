package gozlog_test

import (
	"context"
	"io"
	"os"

	"github.com/bizvip/go-utils/os/io/gozlog"
)

func ExampleManager() {
	cfg := gozlog.DefaultConfig()
	cfg.Output = "stderr"
	cfg.Format = "json"
	cfg.Caller = false

	mgr, err := gozlog.NewManager(cfg)
	if err != nil {
		panic(err)
	}

	log := mgr.Service("example")
	ctx := gozlog.IntoContext(context.Background(), log)
	gozlog.FromContext(ctx).Info().Str("component", "demo").Msg("ready")

	_, _ = io.Copy(io.Discard, os.Stderr)
}
