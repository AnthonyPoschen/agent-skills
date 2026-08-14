// Command ticket-orchestrator runs a durable, dependency-gated implementation
// queue. It deliberately keeps tracker publication authority outside workers.
package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const stateVersion = 1

type stringList []string

func (s *stringList) String() string { return strings.Join(*s, ",") }
func (s *stringList) Set(v string) error {
	*s = append(*s, v)
	return nil
}

type intList []int

func (s *intList) String() string {
	values := make([]string, 0, len(*s))
	for _, v := range *s {
		values = append(values, strconv.Itoa(v))
	}
	return strings.Join(values, ",")
}
func (s *intList) Set(v string) error {
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		return fmt.Errorf("issue number must be a positive integer: %q", v)
	}
	*s = append(*s, n)
	return nil
}

type Config struct {
	Repo           string   `json:"repo"`
	RepoSlug       string   `json:"repo_slug"`
	Source         string   `json:"source"`
	Harness        string   `json:"harness"`
	Target         string   `json:"target"`
	Remote         string   `json:"remote"`
	ReadyLabel     string   `json:"ready_label"`
	Concurrency    int      `json:"concurrency"`
	PollSeconds    int      `json:"poll_seconds"`
	StallMinutes   int      `json:"stall_minutes"`
	WorktreeRoot   string   `json:"worktree_root"`
	RunDir         string   `json:"run_dir"`
	WorkerAgent    string   `json:"worker_agent,omitempty"`
	WorkerModel    string   `json:"worker_model,omitempty"`
	ManagerModel   string   `json:"manager_model,omitempty"`
	LauncherArgs   []string `json:"launcher_args,omitempty"`
	VerifyCommands []string `json:"verify_commands,omitempty"`
	Invocation     string   `json:"implementation_invocation,omitempty"`
	Publish        bool     `json:"publish"`
	IssueNumbers   []int    `json:"issue_numbers,omitempty"`
	IssueQuery     string   `json:"issue_query,omitempty"`
}

// ProjectSetup is the repository-owned orchestration contract. Command-line
// flags override it; tool and documentation discovery fill only missing data.
type ProjectSetup struct {
	Version      int      `json:"version"`
	Source       string   `json:"source,omitempty"`
	Target       string   `json:"target,omitempty"`
	Remote       string   `json:"remote,omitempty"`
	ReadyLabel   string   `json:"ready_label,omitempty"`
	Harness      string   `json:"harness,omitempty"`
	Concurrency  int      `json:"concurrency,omitempty"`
	Verification []string `json:"verification,omitempty"`
	Invocation   string   `json:"implementation_invocation,omitempty"`
	WorkerAgent  string   `json:"worker_agent,omitempty"`
	WorkerModel  string   `json:"worker_model,omitempty"`
	ManagerModel string   `json:"manager_model,omitempty"`
}

type State struct {
	Version        int              `json:"version"`
	RunID          string           `json:"run_id"`
	CreatedAt      string           `json:"created_at"`
	UpdatedAt      string           `json:"updated_at"`
	Login          string           `json:"login,omitempty"`
	TargetHead     string           `json:"target_head,omitempty"`
	StopRequested  bool             `json:"stop_requested,omitempty"`
	CompletedAt    string           `json:"completed_at,omitempty"`
	NextEventSeq   int64            `json:"next_event_seq,omitempty"`
	QueueCondition string           `json:"queue_condition,omitempty"`
	Config         Config           `json:"config"`
	Items          map[string]*Item `json:"items"`
}

type Item struct {
	Number            int       `json:"number"`
	NodeID            string    `json:"node_id,omitempty"`
	Title             string    `json:"title"`
	Body              string    `json:"body"`
	URL               string    `json:"url"`
	State             string    `json:"tracker_state"`
	Labels            []string  `json:"labels,omitempty"`
	Assignees         []string  `json:"assignees,omitempty"`
	Blockers          []Blocker `json:"blockers,omitempty"`
	Status            string    `json:"status"`
	Branch            string    `json:"branch,omitempty"`
	Worktree          string    `json:"worktree,omitempty"`
	BaseHead          string    `json:"base_head,omitempty"`
	Managed           bool      `json:"managed,omitempty"`
	Worker            *Worker   `json:"worker,omitempty"`
	PR                *Pull     `json:"pull_request,omitempty"`
	Verification      []Check   `json:"verification,omitempty"`
	SeenFeedback      []string  `json:"seen_feedback,omitempty"`
	Pending           []Request `json:"pending_requests,omitempty"`
	Error             string    `json:"error,omitempty"`
	UpdatedAt         string    `json:"updated_at,omitempty"`
	CheckoutRemovedAt string    `json:"checkout_removed_at,omitempty"`
}

type Blocker struct {
	Number     int    `json:"number"`
	State      string `json:"state"`
	Repository string `json:"repository,omitempty"`
	Integrated bool   `json:"integrated"`
	Reason     string `json:"reason,omitempty"`
}

type Worker struct {
	Token       string   `json:"token"`
	PID         int      `json:"pid"`
	PIDStart    string   `json:"pid_start,omitempty"`
	StartedAt   string   `json:"started_at"`
	PromptPath  string   `json:"prompt_path"`
	LogPath     string   `json:"log_path"`
	LastMessage string   `json:"last_message_path,omitempty"`
	ExitPath    string   `json:"exit_path"`
	Attempt     int      `json:"attempt"`
	RequestIDs  []string `json:"request_ids,omitempty"`
}

type Pull struct {
	Number     int    `json:"number"`
	URL        string `json:"url"`
	State      string `json:"state"`
	Draft      bool   `json:"draft"`
	HeadOID    string `json:"head_oid,omitempty"`
	MergeOID   string `json:"merge_oid,omitempty"`
	MergeState string `json:"merge_state,omitempty"`
	Checks     string `json:"checks,omitempty"`
}

type Check struct {
	Command  string `json:"command"`
	ExitCode int    `json:"exit_code"`
	Output   string `json:"output,omitempty"`
}

type Request struct {
	ID        string `json:"id"`
	Action    string `json:"action"`
	Message   string `json:"message,omitempty"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
	Source    string `json:"source,omitempty"`
	URL       string `json:"url,omitempty"`
	Path      string `json:"path,omitempty"`
	Line      int    `json:"line,omitempty"`
	HeadOID   string `json:"head_oid,omitempty"`
	ThreadID  string `json:"thread_id,omitempty"`
	ReplySent bool   `json:"reply_sent,omitempty"`
}

type event struct {
	Seq      int64  `json:"seq"`
	At       string `json:"at"`
	Type     string `json:"type"`
	Priority string `json:"priority"`
	Item     int    `json:"item,omitempty"`
	Message  string `json:"message,omitempty"`
}

type feedbackInput struct {
	ID, Source, RemoteID, ThreadID, Body, Author, URL, Path, HeadOID string
	Line                                                             int
	Inline, RequestedChanges                                         bool
}

type reviewSnapshot struct {
	Review   *Pull
	Feedback []feedbackInput
	Failures []string
}

var printSupervisorEvents bool

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	var err error
	switch os.Args[1] {
	case "init":
		err = cmdInit(os.Args[2:])
	case "sync":
		err = withStateCommand(os.Args[2:], syncState)
	case "dispatch":
		err = withStateCommand(os.Args[2:], dispatch)
	case "supervise":
		err = cmdSupervise(os.Args[2:])
	case "status":
		err = cmdStatus(os.Args[2:])
	case "events":
		err = cmdEvents(os.Args[2:])
	case "inbox":
		err = cmdInbox(os.Args[2:])
	case "stop":
		err = cmdStop(os.Args[2:])
	case "publish":
		err = cmdPublish(os.Args[2:])
	case "_worker":
		err = cmdWorker(os.Args[2:])
	case "help", "-h", "--help":
		usage()
	default:
		err = fmt.Errorf("unknown command %q", os.Args[1])
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "ticket-orchestrator:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: ticket-orchestrator <command> [options]

commands:
  init       create durable run state
  sync       refresh target, issues, reviews, checks, and feedback
  dispatch   launch ready workers up to the concurrency cap
  supervise  run sync/reap/dispatch loop (use --once for one cycle)
  status     print a compact run table
  events     print important progress events (optionally follow)
  inbox      add feedback, stop, or resume request for an item
  publish    verify and publish one completed item
  stop       request supervisor stop, optionally stopping workers`)
	os.Exit(2)
}

