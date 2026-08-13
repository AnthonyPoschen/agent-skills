package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseRepoSlug(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"git@github.com:owner/repo.git", "owner/repo"},
		{"https://github.com/owner/repo.git", "owner/repo"},
		{"ssh://git@github.com/owner/repo", "owner/repo"},
	} {
		got, err := parseRepoSlug(tc.in)
		if err != nil || got != tc.want {
			t.Fatalf("parseRepoSlug(%q) = %q, %v", tc.in, got, err)
		}
	}
	if _, err := parseRepoSlug("https://gitlab.com/owner/repo"); err == nil {
		t.Fatal("expected non-GitHub remote to fail closed")
	}
}

func TestParseRemoteSupportsGitHubAndGitLab(t *testing.T) {
	for _, tc := range []struct{ in, provider, slug string }{
		{"git@github.com:owner/repo.git", "github", "owner/repo"},
		{"https://gitlab.com/group/subgroup/repo.git", "gitlab", "group/subgroup/repo"},
	} {
		provider, slug, err := parseRemote(tc.in)
		if err != nil || provider != tc.provider || slug != tc.slug {
			t.Fatalf("parseRemote(%q) = %q, %q, %v", tc.in, provider, slug, err)
		}
	}
}

func TestReadyRequiresIntegratedBlockers(t *testing.T) {
	i := &Item{State: "OPEN", Labels: []string{"ready-for-agent"}, Blockers: []Blocker{{Number: 1, Integrated: false}}}
	if ready(i, "me", "ready-for-agent") {
		t.Fatal("item with unintegrated blocker was ready")
	}
	i.Blockers[0].Integrated = true
	if !ready(i, "me", "ready-for-agent") {
		t.Fatal("item with integrated blockers was not ready")
	}
	i.Assignees = []string{"someone-else"}
	if ready(i, "me", "ready-for-agent") {
		t.Fatal("item assigned to another user was ready")
	}
}

func TestOpenBlockerNeverCountsAsIntegrated(t *testing.T) {
	i := &Item{State: "OPEN", Labels: []string{"ready-for-agent"}, Blockers: []Blocker{{Number: 1, State: "OPEN", Integrated: false}}}
	if ready(i, "me", "ready-for-agent") {
		t.Fatal("open blocker was accepted")
	}
}

func TestWorkerPromptPreservesBodyAndAuthorityBoundary(t *testing.T) {
	s := &State{Config: Config{Invocation: "$implement", Repo: "/repo", Remote: "origin", Target: "main"}}
	i := &Item{Title: "Add a thing", Body: "Line one\n\n- exact acceptance criterion", Branch: "issue/2-add", Worktree: "/worktree", BaseHead: "abc"}
	prompt := workerPrompt(s, i)
	for _, required := range []string{i.Body, "$implement", "Do not push", "leave the tested changes unstaged", "Commit subject:"} {
		if !strings.Contains(prompt, required) {
			t.Fatalf("prompt missing %q", required)
		}
	}
}

