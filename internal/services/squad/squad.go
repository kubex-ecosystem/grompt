package squad

import (
	"os"
	"regexp"
	"strings"
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