func cmdInit(args []string) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	repo := fs.String("repo", ".", "target repository")
	source := fs.String("source", "github", "work source (github)")
	harness := fs.String("harness", "", "worker harness (codex or opencode); discovered when unambiguous")
	target := fs.String("target", "", "target branch; defaults to the repository default branch")
	remote := fs.String("remote", "origin", "Git remote")
	label := fs.String("ready-label", "", "label that explicitly makes an issue schedulable")
	concurrency := fs.Int("concurrency", 3, "maximum active workers")
	poll := fs.Int("poll-seconds", 30, "supervisor polling interval")
	stall := fs.Int("stall-minutes", 20, "worker log stall threshold")
	worktrees := fs.String("worktree-root", "", "run-owned worktree parent")
	stateDir := fs.String("state-dir", "", "durable state parent")
	runID := fs.String("run-id", "", "stable run identifier")
	agent := fs.String("worker-agent", "worker", "OpenCode agent name")
	workerModel := fs.String("worker-model", "", "explicit worker model")
	managerModel := fs.String("manager-model", "", "recorded manager model")
	invocation := fs.String("implementation-invocation", "$implement", "worker prompt prefix")
	projectConfig := fs.String("project-config", ".github/implement-tickets.json", "repository-owned orchestration contract")
	publish := fs.Bool("publish", false, "push and create draft PRs after verification")
	issueQuery := fs.String("issue-query", "", "GitHub issue search query selecting this run")
	var launcher, verify stringList
	var issues intList
	var issueRanges stringList
	fs.Var(&launcher, "launcher-arg", "extra harness argument (repeatable)")
	fs.Var(&verify, "verify-command", "trusted verification shell command (repeatable)")
	fs.Var(&issues, "issue", "specific issue number (repeatable)")
	fs.Var(&issueRanges, "issue-range", "inclusive issue range such as 10-20 (repeatable)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	explicit := map[string]bool{}
	fs.Visit(func(f *flag.Flag) { explicit[f.Name] = true })
	for _, value := range issueRanges {
		expanded, rangeErr := parseIssueRange(value)
		if rangeErr != nil {
			return rangeErr
		}
		issues = append(issues, expanded...)
	}
	issues = uniqueInts(issues)
	if *issueQuery != "" && len(issues) > 0 {
		return errors.New("--issue-query cannot be combined with --issue or --issue-range")
	}
	absRepo, err := filepath.Abs(*repo)
	if err != nil {
		return err
	}
	root, err := git(absRepo, "rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("resolve repository: %w", err)
	}
	absRepo = strings.TrimSpace(root)
	setup, err := loadProjectSetup(absRepo, *projectConfig)
	if err != nil {
		return err
	}
	if setup != nil {
		if !explicit["source"] && setup.Source != "" {
			*source = setup.Source
		}
		if !explicit["target"] && setup.Target != "" {
			*target = setup.Target
		}
		if !explicit["remote"] && setup.Remote != "" {
			*remote = setup.Remote
		}
		if !explicit["ready-label"] && setup.ReadyLabel != "" {
			*label = setup.ReadyLabel
		}
		if !explicit["harness"] && setup.Harness != "" {
			*harness = setup.Harness
		}
		if !explicit["concurrency"] && setup.Concurrency != 0 {
			*concurrency = setup.Concurrency
		}
		if !explicit["verify-command"] && len(setup.Verification) > 0 {
			verify = append(verify, setup.Verification...)
		}
		if !explicit["implementation-invocation"] && setup.Invocation != "" {
			*invocation = setup.Invocation
		}
		if !explicit["worker-agent"] && setup.WorkerAgent != "" {
			*agent = setup.WorkerAgent
		}
		if !explicit["worker-model"] && setup.WorkerModel != "" {
			*workerModel = setup.WorkerModel
		}
		if !explicit["manager-model"] && setup.ManagerModel != "" {
			*managerModel = setup.ManagerModel
		}
	}
	if *source != "github" && *source != "gitlab" && *source != "jira" {
		return fmt.Errorf("unsupported work source %q; choose github, gitlab, or jira", *source)
	}
	if *source == "jira" {
		if _, lookErr := exec.LookPath("acli"); lookErr != nil {
			return errors.New("Jira adapter is beta and not executable yet; install/authenticate Atlassian acli for future adapter testing")
		}
		return errors.New("Jira adapter is beta: acli was found, but repository/review linkage and polling require live contract testing before workers may launch")
	}
	if *harness == "opencode" && contains(launcher, "--auto") {
		return errors.New("--launcher-arg --auto is forbidden for headless OpenCode workers")
	}
	if *concurrency < 1 || *concurrency > 32 {
		return errors.New("concurrency must be between 1 and 32")
	}
	dirty, err := git(absRepo, "status", "--porcelain")
	if err != nil {
		return err
	}
	if strings.TrimSpace(dirty) != "" {
		return errors.New("primary worktree has uncommitted changes; preserve or commit them before init")
	}
	remoteURL, err := git(absRepo, "remote", "get-url", *remote)
	if err != nil {
		return fmt.Errorf("read remote: %w", err)
	}
	detectedSource, slug, err := parseRemote(strings.TrimSpace(remoteURL))
	if err != nil {
		return err
	}
	if !explicit["source"] {
		*source = detectedSource
	} else if *source != detectedSource {
		return fmt.Errorf("--source %s does not match remote provider %s", *source, detectedSource)
	}
	loginOut, err := sourceLogin(*source)
	if err != nil {
		return fmt.Errorf("read GitHub login: %w", err)
	}
	projectText, err := trackedProjectText(absRepo)
	if err != nil {
		return fmt.Errorf("read project setup documentation: %w", err)
	}
	labels, err := sourceLabels(*source, slug)
	if err != nil {
		return err
	}
	var questions []string
	if *target == "" {
		out, targetErr := sourceDefaultBranch(*source, slug)
		if targetErr == nil {
			*target = strings.TrimSpace(out)
		}
		if *target == "" {
			questions = append(questions, "Which branch is the integration target? Pass --target <branch>.")
		}
	}
	if *label == "" {
		*label = discoverReadyLabel(projectText, labels)
		if *label == "" {
			questions = append(questions, fmt.Sprintf("Which GitHub label makes an issue ready for an agent? Pass --ready-label <label>. Available labels: %s.", printable(labels)))
		}
	} else if !contains(labels, *label) {
		questions = append(questions, fmt.Sprintf("The ready label %q does not exist. Create it or choose one of: %s.", *label, printable(labels)))
	}
	if *harness == "" {
		*harness = discoverHarness()
		if *harness == "" {
			questions = append(questions, "Which worker harness should run: codex or opencode? Pass --harness <name>.")
		}
	} else if *harness != "codex" && *harness != "opencode" {
		questions = append(questions, fmt.Sprintf("Worker harness %q is unsupported; choose codex or opencode.", *harness))
	}
	if len(verify) == 0 {
		verify = discoverVerification(absRepo)
		if len(verify) == 0 {
			questions = append(questions, "Which independent verification command proves each worker result? Pass --verify-command <command> (repeatable).")
		}
	}
	if *label != "" && contains(labels, *label) {
		candidates, issueErr := sourceCandidateNumbers(*source, slug, *label, issues, *issueQuery)
		if issueErr != nil {
			return fmt.Errorf("inspect selected issues: %w", issueErr)
		}
		if len(candidates) == 0 {
			questions = append(questions, fmt.Sprintf("The run selection contains no issue carrying label %q. Confirm the issue selection/query and label the intended stories before initialization.", *label))
		} else {
			invalid, validationErr := invalidCandidates(*source, slug, *label, candidates)
			if validationErr != nil {
				return validationErr
			}
			if len(invalid) > 0 {
				questions = append(questions, fmt.Sprintf("Selected issues %v are not open with label %q. Update the stories or change the run selection.", invalid, *label))
			} else {
				issues = candidates
			}
		}
	}
	if len(questions) > 0 {
		return setupError(questions)
	}
	if _, err := run("", nil, harnessBinary(*harness), "--version"); err != nil {
		return fmt.Errorf("%s preflight: %w", *harness, err)
	}
	if *runID == "" {
		*runID = time.Now().UTC().Format("20060102T150405Z")
	}
	if *stateDir == "" {
		base := os.Getenv("XDG_STATE_HOME")
		if base == "" {
			home, _ := os.UserHomeDir()
			base = filepath.Join(home, ".local", "state")
		}
		h := sha256.Sum256([]byte(absRepo))
		*stateDir = filepath.Join(base, "agent-orchestration", filepath.Base(absRepo)+"-"+hex.EncodeToString(h[:4]))
	}
	runDir := filepath.Join(*stateDir, *runID)
	if *worktrees == "" {
		*worktrees = filepath.Join(*stateDir, "worktrees", *runID)
	}
	for _, p := range []string{runDir, filepath.Join(runDir, "workers"), *worktrees} {
		if err := os.MkdirAll(p, 0o700); err != nil {
			return err
		}
	}
	s := &State{
		Version: stateVersion, RunID: *runID, CreatedAt: now(), UpdatedAt: now(),
		Login: strings.TrimSpace(loginOut), Items: map[string]*Item{},
		Config: Config{Repo: absRepo, RepoSlug: slug, Source: *source, Harness: *harness,
			Target: *target, Remote: *remote, ReadyLabel: *label, Concurrency: *concurrency,
			PollSeconds: *poll, StallMinutes: *stall, WorktreeRoot: *worktrees, RunDir: runDir,
			WorkerAgent: *agent, WorkerModel: *workerModel, ManagerModel: *managerModel,
			LauncherArgs: launcher, VerifyCommands: verify, Invocation: *invocation, Publish: *publish,
			IssueNumbers: issues, IssueQuery: *issueQuery},
	}
	statePath := filepath.Join(runDir, "state.json")
	appendEvent(s, event{At: now(), Type: "run_initialized", Message: slug})
	if err := saveState(statePath, s); err != nil {
		return err
	}
	fmt.Println(statePath)
	return nil
}

func harnessBinary(h string) string {
	if h == "opencode" {
		return "opencode"
	}
	return "codex"
}

func loadProjectSetup(repo, rel string) (*ProjectSetup, error) {
	if filepath.IsAbs(rel) || rel == "" || strings.HasPrefix(filepath.Clean(rel), "..") {
		return nil, errors.New("--project-config must be a repository-relative path")
	}
	path := filepath.Join(repo, rel)
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read project orchestration contract: %w", err)
	}
	var setup ProjectSetup
	decoder := json.NewDecoder(strings.NewReader(string(b)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&setup); err != nil {
		return nil, fmt.Errorf("decode %s: %w", rel, err)
	}
	if setup.Version != 1 {
		return nil, fmt.Errorf("%s has unsupported version %d", rel, setup.Version)
	}
	return &setup, nil
}

func setupError(questions []string) error {
	var b strings.Builder
	b.WriteString("setup is incomplete; answer these questions before workers are launched:")
	for _, q := range questions {
		b.WriteString("\n- ")
		b.WriteString(q)
	}
	return errors.New(b.String())
}

func githubLabels(slug string) ([]string, error) {
	out, err := run("", nil, "gh", "label", "list", "--repo", slug, "--limit", "100", "--json", "name")
	if err != nil {
		return nil, fmt.Errorf("discover GitHub labels: %w", err)
	}
	var rows []struct{ Name string }
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		return nil, err
	}
	labels := make([]string, 0, len(rows))
	for _, row := range rows {
		labels = append(labels, row.Name)
	}
	sort.Strings(labels)
	return labels, nil
}

func sourceLogin(source string) (string, error) {
	if source == "gitlab" {
		out, err := run("", nil, "glab", "api", "user")
		if err != nil {
			return "", err
		}
		var user struct{ Username string }
		if err := json.Unmarshal([]byte(out), &user); err != nil {
			return "", err
		}
		return user.Username, nil
	}
	return run("", nil, "gh", "api", "user", "--jq", ".login")
}

func sourceLabels(source, slug string) ([]string, error) {
	if source == "github" {
		return githubLabels(slug)
	}
	out, err := run("", nil, "glab", "api", "projects/"+url.PathEscape(slug)+"/labels?per_page=100")
	if err != nil {
		return nil, fmt.Errorf("discover GitLab labels: %w", err)
	}
	var rows []struct{ Name string }
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		return nil, err
	}
	labels := make([]string, 0, len(rows))
	for _, row := range rows {
		labels = append(labels, row.Name)
	}
	sort.Strings(labels)
	return labels, nil
}

func sourceDefaultBranch(source, slug string) (string, error) {
	if source == "gitlab" {
		out, err := run("", nil, "glab", "api", "projects/"+url.PathEscape(slug))
		if err != nil {
			return "", err
		}
		var project struct {
			DefaultBranch string `json:"default_branch"`
		}
		if err := json.Unmarshal([]byte(out), &project); err != nil {
			return "", err
		}
		return project.DefaultBranch, nil
	}
	return run("", nil, "gh", "repo", "view", slug, "--json", "defaultBranchRef", "--jq", ".defaultBranchRef.name")
}

func trackedProjectText(repo string) (string, error) {
	out, err := git(repo, "ls-files")
	if err != nil {
		return "", err
	}
	var b strings.Builder
	for _, rel := range strings.Split(out, "\n") {
		lower := strings.ToLower(rel)
		if rel == "" || !(strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".json") || strings.HasSuffix(lower, ".yaml") || strings.HasSuffix(lower, ".yml") || strings.HasSuffix(lower, ".toml") || filepath.Base(lower) == "makefile") {
			continue
		}
		info, statErr := os.Stat(filepath.Join(repo, rel))
		if statErr != nil || info.Size() > 1<<20 {
			continue
		}
		content, readErr := os.ReadFile(filepath.Join(repo, rel))
		if readErr != nil {
			return "", readErr
		}
		b.WriteString("\n")
		b.Write(content)
	}
	return strings.ToLower(b.String()), nil
}

func discoverReadyLabel(projectText string, labels []string) string {
	var matches []string
	for _, label := range labels {
		lower := strings.ToLower(label)
		looksReady := strings.Contains(lower, "ready") && (strings.Contains(lower, "agent") || strings.Contains(lower, "status") || lower == "ready")
		if looksReady && strings.Contains(projectText, lower) {
			matches = append(matches, label)
		}
	}
	if len(matches) == 1 {
		return matches[0]
	}
	return ""
}

func discoverHarness() string {
	_, codexErr := exec.LookPath("codex")
	_, openCodeErr := exec.LookPath("opencode")
	if codexErr == nil && openCodeErr != nil {
		return "codex"
	}
	if openCodeErr == nil && codexErr != nil {
		return "opencode"
	}
	return ""
}

func discoverVerification(repo string) []string {
	if fileContains(filepath.Join(repo, "Makefile"), "test:") {
		return []string{"make test"}
	}
	if _, err := os.Stat(filepath.Join(repo, "go.mod")); err == nil {
		return []string{"go test ./..."}
	}
	if fileContains(filepath.Join(repo, "package.json"), `"test"`) {
		return []string{"npm test"}
	}
	return nil
}

func fileContains(path, needle string) bool {
	b, err := os.ReadFile(path)
	return err == nil && strings.Contains(string(b), needle)
}

func printable(values []string) string {
	if len(values) == 0 {
		return "(none)"
	}
	return strings.Join(values, ", ")
}

func githubCandidateNumbers(slug, label string, selected []int, query string) ([]int, error) {
	if len(selected) > 0 {
		return uniqueInts(selected), nil
	}
	args := []string{"issue", "list", "--repo", slug, "--state", "open", "--label", label, "--limit", "1000", "--json", "number"}
	if query != "" {
		args = append(args, "--search", query)
	}
	out, err := run("", nil, "gh", args...)
	if err != nil {
		return nil, err
	}
	var rows []struct{ Number int }
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		return nil, err
	}
	numbers := make([]int, 0, len(rows))
	for _, row := range rows {
		numbers = append(numbers, row.Number)
	}
	return uniqueInts(numbers), nil
}

func sourceCandidateNumbers(source, slug, label string, selected []int, query string) ([]int, error) {
	if source == "github" {
		return githubCandidateNumbers(slug, label, selected, query)
	}
	if len(selected) > 0 {
		return uniqueInts(selected), nil
	}
	endpoint := "projects/" + url.PathEscape(slug) + "/issues?state=opened&per_page=100&labels=" + url.QueryEscape(label)
	if query != "" {
		endpoint += "&search=" + url.QueryEscape(query)
	}
	out, err := run("", nil, "glab", "api", endpoint)
	if err != nil {
		return nil, err
	}
	var rows []struct {
		IID int `json:"iid"`
	}
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		return nil, err
	}
	numbers := make([]int, 0, len(rows))
	for _, row := range rows {
		numbers = append(numbers, row.IID)
	}
	return uniqueInts(numbers), nil
}

