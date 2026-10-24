package cms

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// publishCommitMessage is the commit message used when publishing the site.
const publishCommitMessage = "Doolox update."

// PublishSite generates the static website and then commits and pushes the
// changes with git in the main website's folder. It shows a native system
// notification with the result and returns a short human-readable summary or
// an error.
func PublishSite() (string, error) {
	msg, err := publish()
	if err != nil {
		log.Printf("[publish] error: %v", err)
		notify("Doolox", fmt.Sprintf("Publish failed: %v", err))
		return msg, err
	}

	log.Printf("[publish] %s", msg)
	notify("Doolox", msg)
	return msg, nil
}

// publish generates the site (without the generation notification, since the
// publish notification covers the whole flow), then commits and pushes it.
func publish() (string, error) {
	generated, err := generate()
	if err != nil {
		return "", fmt.Errorf("generation failed: %w", err)
	}

	if err := gitCommitAndPush(websiteDir(Website())); err != nil {
		return "", err
	}

	return fmt.Sprintf("Site published: %d pages generated, committed and pushed", generated), nil
}

// gitCommitAndPush stages all changes (including new untracked files and
// deletions) with `git add -A`, commits them with `git commit -m <message>`,
// and then runs `git push` in the given directory. A "nothing to commit"
// result is not treated as an error: the tree is already up to date, so the
// push still runs.
func gitCommitAndPush(dir string) error {
	if out, err := runGit(dir, "add", "-A"); err != nil {
		return fmt.Errorf("git add: %v: %s", err, strings.TrimSpace(out))
	}

	if out, err := runGit(dir, "commit", "-m", publishCommitMessage); err != nil {
		if !strings.Contains(strings.ToLower(out), "nothing to commit") {
			return fmt.Errorf("git commit: %v: %s", err, strings.TrimSpace(out))
		}
	}

	if out, err := runGit(dir, "push"); err != nil {
		return fmt.Errorf("git push: %v: %s", err, strings.TrimSpace(out))
	}

	return nil
}

// runGit runs a git command in the given directory and returns its combined
// output and error.
func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// publishSite generates the static website and pushes it with git. It is
// called from the system tray menu.
func publishSite() {
	PublishSite()
}
