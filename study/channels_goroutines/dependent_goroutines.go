package main

import (
	"fmt"
	"math/rand"
	"time" // For simulating delays
)

// simulateFetchUserID simulates fetching a user ID from a database or API.
// It returns the fetched ID and sends a signal on the 'done' channel when complete.
func simulateFetchUserID(done chan<- map[int]int) { // `chan<- int` means this channel can only send integers.
	fmt.Println("Task A: Fetching user ID...")

	userID := rand.Intn(1000) // Simulate fetching a user ID (random number for demonstration)
	userIDs := make(map[int]int)
	for i := 1; i <= 5; i++ {
		time.Sleep(300 * time.Millisecond) // Simulate network/database latency
		userIDs[i] = userID + i
	}
	fmt.Printf("Task A: User IDs fetched: %v\n", userIDs)
	done <- userIDs // Send the user IDs to the channel. This "signals" completion and passes data.
}

// simulateFetchUserDetails simulates fetching user details using a user ID.
// It receives the user ID from the 'userIDChan' channel.
func simulateFetchUserDetails(userIDChan <-chan map[int]int, done chan<- bool) { // `<-chan int` means this channel can only receive integers.
	fmt.Println("Task B: Waiting for user IDs...")
	userIDs := <-userIDChan // Receive the user IDs from the channel. This operation blocks until a value is sent.
	userDetails := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		time.Sleep(1 * time.Second) // Simulate more latency
		fmt.Printf("Task B: Received user ID %d. Fetching user details...\n", userID)
		userDetails = append(userDetails, fmt.Sprintf("Details for User %d: Name 'User-%d', Email 'user-%d@example.com'", userID, userID, userID))
	}
	fmt.Printf("Task B: User details fetched: %s\n", userDetails)
	done <- true // Signal that Task B is complete.
}

func main() {
	fmt.Println("Starting main program...")

	// 1. Create channels for communication and signaling
	//    userIDChannel: To send the user ID from Task A to Task B. (unbuffered channel)
	//    taskBDoneChannel: To signal when Task B has completed. (unbuffered channel)
	userIDsChannel := make(chan map[int]int)
	taskBDoneChannel := make(chan bool)

	// 2. Launch Goroutine for Task A
	// This goroutine will run concurrently. It sends the userID to userIDChannel when done.
	go simulateFetchUserID(userIDsChannel)

	// 3. Launch Goroutine for Task B
	// This goroutine also runs concurrently. It *waits* for a value on userIDsChannel.
	// Once it receives the userIDs, it proceeds.
	go simulateFetchUserDetails(userIDsChannel, taskBDoneChannel)

	// 4. Main goroutine waits for Task B to finish
	// This line will block until Task B sends a value on taskBDoneChannel,
	// ensuring that both tasks have completed before the main program exits.
	<-taskBDoneChannel
	fmt.Println("Main program: Both dependent tasks completed.")
}