func invalidCandidates(source, slug, label string, numbers []int) ([]int, error) {
	var invalid []int
	for _, number := range numbers {
		var out string
		var err error
		if source == "gitlab" {
			out, err = run("", nil, "glab", "api", "projects/"+url.PathEscape(slug)+"/issues/"+strconv.Itoa(number))
		} else {
			out, err = run("", nil, "gh", "issue", "view", strconv.Itoa(number), "--repo", slug, "--json", "state,labels")
		}
		if err != nil {
			return nil, fmt.Errorf("inspect selected issue #%d: %w", number, err)
		}
		hasLabel, open, decodeErr := candidateState(source, []byte(out), label)
		if decodeErr != nil {
			return nil, decodeErr
		}
		if !open || !hasLabel {
			invalid = append(invalid, number)
		}
	}
	return invalid, nil
}

func candidateState(source string, raw []byte, label string) (bool, bool, error) {
	if source == "gitlab" {
		var row struct {
			State  string
			Labels []string
		}
		if err := json.Unmarshal(raw, &row); err != nil {
			return false, false, err
		}
		return contains(row.Labels, label), row.State == "opened", nil
	}
	var row struct {
		State  string
		Labels []struct{ Name string }
	}
	if err := json.Unmarshal(raw, &row); err != nil {
		return false, false, err
	}
	labels := make([]string, 0, len(row.Labels))
	for _, itemLabel := range row.Labels {
		labels = append(labels, itemLabel.Name)
	}
	return contains(labels, label), row.State == "OPEN", nil
}

func parseIssueRange(value string) ([]int, error) {
	parts := strings.Split(value, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("issue range must use START-END: %q", value)
	}
	start, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	end, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil || start < 1 || end < start || end-start > 999 {
		return nil, fmt.Errorf("invalid or overly broad issue range %q", value)
	}
	values := make([]int, 0, end-start+1)
	for n := start; n <= end; n++ {
		values = append(values, n)
	}
	return values, nil
}

func uniqueInts(values []int) []int {
	seen := map[int]bool{}
	result := make([]int, 0, len(values))
	for _, value := range values {
		if value > 0 && !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	sort.Ints(result)
	return result
}

type stateFn func(string, *State) error

func withStateCommand(args []string, fn stateFn) error {
	fs := flag.NewFlagSet("state command", flag.ContinueOnError)
	path := fs.String("state", "", "path to state.json")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *path == "" {
		return errors.New("--state is required")
	}
	return locked(*path, func(s *State) error {
		if err := fn(*path, s); err != nil {
			return err
		}
		return saveState(*path, s)
	})
}

func cmdSupervise(args []string) error {
	fs := flag.NewFlagSet("supervise", flag.ContinueOnError)
	path := fs.String("state", "", "path to state.json")
	once := fs.Bool("once", false, "run one cycle")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *path == "" {
		return errors.New("--state is required")
	}
	owner, err := acquireSupervisorOwner(*path)
	if err != nil {
		return err
	}
	defer releaseSupervisorOwner(owner)
	printSupervisorEvents = true
	defer func() { printSupervisorEvents = false }()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	for {
		err := locked(*path, func(s *State) error {
			if err := processInbox(s); err != nil {
				return err
			}
			if s.StopRequested {
				return errStop
			}
			syncErr := syncState(*path, s)
			if syncErr != nil {
				appendEvent(s, event{At: now(), Type: "sync_failed", Message: syncErr.Error()})
			}
			if err := reapWorkers(*path, s); err != nil {
				return err
			}
			if syncErr == nil {
				if err := dispatch(*path, s); err != nil {
					return err
				}
			}
			cleanupCompletedRun(s)
			emitQueueCondition(s)
			complete := runComplete(s)
			if complete && s.CompletedAt == "" {
				s.CompletedAt = now()
				appendEvent(s, event{At: s.CompletedAt, Type: "run_completed", Message: terminalSummary(s)})
			}
			if err := saveState(*path, s); err != nil {
				return err
			}
			if complete {
				return errRunComplete
			}
			return nil
		})
		if errors.Is(err, errStop) || errors.Is(err, errRunComplete) {
			return nil
		}
		if err != nil {
			return err
		}
		if *once {
			return nil
		}
		s, err := loadState(*path)
		if err != nil {
			return err
		}
		select {
		case <-stop:
			return nil
		case <-time.After(time.Duration(s.Config.PollSeconds) * time.Second):
		}
	}
}

func acquireSupervisorOwner(statePath string) (*os.File, error) {
	path := statePath + ".supervisor.lock"
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	lock, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		lock.Close()
		return nil, errors.New("another supervisor owns the run for its full lifetime")
	}
	return lock, nil
}

func releaseSupervisorOwner(lock *os.File) {
	if lock == nil {
		return
	}
	_ = syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
	_ = lock.Close()
}

func emitQueueCondition(s *State) {
	active, runnable, input, waiting, done := 0, 0, 0, 0, 0
	for _, i := range s.Items {
		switch i.Status {
		case "running":
			active++
		case "ready", "feedback", "ci_failed", "reconcile_pending":
			runnable++
		case "needs_input", "needs_attention", "failed", "publication_failed":
			input++
		case "integrated", "review_closed", "skipped":
			done++
		default:
			waiting++
		}
	}
	condition := fmt.Sprintf("active=%d runnable=%d input=%d waiting=%d done=%d", active, runnable, input, waiting, done)
	if condition == s.QueueCondition {
		return
	}
	s.QueueCondition = condition
	typeName := "queue_progress"
	if active == 0 && runnable == 0 && input > 0 {
		typeName = "queue_blocked"
	} else if active == 0 && runnable == 0 && waiting == 0 {
		typeName = "queue_idle"
	}
	appendEvent(s, event{At: now(), Type: typeName, Message: condition})
}

func cleanupCompletedRun(s *State) {
	if !allItemsTerminal(s) {
		return
	}
	stopping := false
	for _, i := range s.Items {
		if i.Worker != nil && processAlive(i.Worker.PID, i.Worker.PIDStart) {
			_ = syscall.Kill(-i.Worker.PID, syscall.SIGTERM)
			appendEvent(s, event{At: now(), Type: "terminal_worker_stopped", Item: i.Number, Message: i.Worker.Token})
			stopping = true
		}
	}
	if stopping {
		return
	}
	for _, i := range s.Items {
		if i.CheckoutRemovedAt != "" || !i.Managed || i.Worktree == "" {
			continue
		}
		if err := removeManagedCheckout(s, i); err != nil {
			i.Error = "terminal checkout cleanup: " + err.Error()
			appendEvent(s, event{At: now(), Type: "cleanup_failed", Item: i.Number, Message: i.Error})
			continue
		}
		i.CheckoutRemovedAt, i.Error = now(), ""
		appendEvent(s, event{At: i.CheckoutRemovedAt, Type: "checkout_removed", Item: i.Number, Message: i.Worktree})
	}
	if err := os.Remove(s.Config.WorktreeRoot); err == nil {
		appendEvent(s, event{At: now(), Type: "worktree_root_removed", Message: s.Config.WorktreeRoot})
	} else if !os.IsNotExist(err) {
		appendEvent(s, event{At: now(), Type: "cleanup_failed", Message: "remove empty worktree root: " + err.Error()})
	}
}

func removeManagedCheckout(s *State, i *Item) error {
	root, err := filepath.Abs(s.Config.WorktreeRoot)
	if err != nil {
		return err
	}
	checkout, err := filepath.Abs(i.Worktree)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(root, checkout)
	if err != nil || rel == "." || rel == "" || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("refusing checkout outside run-owned root: %s", checkout)
	}
	info, err := os.Stat(filepath.Join(checkout, ".git"))
	if err != nil || !info.IsDir() {
		return errors.New("managed checkout does not contain self-contained Git metadata")
	}
	status, err := git(checkout, "status", "--porcelain")
	if err != nil || strings.TrimSpace(status) != "" {
		return errors.New("managed checkout is not clean")
	}
	head, err := git(checkout, "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	if i.PR == nil || i.PR.HeadOID == "" || strings.TrimSpace(head) != i.PR.HeadOID {
		return errors.New("managed checkout head is not the published review head")
	}
	return os.RemoveAll(checkout)
}

func terminalItem(i *Item) bool {
	return i.Status == "integrated" || i.Status == "review_closed" || i.Status == "skipped"
}

func allItemsTerminal(s *State) bool {
	if len(s.Items) == 0 {
		return false
	}
	for _, i := range s.Items {
		if !terminalItem(i) {
			return false
		}
	}
	return true
}

func runComplete(s *State) bool {
	if !allItemsTerminal(s) {
		return false
	}
	for _, i := range s.Items {
		if i.Managed && i.Worktree != "" && i.CheckoutRemovedAt == "" {
			return false
		}
		if i.Worker != nil && processAlive(i.Worker.PID, i.Worker.PIDStart) {
			return false
		}
	}
	return true
}

func terminalSummary(s *State) string {
	counts := map[string]int{}
	for _, i := range s.Items {
		counts[i.Status]++
	}
	return fmt.Sprintf("integrated=%d review_closed=%d skipped=%d", counts["integrated"], counts["review_closed"], counts["skipped"])
}

var errStop = errors.New("stop requested")
var errRunComplete = errors.New("run complete")

func syncState(_ string, s *State) error {
	if _, err := run("", nil, "git", "-C", s.Config.Repo, "fetch", "--prune", s.Config.Remote, s.Config.Target); err != nil {
		return fmt.Errorf("fetch target: %w", err)
	}
	head, err := git(s.Config.Repo, "rev-parse", s.Config.Remote+"/"+s.Config.Target)
	if err != nil {
		return err
	}
	s.TargetHead = strings.TrimSpace(head)
	issues, err := sourceIssues(s)
	if err != nil {
		return err
	}
	seen := map[string]bool{}
	for _, fresh := range issues {
		key := strconv.Itoa(fresh.Number)
		seen[key] = true
		old := s.Items[key]
		if old != nil {
			fresh.Status, fresh.Branch, fresh.Worktree, fresh.BaseHead, fresh.Managed = old.Status, old.Branch, old.Worktree, old.BaseHead, old.Managed
			fresh.Worker, fresh.PR, fresh.Verification, fresh.SeenFeedback, fresh.Pending, fresh.Error = old.Worker, old.PR, old.Verification, old.SeenFeedback, old.Pending, old.Error
		}
		if fresh.Status == "" {
			fresh.Status = "waiting"
		}
		if fresh.Worker == nil && fresh.PR == nil && (fresh.Status == "waiting" || fresh.Status == "ready") {
			if ready(fresh, s.Login, s.Config.ReadyLabel) {
				fresh.Status = "ready"
			} else {
				fresh.Status = "waiting"
			}
		}
		if err := syncPull(s, fresh); err != nil {
			fresh.Error = err.Error()
		}
		fresh.UpdatedAt = now()
		s.Items[key] = fresh
	}
	for key, item := range s.Items {
		if !seen[key] && item.State == "OPEN" {
			if err := refreshMissingItem(s, item); err != nil {
				item.Error = err.Error()
			}
		}
	}
	appendEvent(s, event{At: now(), Type: "sync", Message: fmt.Sprintf("%d open issues", len(issues))})
	return nil
}

type githubIssueNode struct {
	ID, Title, Body, URL, State string
	Number                      int
	Labels                      struct{ Nodes []struct{ Name string } }
	Assignees                   struct{ Nodes []struct{ Login string } }
	BlockedBy                   struct {
		Nodes []struct {
			Number     int
			State      string
			Repository struct{ NameWithOwner string }
			Closing    struct {
				Nodes []struct {
					Number      int
					BaseRefName string
					MergedAt    string
					MergeCommit *struct{ OID string }
				} `json:"nodes"`
			} `json:"closedByPullRequestsReferences"`
		} `json:"nodes"`
	} `json:"blockedBy"`
}

func githubIssues(s *State) ([]*Item, error) {
	numbers, err := githubCandidateNumbers(s.Config.RepoSlug, s.Config.ReadyLabel, s.Config.IssueNumbers, s.Config.IssueQuery)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(s.Config.RepoSlug, "/")
	query := `query($owner:String!,$name:String!,$number:Int!){repository(owner:$owner,name:$name){issue(number:$number){id number title body url state labels(first:50){nodes{name}} assignees(first:20){nodes{login}} blockedBy(first:50){nodes{number state url repository{nameWithOwner} closedByPullRequestsReferences(first:20){nodes{number mergedAt baseRefName mergeCommit{oid}}}}}}}}`
	result := make([]*Item, 0, len(numbers))
	for _, number := range numbers {
		out, queryErr := run("", nil, "gh", "api", "graphql", "-f", "query="+query, "-F", "owner="+parts[0], "-F", "name="+parts[1], "-F", "number="+strconv.Itoa(number))
		if queryErr != nil {
			return nil, fmt.Errorf("query GitHub issue #%d: %w", number, queryErr)
		}
		var response struct {
			Data struct {
				Repository struct {
					Issue githubIssueNode `json:"issue"`
				} `json:"repository"`
			} `json:"data"`
		}
		if err := json.Unmarshal([]byte(out), &response); err != nil {
			return nil, fmt.Errorf("decode GitHub issue #%d: %w", number, err)
		}
		n := response.Data.Repository.Issue
		if n.Number == 0 {
			return nil, fmt.Errorf("selected issue #%d does not exist", number)
		}
		i := &Item{Number: n.Number, NodeID: n.ID, Title: n.Title, Body: n.Body, URL: n.URL, State: n.State}
		for _, label := range n.Labels.Nodes {
			i.Labels = append(i.Labels, label.Name)
		}
		for _, a := range n.Assignees.Nodes {
			i.Assignees = append(i.Assignees, a.Login)
		}
		for _, b := range n.BlockedBy.Nodes {
			blocker := Blocker{Number: b.Number, State: b.State, Repository: b.Repository.NameWithOwner}
			for _, p := range b.Closing.Nodes {
				if b.State == "CLOSED" && p.MergedAt != "" && p.BaseRefName == s.Config.Target && p.MergeCommit != nil && isAncestor(s.Config.Repo, p.MergeCommit.OID, s.Config.Remote+"/"+s.Config.Target) {
					blocker.Integrated = true
					break
				}
			}
			if !blocker.Integrated {
				blocker.Reason = "no merged closing PR is present on the configured target"
			}
			i.Blockers = append(i.Blockers, blocker)
		}
		result = append(result, i)
	}
	return result, nil
}

