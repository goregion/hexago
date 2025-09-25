package grexit

import (
	"context"
	"os"
	"runtime"
	"syscall"
	"testing"
	"time"
)

func TestWithGrexitContext(t *testing.T) {
	ctx := context.Background()
	grexitCtx := WithGrexitContext(ctx)

	// Context should not be done initially
	select {
	case <-grexitCtx.Done():
		t.Fatal("Context should not be done initially")
	default:
	}
}

func TestWithGrexitCancelContext(t *testing.T) {
	ctx := context.Background()
	grexitCtx, cancel := WithGrexitCancelContext(ctx)
	defer cancel()

	// Context should not be done initially
	select {
	case <-grexitCtx.Done():
		t.Fatal("Context should not be done initially")
	default:
	}

	// Cancel the context manually
	cancel()

	// Context should be done after cancel
	select {
	case <-grexitCtx.Done():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Context should be done after cancel")
	}
}

func TestWithGrexitCancelContextSignal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Signal testing is complex on Windows")
	}

	ctx := context.Background()
	grexitCtx, cancel := WithGrexitCancelContext(ctx)
	defer cancel()

	// Send SIGTERM to self
	go func() {
		time.Sleep(50 * time.Millisecond)
		process, _ := os.FindProcess(os.Getpid())
		process.Signal(syscall.SIGTERM)
	}()

	// Context should be done after signal
	select {
	case <-grexitCtx.Done():
		// Expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Context should be done after signal")
	}
}

func TestNoGoroutineLeak(t *testing.T) {
	initialGoroutines := runtime.NumGoroutine()

	// Create and immediately cancel context multiple times
	for i := 0; i < 10; i++ {
		ctx := context.Background()
		grexitCtx, cancel := WithGrexitCancelContext(ctx)
		cancel()

		// Wait for context to be done
		<-grexitCtx.Done()
	}

	// Give some time for goroutines to cleanup
	time.Sleep(100 * time.Millisecond)
	runtime.GC()
	runtime.Gosched()

	finalGoroutines := runtime.NumGoroutine()

	// Allow for some variance in goroutine count
	if finalGoroutines > initialGoroutines+2 {
		t.Errorf("Potential goroutine leak: initial=%d, final=%d", initialGoroutines, finalGoroutines)
	}
}

func TestWithGrexitTimeout(t *testing.T) {
	ctx := context.Background()
	grexitCtx := WithGrexitTimeout(ctx)

	// Context should not be done initially
	select {
	case <-grexitCtx.Done():
		t.Fatal("Context should not be done initially")
	default:
	}
}

func TestWithGrexitTimeoutDuration(t *testing.T) {
	ctx := context.Background()
	grexitCtx := WithGrexitTimeoutDuration(ctx, 100*time.Millisecond)

	// Context should not be done initially
	select {
	case <-grexitCtx.Done():
		t.Fatal("Context should not be done initially")
	default:
	}
}

func TestWithGrexitTimeoutCancel(t *testing.T) {
	ctx := context.Background()
	grexitCtx, cancel := WithGrexitTimeoutCancel(ctx)
	defer cancel()

	// Context should not be done initially
	select {
	case <-grexitCtx.Done():
		t.Fatal("Context should not be done initially")
	default:
	}

	// Cancel the context manually
	cancel()

	// Context should be done after cancel
	select {
	case <-grexitCtx.Done():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Context should be done after cancel")
	}
}

func TestWithGrexitTimeoutCancelContext(t *testing.T) {
	ctx := context.Background()
	grexitCtx, cancel := WithGrexitTimeoutCancelContext(ctx, 100*time.Millisecond)
	defer cancel()

	// Context should not be done initially
	select {
	case <-grexitCtx.Done():
		t.Fatal("Context should not be done initially")
	default:
	}

	// Cancel the context manually
	cancel()

	// Context should be done after cancel
	select {
	case <-grexitCtx.Done():
		// Expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Context should be done after cancel")
	}
}

func TestTimeoutFunctionality(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Signal testing is complex on Windows")
	}

	ctx := context.Background()
	shortTimeout := 50 * time.Millisecond
	grexitCtx, cancel := WithGrexitTimeoutCancelContext(ctx, shortTimeout)
	defer cancel()

	start := time.Now()

	// Send signal to trigger timeout
	go func() {
		time.Sleep(10 * time.Millisecond)
		process, _ := os.FindProcess(os.Getpid())
		process.Signal(syscall.SIGTERM)
	}()

	// Wait for timeout
	<-grexitCtx.Done()
	elapsed := time.Since(start)

	// Should complete within timeout + some margin
	if elapsed > shortTimeout+100*time.Millisecond {
		t.Errorf("Timeout took too long: %v, expected around %v", elapsed, shortTimeout)
	}

	// Should take at least the signal delay
	if elapsed < 10*time.Millisecond {
		t.Errorf("Completed too quickly: %v", elapsed)
	}
}
