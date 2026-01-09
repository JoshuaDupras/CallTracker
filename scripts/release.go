package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginLeft(2)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF0000"))

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF00"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D7D7D"))
)

type WailsConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type step int

const (
	stepVersionChoice step = iota
	stepCustomVersion
	stepReleaseNotes
	stepConfirm
	stepExecuting
	stepComplete
)

type versionItem struct {
	title, desc string
	version     string
}

func (i versionItem) Title() string       { return i.title }
func (i versionItem) Description() string { return i.desc }
func (i versionItem) FilterValue() string { return i.title }

type model struct {
	currentStep    step
	currentVersion string
	newVersion     string
	releaseNotes   string
	versionList    list.Model
	customInput    textinput.Model
	notesInput     textarea.Model
	confirmed      bool
	err            error
	output         []string
}

func initialModel() model {
	// Get current version
	currentVer, _ := getCurrentVersion()

	major, minor, patch, _ := parseVersion(currentVer)

	items := []list.Item{
		versionItem{
			title:   "Patch",
			desc:    fmt.Sprintf("Bug fixes → v%d.%d.%d", major, minor, patch+1),
			version: fmt.Sprintf("%d.%d.%d", major, minor, patch+1),
		},
		versionItem{
			title:   "Minor",
			desc:    fmt.Sprintf("New features → v%d.%d.0", major, minor+1),
			version: fmt.Sprintf("%d.%d.0", major, minor+1),
		},
		versionItem{
			title:   "Major",
			desc:    fmt.Sprintf("Breaking changes → v%d.0.0", major+1),
			version: fmt.Sprintf("%d.0.0", major+1),
		},
		versionItem{
			title:   "Custom",
			desc:    "Enter version manually",
			version: "custom",
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select version bump type"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	ti := textinput.New()
	ti.Placeholder = "e.g., 1.0.0"
	ti.CharLimit = 20

	ta := textarea.New()
	ta.Placeholder = "Enter release notes (what's new, fixed, changed)..."
	ta.SetHeight(5)

	return model{
		currentStep:    stepVersionChoice,
		currentVersion: currentVer,
		versionList:    l,
		customInput:    ti,
		notesInput:     ta,
		output:         []string{},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.versionList.SetSize(msg.Width, msg.Height-10)
		return m, nil

	case tea.KeyMsg:
		switch m.currentStep {
		case stepVersionChoice:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				selected := m.versionList.SelectedItem().(versionItem)
				if selected.version == "custom" {
					m.currentStep = stepCustomVersion
					m.customInput.Focus()
					return m, textinput.Blink
				}
				m.newVersion = selected.version
				m.currentStep = stepReleaseNotes
				m.notesInput.Focus()
				return m, textarea.Blink
			}

		case stepCustomVersion:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "enter":
				m.newVersion = m.customInput.Value()
				if !isValidVersion(m.newVersion) {
					m.err = fmt.Errorf("invalid version format")
					return m, nil
				}
				m.currentStep = stepReleaseNotes
				m.notesInput.Focus()
				return m, textarea.Blink
			}

		case stepReleaseNotes:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "ctrl+d", "esc":
				m.releaseNotes = m.notesInput.Value()
				m.currentStep = stepConfirm
				return m, nil
			}

		case stepConfirm:
			switch msg.String() {
			case "ctrl+c", "n", "N":
				return m, tea.Quit
			case "y", "Y":
				m.currentStep = stepExecuting
				return m, m.executeRelease
			}

		case stepComplete:
			return m, tea.Quit
		}
	}

	// Update components
	var cmd tea.Cmd
	switch m.currentStep {
	case stepVersionChoice:
		m.versionList, cmd = m.versionList.Update(msg)
	case stepCustomVersion:
		m.customInput, cmd = m.customInput.Update(msg)
	case stepReleaseNotes:
		m.notesInput, cmd = m.notesInput.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("🚀 Fire Department Call Log - Release Tool"))
	s.WriteString("\n\n")

	if m.currentVersion != "" {
		s.WriteString(infoStyle.Render(fmt.Sprintf("Current version: %s", m.currentVersion)))
		s.WriteString("\n\n")
	}

	switch m.currentStep {
	case stepVersionChoice:
		s.WriteString(m.versionList.View())

	case stepCustomVersion:
		s.WriteString("Enter custom version:\n\n")
		s.WriteString(m.customInput.View())
		s.WriteString("\n\n")
		if m.err != nil {
			s.WriteString(errorStyle.Render(fmt.Sprintf("❌ %v", m.err)))
			s.WriteString("\n\n")
		}
		s.WriteString(infoStyle.Render("Press Enter to continue"))

	case stepReleaseNotes:
		s.WriteString("Release notes (press Ctrl+D or Esc when done):\n\n")
		s.WriteString(m.notesInput.View())

	case stepConfirm:
		s.WriteString(titleStyle.Render("📋 Release Summary"))
		s.WriteString("\n\n")
		s.WriteString(fmt.Sprintf("Old version: %s\n", m.currentVersion))
		s.WriteString(fmt.Sprintf("New version: %s\n", selectedStyle.Render(m.newVersion)))
		s.WriteString(fmt.Sprintf("Tag: %s\n\n", selectedStyle.Render("v"+m.newVersion)))
		if m.releaseNotes != "" {
			s.WriteString("Release notes:\n")
			s.WriteString(infoStyle.Render(m.releaseNotes))
			s.WriteString("\n\n")
		}
		s.WriteString("Proceed with release? (y/N): ")

	case stepExecuting:
		s.WriteString(titleStyle.Render("🔨 Executing Release"))
		s.WriteString("\n\n")
		for _, line := range m.output {
			s.WriteString(line)
			s.WriteString("\n")
		}

	case stepComplete:
		s.WriteString(successStyle.Render("✅ Release completed successfully!"))
		s.WriteString("\n\n")
		s.WriteString(fmt.Sprintf("Version %s has been released.\n", selectedStyle.Render("v"+m.newVersion)))
		s.WriteString("\n")
		s.WriteString(infoStyle.Render("GitHub Actions will now build and publish the release."))
		s.WriteString("\n\n")
		s.WriteString("Press any key to exit...")
	}

	return s.String()
}