func sourceIssues(s *State) ([]*Item, error) {
	if s.Config.Source == "gitlab" {
		return gitlabIssues(s)
	}
	return githubIssues(s)
}

func gitlabIssues(s *State) ([]*Item, error) {
	numbers, err := sourceCandidateNumbers("gitlab", s.Config.RepoSlug, s.Config.ReadyLabel, s.Config.IssueNumbers, s.Config.IssueQuery)
	if err != nil {
		return nil, err
	}
	project := "projects/" + url.PathEscape(s.Config.RepoSlug)
	result := make([]*Item, 0, len(numbers))
	for _, number := range numbers {
		out, err := run("", nil, "glab", "api", project+"/issues/"+strconv.Itoa(number))
		if err != nil {
			return nil, fmt.Errorf("query GitLab issue #%d: %w", number, err)
		}
		var row struct {
			ID          int      `json:"id"`
			IID         int      `json:"iid"`
			Title       string   `json:"title"`
			Description string   `json:"description"`
			WebURL      string   `json:"web_url"`
			State       string   `json:"state"`
			Labels      []string `json:"labels"`
			Assignees   []struct {
				Username string `json:"username"`
			} `json:"assignees"`
		}
		if err := json.Unmarshal([]byte(out), &row); err != nil {
			return nil, err
		}
		i := &Item{Number: row.IID, NodeID: strconv.Itoa(row.ID), Title: row.Title, Body: row.Description, URL: row.WebURL, State: normalizedGitLabState(row.State), Labels: row.Labels}
		for _, assignee := range row.Assignees {
			i.Assignees = append(i.Assignees, assignee.Username)
		}
		links, err := gitlabBlockers(s, number, row.Description)
		if err != nil {
			return nil, err
		}
		i.Blockers = links
		result = append(result, i)
	}
	return result, nil
}

func normalizedGitLabState(state string) string {
	if strings.EqualFold(state, "opened") {
		return "OPEN"
	}
	return strings.ToUpper(state)
}

func gitlabBlockers(s *State, number int, description string) ([]Blocker, error) {
	project := "projects/" + url.PathEscape(s.Config.RepoSlug)
	out, err := run("", nil, "glab", "api", project+"/issues/"+strconv.Itoa(number)+"/links?per_page=100")
	if err != nil {
		return nil, fmt.Errorf("query GitLab blockers for #%d: %w", number, err)
	}
	var rows []struct {
		IID      int    `json:"iid"`
		State    string `json:"state"`
		WebURL   string `json:"web_url"`
		LinkType string `json:"link_type"`
	}
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		return nil, err
	}
	var result []Blocker
	for _, row := range rows {
		if row.LinkType != "is_blocked_by" {
			continue
		}
		blocker := Blocker{Number: row.IID, State: strings.ToUpper(row.State), Repository: s.Config.RepoSlug}
		if blocker.State == "CLOSED" {
			blocker.Integrated = gitlabIssueIntegrated(s, row.IID)
		}
		if !blocker.Integrated {
			blocker.Reason = "no merged related merge request is present on the configured target"
		}
		result = append(result, blocker)
	}
	known := make(map[int]bool, len(result))
	for _, blocker := range result {
		known[blocker.Number] = true
	}
	for _, blockerNumber := range explicitBlockerNumbers(description) {
		if known[blockerNumber] {
			continue
		}
		issue, issueErr := run("", nil, "glab", "api", project+"/issues/"+strconv.Itoa(blockerNumber))
		if issueErr != nil {
			return nil, fmt.Errorf("query explicit GitLab blocker #%d: %w", blockerNumber, issueErr)
		}
		var row struct {
			State string `json:"state"`
		}
		if err := json.Unmarshal([]byte(issue), &row); err != nil {
			return nil, err
		}
		blocker := Blocker{Number: blockerNumber, State: strings.ToUpper(row.State), Repository: s.Config.RepoSlug}
		if blocker.State == "CLOSED" {
			blocker.Integrated = gitlabIssueIntegrated(s, blockerNumber)
		}
		if !blocker.Integrated {
			blocker.Reason = "no merged related merge request is present on the configured target"
		}
		result = append(result, blocker)
	}
	return result, nil
}

var issueReference = regexp.MustCompile(`(?m)(?:^|\s)-?\s*#(\d+)\b`)

// explicitBlockerNumbers reads the documented Blocked by section without
// treating incidental issue references elsewhere in the description as edges.
func explicitBlockerNumbers(description string) []int {
	seen := map[int]bool{}
	var numbers []int
	inSection := false
	for _, line := range strings.Split(description, "\n") {
		heading := strings.TrimSpace(line)
		if strings.EqualFold(heading, "## Blocked by") {
			inSection = true
			continue
		}
		if inSection == false {
			continue
		}
		if strings.HasPrefix(heading, "## ") {
			break
		}
		if strings.Contains(strings.ToLower(line), "none") {
			return nil
		}
		for _, match := range issueReference.FindAllStringSubmatch(line, -1) {
			number, err := strconv.Atoi(match[1])
			if err != nil || seen[number] {
				continue
			}
			seen[number] = true
			numbers = append(numbers, number)
		}
	}
	sort.Ints(numbers)
	return numbers
}

func gitlabIssueIntegrated(s *State, number int) bool {
	endpoint := "projects/" + url.PathEscape(s.Config.RepoSlug) + "/issues/" + strconv.Itoa(number) + "/related_merge_requests"
	out, err := run("", nil, "glab", "api", endpoint)
	if err != nil {
		return false
	}
	var rows []struct {
		State          string `json:"state"`
		TargetBranch   string `json:"target_branch"`
		MergeCommitSHA string `json:"merge_commit_sha"`
	}
	if json.Unmarshal([]byte(out), &rows) != nil {
		return false
	}
	for _, row := range rows {
		if row.State == "merged" && row.TargetBranch == s.Config.Target && row.MergeCommitSHA != "" && isAncestor(s.Config.Repo, row.MergeCommitSHA, s.Config.Remote+"/"+s.Config.Target) {
			return true
		}
	}
	return false
}

func refreshMissingItem(s *State, i *Item) error {
	if s.Config.Source == "gitlab" {
		out, err := run("", nil, "glab", "api", "projects/"+url.PathEscape(s.Config.RepoSlug)+"/issues/"+strconv.Itoa(i.Number))
		if err != nil {
			return err
		}
		var row struct {
			State     string
			Labels    []string
			Assignees []struct{ Username string }
		}
		if err := json.Unmarshal([]byte(out), &row); err != nil {
			return err
		}
		i.State, i.Labels, i.Assignees = normalizedGitLabState(row.State), row.Labels, nil
		for _, assignee := range row.Assignees {
			i.Assignees = append(i.Assignees, assignee.Username)
		}
		return nil
	}
	out, err := run("", nil, "gh", "issue", "view", strconv.Itoa(i.Number), "--repo", s.Config.RepoSlug, "--json", "state,labels,assignees")
	if err != nil {
		return fmt.Errorf("refresh issue no longer returned by ready-label query: %w", err)
	}
	var row struct {
		State     string
		Labels    []struct{ Name string }
		Assignees []struct{ Login string }
	}
	if err := json.Unmarshal([]byte(out), &row); err != nil {
		return err
	}
	i.State, i.Labels, i.Assignees = row.State, nil, nil
	for _, label := range row.Labels {
		i.Labels = append(i.Labels, label.Name)
	}
	for _, a := range row.Assignees {
		i.Assignees = append(i.Assignees, a.Login)
	}
	if err := syncPull(s, i); err != nil {
		return err
	}
	if i.State == "CLOSED" {
		if mergedIntoTarget(s, i) {
			i.Status = "integrated"
		} else {
			i.Status = "closed_unverified"
		}
	} else if i.Worker == nil && i.PR == nil && (i.Status == "ready" || i.Status == "waiting") {
		i.Status = "waiting"
	}
	return nil
}

func ready(i *Item, login, label string) bool {
	if i.State != "OPEN" || !contains(i.Labels, label) {
		return false
	}
	if len(i.Assignees) > 0 && !contains(i.Assignees, login) {
		return false
	}
	for _, b := range i.Blockers {
		if !b.Integrated {
			return false
		}
	}
	return true
}

func dispatch(statePath string, s *State) error {
	active := 0
	for _, i := range s.Items {
		if i.Status == "running" && i.Worker != nil && processAlive(i.Worker.PID, i.Worker.PIDStart) {
			active++
		}
	}
	if active >= s.Config.Concurrency {
		return nil
	}
	items := sortedItems(s)
	for _, i := range items {
		if active >= s.Config.Concurrency {
			break
		}
		if i.Status != "ready" && i.Status != "feedback" && i.Status != "ci_failed" && i.Status != "reconcile_pending" {
			continue
		}
		if err := prepareWorktree(s, i); err != nil {
			i.Status, i.Error = "needs_attention", err.Error()
			continue
		}
		if err := claimIssue(s, i); err != nil {
			i.Status, i.Error = "needs_attention", err.Error()
			continue
		}
		if err := launchWorker(statePath, s, i); err != nil {
			i.Status, i.Error = "failed", err.Error()
			continue
		}
		active++
	}
	return nil
}

func prepareWorktree(s *State, i *Item) error {
	if i.Branch == "" {
		i.Branch = fmt.Sprintf("issue/%d-%s", i.Number, slugify(i.Title))
	}
	if i.Worktree == "" {
		i.Worktree = filepath.Join(s.Config.WorktreeRoot, fmt.Sprintf("%d-%s", i.Number, slugify(i.Title)))
	}
	if i.Managed {
		info, err := os.Stat(filepath.Join(i.Worktree, ".git"))
		if err != nil {
			return fmt.Errorf("managed worktree is missing: %s", i.Worktree)
		}
		if !info.IsDir() {
			return errors.New("managed checkout uses shared Git metadata; preserve it and start a fresh run with self-contained checkouts")
		}
		return nil
	}
	if _, err := os.Stat(i.Worktree); err == nil {
		return fmt.Errorf("refusing existing unmanaged path %s", i.Worktree)
	} else if !os.IsNotExist(err) {
		return err
	}
	remoteURL, err := git(s.Config.Repo, "remote", "get-url", s.Config.Remote)
	if err != nil {
		return err
	}
	if _, err := run("", nil, "git", "clone", "--no-checkout", "--no-local", s.Config.Repo, i.Worktree); err != nil {
		return err
	}
	for _, key := range []string{"user.name", "user.email"} {
		value, err := git(s.Config.Repo, "config", "--get", key)
		if err != nil || strings.TrimSpace(value) == "" {
			return fmt.Errorf("read manager checkout %s: %w", key, err)
		}
		if _, err := git(i.Worktree, "config", key, strings.TrimSpace(value)); err != nil {
			return fmt.Errorf("set worker checkout %s: %w", key, err)
		}
	}
	if _, err := git(i.Worktree, "remote", "set-url", s.Config.Remote, strings.TrimSpace(remoteURL)); err != nil {
		return err
	}
	if _, err := git(i.Worktree, "fetch", s.Config.Remote, s.Config.Target); err != nil {
		return err
	}
	if _, err := git(i.Worktree, "checkout", "-b", i.Branch, s.Config.Remote+"/"+s.Config.Target); err != nil {
		return err
	}
	_, _ = git(i.Worktree, "branch", "--unset-upstream")
	i.Managed, i.BaseHead = true, s.TargetHead
	return nil
}

