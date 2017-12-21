package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_snippetCmd(t *testing.T) {
	var snipID string
	t.Run("create_personal", func(t *testing.T) {
		repo := copyTestRepo(t)
		cmd := exec.Command("../lab_bin", "snippet", "create", "-g",
			"-m", "personal snippet title",
			"-m", "personal snippet description")
		cmd.Dir = repo

		rc, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}

		_, err = rc.Write([]byte("personal snippet contents"))
		if err != nil {
			t.Fatal(err)
		}
		err = rc.Close()
		if err != nil {
			t.Fatal(err)
		}

		b, err := cmd.CombinedOutput()
		if err != nil {
			t.Log(string(b))
			t.Fatal(err)
		}

		out := string(b)
		require.Contains(t, out, "https://gitlab.com/snippets/")

		i := strings.Index(out, "\n")
		snipID = strings.TrimPrefix(out[:i], "https://gitlab.com/snippets/")
		t.Log(snipID)
	})
	t.Run("delete_personal", func(t *testing.T) {
		if snipID == "" {
			t.Skip("snipID is empty, create likely failed")
		}
		repo := copyTestRepo(t)
		cmd := exec.Command("../lab_bin", "snippet", "-g", "-d", snipID)
		cmd.Dir = repo

		b, err := cmd.CombinedOutput()
		if err != nil {
			t.Log(string(b))
			t.Fatal(err)
		}
		require.Contains(t, string(b), fmt.Sprintf("Snippet #%s deleted", snipID))
	})
}

func Test_snippetCmd_noArgs(t *testing.T) {
	repo := copyTestRepo(t)
	cmd := exec.Command("../lab_bin", "snippet")
	cmd.Dir = repo

	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}
	require.Contains(t, string(b), `Usage:
  lab snippet [flags]
  lab snippet [command]`)
}