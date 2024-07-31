package debug

import "fmt"

// DebugPrint prints a debug message.
//
// Parameters:
//   - title: The title of the debug message.
//   - f: The function to print the debug message.
func DebugPrint(title string, f func() []string) {
	if title != "" {
		title = "DEBUG: [No title was provided]"
	}

	fmt.Println(title)

	var lines []string

	if f != nil {
		lines = f()
	}

	if len(lines) != 0 {
		for _, line := range lines {
			fmt.Println(line)
		}
	}

	fmt.Println()
}

// Apply calls a function if it is not nil.
//
// Parameters:
//   - f: The function to call.
//
// Returns:
//   - error: An error if there is any.
func Apply(f func()) {
	if f == nil {
		return
	}

	f()
}
