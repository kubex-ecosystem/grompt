package squad

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// Agent represents an AI agent specification for AGENTS.md
// Each agent has a title, role description, skills, restrictions and a prompt example.
type Agent struct {
	Title         string
	Role          string
	Skills        []string
	Restrictions  []string
	PromptExample string
}

// ParseRequirements parses a free text requirement string and returns a slice of Agent definitions.
// The parsing uses a simple heuristic based on keywords.
func ParseRequirements(req string) []Agent {
	lower := strings.ToLower(req)
	agents := []Agent{}

	langs := []string{}
	langPrefs := map[string]bool{}

	for _, lang := range []string{"go", "python", "javascript", "typescript", "java", "rust", "ruby"} {
		if strings.Contains(lower, lang) {
			langs = append(langs, strings.Title(lang))
			langPrefs[lang] = true
		}
	}

	restricts := []string{}
	rJava := regexp.MustCompile(`(?i)(sem|without|no)\s+java`)
	if rJava.MatchString(lower) {
		restricts = append(restricts, "No Java")
	}

	if strings.Contains(lower, "backend") || strings.Contains(lower, "microservi") || strings.Contains(lower, "api") {
		title := "Backend Developer"
		if len(langs) > 0 {
			title += " (" + strings.Join(langs, "/") + ")"
		}
		agents = append(agents, Agent{
			Title:         title,
			Role:          "Implement backend services",
			Skills:        append([]string{"REST"}, langs...),
			Restrictions:  restricts,
			PromptExample: "Design and implement backend APIs following the requirements.",
		})
	}

	if strings.Contains(lower, "docker") || strings.Contains(lower, "deploy") || strings.Contains(lower, "ci/cd") || strings.Contains(lower, "kubernetes") {
		agents = append(agents, Agent{
			Title:         "DevOps Engineer",
			Role:          "Setup Docker and CI/CD pipelines",
			Skills:        []string{"Docker", "CI/CD", "Cloud"},
			PromptExample: "Configure Docker containers and deployment pipeline.",
		})
	}

	if strings.Contains(lower, "test") || strings.Contains(lower, "qa") {
		skills := []string{"Automated Testing"}
		if langPrefs["go"] {
			skills = append(skills, "go test")
		}
		if langPrefs["python"] {
			skills = append(skills, "pytest")
		}
		agents = append(agents, Agent{
			Title:         "QA Engineer",
			Role:          "Write and execute automated tests",
			Skills:        skills,
			PromptExample: "Create test suites ensuring high coverage.",
		})
	}

	if len(agents) == 0 {
		agents = append(agents, Agent{
			Title:         "Software Engineer",
			Role:          "Implement requested features",
			Skills:        langs,
			Restrictions:  restricts,
			PromptExample: "Develop the project according to the requirements.",
		})
	}

	return agents
}

// GenerateMarkdown creates the AGENTS.md content based on agents slice.
func GenerateMarkdown(agents []Agent) string {
	b := &strings.Builder{}
	b.WriteString("# Agents\n\n")
	for _, a := range agents {
		b.WriteString("## " + a.Title + "\n")
		if a.Role != "" {
			b.WriteString("- Role: " + a.Role + "\n")
		}
		if len(a.Skills) > 0 {
			b.WriteString("- Skills: " + strings.Join(a.Skills, ", ") + "\n")
		}
		if len(a.Restrictions) > 0 {
			b.WriteString("- Restrictions: " + strings.Join(a.Restrictions, ", ") + "\n")
		}
		if a.PromptExample != "" {
			b.WriteString("- Prompt Example: " + a.PromptExample + "\n")
		}
		b.WriteString("\n")
	}
	return b.String()
}

// WriteFile writes the AGENTS.md content to the provided path.
func WriteFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}

// BuildAndSave generates the agents from requirements and writes AGENTS.md at path.
func BuildAndSave(req, path string) (string, error) {
	agents := ParseRequirements(req)
	md := GenerateMarkdown(agents)
	if err := WriteFile(path, md); err != nil {
		return "", err
	}
	return md, nil
}

