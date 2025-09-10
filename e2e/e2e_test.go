package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGithubReleasesReal hits real GitHub release API endpoints and asserts
// that expected asset names or URLs exist for those releases. Skipped when
// -short is used or E2E_SKIP=1 is set.
func TestGithubReleasesReal(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}
	if os.Getenv("E2E_SKIP") == "1" {
		t.Skip("E2E_SKIP=1: skipping external github tests")
	}

	cases := []struct {
		owner string
		repo  string
		tag   string
		want  string
	}{
		// uv 0.8.16
		{"astral-sh", "uv", "0.8.16", "uv-x86_64-unknown-linux-musl.tar.gz"},
		{"astral-sh", "uv", "0.8.16", "uv-x86_64-apple-darwin.tar.gz"},
		{"astral-sh", "uv", "0.8.16", "uv-x86_64-pc-windows-msvc.zip"},
		{"astral-sh", "uv", "0.8.16", "uv-loongarch64-unknown-linux-gnu.tar.gz"},
		{"astral-sh", "uv", "0.8.16", "uv-powerpc64-unknown-linux-gnu.tar.gz"},
		{"astral-sh", "uv", "0.8.16", "uv-powerpc64le-unknown-linux-gnu.tar.gz"},
		{"astral-sh", "uv", "0.8.16", "uv-riscv64gc-unknown-linux-gnu.tar.gz"},
		{"astral-sh", "uv", "0.8.16", "uv-s390x-unknown-linux-gnu.tar.gz"},
		{"astral-sh", "uv", "0.8.16", "uv-armv7-unknown-linux-musleabihf.tar.gz"},
		{"astral-sh", "uv", "0.8.16", "uv-i686-unknown-linux-musl.tar.gz"},
		{"astral-sh", "uv", "0.8.16", "uv-aarch64-unknown-linux-musl.tar.gz"},
		{"astral-sh", "uv", "0.8.16", "uv-aarch64-apple-darwin.tar.gz"},
		{"astral-sh", "uv", "0.8.16", "uv-aarch64-pc-windows-msvc.zip"},
		{"astral-sh", "uv", "0.8.16", "uv-i686-pc-windows-msvc.zip"},

		// tailspin 5.5.0
		{"bensadeh", "tailspin", "5.5.0", "tailspin-x86_64-apple-darwin.tar.gz"},
		{"bensadeh", "tailspin", "5.5.0", "tailspin-aarch64-apple-darwin.tar.gz"},
		{"bensadeh", "tailspin", "5.5.0", "tailspin-x86_64-unknown-linux-musl.tar.gz"},
		{"bensadeh", "tailspin", "5.5.0", "tailspin-aarch64-unknown-linux-musl.tar.gz"},

		// sd v1.0.0
		{"chmln", "sd", "v1.0.0", "sd-v1.0.0-x86_64-apple-darwin.tar.gz"},
		{"chmln", "sd", "v1.0.0", "sd-v1.0.0-aarch64-apple-darwin.tar.gz"},
		{"chmln", "sd", "v1.0.0", "sd-v1.0.0-x86_64-unknown-linux-musl.tar.gz"},
		{"chmln", "sd", "v1.0.0", "sd-v1.0.0-arm-unknown-linux-gnueabihf.tar.gz"},
		{"chmln", "sd", "v1.0.0", "sd-v1.0.0-aarch64-unknown-linux-musl.tar.gz"},

		// delta 0.18.2
		{"dandavison", "delta", "0.18.2", "delta-0.18.2-x86_64-apple-darwin.tar.gz"},
		{"dandavison", "delta", "0.18.2", "delta-0.18.2-aarch64-apple-darwin.tar.gz"},
		{"dandavison", "delta", "0.18.2", "delta-0.18.2-i686-unknown-linux-gnu.tar.gz"},
		{"dandavison", "delta", "0.18.2", "delta-0.18.2-x86_64-unknown-linux-musl.tar.gz"},
		{"dandavison", "delta", "0.18.2", "delta-0.18.2-arm-unknown-linux-gnueabihf.tar.gz"},
		{"dandavison", "delta", "0.18.2", "delta-0.18.2-aarch64-unknown-linux-gnu.tar.gz"},

		// eza v0.23.2
		{"eza-community", "eza", "v0.23.2", "eza_x86_64-unknown-linux-musl.tar.gz"},
		{"eza-community", "eza", "v0.23.2", "eza_aarch64-unknown-linux-gnu.tar.gz"},

		// age v1.2.1
		{"FiloSottile", "age", "v1.2.1", "age-v1.2.1-darwin-amd64.tar.gz"},
		{"FiloSottile", "age", "v1.2.1", "age-v1.2.1-darwin-arm64.tar.gz"},
		{"FiloSottile", "age", "v1.2.1", "age-v1.2.1-freebsd-amd64.tar.gz"},
		{"FiloSottile", "age", "v1.2.1", "age-v1.2.1-linux-amd64.tar.gz"},
		{"FiloSottile", "age", "v1.2.1", "age-v1.2.1-linux-arm.tar.gz"},
		{"FiloSottile", "age", "v1.2.1", "age-v1.2.1-linux-arm64.tar.gz"},

		// gitleaks v8.28.0
		{"gitleaks", "gitleaks", "v8.28.0", "gitleaks_8.28.0_darwin_x64.tar.gz"},
		{"gitleaks", "gitleaks", "v8.28.0", "gitleaks_8.28.0_darwin_arm64.tar.gz"},
		{"gitleaks", "gitleaks", "v8.28.0", "gitleaks_8.28.0_linux_x32.tar.gz"},
		{"gitleaks", "gitleaks", "v8.28.0", "gitleaks_8.28.0_linux_x64.tar.gz"},
		{"gitleaks", "gitleaks", "v8.28.0", "gitleaks_8.28.0_linux_armv6.tar.gz"},
		{"gitleaks", "gitleaks", "v8.28.0", "gitleaks_8.28.0_linux_arm64.tar.gz"},

		// gitui v0.27.0
		{"gitui-org", "gitui", "v0.27.0", "gitui-mac-x86.tar.gz"},
		{"gitui-org", "gitui", "v0.27.0", "gitui-linux-x86_64.tar.gz"},
		{"gitui-org", "gitui", "v0.27.0", "gitui-linux-arm.tar.gz"},
		{"gitui-org", "gitui", "v0.27.0", "gitui-linux-aarch64.tar.gz"},

		// piknik 0.10.2 (no linux/arm64 asset)
		{"jedisct1", "piknik", "0.10.2", "piknik-macos-0.10.2.tar.gz"},
		{"jedisct1", "piknik", "0.10.2", "piknik-dragonflybsd_amd64-0.10.2.tar.gz"},
		{"jedisct1", "piknik", "0.10.2", "piknik-freebsd_i386-0.10.2.tar.gz"},
		{"jedisct1", "piknik", "0.10.2", "piknik-freebsd_amd64-0.10.2.tar.gz"},
		{"jedisct1", "piknik", "0.10.2", "piknik-linux_i386-0.10.2.tar.gz"},
		{"jedisct1", "piknik", "0.10.2", "piknik-linux_x86_64-0.10.2.tar.gz"},
		{"jedisct1", "piknik", "0.10.2", "piknik-linux_arm-0.10.2.tar.gz"},
		{"jedisct1", "piknik", "0.10.2", "piknik-netbsd_i386-0.10.2.tar.gz"},
		{"jedisct1", "piknik", "0.10.2", "piknik-netbsd_amd64-0.10.2.tar.gz"},
		{"jedisct1", "piknik", "0.10.2", "piknik-openbsd_i386-0.10.2.tar.gz"},
		{"jedisct1", "piknik", "0.10.2", "piknik-openbsd_amd64-0.10.2.tar.gz"},

		// lazygit v0.55.0
		{"jesseduffield", "lazygit", "v0.55.0", "lazygit_0.55.0_darwin_x86_64.tar.gz"},
		{"jesseduffield", "lazygit", "v0.55.0", "lazygit_0.55.0_darwin_arm64.tar.gz"},
		{"jesseduffield", "lazygit", "v0.55.0", "lazygit_0.55.0_freebsd_32-bit.tar.gz"},
		{"jesseduffield", "lazygit", "v0.55.0", "lazygit_0.55.0_freebsd_x86_64.tar.gz"},
		{"jesseduffield", "lazygit", "v0.55.0", "lazygit_0.55.0_freebsd_armv6.tar.gz"},
		{"jesseduffield", "lazygit", "v0.55.0", "lazygit_0.55.0_freebsd_arm64.tar.gz"},
		{"jesseduffield", "lazygit", "v0.55.0", "lazygit_0.55.0_linux_32-bit.tar.gz"},
		{"jesseduffield", "lazygit", "v0.55.0", "lazygit_0.55.0_linux_x86_64.tar.gz"},
		{"jesseduffield", "lazygit", "v0.55.0", "lazygit_0.55.0_linux_armv6.tar.gz"},
		{"jesseduffield", "lazygit", "v0.55.0", "lazygit_0.55.0_linux_arm64.tar.gz"},

		// codex
		{"openai", "codex", "rust-v0.31.0", "codex-x86_64-apple-darwin.tar.gz"},
		{"openai", "codex", "rust-v0.31.0", "codex-aarch64-apple-darwin.tar.gz"},
		{"openai", "codex", "rust-v0.31.0", "codex-x86_64-unknown-linux-musl.tar.gz"},
		{"openai", "codex", "rust-v0.31.0", "codex-aarch64-unknown-linux-musl.tar.gz"},

		// ouch 0.6.1 (no darwin/arm64)
		{"ouch-org", "ouch", "0.6.1", "ouch-x86_64-apple-darwin.tar.gz"},
		{"ouch-org", "ouch", "0.6.1", "ouch-x86_64-unknown-linux-musl.tar.gz"},
		{"ouch-org", "ouch", "0.6.1", "ouch-armv7-unknown-linux-musleabihf.tar.gz"},
		{"ouch-org", "ouch", "0.6.1", "ouch-aarch64-unknown-linux-musl.tar.gz"},

		// croc v10.2.4
		{"schollz", "croc", "v10.2.4", "croc_v10.2.4_macOS-64bit.tar.gz"},
		{"schollz", "croc", "v10.2.4", "croc_v10.2.4_macOS-ARM64.tar.gz"},
		{"schollz", "croc", "v10.2.4", "croc_v10.2.4_DragonFlyBSD-64bit.tar.gz"},
		{"schollz", "croc", "v10.2.4", "croc_v10.2.4_FreeBSD-64bit.tar.gz"},
		{"schollz", "croc", "v10.2.4", "croc_v10.2.4_FreeBSD-ARM64.tar.gz"},
		{"schollz", "croc", "v10.2.4", "croc_v10.2.4_Linux-32bit.tar.gz"},
		{"schollz", "croc", "v10.2.4", "croc_v10.2.4_Linux-64bit.tar.gz"},
		{"schollz", "croc", "v10.2.4", "croc_v10.2.4_Linux-ARM.tar.gz"},
		{"schollz", "croc", "v10.2.4", "croc_v10.2.4_Linux-ARM64.tar.gz"},
		{"schollz", "croc", "v10.2.4", "croc_v10.2.4_NetBSD-32bit.tar.gz"},
		{"schollz", "croc", "v10.2.4", "croc_v10.2.4_NetBSD-64bit.tar.gz"},
		{"schollz", "croc", "v10.2.4", "croc_v10.2.4_NetBSD-ARM64.tar.gz"},
		{"schollz", "croc", "v10.2.4", "croc_v10.2.4_OpenBSD-64bit.tar.gz"},
		{"schollz", "croc", "v10.2.4", "croc_v10.2.4_OpenBSD-ARM64.tar.gz"},

		// jid v0.7.6 (no darwin/arm64)
		{"simeji", "jid", "v0.7.6", "jid_darwin_386.zip"},
		{"simeji", "jid", "v0.7.6", "jid_darwin_amd64.zip"},
		{"simeji", "jid", "v0.7.6", "jid_freebsd_386.zip"},
		{"simeji", "jid", "v0.7.6", "jid_freebsd_amd64.zip"},
		{"simeji", "jid", "v0.7.6", "jid_linux_386.zip"},
		{"simeji", "jid", "v0.7.6", "jid_linux_amd64.zip"},
		{"simeji", "jid", "v0.7.6", "jid_linux_arm64.zip"},
		{"simeji", "jid", "v0.7.6", "jid_netbsd_386.zip"},
		{"simeji", "jid", "v0.7.6", "jid_netbsd_amd64.zip"},
		{"simeji", "jid", "v0.7.6", "jid_openbsd_386.zip"},
		{"simeji", "jid", "v0.7.6", "jid_openbsd_amd64.zip"},

		// gtrash v0.0.6
		{"umlx5h", "gtrash", "v0.0.6", "gtrash_Darwin_x86_64.tar.gz"},
		{"umlx5h", "gtrash", "v0.0.6", "gtrash_Darwin_arm64.tar.gz"},
		{"umlx5h", "gtrash", "v0.0.6", "gtrash_Linux_i386.tar.gz"},
		{"umlx5h", "gtrash", "v0.0.6", "gtrash_Linux_x86_64.tar.gz"},
		{"umlx5h", "gtrash", "v0.0.6", "gtrash_Linux_arm64.tar.gz"},

		// dive v0.13.1
		{"wagoodman", "dive", "v0.13.1", "dive_0.13.1_darwin_amd64.tar.gz"},
		{"wagoodman", "dive", "v0.13.1", "dive_0.13.1_darwin_arm64.tar.gz"},
		{"wagoodman", "dive", "v0.13.1", "dive_0.13.1_linux_amd64.tar.gz"},
		{"wagoodman", "dive", "v0.13.1", "dive_0.13.1_linux_arm64.tar.gz"},

		// yt-dlp 2025.09.05
		{"yt-dlp", "yt-dlp", "2025.09.05", "yt-dlp_macos"},
		{"yt-dlp", "yt-dlp", "2025.09.05", "yt-dlp_linux"},
		{"yt-dlp", "yt-dlp", "2025.09.05", "yt-dlp_linux_armv7l.zip"},
		{"yt-dlp", "yt-dlp", "2025.09.05", "yt-dlp_linux_aarch64"},
	}

	for _, c := range cases {
		rel, err := getRelease(t, c.owner, c.repo, c.tag)
		if err != nil {
			t.Fatalf("getRelease %s/%s@%s: %v", c.owner, c.repo, c.tag, err)
		}
		found := false
		for _, a := range rel.Assets {
			if strings.Contains(a.Name, c.want) ||
				strings.Contains(a.BrowserDownloadURL, c.want) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected to find %q in assets for %s/%s@%s", c.want, c.owner, c.repo, c.tag)
		}
	}
}

