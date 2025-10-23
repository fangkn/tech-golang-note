package version

import "testing"

func TestVersion(t *testing.T) {
	Info.Version = "v0.0.1"
	Info.BuildTime = "2017-11-11T00:00:00+0800"
	Info.GitCommit = "2c5cb13"
	Info.GitBranch = "test"
	Info.GoVersion = "1.8"

	want := "Version: \033[32mv0.0.1\033[0m\nGitCommit: \033[33m2c5cb13\033[0m\nGitBranch: \033[36mtest\033[0m\nGoVersion: 1.8\nBuildTime: 2017-11-11T00:00:00+0800\n"

	if got := Version(); got != want {
		t.Errorf("Version() = %q; want %q", got, want)
	}
}

func TestBuildInfoString(t *testing.T) {
	Info.Version = "v0.0.1"
	Info.GitCommit = "2c5cb13"
	Info.GitBranch = "test"

	want := `v0.0.1 2c5cb13`

	if got := Info.String(); got != want {
		t.Errorf("BuildInfo.String() = %q; want %q", got, want)
	}
}
