package cms

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// runGitTest runs a git command in dir and fails the test on error.
func runGitTest(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func TestGitCommitAndPush(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	// A bare repository acting as the remote.
	remote := t.TempDir()
	runGitTest(t, remote, "init", "--bare", "--initial-branch=main")

	// A working repository with an initial commit pushed to the remote.
	work := t.TempDir()
	runGitTest(t, work, "init", "--initial-branch=main")
	runGitTest(t, work, "config", "user.email", "test@doolox.example")
	runGitTest(t, work, "config", "user.name", "Doolox Test")
	if err := os.WriteFile(filepath.Join(work, "index.html"), []byte("v1"), 0644); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, work, "add", ".")
	runGitTest(t, work, "commit", "-m", "initial")
	runGitTest(t, work, "remote", "add", "origin", remote)
	runGitTest(t, work, "push", "-u", "origin", "main")

	// Change a tracked file, then commit and push it.
	if err := os.WriteFile(filepath.Join(work, "index.html"), []byte("v2"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := gitCommitAndPush(work); err != nil {
		t.Fatalf("gitCommitAndPush() error = %v", err)
	}

	// The remote HEAD should now be the "Doolox update." commit.
	if got := runGitTest(t, remote, "log", "-1", "--format=%s"); strings.TrimSpace(got) != publishCommitMessage {
		t.Errorf("remote HEAD subject = %q, want %q", strings.TrimSpace(got), publishCommitMessage)
	}

	// With no changes, "nothing to commit" must not be treated as an error.
	if err := gitCommitAndPush(work); err != nil {
		t.Errorf("gitCommitAndPush() with clean tree error = %v, want nil", err)
	}
}
