package comparable

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpecTestsDirFromAnchorsToGoModRoot(t *testing.T) {
	repoRoot := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(repoRoot, "go.mod"), []byte("module example.com/test\n"), 0600))

	specTestsDir, err := SpecTestsDirFrom(filepath.Join(repoRoot, "qbft", "spectest"))
	require.NoError(t, err)
	require.Equal(t, filepath.Join(repoRoot, "..", "spec-tests", "qbft"), specTestsDir)

	_, err = SpecTestsDirFrom(repoRoot)
	require.Error(t, err)
}

func TestEnsureSpecTestsSubdir(t *testing.T) {
	repoRoot := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(repoRoot, "go.mod"), []byte("module example.com/test\n"), 0600))
	wd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(repoRoot))
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(wd))
	})
	repoRoot, err = os.Getwd()
	require.NoError(t, err)

	require.NoError(t, EnsureSpecTestsSubdir("qbft", filepath.Join(repoRoot, "..", "spec-tests", "qbft", "tests")))
	require.Error(t, EnsureSpecTestsSubdir("qbft", filepath.Join(repoRoot, "..", "spec-tests", "qbft")))
	require.Error(t, EnsureSpecTestsSubdir("qbft", filepath.Join(repoRoot, "..", "spec-tests", "ssv", "tests")))
}