// ParseRequirementsWithLLM generates agents using LLM analysis of requirements
func ParseRequirementsWithLLM(req string, llmFunc func(string) (string, error)) ([]Agent, error) {
	prompt := `Analyze the following project requirements and generate a JSON array of AI agents needed for this project.

Requirements: ` + req + `

For each agent, provide:
- Title: A descriptive title for the agent role
- Role: A clear description of what this agent does
- Skills: Array of technical skills and technologies
- Restrictions: Array of things this agent should avoid or not do
- PromptExample: A detailed example of how to use this agent in a prompt

Return ONLY a valid JSON array with this structure:
[
  {
    "Title": "Backend Developer (Python/FastAPI)",
    "Role": "Implement scalable REST APIs and database integrations",
    "Skills": ["Python", "FastAPI", "PostgreSQL", "SQLAlchemy", "Pydantic"],
    "Restrictions": ["No Java", "Avoid synchronous database calls"],
    "PromptExample": "You are a Python backend developer specializing in FastAPI. Create a REST API endpoint for user authentication with JWT tokens, including proper error handling and async database operations."
  }
]

Generate 3-6 agents that would be needed for this project. Be specific about technologies mentioned in the requirements.`

	response, err := llmFunc(prompt)
	if err != nil {
		// Fallback to basic parsing if LLM fails
		return ParseRequirements(req), nil
	}

	// Try to parse the JSON response
	var agents []Agent
	err = json.Unmarshal([]byte(response), &agents)
	if err != nil {
		// If JSON parsing fails, fallback to basic parsing
		return ParseRequirements(req), nil
	}

	return agents, nil
}

// ParseError represents an error that occurred during parsing with location information
type ParseError struct {
	Line    int    `json:"line"`
	Section string `json:"section"`
	Message string `json:"message"`
	Content string `json:"content"`
}

// ParseResult contains the parsed agents and any errors that occurred
type ParseResult struct {
	Agents []Agent      `json:"agents"`
	Errors []ParseError `json:"errors"`
}

// ParseAgentsMarkdown parses an AGENTS.md file content and returns agents with detailed error reporting
func ParseAgentsMarkdown(content string) ParseResult {
	result := ParseResult{
		Agents: []Agent{},
		Errors: []ParseError{},
	}

	lines := strings.Split(content, "\n")
	var currentAgent *Agent
	var currentSection string
	var lineNum int

	for i, line := range lines {
		lineNum = i + 1
		trimmed := strings.TrimSpace(line)

		// Skip empty lines and main title
		if trimmed == "" || strings.HasPrefix(trimmed, "# Agents") {
			continue
		}

		// New agent section (## Title)
		if strings.HasPrefix(trimmed, "## ") {
			// Save previous agent if exists
			if currentAgent != nil {
				if currentAgent.Title == "" {
					result.Errors = append(result.Errors, ParseError{
						Line:    lineNum - 1,
						Section: "Agent Title",
						Message: "Agent title is required",
						Content: line,
					})
				} else {
					result.Agents = append(result.Agents, *currentAgent)
				}
			}

			// Start new agent
			title := strings.TrimSpace(strings.TrimPrefix(trimmed, "##"))
			if title == "" {
				result.Errors = append(result.Errors, ParseError{
					Line:    lineNum,
					Section: "Agent Title",
					Message: "Empty agent title found",
					Content: line,
				})
				currentAgent = &Agent{Title: "Untitled Agent"}
			} else {
				currentAgent = &Agent{Title: title}
			}
			currentSection = "title"
			continue
		}

		// Agent properties (- Property: Value)
		if strings.HasPrefix(trimmed, "- ") && currentAgent != nil {
			propertyLine := strings.TrimPrefix(trimmed, "- ")
			parts := strings.SplitN(propertyLine, ":", 2)

			if len(parts) != 2 {
				result.Errors = append(result.Errors, ParseError{
					Line:    lineNum,
					Section: currentSection,
					Message: "Invalid property format, expected '- Property: Value'",
					Content: line,
				})
				continue
			}

			property := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch strings.ToLower(property) {
			case "role":
				currentAgent.Role = value
				currentSection = "role"
			case "skills":
				skills := parseListValue(value)
				if len(skills) == 0 {
					result.Errors = append(result.Errors, ParseError{
						Line:    lineNum,
						Section: "skills",
						Message: "Empty skills list",
						Content: line,
					})
				}
				currentAgent.Skills = skills
				currentSection = "skills"
			case "restrictions":
				restrictions := parseListValue(value)
				currentAgent.Restrictions = restrictions
				currentSection = "restrictions"
			case "prompt example":
				currentAgent.PromptExample = value
				currentSection = "prompt_example"
			default:
				result.Errors = append(result.Errors, ParseError{
					Line:    lineNum,
					Section: currentSection,
					Message: fmt.Sprintf("Unknown property '%s'", property),
					Content: line,
				})
			}
		}
	}

	// Save last agent
	if currentAgent != nil {
		if currentAgent.Title == "" {
			result.Errors = append(result.Errors, ParseError{
				Line:    lineNum,
				Section: "Agent Title",
				Message: "Agent title is required",
				Content: "End of file",
			})
		} else {
			result.Agents = append(result.Agents, *currentAgent)
		}
	}

	return result
}

