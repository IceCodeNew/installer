package handler

import "testing"

func TestFilExt(t *testing.T) {
	tests := []struct {
		file, ext string
	}{
		{"my.file.tar.gz", ".tar.gz"},
		{"my.file.tar.bz2", ".tar.bz2"},
		{"my.file.tar.bz", ".tar.bz"},
		{"my.file.bz2", ".bz2"},
		{"my.file.gz", ".gz"},
		{"my.file.tar.zip", ".tar.zip"}, // :(
	}
	for _, tc := range tests {
		ext := getFileExt(tc.file)
		if ext != tc.ext {
			t.Fatalf("getFileExt(%s) = %s, want %s", tc.file, ext, tc.ext)
		}
	}
}

func TestOSArch(t *testing.T) {
	for _, tc := range []struct {
		name, os, arch string
	}{
		// m"arm"ite
		{"marmite-0.2.5-x86_64-unknown-linux-musl.tar.gz", "linux", "amd64"},
		{"yt-dlp_linux", "linux", "amd64"},
		{"yt-dlp_linux_armv7l", "linux", "arm"},
		{"gitleaks_8.24.0_linux_x64.tar.gz", "linux", "amd64"},
		{"gitleaks_8.24.0_linux_x32.tar.gz", "linux", "386"},
		{"gitleaks_8.24.0_linux_armv6.tar.gz", "linux", "arm"},
		{"gitleaks_8.24.0_linux_arm64.tar.gz", "linux", "arm64"},
		{"gitleaks_8.24.0_darwin_x64.tar.gz", "darwin", "amd64"},
		{"gitleaks_8.24.0_darwin_arm64.tar.gz", "darwin", "arm64"},
		{"gitleaks_8.24.0_windows_armv6.zip", "windows", "arm"},
		{"gitleaks_8.24.0_windows_x32.zip", "windows", "386"},
		{"gitleaks_8.24.0_windows_x64.zip", "windows", "amd64"},
		{"gitui-linux-x86_64.tar.gz", "linux", "amd64"},
		{"gitui-linux-arm.tar.gz", "linux", "arm"},
		{"gitui-linux-aarch64.tar.gz", "linux", "arm64"},
		{"gg-linux-x86_64", "linux", "amd64"},
		{"gg-linux-armv5", "linux", "arm"},
		{"gg-linux-arm64", "linux", "arm64"},
		{"croc_v10.2.1_Linux-64bit.tar.gz", "linux", "amd64"},
		{"croc_v10.2.1_Linux-32bit.tar.gz", "linux", "386"},
		{"croc_v10.2.1_Linux-ARM.tar.gz", "linux", "arm"},
		{"croc_v10.2.1_Linux-ARM64.tar.gz", "linux", "arm64"},
		{"ouch-x86_64-unknown-linux-musl.tar.gz", "linux", "amd64"},
		{"ouch-armv7-unknown-linux-musleabihf.tar.gz", "linux", "arm"},
		{"ouch-aarch64-unknown-linux-musl.tar.gz", "linux", "arm64"},
		{"crun-1.20-linux-ppc64le", "linux", "ppc64le"},
		{"crun-1.20-linux-riscv64", "linux", "riscv64"},
		{"crun-1.20-linux-s390x", "linux", "s390x"},
		{"libtree_aarch64", "linux", "arm64"},
		{"libtree_armv6l", "linux", "arm"},
		{"libtree_i686", "linux", "386"},
		{"libtree_x86_64", "linux", "amd64"},
		{"runc.armel", "linux", "arm"},
		{"runc.armhf", "linux", "arm"},
		{"runc.ppc64le", "linux", "ppc64le"},
		{"runc.riscv64", "linux", "riscv64"},
		{"runc.s390x", "linux", "s390x"},
	} {
		os := getOS(tc.name)
		arch := getArch(tc.name)
		if os != tc.os || arch != tc.arch {
			t.Fatalf("file '%s' results in %s/%s, expected %s/%s", tc.name, os, arch, tc.os, tc.arch)
		}
	}
}