func (m *model) executeRelease() tea.Msg {
	m.output = append(m.output, "📝 Updating wails.json...")
	if err := updateWailsVersion(m.newVersion); err != nil {
		m.output = append(m.output, errorStyle.Render(fmt.Sprintf("❌ Error: %v", err)))
		return err
	}
	m.output = append(m.output, successStyle.Render("✅ wails.json updated"))

	if m.releaseNotes != "" {
		m.output = append(m.output, "📝 Updating CHANGELOG.md...")
		if err := updateChangelog(m.newVersion, m.releaseNotes); err != nil {
			m.output = append(m.output, infoStyle.Render(fmt.Sprintf("⚠️  Warning: %v", err)))
		} else {
			m.output = append(m.output, successStyle.Render("✅ CHANGELOG.md updated"))
		}
	}

	m.output = append(m.output, "🧪 Running tests...")
	if err := runCommandSilent("go", "test", "./...", "-v"); err != nil {
		m.output = append(m.output, infoStyle.Render("⚠️  Tests failed (continuing anyway)"))
	} else {
		m.output = append(m.output, successStyle.Render("✅ Tests passed"))
	}

	m.output = append(m.output, "🔨 Building application...")
	if err := runCommandSilent("wails", "build"); err != nil {
		m.output = append(m.output, errorStyle.Render(fmt.Sprintf("❌ Build failed: %v", err)))
		return err
	}
	m.output = append(m.output, successStyle.Render("✅ Build successful"))

	m.output = append(m.output, "📦 Creating git commit...")
	runCommandSilent("git", "add", "wails.json", "CHANGELOG.md")
	commitMsg := fmt.Sprintf("Release v%s", m.newVersion)
	if err := runCommandSilent("git", "commit", "-m", commitMsg); err != nil {
		m.output = append(m.output, errorStyle.Render(fmt.Sprintf("❌ Commit failed: %v", err)))
		return err
	}
	m.output = append(m.output, successStyle.Render("✅ Changes committed"))

	tagName := fmt.Sprintf("v%s", m.newVersion)
	tagMsg := fmt.Sprintf("Release version %s", m.newVersion)
	m.output = append(m.output, fmt.Sprintf("🏷️  Creating tag %s...", tagName))
	if err := runCommandSilent("git", "tag", "-a", tagName, "-m", tagMsg); err != nil {
		m.output = append(m.output, errorStyle.Render(fmt.Sprintf("❌ Tag failed: %v", err)))
		return err
	}
	m.output = append(m.output, successStyle.Render("✅ Tag created"))

	m.output = append(m.output, "🚀 Pushing to GitHub...")
	if err := runCommandSilent("git", "push", "origin", "main"); err != nil {
		m.output = append(m.output, errorStyle.Render(fmt.Sprintf("❌ Push failed: %v", err)))
		return err
	}
	if err := runCommandSilent("git", "push", "origin", tagName); err != nil {
		m.output = append(m.output, errorStyle.Render(fmt.Sprintf("❌ Push tag failed: %v", err)))
		return err
	}
	m.output = append(m.output, successStyle.Render("✅ Pushed to GitHub"))

	m.currentStep = stepComplete
	return nil
}

