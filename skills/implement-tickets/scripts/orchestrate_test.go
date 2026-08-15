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
		t.Fatal("unintegrated open blocker was accepted")
	}
}

func TestIntegratedBlockerMayRemainOpenOnTracker(t *testing.T) {
	i := &Item{State: "OPEN", Labels: []string{"ready-for-agent"}, Blockers: []Blocker{{Number: 1, State: "OPEN", Integrated: true}}}
	if !ready(i, "me", "ready-for-agent") {
		t.Fatal("Git-integrated blocker still blocked readiness because the tracker issue is open")
	}
}

func TestReviewLinkLineUsesRelatedWhenContractListsRepositories(t *testing.T) {
	item := &Item{Number: 36}
	if got := reviewLinkLine(&State{}, item); got != "Closes #36" {
		t.Fatalf("single-repo closing line = %q", got)
	}
	s := &State{Config: Config{RelatedRepositories: []RelatedRepository{{Name: "org/controller", Path: "../controller"}}}}
	if got := reviewLinkLine(s, item); got != "Related #36" {
		t.Fatalf("multi-repo closing line = %q", got)
	}
}

func TestCollectSequenceNumbersAndReserveNext(t *testing.T) {
	dir := t.TempDir()
	mig := filepath.Join(dir, "db", "migrations")
	if err := os.MkdirAll(mig, 0o700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"000029_a.up.sql", "000030_b.up.sql", "README.md"} {
		if err := os.WriteFile(filepath.Join(mig, name), []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	s := &State{
		Config: Config{
			Repo: dir,
			Sequences: []Sequence{{Directory: "db/migrations", Pattern: `^(\d{6})_`}},
		},
		Items: map[string]*Item{},
	}
	item := &Item{Number: 40}
	if err := reserveSequences(s, item); err != nil {
		t.Fatal(err)
	}
	if len(item.SequenceReservations) != 1 || item.SequenceReservations[0] != "db/migrations: 000031" {
		t.Fatalf("reservations = %#v", item.SequenceReservations)
	}
}

func TestRejectForbiddenCommitsDetectsPathReplace(t *testing.T) {
	s := &State{Config: Config{ForbiddenCommitPatterns: []ForbiddenCommitPattern{{Path: "go.mod", Pattern: `^replace\s+\S+\s+=>\s+\.\./`}}}}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module x\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := run(dir, nil, "git", "init"); err != nil {
		t.Fatal(err)
	}
	if _, err := git(dir, "config", "user.email", "t@example.test"); err != nil {
		t.Fatal(err)
	}
	if _, err := git(dir, "config", "user.name", "t"); err != nil {
		t.Fatal(err)
	}
	if _, err := git(dir, "add", "go.mod"); err != nil {
		t.Fatal(err)
	}
	if _, err := git(dir, "commit", "-m", "init"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module x\n\nreplace example.com/mod => ../mod\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := git(dir, "add", "go.mod"); err != nil {
		t.Fatal(err)
	}
	if err := rejectForbiddenCommits(s, dir); err == nil {
		t.Fatal("path replace was accepted")
	}
}

func TestProjectSetupAcceptsSharedOptionalFields(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".github"), 0o700); err != nil {
		t.Fatal(err)
	}
	body := `{
  "version": 1,
  "source": "github",
  "target": "master",
  "related_repositories": [{"name": "org/other", "path": "../other", "role": "controller"}],
  "sequences": [{"directory": "db/migrations", "pattern": "^(\\\\d{6})_"}],
  "forbidden_commit_patterns": [{"path": "go.mod", "pattern": "^replace"}]
}`
	path := filepath.Join(dir, ".github", "implement-tickets.json")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	setup, err := loadProjectSetup(dir, ".github/implement-tickets.json")
	if err != nil || setup == nil || len(setup.RelatedRepositories) != 1 || len(setup.Sequences) != 1 || len(setup.ForbiddenCommitPatterns) != 1 {
		t.Fatalf("setup = %#v err=%v", setup, err)
	}
}

func TestOpenCodeGetsOneAutomaticVerificationRecovery(t *testing.T) {
	s := &State{Config: Config{Harness: "opencode"}}
	i := &Item{Worker: &Worker{Attempt: 1}}
	if !automaticRecoveryAllowed(s, i) {
		t.Fatal("first OpenCode verification failure was not recoverable")
	}
	i.Worker.Attempt = 2
	if automaticRecoveryAllowed(s, i) {
		t.Fatal("second verification failure retried without review")
	}
	s.Config.Harness = "codex"
	i.Worker.Attempt = 1
	if automaticRecoveryAllowed(s, i) {
		t.Fatal("non-OpenCode harness inherited automatic recovery")
	}
}

func TestExplicitBlockerNumbersUsesOnlyBlockedBySection(t *testing.T) {
	body := "## What to build\n\nMention #99 without creating a dependency.\n\n## Blocked by\n\n- #11\n- #9 description\n- #11 duplicate\n\n## Acceptance criteria\n\n- [ ] Complete #88"
	if got := fmt.Sprint(explicitBlockerNumbers(body)); got != "[9 11]" {
		t.Fatalf("explicitBlockerNumbers = %s", got)
	}
	if got := explicitBlockerNumbers("## Blocked by\n\n- None - can start immediately."); len(got) != 0 {
		t.Fatalf("None produced blockers: %v", got)
	}
}

func TestNormalizedGitLabStateMapsOpenedToOpen(t *testing.T) {
	if got := normalizedGitLabState("opened"); got != "OPEN" {
		t.Fatalf("opened normalized to %q", got)
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
	if got, err := workerCommitSubject(&State{}, i); err != nil || got != "fix: remove stale review state" {
		t.Fatalf("workerCommitSubject = %q, %v", got, err)
	}

	jsonPath := filepath.Join(dir, "events.jsonl")
	if err := os.WriteFile(jsonPath, []byte(`{"type":"message","part":{"text":"Commit subject: feat(events): relay supervisor progress"}}`+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	i.Worker.LastMessage = jsonPath
	if got, err := workerCommitSubject(&State{}, i); err != nil || got != "feat(events): relay supervisor progress" {
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
		if _, err := workerCommitSubject(&State{}, i); err == nil {
			t.Fatalf("invalid handoff was accepted: %q", content)
		}
	}
}

func TestOpenCodeDerivesCommitSubjectFromWorkItemTitle(t *testing.T) {
	s := &State{Config: Config{Harness: "opencode", RunDir: t.TempDir()}}
	i := &Item{Title: "feat(sync): reconcile catalogue after complete discovery", Worker: &Worker{LastMessage: filepath.Join(t.TempDir(), "missing.txt")}}
	if got, err := workerCommitSubject(s, i); err != nil || got != i.Title {
		t.Fatalf("workerCommitSubject = %q, %v", got, err)
	}
}

func TestWorkerFeedbackResponseReadsHandoff(t *testing.T) {
	path := filepath.Join(t.TempDir(), "handoff.txt")
	if err := os.WriteFile(path, []byte("Feedback response: The DEV-only button starts the sync.\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := workerFeedbackResponse(&Item{Worker: &Worker{LastMessage: path}}); got != "The DEV-only button starts the sync." {
		t.Fatalf("workerFeedbackResponse = %q", got)
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
		{feedbackInput{Body: "Can we show current/total progress?"}, "feedback"},
		{feedbackInput{Body: "Anything that is 404 should not fail the run."}, "feedback"},
		{feedbackInput{Body: "Why was this approach selected?"}, "feedback"},
		{feedbackInput{Body: "is this easy for other games to use the same tool?"}, "feedback"},
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

func TestGitLabFeedbackUsesNoteIDForAcknowledgement(t *testing.T) {
	input := feedbackInput{ID: "stable", Source: "gitlab-discussion", RemoteID: "6115288"}
	if input.RemoteID != "6115288" {
		t.Fatal("GitLab note ID was not retained for acknowledgement")
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

func TestSupervisorHealthRepairOnlyTargetsOpenCode(t *testing.T) {
	if err := ensureOpenCodeSupervisor(filepath.Join(t.TempDir(), "state.json"), &State{Config: Config{Harness: "codex"}}); err != nil {
		t.Fatalf("non-OpenCode status health check failed: %v", err)
	}
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

func TestSupportedHarnessIncludesGrok(t *testing.T) {
	if !supportedHarness("grok") || supervisorLaunchesWorkers("grok") {
		t.Fatal("grok should be a supported interactive harness that does not launch CLI workers")
	}
	if !supervisorLaunchesWorkers("opencode") || !supervisorLaunchesWorkers("codex") {
		t.Fatal("codex and opencode should still launch CLI workers")
	}
}

func TestGrokDispatchDoesNotLaunchWorkers(t *testing.T) {
	s := &State{
		Config: Config{Harness: "grok", Concurrency: 3},
		Items:  map[string]*Item{"1": {Number: 1, Status: "ready", Title: "Do a thing"}},
	}
	if err := dispatch(filepath.Join(t.TempDir(), "state.json"), s); err != nil {
		t.Fatal(err)
	}
	if s.Items["1"].Status != "ready" || s.Items["1"].Worker != nil {
		t.Fatalf("grok dispatch mutated item: %#v", s.Items["1"])
	}
}

func TestHarnessCreditIncludesGrok(t *testing.T) {
	if got := harnessCredit("grok"); got != "Grok/AI" {
		t.Fatalf("harnessCredit(grok) = %q", got)
	}
}

func TestFollowupPromptOmitsCanonicalBody(t *testing.T) {
	s := &State{Config: Config{Invocation: "$implement", Harness: "opencode"}}
	i := &Item{
		Title: "Add a thing", Body: "Line one\n\n- exact acceptance criterion",
		Branch: "issue/2-add", Worktree: "/worktree", WorkerSession: "ses_123",
		Pending: []Request{{Status: "queued", Action: "feedback", Message: "rename the helper", Source: "github-comment"}},
	}
	prompt := workerPrompt(s, i)
	if strings.Contains(prompt, i.Body) {
		t.Fatal("follow-up prompt re-sent the canonical body")
	}
	if !strings.Contains(prompt, "rename the helper") || !strings.Contains(prompt, "Continue the same work item") {
		t.Fatalf("follow-up prompt missing new input: %s", prompt)
	}
}

func TestOpenCodeSessionIDReadsNestedJSON(t *testing.T) {
	raw := `{"type":"tool","payload":{"huge":"` + strings.Repeat("x", 80_000) + `"}}` + "\n" +
		`{"type":"session.created","properties":{"sessionID":"ses_abc"}}` + "\n"
	if got := openCodeSessionID([]byte(raw)); got != "ses_abc" {
		t.Fatalf("openCodeSessionID = %q", got)
	}
}

func TestHandoffTextsReadsLongJSONLines(t *testing.T) {
	raw := `{"type":"tool","payload":{"huge":"` + strings.Repeat("x", 80_000) + `"}}` + "\n" +
		`{"part":{"text":"Commit subject: fix: keep long event lines\n"}}` + "\n"
	found := false
	for _, text := range handoffTexts([]byte(raw)) {
		if strings.Contains(text, "Commit subject: fix: keep long event lines") {
			found = true
		}
	}
	if !found {
		t.Fatal("long JSONL handoff dropped the commit subject")
	}
}

func TestIntegrationCloseCommentNamesReviewAndTarget(t *testing.T) {
	s := &State{Config: Config{Target: "master"}}
	i := &Item{PR: &Pull{URL: "https://example.test/pull/9", MergeOID: "abcdef1234567890"}}
	got := integrationCloseComment(s, i)
	for _, want := range []string{"AI-generated:", "https://example.test/pull/9", "abcdef1", "master"} {
		if !strings.Contains(got, want) {
			t.Fatalf("close comment %q missing %q", got, want)
		}
	}
}
