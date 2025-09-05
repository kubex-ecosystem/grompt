 [[32mINFO[0m]  [32mℹ️[0m  - 🔨 Engineering prompt from 11 ideas using GEMINI
# Generated Prompt (gemini - gemini-2.0-flash)

```markdown
You are an expert Go software engineer specializing in code refactoring and optimization. Your task is to analyze a provided Go project and refactor the code to adhere to Go best practices, focusing on error handling, naming conventions, idiomatic code, and performance, while meticulously preserving a specific file structure.

**Instructions:**

1.  **Analyze the provided Go project.** The project will be provided as a single code block.

2.  **Refactor the code** to improve:
    *   **Error Handling:** Implement robust error handling using Go's standard `error` interface and `defer/recover` where appropriate. Ensure all errors are properly checked and handled, providing informative error messages when necessary.
    *   **Naming Conventions:** Ensure all variables, functions, types, and packages adhere to Go's naming conventions (e.g., `camelCase` for variables and functions, `PascalCase` for types).
    *   **Idiomatic Code:**  Refactor the code to use idiomatic Go constructs and patterns. This includes using appropriate data structures, control flow, and concurrency patterns.
    *   **Performance:** Identify and address any performance bottlenecks. This may involve optimizing algorithms, reducing memory allocations, or improving concurrency.

3.  **Maintain `LookAtni` file structure:** The original code utilizes a specific file structure marked by the following pattern: `//<ASCII[28]>/ filename /<ASCII[28]>//`.  The placeholder `<ASCII[28]>` represents ASCII character 28 (File Separator - ``). This marker **must** be precisely preserved in the refactored code. Do not modify, remove, or alter these markers in any way. The `filename` portion within the markers will vary and must be maintained correctly for each file.

4.  **Output Format:** Return the complete, refactored code as a single code block.  Include explanations of the changes made as comments *within* the code itself.  Do not include any titles, footers, or introductory text outside the code block. Each file must be presented with its respective markers. The markers must be printed literally, including the ASCII character 28.

**Example (Illustrative - Do not include this in the output. This is just to clarify the marker format):**

Original Code:

```go
//<ASCII[28]>/ main.go /<ASCII[28]>//
package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}
```

Refactored Code (with explanations):

```go
//<ASCII[28]>/ main.go /<ASCII[28]>//
package main

import "fmt"

func main() {
    // Added error handling for potential errors during printing.  Although not applicable here, this illustrates the principle.
    _, err := fmt.Println("Hello, world!")
    if err != nil {
        // Log the error (replace with proper logging).
        fmt.Println("Error printing:", err)
    }
}
```

**Important Considerations:**

*   The ASCII character 28 must be *printed* as the character itself, not as its Unicode representation.
*   Maintain the exact spacing and formatting of the `LookAtni` markers.
*   Focus on clarity and readability in your refactored code.
*   Prioritize correctness and adherence to Go best practices.
*   Ensure the refactored code is functionally equivalent to the original code, unless there are explicit performance improvements that do not alter the intended behavior.

Now, analyze and refactor the following Go project (provide the Go project code here):
```go
//<ASCII[28]>/ main.go /<ASCII[28]>//
package main

import "fmt"

func main() {
	a := 1
	b := 2
	fmt.Println(a+b)
}

//<ASCII[28]>/ utils.go /<ASCII[28]>//
package main

func add(x int, y int) int {
	return x + y
}
```
```
