package launcher

import (
	"context"

	"github.com/goregion/hexago/pkg/goture"
	"github.com/goregion/hexago/pkg/grexit"
	"github.com/goregion/hexago/pkg/log"
)

type appResult struct {
	err error
}

func (r appResult) LogIfError(logger *log.Logger, messages ...any) {
	if r.err != nil {
		logger.LogIfError(r.err, messages...)
	}
}

func (r appResult) Error() error {
	return r.err
}

type appLauncher struct {
	context.Context
}

func NewAppLauncher() *appLauncher {
	return NewAppLauncherWithContext(context.Background())
}

func NewAppLauncherWithContext(ctx context.Context) *appLauncher {
	return &appLauncher{
		Context: ctx,
	}
}

func (a *appLauncher) WithLoggerContext(logger *log.Logger) *appLauncher {
	a.Context = log.WithLoggerContext(a.Context, logger)
	return a
}

func (a *appLauncher) WithGrexitContext() *appLauncher {
	a.Context = grexit.WithGrexitContext(a.Context)
	return a
}

func (a *appLauncher) WaitApplication(task goture.Task) appResult {
	return appResult{
		err: goture.NewGoture(a.Context, task).Wait(),
	}
}

func (a *appLauncher) WaitApplications(task ...goture.Task) appResult {
	return appResult{
		err: goture.NewParallelGoture(a.Context, task...).Wait(),
	}
}
