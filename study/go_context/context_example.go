package main

import (
	"context" // The context package
	"fmt"
	"time" // For time.Sleep and time.Duration
)

// performLongRunningTask simulates a task that takes time and can be cancelled.
// It takes a context.Context to listen for cancellation signals.
func performLongRunningTask(ctx context.Context, taskName string) {
	fmt.Printf("%s: Starting long-running task...\n", taskName)

	select {
	case <-time.After(5 * time.Second): // Simulate 5 seconds of work.
		// This case will be selected if the 5-second timer fires before the context is cancelled.
		fmt.Printf("%s: Task completed successfully after 5 seconds.\n", taskName)
	case <-ctx.Done(): // This case will be selected if the context is cancelled.
		// ctx.Err() returns the reason for cancellation (e.g., context.Canceled, context.DeadlineExceeded).
		fmt.Printf("%s: Task cancelled! Reason: %v\n", taskName, ctx.Err())
	}
}

// main function to demonstrate different context types.
func main() {
	fmt.Println("--- Context Example: Cancellation & Timeout ---")

	// --- Scenario 1: Manual Cancellation ---
	fmt.Println("\n--- Scenario 1: Manual Cancellation ---")
	// context.WithCancel returns a new context and a cancel function.
	// We pass ctxCancel to our task, and can call cancelFunc anytime to cancel it.
	ctxCancel, cancelFunc := context.WithCancel(context.Background())
	go performLongRunningTask(ctxCancel, "Task_A (Manual)")

	// Simulate some other work, then decide to cancel Task_A
	time.Sleep(2 * time.Second)
	fmt.Println("Main: Deciding to cancel Task_A...")
	cancelFunc()                // Call the cancel function to signal cancellation.
	time.Sleep(1 * time.Second) // Give some time for Task_A to react and print its message.

	// --- Scenario 2: Timeout ---
	fmt.Println("\n--- Scenario 2: Timeout ---")
	// context.WithTimeout returns a context that automatically cancels after a duration.
	ctxTimeout, cancelTimeoutFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelTimeoutFunc() // It's good practice to defer the cancel function for timeouts/deadlines
	// to release resources associated with the context early if the task finishes before timeout.
	go performLongRunningTask(ctxTimeout, "Task_B (Timeout)")

	// Let Task_B run for its duration + a little extra to see it react to timeout.
	time.Sleep(4 * time.Second)

	// --- Scenario 3: Task completes before timeout (Still defer cancelFunc!) ---
	fmt.Println("\n--- Scenario 3: Task completes before Timeout ---")
	ctxCompleteBeforeTimeout, cancelCompleteFunc := context.WithTimeout(context.Background(), 10*time.Second) // Long timeout
	defer cancelCompleteFunc()
	go performLongRunningTask(ctxCompleteBeforeTimeout, "Task_C (Completes)")
	time.Sleep(6 * time.Second) // Task_C finishes in 5s, we wait 6s to see the completion message.

	fmt.Println("\nMain: All scenarios demonstrated. Exiting.")
}