func TestWorkerCommitSubjectUsesAIHandoff(t *testing.T) {
	dir := t.TempDir()
	plain := filepath.Join(dir, "plain.txt")
	if err := os.WriteFile(plain, []byte("Done.\nCommit subject: fix: remove stale review state\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	i := &Item{Worker: &Worker{LastMessage: plain}}
	if got, err := workerCommitSubject(i); err != nil || got != "fix: remove stale review state" {
		t.Fatalf("workerCommitSubject = %q, %v", got, err)
	}

	jsonPath := filepath.Join(dir, "events.jsonl")
	if err := os.WriteFile(jsonPath, []byte(`{"type":"message","part":{"text":"Commit subject: feat(events): relay supervisor progress"}}`+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	i.Worker.LastMessage = jsonPath
	if got, err := workerCommitSubject(i); err != nil || got != "feat(events): relay supervisor progress" {
		t.Fatalf("JSON workerCommitSubject = %q, %v", got, err)
	}
}

func TestWorkerCommitSubjectRejectsGenericOrInvalidHandoff(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "last.txt")
	i := &Item{Worker: &Worker{LastMessage: path}}
	for _, content := range []string{
		"Implemented issue 6.",
		"Commit subject: implement issue 6",
		"Commit subject: fix: Remove stale state.",
	} {
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := workerCommitSubject(i); err == nil {
			t.Fatalf("invalid handoff was accepted: %q", content)
		}
	}
}

func TestStateRoundTripIsAtomicAndPrivate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")
	want := &State{Version: stateVersion, RunID: "test", Items: map[string]*Item{"7": {Number: 7, Status: "ready"}}}
	if err := saveState(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := loadState(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Items["7"].Status != "ready" {
		t.Fatalf("unexpected state: %#v", got.Items["7"])
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("state mode = %o", info.Mode().Perm())
	}
}

func TestSlugifyIsStableAndBounded(t *testing.T) {
	got := slugify("  Lifecycle: Intent & Actions!  ")
	if got != "lifecycle-intent-actions" {
		t.Fatalf("unexpected slug %q", got)
	}
	if len(slugify(strings.Repeat("long title ", 20))) > 48 {
		t.Fatal("slug exceeded branch-name bound")
	}
}

func TestReadyLabelDiscoveryRequiresDocumentation(t *testing.T) {
	labels := []string{"bug", "ready-for-agent", "status: ready"}
	if got := discoverReadyLabel("use the `ready-for-agent` label for executable work.", labels); got != "ready-for-agent" {
		t.Fatalf("discovered %q", got)
	}
	if got := discoverReadyLabel("no workflow is documented", labels); got != "" {
		t.Fatalf("silently guessed undocumented label %q", got)
	}
}

func TestSetupErrorListsEveryQuestion(t *testing.T) {
	err := setupError([]string{"Which target?", "Which ready label?"})
	for _, want := range []string{"Which target?", "Which ready label?"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("missing setup question %q", want)
		}
	}
}

func TestProjectSetupRejectsUnknownFields(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".github"), 0o700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, ".github", "implement-tickets.json")
	if err := os.WriteFile(path, []byte(`{"version":1,"ready_lable":"typo"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadProjectSetup(dir, ".github/implement-tickets.json"); err == nil {
		t.Fatal("unknown contract field was accepted")
	}
}

func TestIssueRangeAndSelection(t *testing.T) {
	got, err := parseIssueRange("7-10")
	if err != nil || len(got) != 4 || got[0] != 7 || got[3] != 10 {
		t.Fatalf("parseIssueRange = %v, %v", got, err)
	}
	if _, err := parseIssueRange("10-7"); err == nil {
		t.Fatal("descending range was accepted")
	}
	if unique := uniqueInts([]int{9, 7, 9, 8}); fmt.Sprint(unique) != "[7 8 9]" {
		t.Fatalf("unique selection = %v", unique)
	}
}

func TestOpenCodePolicyIsDenyByDefault(t *testing.T) {
	s := &State{Config: Config{WorkerAgent: "worker", VerifyCommands: []string{"make test"}}}
	raw, err := opencodeWorkerPolicy(s)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"external_directory":"deny"`, `"webfetch":"deny"`, `"make test":"allow"`} {
		if !strings.Contains(raw, want) {
			t.Fatalf("policy missing %s: %s", want, raw)
		}
	}
	if strings.Contains(raw, `"git push *":"allow"`) {
		t.Fatal("policy allowed git push")
	}
}

func TestWorkerEnvironmentDropsPublicationCredentials(t *testing.T) {
	got := scrubWorkerEnv([]string{"PATH=/bin", "GH_TOKEN=secret", "AWS_PROFILE=prod", "OPENAI_API_KEY=provider"})
	joined := strings.Join(got, "\n")
	if strings.Contains(joined, "secret") || strings.Contains(joined, "AWS_PROFILE") {
		t.Fatalf("sensitive publication credentials remained: %s", joined)
	}
	if !strings.Contains(joined, "OPENAI_API_KEY") {
		t.Fatal("model provider credential was removed")
	}
}

func TestFeedbackClassificationIsConservative(t *testing.T) {
	cases := []struct {
		input feedbackInput
		want  string
	}{
		{feedbackInput{Body: "Please rename this helper."}, "feedback"},
		{feedbackInput{Body: "Can you add a regression test?"}, "feedback"},
		{feedbackInput{Body: "Why was this approach selected?"}, "needs_input"},
		{feedbackInput{Body: "This requires a product decision before implementation."}, "needs_input"},
		{feedbackInput{Body: "The guard is inverted.", Inline: true}, "feedback"},
		{feedbackInput{Body: "Looks good to me."}, "observed"},
	}
	for _, tc := range cases {
		if got := classifyFeedback(tc.input); got != tc.want {
			t.Errorf("classifyFeedback(%q) = %q, want %q", tc.input.Body, got, tc.want)
		}
	}
}

func TestFeedbackAuthorizationUsesAuthenticatedAccount(t *testing.T) {
	inputs := []feedbackInput{
		{ID: "own", Author: "AnthonyPoschen", Body: "Please add a test."},
		{ID: "ai", Author: "AnthonyPoschen", Body: "AI-generated: updated smoke/feedback.txt"},
		{ID: "other", Author: "reviewer", Body: "Please rename this helper."},
		{ID: "scanner", Author: "codecov[bot]", Body: "Please cover the error path."},
	}
	var eligible []feedbackInput
	for _, input := range inputs {
		if eligibleFeedback(input.Author, input.Body) {
			eligible = append(eligible, input)
		}
	}
	got := authorizeFeedback(eligible, "anTHonyposchen")
	if fmt.Sprint(feedbackIDs(got)) != "[own scanner]" {
		t.Fatalf("authorized feedback = %v, want own and scanner feedback", feedbackIDs(got))
	}
}

func TestThirdPartyFeedbackRequiresExactReactionOrHumanEndorsement(t *testing.T) {
	inputs := []feedbackInput{
		{ID: "reaction", Author: "reviewer", Body: "Please add a test.", PositiveReactors: []string{"AnthonyPoschen"}},
		{ID: "negative", Author: "reviewer", Body: "Please rename this.", PositiveReactors: []string{"someone-else"}},
		{ID: "threaded", Author: "reviewer-2", Body: "Use the shared helper.", Group: "thread-1", Threaded: true},
		{ID: "endorsement", Author: "AnthonyPoschen", Body: "Yes, please implement this.", Group: "thread-1", Threaded: true},
	}
	got := authorizeFeedback(inputs, "AnthonyPoschen")
	if fmt.Sprint(feedbackIDs(got)) != "[reaction threaded]" {
		t.Fatalf("authorized feedback = %v", feedbackIDs(got))
	}
	if !strings.Contains(got[1].Body, "Use the shared helper.") || !strings.Contains(got[1].Body, "Yes, please implement this.") {
		t.Fatalf("comment endorsement did not combine both messages: %q", got[1].Body)
	}
}

func TestUnthreadedEndorsementMustIdentifyFeedback(t *testing.T) {
	feedback := feedbackInput{ID: "other", Author: "alice", Body: "Please add a guard.", URL: "https://example.test/comment/1", Group: "pr"}
	generic := feedbackInput{ID: "generic", Author: "owner", Body: "Yes, please implement this.", Group: "pr"}
	if got := authorizeFeedback([]feedbackInput{feedback, generic}, "owner"); len(got) != 0 {
		t.Fatalf("generic unthreaded endorsement authorized unrelated feedback: %v", feedbackIDs(got))
	}
	explicit := feedbackInput{ID: "explicit", Author: "owner", Body: "Please implement @alice's feedback.", Group: "pr"}
	if got := authorizeFeedback([]feedbackInput{feedback, explicit}, "owner"); len(got) != 1 || got[0].ID != "other" {
		t.Fatalf("explicit unthreaded endorsement was not applied: %v", feedbackIDs(got))
	}
}

func TestThreadedEndorsementTargetsNearestPrecedingFeedback(t *testing.T) {
	inputs := []feedbackInput{
		{ID: "first", Author: "alice", Body: "Rename the type.", Group: "thread", Threaded: true},
		{ID: "second", Author: "bob", Body: "Add a regression test.", Group: "thread", Threaded: true},
		{ID: "endorsement", Author: "owner", Body: "Yes, please implement this.", Group: "thread", Threaded: true},
		{ID: "later", Author: "carol", Body: "Also change the API.", Group: "thread", Threaded: true},
	}
	got := authorizeFeedback(inputs, "owner")
	if fmt.Sprint(feedbackIDs(got)) != "[second]" {
		t.Fatalf("thread endorsement did not target nearest preceding feedback: %v", feedbackIDs(got))
	}
}

func TestPositiveReactionSetExcludesAmbiguousEmoji(t *testing.T) {
	for _, reaction := range []string{"THUMBS_UP", "hooray", "heart", "rocket", "tada"} {
		if !positiveReaction(reaction) {
			t.Errorf("positive reaction %q was rejected", reaction)
		}
	}
	for _, reaction := range []string{"THUMBS_DOWN", "CONFUSED", "EYES", "LAUGH"} {
		if positiveReaction(reaction) {
			t.Errorf("ambiguous or negative reaction %q was accepted", reaction)
		}
	}
}

func feedbackIDs(inputs []feedbackInput) []string {
	result := make([]string, 0, len(inputs))
	for _, input := range inputs {
		result = append(result, input.ID)
	}
	return result
}

func TestFeedbackArrivingDuringWorkerRemainsQueued(t *testing.T) {
	i := &Item{Number: 4, Status: "running", Worker: &Worker{RequestIDs: []string{"old"}}, Pending: []Request{{ID: "old", Action: "feedback", Status: "queued"}}}
	s := &State{Config: Config{RunDir: t.TempDir()}, Items: map[string]*Item{"4": i}}
	queueRequest(s, i, Request{ID: "new", Action: "feedback", Message: "Please add a test.", Status: "queued"})
	completeWorkerRequests(i)
	applyPendingStatus(i)
	if i.Status != "feedback" || i.Pending[1].Status != "queued" {
		t.Fatalf("later feedback was consumed: status=%s pending=%#v", i.Status, i.Pending)
	}
}

func TestEventsAreSequencedAndPrioritized(t *testing.T) {
	dir := t.TempDir()
	s := &State{Config: Config{RunDir: dir}}
	appendEvent(s, event{Type: "sync"})
	appendEvent(s, event{Type: "feedback_detected", Item: 2})
	events, err := readEvents(filepath.Join(dir, "events.jsonl"), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 || events[0].Seq != 1 || events[1].Seq != 2 || events[1].Priority != "important" {
		t.Fatalf("unexpected events: %#v", events)
	}
}

func TestPublishedPullMustBeDraftAtExpectedHead(t *testing.T) {
	head := "abc123"
	if !publishedPullMatches(&Pull{Draft: true, HeadOID: head}, head) {
		t.Fatal("matching draft pull was rejected")
	}
	if publishedPullMatches(&Pull{Draft: true, HeadOID: "stale"}, head) {
		t.Fatal("stale pull head was accepted")
	}
	if publishedPullMatches(&Pull{Draft: false, HeadOID: head}, head) {
		t.Fatal("non-draft pull was accepted")
	}
}

func TestSupervisorOwnershipCoversProcessLifetime(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "state.json")
	first, err := acquireSupervisorOwner(statePath)
	if err != nil {
		t.Fatal(err)
	}
	if second, err := acquireSupervisorOwner(statePath); err == nil {
		releaseSupervisorOwner(second)
		t.Fatal("a second supervisor acquired the same run")
	}
	releaseSupervisorOwner(first)
	third, err := acquireSupervisorOwner(statePath)
	if err != nil {
		t.Fatalf("ownership was not released: %v", err)
	}
	releaseSupervisorOwner(third)
}

func TestClosedReviewBecomesTerminal(t *testing.T) {
	s := &State{Config: Config{RunDir: t.TempDir()}}
	i := &Item{Number: 7, Status: "in_review"}
	pr := &Pull{State: "CLOSED", URL: "https://example.test/pull/7"}
	if !applyClosedReview(s, i, pr, i.Status) {
		t.Fatal("closed review was not handled")
	}
	if i.Status != "review_closed" || !terminalItem(i) {
		t.Fatalf("closed review status = %q", i.Status)
	}
}

func TestCleanupWaitsForWholeRunAndKeepsAuditLogs(t *testing.T) {
	root := filepath.Join(t.TempDir(), "worktrees")
	checkout := filepath.Join(root, "1-item")
	audit := filepath.Join(t.TempDir(), "run")
	if err := os.MkdirAll(checkout, 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := run("", nil, "git", "init", checkout); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(checkout, "fixture.txt"), []byte("kept\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := git(checkout, "add", "fixture.txt"); err != nil {
		t.Fatal(err)
	}
	if _, err := git(checkout, "-c", "user.name=Test", "-c", "user.email=test@example.test", "commit", "-m", "test: add fixture"); err != nil {
		t.Fatal(err)
	}
	head, err := git(checkout, "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(audit, 0o700); err != nil {
		t.Fatal(err)
	}
	logPath := filepath.Join(audit, "events.jsonl")
	if err := os.WriteFile(logPath, []byte("audit\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	i := &Item{Number: 1, Status: "integrated", Managed: true, Worktree: checkout, PR: &Pull{HeadOID: strings.TrimSpace(head)}}
	waiting := &Item{Number: 2, Status: "in_review"}
	s := &State{Config: Config{WorktreeRoot: root, RunDir: audit}, Items: map[string]*Item{"1": i, "2": waiting}}

	cleanupCompletedRun(s)
	if _, err := os.Stat(checkout); err != nil {
		t.Fatalf("checkout removed before run completion: %v", err)
	}

	waiting.Status = "review_closed"
	cleanupCompletedRun(s)
	if _, err := os.Stat(checkout); !os.IsNotExist(err) {
		t.Fatalf("terminal checkout still exists: %v", err)
	}
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		t.Fatalf("empty worktree root still exists: %v", err)
	}
	if _, err := os.Stat(logPath); err != nil {
		t.Fatalf("audit log was removed: %v", err)
	}
	if !runComplete(s) {
		t.Fatal("cleaned terminal run was not complete")
	}
}
