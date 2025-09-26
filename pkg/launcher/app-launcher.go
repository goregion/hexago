// Package launcher provides a fluent API for launching applications
// with proper context management, logging, and graceful shutdown capabilities.
package launcher

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/goregion/goture"
	"github.com/goregion/grexit"
	"github.com/goregion/hexago/pkg/log"
)

// AppResult represents the result of application execution with error handling capabilities
type AppResult struct {
	Err error // Public error field for easy access
}

// LogIfError logs the error if it exists using the provided logger.
// It safely handles nil logger and only logs when there's an actual error.
func (r *AppResult) LogIfError(logger *log.Logger, messages ...any) {
	if r.Err != nil && logger != nil {
		logger.LogIfError(r.Err, messages...)
	}
}

// Error returns the underlying error from application execution
func (r *AppResult) Error() error {
	return r.Err
}

// AppLauncher provides a fluent API for launching applications
// with proper context management, logging, and graceful shutdown capabilities.
// It manages application context and provides methods for launching applications
// with proper context enrichment (logging, graceful shutdown, etc.)
type AppLauncher struct {
	context.Context
	parentContext context.Context    // Store parent context before timeout
	cancelFunc    context.CancelFunc // Store cancel function for timeout cleanup - keep private for safety
}

// NewAppLauncher creates a new application launcher with background context.
// This is the entry point for building application launch configuration.
func NewAppLauncher() *AppLauncher {
	return NewAppLauncherWithContext(context.Background())
}

// NewAppLauncherWithContext creates a new application launcher with the provided context.
// If the provided context is nil, it defaults to background context for safety.
func NewAppLauncherWithContext(ctx context.Context) *AppLauncher {
	if ctx == nil {
		ctx = context.Background()
	}
	return &AppLauncher{
		Context:       ctx,
		parentContext: ctx,
	}
}

// WithLoggerContext enriches the launcher context with a logger instance.
// The logger will be available to all launched applications through the context.
// Returns the same launcher instance for method chaining (fluent API).
func (a *AppLauncher) WithLoggerContext(logger *log.Logger) *AppLauncher {
	if logger != nil {
		a.Context = log.WithLoggerContext(a.Context, logger)
		a.parentContext = log.WithLoggerContext(a.parentContext, logger)
	}
	return a
}

// WithGrexitContext enriches the launcher context with graceful exit capabilities.
// This enables automatic handling of system signals (SIGINT, SIGTERM) for clean shutdown.
// Returns the same launcher instance for method chaining (fluent API).
func (a *AppLauncher) WithGrexitContext() *AppLauncher {
	a.Context = grexit.WithGrexitContext(a.Context)
	a.parentContext = grexit.WithGrexitContext(a.parentContext)
	return a
}

// WithContext enriches the launcher context with a custom key-value pair.
// This allows passing custom configuration or dependencies to launched applications.
// Returns the same launcher instance for method chaining (fluent API).
func (a *AppLauncher) WithContext(key, value any) *AppLauncher {
	if key != nil {
		a.Context = context.WithValue(a.Context, key, value)
		a.parentContext = context.WithValue(a.parentContext, key, value)
	}
	return a
}

// WithTimeout sets a timeout for application execution.
// If the timeout expires, the context will be canceled and applications will receive cancellation signal.
// If called multiple times, the previous timeout will be canceled and replaced with the new one.
// Returns the same launcher instance for method chaining (fluent API).
func (a *AppLauncher) WithTimeout(timeout time.Duration) *AppLauncher {
	// Cancel previous timeout if it exists to prevent resource leaks
	if a.cancelFunc != nil {
		a.cancelFunc()
		a.cancelFunc = nil
	}

	if timeout > 0 {
		ctx, cancel := context.WithTimeout(a.parentContext, timeout)
		a.Context = ctx
		a.cancelFunc = cancel
	}
	return a
}

// WithoutTimeout removes any previously set timeout.
// This cancels the existing timeout and restores the context to its state before timeout was applied.
// Returns the same launcher instance for method chaining (fluent API).
func (a *AppLauncher) WithoutTimeout() *AppLauncher {
	if a.cancelFunc != nil {
		a.cancelFunc()
		a.cancelFunc = nil
		a.Context = a.parentContext
	}
	return a
}

// WaitApplication launches a single application task and waits for its completion.
// The task receives the enriched context with all configured dependencies (logger, graceful exit, etc.)
// Returns AppResult which can be used to check for errors and log them if needed.
func (a *AppLauncher) WaitApplication(task goture.Task) *AppResult {
	if task == nil {
		return &AppResult{Err: errors.New("task cannot be nil")}
	}

	// Ensure cleanup if timeout was set
	if a.cancelFunc != nil {
		defer a.cancelFunc()
	}

	return &AppResult{
		Err: goture.NewGoture(a.Context, task).Wait(),
	}
}

// WaitApplications launches multiple application tasks in parallel and waits for their completion.
// All tasks receive the same enriched context and run concurrently.
// If any task fails, the error will be returned in the result.
// Returns AppResult which can be used to check for errors and log them if needed.
func (a *AppLauncher) WaitApplications(tasks ...goture.Task) *AppResult {
	if len(tasks) == 0 {
		return &AppResult{Err: errors.New("at least one task must be provided")}
	}

	// Validate all tasks before launching
	for i, task := range tasks {
		if task == nil {
			return &AppResult{Err: fmt.Errorf("task at index %d cannot be nil", i)}
		}
	}

	// Ensure cleanup if timeout was set
	if a.cancelFunc != nil {
		defer a.cancelFunc()
	}

	return &AppResult{
		Err: goture.NewParallelGoture(a.Context, tasks...).Wait(),
	}
}