// parseListValue parses comma-separated values or already formatted list
func parseListValue(value string) []string {
	// If it looks like it's already a list format, try to parse it
	if strings.Contains(value, ",") {
		items := strings.Split(value, ",")
		var result []string
		for _, item := range items {
			trimmed := strings.TrimSpace(item)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}

	// Single value
	if strings.TrimSpace(value) != "" {
		return []string{strings.TrimSpace(value)}
	}

	return []string{}
}

// ValidateAgent validates an agent structure and returns validation errors
func ValidateAgent(agent Agent) []ParseError {
	var errors []ParseError

	if agent.Title == "" {
		errors = append(errors, ParseError{
			Section: "title",
			Message: "Agent title is required",
		})
	}

	if agent.Role == "" {
		errors = append(errors, ParseError{
			Section: "role",
			Message: "Agent role is recommended",
		})
	}

	if len(agent.Skills) == 0 {
		errors = append(errors, ParseError{
			Section: "skills",
			Message: "At least one skill is recommended",
		})
	}

	return errors
}

// ExportAgentsToMarkdown creates a comprehensive AGENTS.md with metadata
func ExportAgentsToMarkdown(agents []Agent, includeMetadata bool) string {
	b := &strings.Builder{}

	if includeMetadata {
		b.WriteString("# Agents\n\n")
		b.WriteString("<!-- Generated by Grompt Agent Manager -->\n")
		b.WriteString(fmt.Sprintf("<!-- Generated on: %s -->\n", time.Now().Format("2006-01-02 15:04:05")))
		b.WriteString(fmt.Sprintf("<!-- Total agents: %d -->\n\n", len(agents)))
	} else {
		b.WriteString("# Agents\n\n")
	}

	for _, agent := range agents {
		b.WriteString("## " + agent.Title + "\n")

		if agent.Role != "" {
			b.WriteString("- Role: " + agent.Role + "\n")
		}

		if len(agent.Skills) > 0 {
			b.WriteString("- Skills: " + strings.Join(agent.Skills, ", ") + "\n")
		}

		if len(agent.Restrictions) > 0 {
			b.WriteString("- Restrictions: " + strings.Join(agent.Restrictions, ", ") + "\n")
		}

		if agent.PromptExample != "" {
			b.WriteString("- Prompt Example: " + agent.PromptExample + "\n")
		}

		b.WriteString("\n")
	}

	return b.String()
}