func main() {
	if !checkGitRepo() {
		fmt.Println(errorStyle.Render("❌ Error: Not in a git repository"))
		os.Exit(1)
	}

	if hasUncommittedChanges() {
		fmt.Println(infoStyle.Render("⚠️  Warning: You have uncommitted changes"))
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Continue anyway? (y/N): ")
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" && response != "yes" {
			fmt.Println("Aborted.")
			os.Exit(0)
		}
		fmt.Println()
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func checkGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

func hasUncommittedChanges() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(output) > 0
}

func getCurrentVersion() (string, error) {
	data, err := ioutil.ReadFile("wails.json")
	if err != nil {
		return "", err
	}

	var config WailsConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return "", err
	}

	return config.Version, nil
}

func isValidVersion(version string) bool {
	// Match semantic versioning: major.minor.patch
	matched, _ := regexp.MatchString(`^\d+\.\d+\.\d+$`, version)
	return matched
}

func parseVersion(version string) (major, minor, patch int, err error) {
	// Remove 'v' prefix if present
	version = strings.TrimPrefix(version, "v")
	
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid version format")
	}
	
	_, err = fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)
	return
}

func updateWailsVersion(version string) error {
	data, err := ioutil.ReadFile("wails.json")
	if err != nil {
		return err
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	config["version"] = version

	output, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile("wails.json", output, 0644)
}

func updateChangelog(version, notes string) error {
	data, err := ioutil.ReadFile("CHANGELOG.md")
	if err != nil {
		return err
	}

	content := string(data)
	today := time.Now().Format("2006-01-02")

	newEntry := fmt.Sprintf("## [%s] - %s\n\n%s\n\n", version, today, notes)

	// Insert after "## [Unreleased]" section
	unreleasedIdx := strings.Index(content, "## [Unreleased]")
	if unreleasedIdx == -1 {
		// If no Unreleased section, add after the header
		lines := strings.Split(content, "\n")
		if len(lines) > 5 {
			// Insert after the intro text
			result := strings.Join(lines[:6], "\n") + "\n\n" + newEntry + strings.Join(lines[6:], "\n")
			return ioutil.WriteFile("CHANGELOG.md", []byte(result), 0644)
		}
	} else {
		// Find the next ## or end of file
		nextSection := strings.Index(content[unreleasedIdx+20:], "\n## ")
		if nextSection == -1 {
			// No next section, append at end
			result := content + "\n" + newEntry
			return ioutil.WriteFile("CHANGELOG.md", []byte(result), 0644)
		} else {
			// Insert between Unreleased and next section
			insertPos := unreleasedIdx + 20 + nextSection + 1
			result := content[:insertPos] + newEntry + content[insertPos:]
			return ioutil.WriteFile("CHANGELOG.md", []byte(result), 0644)
		}
	}

	return nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCommandSilent(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func confirm(reader *bufio.Reader, question string) bool {
	fmt.Printf("%s (y/N): ", question)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
