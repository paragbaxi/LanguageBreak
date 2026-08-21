// Command lbtool is a host-side companion to the shell scripts in
// troubleshooting/. It does the two jobs that happen on YOUR computer while the
// Kindle is plugged in over USB:
//
//	lbtool verify  [volume]        did the jailbreak actually take?
//	lbtool payload <src> [volume]  copy the payload correctly, then prove it
//
// The other two fixes (fix-managed-mode.sh, fix-usb-disabled.sh) are NOT here
// and cannot be: they delete files under /var/local as root and reboot, so they
// run ON the device via /mnt/us/emergency.sh. Shell is the right and only tool
// for those.
//
// WHY A GO VERSION EXISTS AT ALL
// One file in the payload has a NAME that is the shell injection carrying the
// exploit — a leading "a;", $(), ${} expansions, no trailing space, 103 bytes
// exactly. Every hop through a shell, a GUI file manager, or an archive tool is
// a chance to alter it, and the failure is SILENT — the jailbreak simply does
// nothing and you are left debugging the device. Go moves filenames as bytes;
// there is no shell in the path. main_test.go round-trips that exact name.
//
// The shell scripts remain the primary, dependency-free path. This is for
// people who would rather run one binary than trust their file manager, and it
// needs no toolchain to READ — only to build.
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const usage = `lbtool — host-side LanguageBreak helpers (see troubleshooting/*.sh for the rest)

  lbtool verify [volume]           check the artefacts the exploit writes
  lbtool payload <src> [volume]    copy the payload to the Kindle, then verify it

[volume] is optional; the Kindle is found automatically by looking for a mounted
volume containing documents/ and system/. Pass it explicitly for unusual mounts.
`

// The five entries that must exist at the volume root after a correct copy.
var payloadRoot = []string{".demo", "documents", "DONT_CHECK_BATTERY", "jb", "patchedUks.sqsh"}

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "verify":
		err = cmdVerify(os.Args[2:])
	case "payload":
		err = cmdPayload(os.Args[2:])
	case "-h", "--help", "help":
		fmt.Print(usage)
		return
	default:
		err = fmt.Errorf("unknown command %q\n\n%s", os.Args[1], strings.TrimSpace(usage))
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

// findVolume locates the mounted Kindle by CONTENT rather than by name: the
// volume label is user-changeable and localised, but a Kindle always carries
// documents/ and system/ at its root.
func findVolume() (string, error) {
	var roots []string
	if runtime.GOOS == "darwin" {
		roots = []string{"/Volumes"}
	} else {
		u := os.Getenv("USER")
		roots = []string{"/media/" + u, "/run/media/" + u, "/media", "/mnt"}
	}
	for _, root := range roots {
		entries, err := os.ReadDir(root)
		if err != nil {
			continue
		}
		for _, e := range entries {
			cand := filepath.Join(root, e.Name())
			if isKindle(cand) {
				return cand, nil
			}
		}
	}
	return "", fmt.Errorf("no mounted Kindle found — plug it in and confirm the volume appears.\n" +
		"If the device enumerates on USB but no volume ever mounts, that is the demo-mode\n" +
		"USB fault: see troubleshooting/fix-usb-disabled.sh")
}

func isKindle(dir string) bool {
	for _, m := range []string{"documents", "system"} {
		st, err := os.Stat(filepath.Join(dir, m))
		if err != nil || !st.IsDir() {
			return false
		}
	}
	return true
}

func volumeArg(args []string) (string, error) {
	if len(args) > 0 && args[0] != "" {
		if !isKindle(args[0]) {
			return "", fmt.Errorf("%s does not look like a Kindle volume (no documents/ and system/)", args[0])
		}
		return args[0], nil
	}
	return findVolume()
}

type checks struct{ pass, fail int }

func (c *checks) ok(f string, a ...any)  { c.pass++; fmt.Printf("  ok      "+f+"\n", a...) }
func (c *checks) bad(f string, a ...any) { c.fail++; fmt.Printf("  FAIL    "+f+"\n", a...) }

func exists(p string) bool { _, err := os.Stat(p); return err == nil }
func isDir(p string) bool  { st, err := os.Stat(p); return err == nil && st.IsDir() }

func fileContains(p, needle string) bool {
	b, err := os.ReadFile(p)
	return err == nil && strings.Contains(string(b), needle)
}

