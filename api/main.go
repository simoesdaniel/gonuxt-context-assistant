package main

import (
	"fmt"                                         // Importing the fmt package for formatted I/O operations.
	"gonuxt-context-assistant/api/internal/tools" // Importing the tools package where our functions are defined.
)

func main() {
	currentTime := tools.GetCurrentDateTime() // Call the GetCurrentDateTime function from the tools package.
	fmt.Println(currentTime)                  // Print the current date and time.

	londonWeather := tools.GetWeather("London") // Call the GetWeather function for London.
	fmt.Println(londonWeather)                  // Print the weather report for London.

	portoWeather := tools.GetWeather("Porto") // Call the GetWeather function for Porto.
	fmt.Println(portoWeather)                 // Print the weather report for Porto.

	capital := tools.GetCapital("Portugal")
	fmt.Println("The capital of Portugal is:", capital) // Print the capital city of Portugal.
}