// releaseModel is a lightweight subset of the GitHub release JSON we need.
type releaseModel struct {
	Assets []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func fixturePath(owner, repo, tag string) string {
	name := fmt.Sprintf("%s_%s_%s.json", owner, repo, tag)
	// sanitize
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, " ", "_")
	return filepath.Join("fixtures", name)
}

// getRelease will either replay from fixture (when E2E_REPLAY=1), record
// to fixture (when E2E_RECORD=1), or perform a live request. If no token is
// present and no fixture exists, it instructs how to record one.
func getRelease(t *testing.T, owner, repo, tag string) (*releaseModel, error) {
	path := fixturePath(owner, repo, tag)
	replay := os.Getenv("E2E_REPLAY") == "1"
	record := os.Getenv("E2E_RECORD") == "1"
	tok := os.Getenv("GITHUB_TOKEN")

	if replay {
		// must exist
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("replay requested but fixture missing: %s (run with E2E_RECORD=1 to create)", path)
		}
		defer func() { _ = f.Close() }()
		var r releaseModel
		if err := json.NewDecoder(f).Decode(&r); err != nil {
			return nil, fmt.Errorf("decode fixture %s: %v", path, err)
		}
		return &r, nil
	}

	// If not recording and token is empty, prefer fixture if exists
	if !record && tok == "" {
		if _, err := os.Stat(path); err == nil {
			f, _ := os.Open(path)
			defer func() { _ = f.Close() }()
			var r releaseModel
			if err := json.NewDecoder(f).Decode(&r); err != nil {
				return nil, fmt.Errorf("decode fixture %s: %v", path, err)
			}
			return &r, nil
		}
		// no fixture and no token -> instruct
		return nil, fmt.Errorf("no GITHUB_TOKEN and no fixture %s; either set GITHUB_TOKEN or run tests with E2E_RECORD=1 to record fixtures", path)
	}

	// perform live request
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", owner, repo, tag)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if tok != "" {
		req.Header.Set("Authorization", "token "+tok)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("release not found %s", url)
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GET %s returned status %d: %s", url, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	// read whole body (we may save it)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// save fixture if requested
	if record {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err == nil {
			_ = os.WriteFile(path, data, 0o644)
		}
	}

	var r releaseModel
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("decode live body: %v", err)
	}
	return &r, nil
}