func claimIssue(s *State, i *Item) error {
	if contains(i.Assignees, s.Login) {
		return nil
	}
	var err error
	if s.Config.Source == "gitlab" {
		_, err = run("", nil, "glab", "issue", "update", strconv.Itoa(i.Number), "--repo", s.Config.RepoSlug, "--assignee", s.Login)
	} else {
		_, err = run("", nil, "gh", "issue", "edit", strconv.Itoa(i.Number), "--repo", s.Config.RepoSlug, "--add-assignee", "@me")
	}
	if err == nil {
		i.Assignees = append(i.Assignees, s.Login)
	}
	return err
}

func launchWorker(statePath string, s *State, i *Item) error {
	attempt := 1
	if i.Worker != nil {
		attempt = i.Worker.Attempt + 1
	}
	token := fmt.Sprintf("%d-a%d-%d", i.Number, attempt, time.Now().UnixNano())
	dir := filepath.Join(s.Config.RunDir, "workers", token)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	prompt := workerPrompt(s, i)
	promptPath := filepath.Join(dir, "prompt.md")
	if err := os.WriteFile(promptPath, []byte(prompt), 0o600); err != nil {
		return err
	}
	logPath, exitPath := filepath.Join(dir, "events.jsonl"), filepath.Join(dir, "exit-code")
	log, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer log.Close()
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(exe, "_worker", "--state", statePath, "--item", strconv.Itoa(i.Number), "--token", token)
	cmd.Stdout, cmd.Stderr = log, log
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	i.Worker = &Worker{Token: token, StartedAt: now(), PromptPath: promptPath, LogPath: logPath, ExitPath: exitPath, LastMessage: filepath.Join(dir, "last-message.txt"), Attempt: attempt, RequestIDs: queuedRequestIDs(i)}
	i.Status, i.Error = "running", ""
	if err := saveState(statePath, s); err != nil {
		return fmt.Errorf("persist worker ownership before launch: %w", err)
	}
	if err := cmd.Start(); err != nil {
		i.Status, i.Error = "failed", err.Error()
		return err
	}
	i.Worker.PID, i.Worker.PIDStart = cmd.Process.Pid, procStart(cmd.Process.Pid)
	if err := saveState(statePath, s); err != nil {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		return fmt.Errorf("persist launched worker PID: %w", err)
	}
	go func() { _ = cmd.Wait() }()
	appendEvent(s, event{At: now(), Type: "worker_started", Item: i.Number, Message: token})
	return nil
}

func workerPrompt(s *State, i *Item) string {
	var context strings.Builder
	for _, r := range i.Pending {
		if r.Status == "queued" && actionableRequest(r) && r.Message != "" {
			location := r.Source
			if r.Path != "" {
				location = fmt.Sprintf("%s %s:%d", location, r.Path, r.Line)
			}
			fmt.Fprintf(&context, "\nOperator/review input (%s; %s; %s):\n%s\n", r.Action, location, r.URL, r.Message)
		}
	}
	return fmt.Sprintf(`%s

Implement exactly one work item in the existing isolated worktree.

Title: %s

Canonical work item body (verbatim):
---
%s
---

Repository: %s
Target ref at dispatch: %s/%s (%s)
Assigned branch: %s
Assigned worktree: %s
%s
Rules:
- Read and obey every applicable project instruction file.
- Preserve unrelated changes and already integrated behavior.
- Implement only this item and its smallest supporting changes.
- Run focused checks regularly and the full applicable suite at the end.
- Review the final diff and leave the tested changes unstaged for the supervisor.
- Do not modify Git metadata or create commits; the supervisor owns the checkpoint commit.
- Do not query GitHub or another tracker.
- Do not push, create or edit a pull request, merge, close issues, or delete worktrees.
- Stop with a precise handoff when authority or product judgment is required.
- End a successful handoff with exactly one Commit subject: <type>: <summary> line.
  Derive that concise Conventional Commit subject from the actual diff, not the issue number.
- When this worker responds to review feedback, also end with exactly one
  Feedback response: <concise answer> line for the supervisor to post in the original thread.
`, s.Config.Invocation, i.Title, i.Body, s.Config.Repo, s.Config.Remote, s.Config.Target, i.BaseHead, i.Branch, i.Worktree, context.String())
}

func queuedRequestIDs(i *Item) []string {
	var ids []string
	for _, request := range i.Pending {
		if request.Status == "queued" && actionableRequest(request) {
			ids = append(ids, request.ID)
		}
	}
	return ids
}

func actionableRequest(request Request) bool {
	return request.Action == "feedback" || request.Action == "ci_failed" || request.Action == "reconcile_pending"
}

func cmdWorker(args []string) error {
	fs := flag.NewFlagSet("_worker", flag.ContinueOnError)
	path := fs.String("state", "", "state")
	number := fs.Int("item", 0, "item")
	token := fs.String("token", "", "worker token")
	if err := fs.Parse(args); err != nil {
		return err
	}
	s, err := loadState(*path)
	if err != nil {
		return err
	}
	i := s.Items[strconv.Itoa(*number)]
	if i == nil || i.Worker == nil || i.Worker.Token != *token {
		return errors.New("worker ownership token no longer matches state")
	}
	promptBytes, err := os.ReadFile(i.Worker.PromptPath)
	if err != nil {
		return err
	}
	var argv []string
	if s.Config.Harness == "codex" {
		argv = []string{"exec", "--json", "--sandbox", "workspace-write", "-c", `approval_policy="never"`, "--ephemeral", "-C", i.Worktree, "-o", i.Worker.LastMessage}
		if s.Config.WorkerModel != "" {
			argv = append(argv, "--model", s.Config.WorkerModel)
		}
		argv = append(argv, s.Config.LauncherArgs...)
		argv = append(argv, "-")
	} else {
		argv = []string{"run", "--agent", s.Config.WorkerAgent, "--format", "json", "--title", "orchestration-" + strconv.Itoa(i.Number), "--dir", i.Worktree}
		if s.Config.WorkerModel != "" {
			argv = append(argv, "--model", s.Config.WorkerModel)
		}
		argv = append(argv, s.Config.LauncherArgs...)
		argv = append(argv, string(promptBytes))
	}
	cmd := exec.Command(harnessBinary(s.Config.Harness), argv...)
	cmd.Dir, cmd.Stdout, cmd.Stderr = i.Worktree, os.Stdout, os.Stderr
	if s.Config.Harness == "codex" {
		cmd.Stdin = strings.NewReader(string(promptBytes))
	} else {
		handoff, handoffErr := os.OpenFile(i.Worker.LastMessage, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
		if handoffErr != nil {
			return handoffErr
		}
		defer handoff.Close()
		cmd.Stdout = io.MultiWriter(os.Stdout, handoff)
		cmd.Env = scrubWorkerEnv(os.Environ())
	}
	err = cmd.Run()
	code := 0
	if err != nil {
		code = exitCode(err)
	}
	if writeErr := atomicWrite(i.Worker.ExitPath, []byte(strconv.Itoa(code)+"\n"), 0o600); writeErr != nil {
		return writeErr
	}
	return err
}

func reapWorkers(statePath string, s *State) error {
	for _, i := range s.Items {
		if i.Status != "running" || i.Worker == nil {
			continue
		}
		exitBytes, err := os.ReadFile(i.Worker.ExitPath)
		if err != nil {
			if os.IsNotExist(err) {
				if !processAlive(i.Worker.PID, i.Worker.PIDStart) {
					i.Status, i.Error = "failed", "worker disappeared without an exit record"
					appendEvent(s, event{At: now(), Type: "worker_lost", Item: i.Number})
				} else if stalled(i.Worker.LogPath, s.Config.StallMinutes) {
					i.Error = "worker log is stale; inspect before stopping or replacing it"
					appendEvent(s, event{At: now(), Type: "worker_stalled", Item: i.Number})
				}
				continue
			}
			return err
		}
		code, _ := strconv.Atoi(strings.TrimSpace(string(exitBytes)))
		if code != 0 {
			i.Status, i.Error = "failed", fmt.Sprintf("worker exited with code %d", code)
			appendEvent(s, event{At: now(), Type: "worker_failed", Item: i.Number, Message: i.Error})
			continue
		}
		if err := verifyItem(s, i); err != nil {
			i.Error = err.Error()
			if automaticRecoveryAllowed(s, i) {
				queueRequest(s, i, Request{ID: stableID("verification-recovery", strconv.Itoa(i.Number), strconv.Itoa(i.Worker.Attempt)), Action: "resume", Message: "Repair the verification failure without changing scope. Evidence: " + err.Error(), CreatedAt: now(), Status: "queued", Source: "supervisor"})
				i.Status = "feedback"
				appendEvent(s, event{At: now(), Type: "verification_recovery_queued", Item: i.Number, Message: err.Error()})
				continue
			}
			i.Status = "needs_attention"
			appendEvent(s, event{At: now(), Type: "verification_failed", Item: i.Number, Message: err.Error()})
			continue
		}
		if s.Config.Publish {
			if err := publishItem(s, i); err != nil {
				i.Status, i.Error = "publication_failed", err.Error()
				appendEvent(s, event{At: now(), Type: "publication_failed", Item: i.Number, Message: err.Error()})
				continue
			}
		} else {
			i.Status = "awaiting_publication"
		}
		if err := replyToFeedback(s, i); err != nil {
			i.Status, i.Error = "needs_attention", err.Error()
			appendEvent(s, event{At: now(), Type: "feedback_reply_failed", Item: i.Number, Message: err.Error()})
			continue
		}
		completeWorkerRequests(i)
		if s.Config.Publish {
			before := i.Status
			applyPendingStatus(i)
			if i.Status != before {
				appendEvent(s, event{At: now(), Type: "followup_queued", Item: i.Number, Message: i.Status})
			}
		}
		appendEvent(s, event{At: now(), Type: "worker_completed", Item: i.Number})
	}
	_ = statePath
	return nil
}

// automaticRecoveryAllowed retries one unambiguous harness-local verification
// failure. Further failures preserve evidence for human product judgment.
func automaticRecoveryAllowed(s *State, i *Item) bool {
	return s.Config.Harness == "opencode" && i.Worker != nil && i.Worker.Attempt == 1
}

func completeWorkerRequests(i *Item) {
	if i.Worker == nil {
		return
	}
	consumed := map[string]bool{}
	for _, id := range i.Worker.RequestIDs {
		consumed[id] = true
	}
	for n := range i.Pending {
		if i.Pending[n].Status == "queued" && consumed[i.Pending[n].ID] {
			i.Pending[n].Status = "applied"
		}
	}
}

// replyToFeedback closes the feedback loop before the supervisor marks the
// request applied. Workers never receive tracker authority.
func replyToFeedback(s *State, i *Item) error {
	if s.Config.Source != "gitlab" || i.PR == nil || i.Worker == nil {
		return nil
	}
	response := workerFeedbackResponse(i)
	for n := range i.Pending {
		req := &i.Pending[n]
		if req.Action != "feedback" || req.ReplySent || req.ThreadID == "" || !contains(i.Worker.RequestIDs, req.ID) {
			continue
		}
		if response == "" {
			return errors.New("feedback worker handoff is missing `Feedback response: <concise answer>`")
		}
		body := "AI-generated: " + response
		if i.PR.HeadOID != "" {
			body += "\n\nUpdated commit: " + short(i.PR.HeadOID)
		}
		if _, err := run("", strings.NewReader(body), "glab", "mr", "note", "create", strconv.Itoa(i.PR.Number), "--repo", s.Config.RepoSlug, "--reply", req.ThreadID); err != nil {
			return fmt.Errorf("reply to GitLab feedback: %w", err)
		}
		req.ReplySent = true
		appendEvent(s, event{At: now(), Type: "feedback_replied", Item: i.Number, Message: req.ThreadID})
	}
	return nil
}

func workerFeedbackResponse(i *Item) string {
	if i.Worker == nil {
		return ""
	}
	b, err := os.ReadFile(i.Worker.LastMessage)
	if err != nil {
		return ""
	}
	for _, text := range handoffTexts(b) {
		for _, line := range strings.Split(text, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "Feedback response:") {
				return strings.TrimSpace(strings.TrimPrefix(line, "Feedback response:"))
			}
		}
	}
	return ""
}

