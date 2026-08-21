package main

import (
	"os"
	"path/filepath"
	"testing"
)

// exploitName is the shape of the payload file whose NAME is the shell
// injection: leading semicolon, $(), ${}, spaces, and a TRAILING SPACE.
//
// The trailing space is not decoration. It is the character most likely to be
// eaten by a shell pipeline, a file manager, or an archive round-trip — and if
// it is lost the jailbreak fails silently. No literal "/" appears because a
// filename cannot contain one; that is precisely why the exploit assigns SLASH.
const exploitName = `;export SLASH=${HOME%${HOME#?}};$(sh ${SLASH}mnt${SLASH}us${SLASH}jb) `

// TestPayloadPreservesExploitFilename is the reason this Go program exists.
// If it ever fails, the copy path has grown a shell somewhere.
func TestPayloadPreservesExploitFilename(t *testing.T) {
	src := t.TempDir()
	dict := filepath.Join(src, "documents", "dictionaries")
	if err := os.MkdirAll(dict, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dict, exploitName), []byte("payload"), 0o644); err != nil {
		t.Fatalf("creating the exploit fixture: %v", err)
	}

	dst := t.TempDir()
	if _, err := copyTree(src, dst); err != nil {
		t.Fatalf("copyTree: %v", err)
	}

	got, found := findExploit(filepath.Join(dst, "documents", "dictionaries"))
	if !found {
		t.Fatal("exploit file absent after copy")
	}
	if got != exploitName {
		t.Errorf("filename altered by the copy\n want %q\n  got %q", exploitName, got)
	}
}

// TestPayloadSkipsAppleDouble — FAT32 cannot hold these, and macOS creates them
// freely. Copying them through is how a payload ends up subtly wrong.
func TestPayloadSkipsAppleDouble(t *testing.T) {
	src := t.TempDir()
	for _, n := range []string{"real", "._sidecar"} {
		if err := os.WriteFile(filepath.Join(src, n), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	dst := t.TempDir()
	if _, err := copyTree(src, dst); err != nil {
		t.Fatal(err)
	}
	if !exists(filepath.Join(dst, "real")) {
		t.Error("real file was not copied")
	}
	if exists(filepath.Join(dst, "._sidecar")) {
		t.Error("AppleDouble sidecar was copied through")
	}
	if s := strayDouble(dst); len(s) != 0 {
		t.Errorf("strayDouble reported %v after a clean copy", s)
	}
}

// TestPayloadCopiesHiddenDemo — .demo/ is hidden, so a drag-copy in Finder
// leaves it behind and the jailbreak does nothing. It must come across.
func TestPayloadCopiesHiddenDemo(t *testing.T) {
	src := t.TempDir()
	if err := os.MkdirAll(filepath.Join(src, ".demo"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, ".demo", "marker"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	dst := t.TempDir()
	if _, err := copyTree(src, dst); err != nil {
		t.Fatal(err)
	}
	if !exists(filepath.Join(dst, ".demo", "marker")) {
		t.Error(".demo/ did not survive the copy — this is the Finder failure mode")
	}
}

// TestIsKindleDetectsByContent — volume labels are user-changeable and
// localised, so detection must not depend on the name.
func TestIsKindleDetectsByContent(t *testing.T) {
	vol := t.TempDir()
	if isKindle(vol) {
		t.Error("empty directory reported as a Kindle")
	}
	for _, d := range []string{"documents", "system"} {
		if err := os.MkdirAll(filepath.Join(vol, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if !isKindle(vol) {
		t.Error("directory with documents/ and system/ not recognised")
	}
}