// cmdVerify answers "did the jailbreak take?" without needing a shell on the
// device.
//
// This exists because the documented check is misleading: `;log` only works
// AFTER the hotfix, so before it a fall-through to library search is EXPECTED
// and proves nothing. Issues #77, #71, #66 and #59 are all people reading a
// working install as a failure.
func cmdVerify(args []string) error {
	vol, err := volumeArg(args)
	if err != nil {
		return err
	}
	fmt.Printf("volume: %s\n\n", vol)
	c := &checks{}

	fmt.Println("exploit ran")
	if exists(filepath.Join(vol, "LanguageBreakRan")) {
		c.ok("LanguageBreakRan present")
	} else {
		c.bad("LanguageBreakRan MISSING — the exploit never ran")
	}
	log := filepath.Join(vol, "languagebreak_log")
	if fileContains(log, "Finished installing jailbreak") {
		c.ok("log reports the install finished")
	} else {
		c.bad("languagebreak_log missing or incomplete")
	}
	if fileContains(log, "I am root") {
		c.ok("log confirms uid=0")
	} else {
		c.bad("log does not confirm root")
	}

	fmt.Println("\nhotfix applied")
	// mkk is the documented success marker. On PW3 the hotfix commonly needs
	// applying TWICE (issue #48): the first pass leaves no mkk, rootfs
	// permission errors, and KUAL apps that launch then silently die.
	if isDir(filepath.Join(vol, "mkk")) {
		c.ok("mkk present")
	} else {
		c.bad("mkk MISSING — apply the hotfix again; twice is normal on PW3, see issue #48")
	}
	if isDir(filepath.Join(vol, "libkh")) {
		c.ok("libkh present")
	} else {
		c.bad("libkh missing")
	}
	if exists(filepath.Join(vol, "libkh", "bin", "fbink")) {
		c.ok("fbink usable")
	} else {
		c.bad("fbink missing")
	}

	fmt.Println("\nroot exec channel")
	for _, f := range []string{"bridge.sh", "gandalf"} {
		if exists(filepath.Join(vol, "mkk", f)) {
			c.ok("mkk/%s present", f)
		} else {
			c.bad("mkk/%s missing", f)
		}
	}
	if fileContains(filepath.Join(vol, "mkk", "bridge.conf"), "BRIDGE_EMERGENCY") {
		c.ok("emergency hook available — a script at /mnt/us/emergency.sh runs as root at boot")
	} else {
		c.bad("no BRIDGE_EMERGENCY in bridge.conf — the recovery path is unavailable")
	}

	fmt.Printf("\n%d ok, %d failed\n", c.pass, c.fail)
	if c.fail > 0 {
		return fmt.Errorf("verification failed")
	}
	fmt.Println("jailbreak is installed and the root exec channel is live.")
	return nil
}

// cmdPayload copies the payload, then proves the copy is usable.
//
// Copying with a GUI file manager fails SILENTLY in two ways: .demo/ is hidden
// so drag-copy leaves it behind, and macOS writes ._ AppleDouble sidecars that
// FAT32 cannot hold. Both produce a payload that looks complete and does
// nothing.
func cmdPayload(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: lbtool payload <LanguageBreak-dir> [volume]")
	}
	src := args[0]
	if !isDir(src) {
		return fmt.Errorf("no such directory: %s", src)
	}
	vol, err := volumeArg(args[1:])
	if err != nil {
		return err
	}
	fmt.Printf("copying %s -> %s\n", src, vol)

	n, err := copyTree(src, vol)
	if err != nil {
		return err
	}
	fmt.Printf("%d files written and flushed\n\n", n)

	c := &checks{}
	fmt.Println("all five root entries must be present")
	for _, f := range payloadRoot {
		if exists(filepath.Join(vol, f)) {
			c.ok("%s", f)
		} else {
			c.bad("%s MISSING", f)
		}
	}

	fmt.Println("\nthe exploit filename must have survived verbatim")
	if name, found := findExploit(filepath.Join(vol, "documents", "dictionaries")); found {
		c.ok("shell-injection dictionary file present")
		fmt.Printf("          %q\n", name)
	} else {
		c.bad("exploit filename did not survive the copy")
	}

	fmt.Println("\nstray AppleDouble sidecars at the volume root")
	if stray := strayDouble(vol); len(stray) == 0 {
		c.ok("none")
	} else {
		c.bad("%d found: %s", len(stray), strings.Join(stray, " "))
	}

	fmt.Printf("\n%d ok, %d failed\n", c.pass, c.fail)
	if c.fail > 0 {
		return fmt.Errorf("copy incomplete — do NOT continue, the jailbreak will fail silently")
	}
	fmt.Println("payload copied correctly. Eject before continuing.")
	return nil
}

// copyTree copies src into dst, skipping AppleDouble sidecars and fsyncing each
// file. Filenames are passed through as bytes — no shell, no globbing, no
// interpolation, which is the entire point.
func copyTree(src, dst string) (int, error) {
	n := 0
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil || rel == "." {
			return err
		}
		if strings.HasPrefix(filepath.Base(rel), "._") {
			return nil // FAT32 cannot hold these and they confuse the payload
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			_ = os.Remove(target)
			return os.Symlink(link, target)
		}
		if err := copyFileSync(path, target); err != nil {
			return fmt.Errorf("%s: %w", rel, err)
		}
		n++
		return nil
	})
	return n, err
}

// copyFileSync copies and fsyncs. The shell version relies on a trailing
// sync(1); per-file fsync is stronger, and on a device people unplug the moment
// the copy "looks done" that difference is the whole point.
func copyFileSync(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	if err := out.Sync(); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// exploitFilename is the payload file whose NAME is the shell injection.
// Compared in full below, not by substring: a name that lost bytes, gained a
// space, or had a metacharacter mangled can still contain "export SLASH"
// while being useless as an exploit.
const exploitFilename = `a; export SLASH=$(awk 'BEGIN {print substr(ARGV[1], 0, 1)}' ${PWD}); sh ${SLASH}mnt${SLASH}us${SLASH}jb`

// findExploit looks for the payload file whose NAME carries the injection.
func findExploit(dir string) (string, bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", false
	}
	for _, e := range entries {
		if e.Name() == exploitFilename {
			return e.Name(), true
		}
	}
	return "", false
}

func strayDouble(vol string) []string {
	var out []string
	entries, err := os.ReadDir(vol)
	if err != nil {
		return out
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "._") {
			out = append(out, e.Name())
		}
	}
	return out
}