func verifyItem(s *State, i *Item) error {
	status, err := git(i.Worktree, "status", "--porcelain")
	if err != nil {
		return err
	}
	if strings.TrimSpace(status) == "" {
		head, headErr := git(i.Worktree, "rev-parse", "HEAD")
		if headErr != nil {
			return headErr
		}
		if strings.TrimSpace(head) == i.BaseHead {
			return errors.New("worker produced no changes")
		}
	} else if strings.Contains(status, "UU ") || strings.Contains(status, "AA ") || strings.Contains(status, "DD ") {
		return errors.New("worker left unresolved conflicts; checkout preserved")
	}
	if !isAncestor(i.Worktree, i.BaseHead, "HEAD") {
		return errors.New("worker branch no longer contains its recorded base")
	}
	i.Verification = nil
	for _, command := range s.Config.VerifyCommands {
		out, err := run(i.Worktree, nil, "/bin/sh", "-lc", command)
		check := Check{Command: command, ExitCode: exitCode(err), Output: truncate(out, 4000)}
		i.Verification = append(i.Verification, check)
		if err != nil {
			return fmt.Errorf("verification command failed: %s", command)
		}
	}
	if strings.TrimSpace(status) != "" {
		if _, err := git(i.Worktree, "add", "--all"); err != nil {
			return fmt.Errorf("stage verified worker diff: %w", err)
		}
		subject, err := workerCommitSubject(s, i)
		if err != nil {
			return err
		}
		message := fmt.Sprintf("%s\n\nAssisted-by: %s", subject, harnessCredit(s.Config.Harness))
		if _, err := git(i.Worktree, "commit", "-m", message); err != nil {
			return fmt.Errorf("create supervisor checkpoint commit: %w", err)
		}
		appendEvent(s, event{At: now(), Type: "checkpoint_created", Item: i.Number, Message: "verified worker changes committed"})
	}
	clean, err := git(i.Worktree, "status", "--porcelain")
	if err != nil || strings.TrimSpace(clean) != "" {
		return errors.New("checkpoint commit did not leave a clean checkout")
	}
	return nil
}

var conventionalSubject = regexp.MustCompile(`^(feat|fix|refactor|test|docs|build|ci|chore|style|perf|revert)(\([a-z0-9._/-]+\))?!?: [a-z0-9][^\r\n]*$`)

func workerCommitSubject(s *State, i *Item) (string, error) {
	if i.Worker == nil || i.Worker.LastMessage == "" {
		return "", errors.New("worker handoff is unavailable for commit-message generation")
	}
	b, err := os.ReadFile(i.Worker.LastMessage)
	if err != nil {
		if s.Config.Harness == "opencode" && conventionalSubject.MatchString(i.Title) {
			appendEvent(s, event{At: now(), Type: "commit_subject_derived", Item: i.Number, Message: i.Title})
			return i.Title, nil
		}
		return "", fmt.Errorf("read worker handoff for commit subject: %w", err)
	}
	for _, text := range handoffTexts(b) {
		for _, line := range strings.Split(text, "\n") {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "Commit subject:") {
				continue
			}
			subject := strings.TrimSpace(strings.TrimPrefix(line, "Commit subject:"))
			if len(subject) > 72 || strings.HasSuffix(subject, ".") || !conventionalSubject.MatchString(subject) {
				return "", fmt.Errorf("worker proposed invalid Conventional Commit subject %q", subject)
			}
			return subject, nil
		}
	}
	if s.Config.Harness == "opencode" {
		if conventionalSubject.MatchString(i.Title) {
			appendEvent(s, event{At: now(), Type: "commit_subject_derived", Item: i.Number, Message: i.Title})
			return i.Title, nil
		}
	}
	return "", errors.New("worker handoff is missing `Commit subject: <type>: <summary>`")
}

func handoffTexts(b []byte) []string {
	result := []string{string(b)}
	scanner := bufio.NewScanner(strings.NewReader(string(b)))
	for scanner.Scan() {
		var value any
		if json.Unmarshal(scanner.Bytes(), &value) == nil {
			collectStrings(value, &result)
		}
	}
	return result
}

func collectStrings(value any, result *[]string) {
	switch typed := value.(type) {
	case string:
		*result = append(*result, typed)
	case []any:
		for _, child := range typed {
			collectStrings(child, result)
		}
	case map[string]any:
		for _, child := range typed {
			collectStrings(child, result)
		}
	}
}

func harnessCredit(harness string) string {
	if harness == "opencode" {
		return "OpenCode/AI"
	}
	return "Codex/AI"
}

