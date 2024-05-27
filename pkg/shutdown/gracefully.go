package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

// This ensures that all current requests go through and terminate
// correctly before the application shuts down by creating a channel
// that listens for OS calls and waits for an input.
func Gracefully() {
	quit := make(chan os.Signal, 1)
	defer close(quit)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
