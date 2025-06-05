package main

import "fmt"

// Filter takes a slice of any type T and a predicate function.
// It returns a new slice containing only the elements for which the predicate returns true.
// `[T any]` is the type parameter list. `any` is an alias for `interface{}` and allows any type.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T // Initialize an empty slice of type T.

	for _, v := range slice {
		if predicate(v) {
			result = append(result, v) // Append elements that satisfy the predicate.
		}
	}
	return result
}

// Map takes a slice of type T and a function that transforms T to R.
// It returns a new slice containing the transformed elements of type R.
func Map[T, R any](slice []T, transform func(T) R) []R {
	result := make([]R, len(slice)) // Pre-allocate slice with known length for efficiency.
	for i, v := range slice {
		result[i] = transform(v)
	}
	return result
}

func Sort[T any](slice []T, less func(a, b T) bool) []T {
	result := make([]T, len(slice))
	copy(result, slice) // Create a copy of the original slice to sort.
	for i := 0; i < len(result)-1; i++ {
		for j := 0; j < len(result)-i-1; j++ {
			if less(result[j+1], result[j]) { // Compare elements using the provided less function.
				result[j], result[j+1] = result[j+1], result[j] // Swap if they are in the wrong order.
			}
		}
	}
	return result // Return the sorted slice.

}

func Slice[T any](slice []T, start, end int) []T {
	// Slice returns a sub-slice of the original slice from start to end indices.
	if start < 0 || end > len(slice) || start > end {
		return nil // Return nil if indices are out of bounds or invalid.
	}
	return slice[start:end] // Return the sub-slice.
}

func Remove[T any](slice []T, index int) []T {
	// Remove removes an element at the specified index from the slice.
	if index < 0 || index >= len(slice) {
		return slice // Return the original slice if index is out of bounds.
	}
	return append(slice[:index], slice[index+1:]...) // Return a new slice without the element at index.
}

func Reduce[T, R any](slice []T, initial R, reducer func(R, T) R) R {
	result := initial // Start with the initial value.
	for _, v := range slice {
		result = reducer(result, v) // Apply the reducer function to each element.
	}
	return result // Return the final accumulated result.
}

// Struct to demonstrate filtering/mapping custom types
type User struct {
	ID   int
	Name string
	Age  int
}

func main() {
	fmt.Println("--- Go Generics Example ---")

	// --- Example 1: Filtering Integers ---
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	evenNumbers := Filter(numbers, func(n int) bool {
		return n%2 == 0
	})
	fmt.Println("Even Numbers:", evenNumbers) // Output: Even Numbers: [2 4 6 8 10]

	// --- Example 2: Filtering Strings ---
	fruits := []string{"apple", "banana", "cherry", "date", "grape"}
	longFruits := Filter(fruits, func(s string) bool {
		return len(s) > 5
	})
	fmt.Println("Long Fruits:", longFruits) // Output: Long Fruits: [banana cherry]

	// --- Example 3: Filtering Custom Structs (Users) ---
	users := []User{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
		{ID: 3, Name: "Charlie", Age: 35},
		{ID: 4, Name: "Diana", Age: 20},
	}
	adultUsers := Filter(users, func(u User) bool {
		return u.Age >= 21
	})
	fmt.Println("Adult Users:", adultUsers)
	// Output: Adult Users: [{1 Alice 30} {2 Bob 25} {3 Charlie 35}]

	// --- Example 4: Mapping Integers to Strings ---
	numbersToWords := Map(numbers, func(n int) string {
		switch n {
		case 1:
			return "one"
		case 2:
			return "two"
		case 3:
			return "three"
		case 4:
			return "four"
		case 5:
			return "five"
		case 6:
			return "six"
		case 7:
			return "seven"
		case 8:
			return "eight"
		case 9:
			return "nine"
		case 10:
			return "ten"
		default:
			return "other"
		}
	})
	fmt.Println("Numbers to Words:", numbersToWords)
	// Output: Numbers to Words: [one two three four five six seven eight nine ten]

	// --- Example 5: Mapping User Structs to their Names ---
	userNames := Map(users, func(u User) string {
		return u.Name
	})
	fmt.Println("User Names:", userNames)
	// Output: User Names: [Alice Bob Charlie Diana]
	// --- Example 6: Sorting Integers ---
	unsortedNumbers := []int{10, 3, 5, 1, 9, 2, 8, 6, 4, 7}
	sortedNumbers := Sort(unsortedNumbers, func(a, b int) bool {
		return a < b // Sort in ascending order.
	})
	fmt.Println("Sorted Numbers:", sortedNumbers) // Output: Sorted Numbers: [1 2 3 4 5 6 7 8 9 10]
	// --- Example 7: Sorting Custom Structs (Users) by Age ---
	sortedUsers := Sort(users, func(a, b User) bool {
		return a.Age < b.Age // Sort users by age in ascending order.
	})
	fmt.Println("Sorted Users by Age:", sortedUsers)
	// Output: Sorted Users by Age: [{4 Diana 20} {2 Bob 25} {1 Alice 30} {3 Charlie 35}]
	// --- Example 8: Sorting Custom Structs (Users) by Name ---

	// --- Example 9: Reducing Integers to their Sum ---
	sum := Reduce(numbers, 0, func(acc int, n int) int {
		return acc + n // Sum all integers in the slice.
	})
	fmt.Println("Sum of Numbers:", sum) // Output: Sum of Numbers: 55
	// --- Example 10: Reducing User to generate a summary string ---
	userSummary := Reduce(users, "", func(acc string, u User) string {
		if acc == "" {
			return fmt.Sprintf("User %d: %s (%d years old)", u.ID, u.Name, u.Age)
		}
		return acc + fmt.Sprintf("\nUser %d: %s (%d years old)", u.ID, u.Name, u.Age)
	})
	fmt.Println("User Summary:\n", userSummary)
	// Output: User Summary:
	// User 1: Alice (30 years old)
	// User 2: Bob (25 years old)
	// User 3: Charlie (35 years old)
	// User 4: Diana (20 years old)

	// Example 11 : Reducing User to generate a map of user IDs to names
	userIDMap := Reduce(users, make(map[int]User), func(acc map[int]User, u User) map[int]User {
		acc[u.ID] = u // Add user ID and name to the map.
		return acc
	})
	fmt.Println("User ID Map:", userIDMap)
	// Output: User ID Map: map[1:{1 Alice 30} 2:{2 Bob 25} 3:{3 Charlie 35} 4:{4 Diana 20}]

	// Example 12: Slicing a slice of integers
	slicedNumbers := Slice(numbers, 2, 5)         // Get a sub-slice from index 2 to 5 (exclusive).
	fmt.Println("Sliced Numbers:", slicedNumbers) // Output: Sliced Numbers: [3 4 5]

	// Example 13: Removing an element from a slice of integers
	removedNumbers := Remove(numbers, 3)            // Remove the element at index 3 (4th element).
	fmt.Println("Removed Numbers:", removedNumbers) // Output: Removed Numbers: [1 2 3 5 6 7 8 9 10]

	// Example 14: Removing an element from a slice of users
	removedUsers := Remove(users, 1) // Remove the user at index 1 (2nd user).
	fmt.Println("Removed Users:", removedUsers)
}