func publishItem(s *State, i *Item) error {
	if err := verifyItem(s, i); err != nil {
		return err
	}
	if _, err := git(i.Worktree, "push", s.Config.Remote, "HEAD:refs/heads/"+i.Branch); err != nil {
		return fmt.Errorf("push explicit branch refspec: %w", err)
	}
	appendEvent(s, event{At: now(), Type: "branch_pushed", Item: i.Number, Message: i.Branch})
	pr, err := findPull(s, i.Branch)
	if err != nil {
		return err
	}
	if pr == nil {
		body := fmt.Sprintf("Closes #%d\n\n## Summary\n\nImplementation for `%s`.\n\n## Verification\n\n", i.Number, i.Title)
		if len(i.Verification) == 0 {
			body += "Independent verification is recorded in the worker handoff.\n"
		} else {
			for _, c := range i.Verification {
				body += fmt.Sprintf("- `%s`\n", c.Command)
			}
		}
		bodyPath := filepath.Join(s.Config.RunDir, fmt.Sprintf("pr-%d.md", i.Number))
		if err := atomicWrite(bodyPath, []byte(body), 0o600); err != nil {
			return err
		}
		if s.Config.Source == "gitlab" {
			endpoint := "projects/" + url.PathEscape(s.Config.RepoSlug) + "/merge_requests"
			if _, err := run("", nil, "glab", "api", "--method", "POST", endpoint,
				"-f", "source_branch="+i.Branch, "-f", "target_branch="+s.Config.Target,
				"-f", "title=Draft: "+i.Title, "-f", "description="+body); err != nil {
				return err
			}
		} else if _, err := run("", nil, "gh", "pr", "create", "--repo", s.Config.RepoSlug, "--draft", "--base", s.Config.Target, "--head", i.Branch, "--title", i.Title, "--body-file", bodyPath); err != nil {
			return err
		}
		pr, err = findPull(s, i.Branch)
		if err != nil || pr == nil {
			return errors.New("draft PR creation did not read back")
		}
	}
	localHead, err := git(i.Worktree, "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	pr, err = awaitPublishedPull(s, i.Branch, strings.TrimSpace(localHead))
	if err != nil {
		return err
	}
	i.PR, i.Status, i.Error = pr, "in_review", ""
	appendEvent(s, event{At: now(), Type: "published", Item: i.Number, Message: pr.URL})
	return nil
}

func awaitPublishedPull(s *State, branch, head string) (*Pull, error) {
	var latest *Pull
	for attempt := 0; attempt < 8; attempt++ {
		pr, err := findPull(s, branch)
		if err != nil {
			return nil, err
		}
		latest = pr
		if publishedPullMatches(pr, head) {
			return pr, nil
		}
		if attempt < 7 {
			time.Sleep(750 * time.Millisecond)
		}
	}
	if latest == nil {
		return nil, errors.New("draft review did not read back after publication")
	}
	return nil, fmt.Errorf("draft review readback stayed stale: expected head %s, observed %s", short(head), short(latest.HeadOID))
}

func publishedPullMatches(pr *Pull, head string) bool {
	return pr != nil && pr.Draft && pr.HeadOID == head
}

func syncPull(s *State, i *Item) error {
	branch := i.Branch
	if branch == "" {
		branch = fmt.Sprintf("issue/%d-%s", i.Number, slugify(i.Title))
	}
	snapshot, err := pollReview(s, branch)
	if err != nil || snapshot == nil || snapshot.Review == nil {
		return err
	}
	pr := snapshot.Review
	previousStatus := i.Status
	i.PR = pr
	if pr.State == "MERGED" {
		if mergedIntoTarget(s, i) {
			i.Status = "integrated"
			if previousStatus != "integrated" {
				appendEvent(s, event{At: now(), Type: "integrated", Item: i.Number, Message: pr.URL})
			}
		} else {
			i.Status = "closed_unverified"
		}
		return nil
	}
	if applyClosedReview(s, i, pr, previousStatus) {
		return nil
	}
	if pr.MergeState == "DIRTY" && i.Status != "running" {
		i.Status = "reconcile_pending"
	}
	if len(snapshot.Failures) > 0 {
		pr.Checks = strings.Join(snapshot.Failures, ", ")
		queueRequest(s, i, Request{ID: stableID("ci", pr.HeadOID, pr.Checks), Action: "ci_failed", Message: "Failing checks: " + pr.Checks, CreatedAt: now(), Status: "queued", Source: s.Config.Source, HeadOID: pr.HeadOID})
	}
	if len(snapshot.Feedback) > 0 {
		appendEvent(s, event{At: now(), Type: "review_feedback_polled", Item: i.Number, Message: fmt.Sprintf("%d human comment(s)", len(snapshot.Feedback))})
	}
	for _, input := range snapshot.Feedback {
		if err := queueFeedback(s, i, input); err != nil {
			appendEvent(s, event{At: now(), Type: "feedback_ack_failed", Item: i.Number, Message: err.Error()})
		}
	}
	return nil
}

func applyClosedReview(s *State, i *Item, pr *Pull, previousStatus string) bool {
	if pr.State != "CLOSED" {
		return false
	}
	if previousStatus != "review_closed" {
		appendEvent(s, event{At: now(), Type: "review_closed", Item: i.Number, Message: pr.URL})
	}
	i.Status, i.Error = "review_closed", "review closed without integration"
	return true
}

func pollReview(s *State, branch string) (*reviewSnapshot, error) {
	if s.Config.Source == "gitlab" {
		return pollGitLabReview(s, branch)
	}
	return pollGitHubReview(s, branch)
}

func pollGitHubReview(s *State, branch string) (*reviewSnapshot, error) {
	pr, err := findGitHubPull(s, branch)
	if err != nil || pr == nil {
		return nil, err
	}
	detail, err := run("", nil, "gh", "pr", "view", strconv.Itoa(pr.Number), "--repo", s.Config.RepoSlug, "--json", "comments,reviews,reviewDecision,statusCheckRollup")
	if err != nil {
		return nil, err
	}
	var d struct {
		Comments []struct {
			ID, Body, URL string
			Author        struct{ Login string }
		}
		Reviews []struct {
			ID, State, Body string
			Author          struct{ Login string }
		}
		Checks []struct{ Name, Conclusion, State string } `json:"statusCheckRollup"`
	}
	if err := json.Unmarshal([]byte(detail), &d); err != nil {
		return nil, err
	}
	snapshot := &reviewSnapshot{Review: pr}
	for _, c := range d.Checks {
		if c.Conclusion == "FAILURE" || c.Conclusion == "TIMED_OUT" || c.Conclusion == "CANCELLED" || c.State == "FAILURE" {
			snapshot.Failures = append(snapshot.Failures, c.Name)
		}
	}
	for _, c := range d.Comments {
		if humanFeedback(c.Author.Login, c.Body) {
			snapshot.Feedback = append(snapshot.Feedback, feedbackInput{ID: stableID("github-comment", c.ID), Source: "github-comment", RemoteID: c.ID, Body: c.Body, Author: c.Author.Login, URL: c.URL, HeadOID: pr.HeadOID})
		}
	}
	for _, r := range d.Reviews {
		if humanFeedback(r.Author.Login, r.Body) && r.State == "CHANGES_REQUESTED" {
			snapshot.Feedback = append(snapshot.Feedback, feedbackInput{ID: stableID("github-review", r.ID), Source: "github-review", RemoteID: r.ID, Body: r.Body, Author: r.Author.Login, RequestedChanges: true, HeadOID: pr.HeadOID})
		}
	}
	threads, err := githubReviewThreads(s, pr)
	if err != nil {
		return nil, err
	}
	snapshot.Feedback = append(snapshot.Feedback, threads...)
	return snapshot, nil
}

func githubReviewThreads(s *State, pr *Pull) ([]feedbackInput, error) {
	parts := strings.SplitN(s.Config.RepoSlug, "/", 2)
	query := `query($owner:String!,$name:String!,$number:Int!){repository(owner:$owner,name:$name){pullRequest(number:$number){reviewThreads(first:100){pageInfo{hasNextPage} nodes{id isResolved isOutdated path line originalLine comments(first:100){pageInfo{hasNextPage} nodes{id body url author{login}}}}}}}}`
	out, err := run("", nil, "gh", "api", "graphql", "-f", "query="+query, "-F", "owner="+parts[0], "-F", "name="+parts[1], "-F", "number="+strconv.Itoa(pr.Number))
	if err != nil {
		return nil, fmt.Errorf("query GitHub review threads: %w", err)
	}
	var response struct {
		Data struct {
			Repository struct {
				PullRequest struct {
					Threads struct {
						PageInfo struct{ HasNextPage bool } `json:"pageInfo"`
						Nodes    []struct {
							ID, Path               string
							IsResolved, IsOutdated bool
							Line, OriginalLine     int
							Comments               struct {
								PageInfo struct{ HasNextPage bool } `json:"pageInfo"`
								Nodes    []struct {
									ID, Body, URL string
									Author        struct{ Login string }
								} `json:"nodes"`
							} `json:"comments"`
						} `json:"nodes"`
					} `json:"reviewThreads"`
				} `json:"pullRequest"`
			} `json:"repository"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &response); err != nil {
		return nil, err
	}
	threads := response.Data.Repository.PullRequest.Threads
	if threads.PageInfo.HasNextPage {
		return nil, errors.New("GitHub PR has more than 100 review threads; refusing incomplete feedback polling")
	}
	var result []feedbackInput
	for _, thread := range threads.Nodes {
		if thread.IsResolved || thread.IsOutdated || thread.Comments.PageInfo.HasNextPage {
			continue
		}
		line := thread.Line
		if line == 0 {
			line = thread.OriginalLine
		}
		for _, comment := range thread.Comments.Nodes {
			if humanFeedback(comment.Author.Login, comment.Body) {
				result = append(result, feedbackInput{ID: stableID("github-thread", thread.ID, comment.ID), Source: "github-thread", RemoteID: comment.ID, Body: comment.Body, Author: comment.Author.Login, URL: comment.URL, Path: thread.Path, Line: line, Inline: true, HeadOID: pr.HeadOID})
			}
		}
	}
	return result, nil
}

func findPull(s *State, branch string) (*Pull, error) {
	if s.Config.Source == "gitlab" {
		return findGitLabMergeRequest(s, branch)
	}
	return findGitHubPull(s, branch)
}

func findGitHubPull(s *State, branch string) (*Pull, error) {
	out, err := run("", nil, "gh", "pr", "list", "--repo", s.Config.RepoSlug, "--state", "all", "--head", branch, "--limit", "1", "--json", "number,url,state,isDraft,headRefOid,baseRefName,mergeStateStatus,mergeCommit")
	if err != nil {
		return nil, err
	}
	var rows []struct {
		Number                                                int
		URL, State, HeadRefOID, BaseRefName, MergeStateStatus string
		IsDraft                                               bool
		MergeCommit                                           *struct{ OID string }
	}
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	r := rows[0]
	if r.BaseRefName != s.Config.Target {
		return nil, fmt.Errorf("PR #%d targets %s, expected %s", r.Number, r.BaseRefName, s.Config.Target)
	}
	p := &Pull{Number: r.Number, URL: r.URL, State: r.State, Draft: r.IsDraft, HeadOID: r.HeadRefOID, MergeState: r.MergeStateStatus}
	if r.MergeCommit != nil {
		p.MergeOID = r.MergeCommit.OID
	}
	return p, nil
}

func findGitLabMergeRequest(s *State, branch string) (*Pull, error) {
	endpoint := "projects/" + url.PathEscape(s.Config.RepoSlug) + "/merge_requests?scope=all&state=all&per_page=1&source_branch=" + url.QueryEscape(branch)
	out, err := run("", nil, "glab", "api", endpoint)
	if err != nil {
		return nil, err
	}
	var rows []struct {
		IID                 int    `json:"iid"`
		WebURL              string `json:"web_url"`
		State               string `json:"state"`
		Draft               bool   `json:"draft"`
		SHA                 string `json:"sha"`
		TargetBranch        string `json:"target_branch"`
		DetailedMergeStatus string `json:"detailed_merge_status"`
		MergeCommitSHA      string `json:"merge_commit_sha"`
	}
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	r := rows[0]
	if r.TargetBranch != s.Config.Target {
		return nil, fmt.Errorf("merge request !%d targets %s, expected %s", r.IID, r.TargetBranch, s.Config.Target)
	}
	state := strings.ToUpper(r.State)
	mergeState := strings.ToUpper(r.DetailedMergeStatus)
	if r.DetailedMergeStatus == "conflict" || r.DetailedMergeStatus == "cannot_be_merged" {
		mergeState = "DIRTY"
	}
	return &Pull{Number: r.IID, URL: r.WebURL, State: state, Draft: r.Draft, HeadOID: r.SHA, MergeOID: r.MergeCommitSHA, MergeState: mergeState}, nil
}

func pollGitLabReview(s *State, branch string) (*reviewSnapshot, error) {
	mr, err := findGitLabMergeRequest(s, branch)
	if err != nil || mr == nil {
		return nil, err
	}
	base := "projects/" + url.PathEscape(s.Config.RepoSlug) + "/merge_requests/" + strconv.Itoa(mr.Number)
	out, err := run("", nil, "glab", "api", base+"/discussions?per_page=100")
	if err != nil {
		return nil, fmt.Errorf("query GitLab merge request discussions: %w", err)
	}
	var discussions []struct {
		ID         string `json:"id"`
		Individual bool   `json:"individual_note"`
		Notes      []struct {
			ID         int                       `json:"id"`
			Body       string                    `json:"body"`
			System     bool                      `json:"system"`
			Resolvable bool                      `json:"resolvable"`
			Resolved   bool                      `json:"resolved"`
			Author     struct{ Username string } `json:"author"`
			Position   *struct {
				NewPath, OldPath string
				NewLine, OldLine int
			} `json:"position"`
		} `json:"notes"`
	}
	if err := json.Unmarshal([]byte(out), &discussions); err != nil {
		return nil, err
	}
	snapshot := &reviewSnapshot{Review: mr}
	for _, discussion := range discussions {
		for _, note := range discussion.Notes {
			if note.System || note.Resolved || !humanFeedback(note.Author.Username, note.Body) {
				continue
			}
			input := feedbackInput{ID: stableID("gitlab-discussion", discussion.ID, strconv.Itoa(note.ID)), Source: "gitlab-discussion", RemoteID: strconv.Itoa(note.ID), ThreadID: discussion.ID, Body: note.Body, Author: note.Author.Username, HeadOID: mr.HeadOID, Inline: !discussion.Individual || note.Position != nil}
			if note.Position != nil {
				input.Path, input.Line = note.Position.NewPath, note.Position.NewLine
				if input.Path == "" {
					input.Path, input.Line = note.Position.OldPath, note.Position.OldLine
				}
			}
			snapshot.Feedback = append(snapshot.Feedback, input)
		}
	}
	pipelines, err := run("", nil, "glab", "api", base+"/pipelines?per_page=1")
	if err != nil {
		return nil, fmt.Errorf("query GitLab merge request pipelines: %w", err)
	}
	var rows []struct {
		ID     int
		Status string
	}
	if err := json.Unmarshal([]byte(pipelines), &rows); err != nil {
		return nil, err
	}
	if len(rows) > 0 && (rows[0].Status == "failed" || rows[0].Status == "canceled") {
		snapshot.Failures = append(snapshot.Failures, fmt.Sprintf("pipeline %d (%s)", rows[0].ID, rows[0].Status))
	}
	return snapshot, nil
}

func mergedIntoTarget(s *State, i *Item) bool {
	return i.PR != nil && i.PR.State == "MERGED" && isAncestor(s.Config.Repo, i.PR.MergeOID, s.Config.Remote+"/"+s.Config.Target)
}

func queueFeedback(s *State, i *Item, input feedbackInput) error {
	action := classifyFeedback(input)
	for n := range i.Pending {
		if i.Pending[n].ID == input.ID && i.Pending[n].Action == "observed" && action == "feedback" {
			i.Pending[n].Action, i.Pending[n].Status = "feedback", "queued"
			appendEvent(s, event{At: now(), Type: "feedback_promoted", Item: i.Number, Message: requestSummary(i.Pending[n])})
			applyPendingStatus(i)
			return acknowledgeFeedback(s, i, input)
		}
	}
	newFeedback := !contains(i.SeenFeedback, input.ID) && !containsRequest(i.Pending, input.ID)
	queueRequest(s, i, Request{
		ID: input.ID, Action: action, Message: strings.TrimSpace(input.Body),
		CreatedAt: now(), Status: "queued", Source: input.Source, URL: input.URL,
		Path: input.Path, Line: input.Line, HeadOID: input.HeadOID, ThreadID: input.ThreadID,
	})
	if action == "feedback" && newFeedback {
		if err := acknowledgeFeedback(s, i, input); err != nil {
			return err
		}
		appendEvent(s, event{At: now(), Type: "feedback_acknowledged", Item: i.Number, Message: input.Source})
	}
	return nil
}

func acknowledgeFeedback(s *State, i *Item, input feedbackInput) error {
	if input.RemoteID == "" || i.PR == nil {
		return errors.New("feedback acknowledgement lacks provider-native comment ID")
	}
	if s.Config.Source == "gitlab" {
		endpoint := "projects/" + url.PathEscape(s.Config.RepoSlug) + "/merge_requests/" + strconv.Itoa(i.PR.Number) + "/notes/" + input.RemoteID + "/award_emoji"
		_, err := run("", nil, "glab", "api", "--method", "POST", endpoint, "-f", "name=eyes")
		return err
	}
	mutation := `mutation($subject:ID!){addReaction(input:{subjectId:$subject,content:EYES}){reaction{content}}}`
	_, err := run("", nil, "gh", "api", "graphql", "-f", "query="+mutation, "-F", "subject="+input.RemoteID)
	return err
}

func classifyFeedback(input feedbackInput) string {
	body := strings.ToLower(strings.TrimSpace(input.Body))
	if body == "" {
		return "needs_input"
	}
	for _, marker := range []string{
		"out of scope", "scope change", "product decision", "design decision",
		"architecture decision", "need to discuss", "let's discuss", "lets discuss",
	} {
		if strings.Contains(body, marker) {
			return "needs_input"
		}
	}
	for _, marker := range []string{
		"please ", "fix ", "change ", "update ", "add ", "remove ", "rename ",
		"use ", "ensure ", "must ", "need to ", "can you ", "can we ", "could you ", "should ", "should not ", "anything that ", "i'd like ",
	} {
		if strings.Contains(body, marker) {
			return "feedback"
		}
	}
	if strings.HasSuffix(body, "?") {
		return "needs_input"
	}
	if input.Inline || input.RequestedChanges {
		return "feedback"
	}
	return "observed"
}

func queueRequest(s *State, i *Item, req Request) {
	if req.ID == "" || contains(i.SeenFeedback, req.ID) || containsRequest(i.Pending, req.ID) {
		return
	}
	if req.CreatedAt == "" {
		req.CreatedAt = now()
	}
	if req.Status == "" {
		req.Status = "queued"
	}
	i.SeenFeedback = append(i.SeenFeedback, req.ID)
	i.Pending = append(i.Pending, req)

	eventType := "feedback_observed"
	if req.Action == "feedback" {
		eventType = "feedback_detected"
	} else if req.Action == "ci_failed" {
		eventType = "ci_failed"
	} else if req.Action == "needs_input" {
		eventType = "input_required"
	}
	appendEvent(s, event{At: now(), Type: eventType, Item: i.Number, Message: requestSummary(req)})
	if i.Status != "running" {
		applyPendingStatus(i)
	}
}

func requestSummary(req Request) string {
	location := ""
	if req.Path != "" {
		location = req.Path
		if req.Line > 0 {
			location += ":" + strconv.Itoa(req.Line)
		}
		location += ": "
	}
	return truncate(location+strings.TrimSpace(req.Message), 500)
}

func applyPendingStatus(i *Item) {
	hasAction, hasInput, hasCI := false, false, false
	for _, req := range i.Pending {
		if req.Status != "queued" {
			continue
		}
		switch req.Action {
		case "needs_input":
			hasInput = true
		case "ci_failed":
			hasCI = true
		case "feedback", "resume":
			hasAction = true
		}
	}
	if hasInput {
		i.Status = "needs_input"
	} else if hasCI {
		i.Status = "ci_failed"
	} else if hasAction {
		i.Status = "feedback"
	}
}

func processInbox(s *State) error {
	path := filepath.Join(s.Config.RunDir, "inbox.jsonl")
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var req struct {
			Item int `json:"item"`
			Request
		}
		if json.Unmarshal(scanner.Bytes(), &req) != nil {
			continue
		}
		i := s.Items[strconv.Itoa(req.Item)]
		if i == nil || containsRequest(i.Pending, req.ID) {
			continue
		}
		if req.Action == "stop" && i.Worker != nil && processAlive(i.Worker.PID, i.Worker.PIDStart) {
			i.Pending = append(i.Pending, req.Request)
			_ = syscall.Kill(-i.Worker.PID, syscall.SIGTERM)
			i.Status = "stopped"
		} else if req.Action == "resume" && i.Status != "running" {
			i.Pending = append(i.Pending, req.Request)
			s.StopRequested = false
			i.Status = "feedback"
		} else if req.Action == "feedback" {
			req.Source = "operator"
			queueRequest(s, i, req.Request)
		}
	}
	return scanner.Err()
}

func cmdInbox(args []string) error {
	fs := flag.NewFlagSet("inbox", flag.ContinueOnError)
	path := fs.String("state", "", "path to state.json")
	item := fs.Int("item", 0, "issue number")
	action := fs.String("action", "feedback", "feedback, stop, or resume")
	message := fs.String("message", "", "request text")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *path == "" || *item == 0 {
		return errors.New("--state and --item are required")
	}
	if *action != "feedback" && *action != "stop" && *action != "resume" {
		return errors.New("--action must be feedback, stop, or resume")
	}
	s, err := loadState(*path)
	if err != nil {
		return err
	}
	req := struct {
		Item int `json:"item"`
		Request
	}{*item, Request{ID: stableID(strconv.Itoa(*item), *action, *message, now()), Action: *action, Message: *message, CreatedAt: now(), Status: "queued"}}
	return appendJSON(filepath.Join(s.Config.RunDir, "inbox.jsonl"), req)
}

func cmdStop(args []string) error {
	fs := flag.NewFlagSet("stop", flag.ContinueOnError)
	path := fs.String("state", "", "path to state.json")
	workers := fs.Bool("workers", false, "also terminate active worker process groups")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return locked(*path, func(s *State) error {
		s.StopRequested = true
		if *workers {
			for _, i := range s.Items {
				if i.Worker != nil && processAlive(i.Worker.PID, i.Worker.PIDStart) {
					_ = syscall.Kill(-i.Worker.PID, syscall.SIGTERM)
					i.Status = "stopped"
				}
			}
		}
		return saveState(*path, s)
	})
}

func cmdPublish(args []string) error {
	fs := flag.NewFlagSet("publish", flag.ContinueOnError)
	path := fs.String("state", "", "path to state.json")
	item := fs.Int("item", 0, "issue number")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return locked(*path, func(s *State) error {
		i := s.Items[strconv.Itoa(*item)]
		if i == nil {
			return errors.New("unknown item")
		}
		if err := publishItem(s, i); err != nil {
			return err
		}
		completeWorkerRequests(i)
		applyPendingStatus(i)
		return saveState(*path, s)
	})
}

func cmdStatus(args []string) error {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	path := fs.String("state", "", "path to state.json")
	asJSON := fs.Bool("json", false, "print state JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	s, err := loadState(*path)
	if err != nil {
		return err
	}
	if err := ensureOpenCodeSupervisor(*path, s); err != nil {
		return err
	}
	if *asJSON {
		return json.NewEncoder(os.Stdout).Encode(s)
	}
	fmt.Printf("RUN %s  TARGET %s@%s  HARNESS %s  LIMIT %d\n", s.RunID, s.Config.Target, short(s.TargetHead), s.Config.Harness, s.Config.Concurrency)
	fmt.Println("ITEM  STATUS                WORKER  PR     BLOCKERS  TITLE")
	for _, i := range sortedItems(s) {
		worker, pr := "-", "-"
		if i.Worker != nil && i.Status == "running" {
			worker = strconv.Itoa(i.Worker.PID)
		}
		if i.PR != nil {
			pr = "#" + strconv.Itoa(i.PR.Number)
		}
		pending := 0
		for _, b := range i.Blockers {
			if !b.Integrated {
				pending++
			}
		}
		fmt.Printf("#%-4d %-21s %-7s %-6s %-9d %s\n", i.Number, i.Status, worker, pr, pending, i.Title)
	}
	return nil
}

// ensureOpenCodeSupervisor restores the durable loop only for OpenCode runs.
// Other harnesses own their lifecycle and are never started from status.
func ensureOpenCodeSupervisor(statePath string, s *State) error {
	if s.Config.Harness != "opencode" || s.StopRequested || runComplete(s) {
		return nil
	}
	owner, err := acquireSupervisorOwner(statePath)
	if err != nil {
		return nil // An existing supervisor owns the run.
	}
	releaseSupervisorOwner(owner)
	cmd := exec.Command(os.Args[0], "supervise", "--state", statePath)
	logPath := filepath.Join(s.Config.RunDir, "supervisor.log")
	log, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	cmd.Stdout, cmd.Stderr = log, log
	if err := cmd.Start(); err != nil {
		log.Close()
		return err
	}
	appendEvent(s, event{At: now(), Type: "supervisor_restarted", Message: strconv.Itoa(cmd.Process.Pid)})
	return saveState(statePath, s)
}

func cmdEvents(args []string) error {
	fs := flag.NewFlagSet("events", flag.ContinueOnError)
	statePath := fs.String("state", "", "path to state.json")
	after := fs.Int64("after", 0, "only events after this sequence")
	follow := fs.Bool("follow", false, "wait for new events")
	all := fs.Bool("all", false, "include debug events")
	asJSON := fs.Bool("json", false, "print JSON lines")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *statePath == "" {
		return errors.New("--state is required")
	}
	s, err := loadState(*statePath)
	if err != nil {
		return err
	}
	path := filepath.Join(s.Config.RunDir, "events.jsonl")
	cursor := *after
	for {
		events, err := readEvents(path, cursor)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		for _, e := range events {
			if e.Seq > cursor {
				cursor = e.Seq
			}
			if !*all && e.Priority == "debug" {
				continue
			}
			printEvent(e, *asJSON)
		}
		if !*follow {
			return nil
		}
		time.Sleep(time.Second)
	}
}

func readEvents(path string, after int64) ([]event, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var result []event
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var e event
		if json.Unmarshal(scanner.Bytes(), &e) == nil && e.Seq > after {
			result = append(result, e)
		}
	}
	return result, scanner.Err()
}

func printEvent(e event, asJSON bool) {
	if asJSON {
		_ = json.NewEncoder(os.Stdout).Encode(e)
		return
	}
	item := ""
	if e.Item > 0 {
		item = " #" + strconv.Itoa(e.Item)
	}
	fmt.Printf("EVENT %d %-9s %-22s%s %s\n", e.Seq, e.Priority, e.Type, item, e.Message)
}

func locked(path string, fn func(*State) error) error {
	if path == "" {
		return errors.New("--state is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	lock, err := os.OpenFile(path+".lock", os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return err
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return errors.New("another supervisor owns the run lock")
	}
	defer syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
	s, err := loadState(path)
	if err != nil {
		return err
	}
	return fn(s)
}

func loadState(path string) (*State, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s State
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	if s.Version != stateVersion {
		return nil, fmt.Errorf("unsupported state version %d", s.Version)
	}
	if s.Items == nil {
		s.Items = map[string]*Item{}
	}
	return &s, nil
}

func saveState(path string, s *State) error {
	s.UpdatedAt = now()
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return atomicWrite(path, append(b, '\n'), 0o600)
}

func atomicWrite(path string, b []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".orchestrator-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if err := f.Chmod(mode); err != nil {
		f.Close()
		return err
	}
	if _, err := f.Write(b); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func appendEvent(s *State, e event) {
	s.NextEventSeq++
	e.Seq = s.NextEventSeq
	if e.At == "" {
		e.At = now()
	}
	if e.Priority == "" {
		e.Priority = eventPriority(e.Type)
	}
	path := filepath.Join(s.Config.RunDir, "events.jsonl")
	if info, err := os.Stat(path); err == nil && info.Size() > 5<<20 {
		_ = os.Rename(path, path+"."+time.Now().UTC().Format("20060102T150405Z"))
	}
	_ = appendJSON(path, e)
	if printSupervisorEvents && e.Priority != "debug" {
		printEvent(e, false)
	}
}

func eventPriority(kind string) string {
	switch kind {
	case "feedback_detected", "input_required", "ci_failed", "worker_lost", "worker_failed",
		"verification_failed", "publication_failed", "checkpoint_created", "branch_pushed",
		"published", "followup_queued", "queue_blocked", "queue_idle", "integrated",
		"review_closed", "terminal_worker_stopped", "checkout_removed", "worktree_root_removed", "run_completed":
		return "important"
	case "worker_started", "worker_completed", "queue_progress":
		return "progress"
	default:
		return "debug"
	}
}

func appendJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(v)
}

func git(dir string, args ...string) (string, error) {
	all := append([]string{"-C", dir}, args...)
	return run("", nil, "git", all...)
}

func run(dir string, stdin io.Reader, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir, cmd.Stdin = dir, stdin
	b, err := cmd.CombinedOutput()
	out := string(b)
	if err != nil {
		return out, fmt.Errorf("%s %s: %w: %s", name, strings.Join(args, " "), err, truncate(strings.TrimSpace(out), 2000))
	}
	return out, nil
}

func parseRepoSlug(remote string) (string, error) {
	provider, slug, err := parseRemote(remote)
	if err != nil {
		return "", err
	}
	if provider != "github" {
		return "", fmt.Errorf("remote is not GitHub: %s", remote)
	}
	return slug, nil
}

func parseRemote(remote string) (string, string, error) {
	remote = strings.TrimSuffix(strings.TrimSpace(remote), ".git")
	if parsed, err := url.Parse(remote); err == nil && parsed.Scheme == "https" && parsed.Host != "" {
		provider := "gitlab"
		if parsed.Host == "github.com" {
			provider = "github"
		}
		slug := strings.TrimPrefix(parsed.Path, "/")
		if len(strings.Split(slug, "/")) >= 2 {
			return provider, slug, nil
		}
	}
	provider := ""
	for _, candidate := range []struct{ prefix, provider string }{
		{"git@github.com:", "github"}, {"ssh://git@github.com/", "github"}, {"https://github.com/", "github"},
		{"git@gitlab.com:", "gitlab"}, {"ssh://git@gitlab.com/", "gitlab"}, {"https://gitlab.com/", "gitlab"},
	} {
		if strings.HasPrefix(remote, candidate.prefix) {
			remote = strings.TrimPrefix(remote, candidate.prefix)
			provider = candidate.provider
			break
		}
	}
	if provider == "" {
		return "", "", fmt.Errorf("remote is not a supported GitHub or GitLab URL: %s", remote)
	}
	if len(strings.Split(remote, "/")) < 2 {
		return "", "", fmt.Errorf("invalid repository slug %q", remote)
	}
	return provider, remote, nil
}

var nonSlug = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(v string) string {
	v = nonSlug.ReplaceAllString(strings.ToLower(v), "-")
	v = strings.Trim(v, "-")
	if len(v) > 48 {
		v = strings.TrimRight(v[:48], "-")
	}
	if v == "" {
		return "work-item"
	}
	return v
}

func isAncestor(repo, ancestor, descendant string) bool {
	if ancestor == "" {
		return false
	}
	_, err := git(repo, "merge-base", "--is-ancestor", ancestor, descendant)
	return err == nil
}

func sortedItems(s *State) []*Item {
	items := make([]*Item, 0, len(s.Items))
	for _, i := range s.Items {
		items = append(items, i)
	}
	sort.Slice(items, func(a, b int) bool { return items[a].Number < items[b].Number })
	return items
}

func contains(values []string, want string) bool {
	for _, v := range values {
		if v == want {
			return true
		}
	}
	return false
}

func containsRequest(values []Request, want string) bool {
	for _, v := range values {
		if v.ID == want {
			return true
		}
	}
	return false
}

func stableID(values ...string) string {
	h := sha256.Sum256([]byte(strings.Join(values, "\x00")))
	return hex.EncodeToString(h[:12])
}

func humanFeedback(login, body string) bool {
	trimmed := strings.TrimSpace(body)
	return login != "" && trimmed != "" && !strings.HasSuffix(login, "[bot]") && !strings.HasPrefix(trimmed, "AI-generated:")
}

func scrubWorkerEnv(values []string) []string {
	denied := []string{"GH_TOKEN=", "GITHUB_TOKEN=", "GITLAB_TOKEN=", "GLAB_TOKEN=", "SSH_AUTH_SOCK=", "GIT_ASKPASS=", "AWS_", "CLERK_", "STRIPE_"}
	result := make([]string, 0, len(values))
	for _, value := range values {
		blocked := false
		for _, prefix := range denied {
			if strings.HasPrefix(value, prefix) {
				blocked = true
				break
			}
		}
		if !blocked && !strings.HasPrefix(value, "OPENCODE_CONFIG_CONTENT=") {
			result = append(result, value)
		}
	}
	return result
}

func procStart(pid int) string {
	b, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return ""
	}
	fields := strings.Fields(string(b))
	if len(fields) < 22 {
		return ""
	}
	return fields[21]
}

func processAlive(pid int, start string) bool {
	if pid <= 0 || syscall.Kill(pid, 0) != nil {
		return false
	}
	return start == "" || procStart(pid) == start
}

func stalled(path string, minutes int) bool {
	info, err := os.Stat(path)
	return err == nil && time.Since(info.ModTime()) > time.Duration(minutes)*time.Minute
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	var e *exec.ExitError
	if errors.As(err, &e) {
		return e.ExitCode()
	}
	return 1
}

func short(v string) string {
	if len(v) > 8 {
		return v[:8]
	}
	return v
}

func truncate(v string, n int) string {
	if len(v) <= n {
		return v
	}
	return v[:n] + "…"
}

func now() string { return time.Now().UTC().Format(time.RFC3339Nano) }
