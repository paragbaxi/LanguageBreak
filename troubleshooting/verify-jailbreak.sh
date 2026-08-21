#!/bin/sh
# verify-jailbreak.sh — did the jailbreak actually take?
#
# ADDRESSES: LanguageBreak issues #77 ("Log not showing after step 18"), #71,
# #66, #59 — all variants of "I followed the steps, did it work?". The current
# answer is "type ;log and see if text appears", which fails for a reason the
# FAQ does not mention: ;log ONLY works AFTER the hotfix. Before it, ;log
# falling through to a library search is EXPECTED and tells you nothing.
#
# This checks the artefacts the exploit itself writes, which are ground truth
# and readable over USB with no shell.
#
# RUN: copy to /mnt/us/emergency.sh and reboot, or run over ssh.
# Writes /mnt/us/verify-jailbreak.txt

OUT=/mnt/us/verify-jailbreak.txt
pass=0; fail=0
ok()   { echo "  PASS  $1" >> "$OUT"; pass=$((pass+1)); }
bad()  { echo "  FAIL  $1" >> "$OUT"; fail=$((fail+1)); }

echo "===== verify-jailbreak $(date) =====" > "$OUT"
echo >> "$OUT"

echo "-- exploit ran --" >> "$OUT"
[ -f /mnt/us/LanguageBreakRan ] && ok "LanguageBreakRan present" || bad "LanguageBreakRan MISSING - exploit never ran"
if grep -q "Finished installing jailbreak" /mnt/us/languagebreak_log 2>/dev/null; then
  ok "languagebreak_log says 'Finished installing jailbreak!'"
else
  bad "languagebreak_log missing or incomplete"
fi
grep -q "I am root" /mnt/us/languagebreak_log 2>/dev/null && ok "log confirms uid=0" || bad "log does not confirm root"

echo >> "$OUT"
echo "-- hotfix applied --" >> "$OUT"
# mkk is the documented success marker. On PW3 the hotfix often needs applying
# TWICE (issue #48): the first pass leaves no mkk, rootfs permission errors,
# and KUAL apps that launch and silently die.
[ -d /mnt/us/mkk ] && ok "mkk present (hotfix applied)" || bad "mkk MISSING - apply the hotfix again; on PW3 twice is normal, see issue #48"
[ -d /mnt/us/libkh ] && ok "libkh present" || bad "libkh missing"
[ -x /mnt/us/libkh/bin/fbink ] && ok "fbink usable" || bad "fbink missing"

echo >> "$OUT"
echo "-- root exec channel --" >> "$OUT"
[ -f /mnt/us/mkk/bridge.sh ]  && ok "mkk/bridge.sh present" || bad "mkk/bridge.sh missing"
[ -f /mnt/us/mkk/gandalf ]    && ok "gandalf (setuid root helper) present" || bad "gandalf missing"
grep -q "BRIDGE_EMERGENCY" /mnt/us/mkk/bridge.conf 2>/dev/null \
  && ok "emergency hook available: a script at /mnt/us/emergency.sh runs as root at boot" \
  || bad "no BRIDGE_EMERGENCY in bridge.conf"

echo >> "$OUT"
echo "-- leftovers that should be gone --" >> "$OUT"
[ -e /mnt/us/jb ] && bad "jb still present - staging file not consumed" || ok "staging files consumed"
[ -e /mnt/us/.demo ] && bad ".demo still present" || ok ".demo consumed"

echo >> "$OUT"
echo "-- state that commonly strands people --" >> "$OUT"
[ -e /var/local/system/DEMO_MODE ] && bad "DEMO_MODE still set - still in demo mode" || ok "DEMO_MODE clear"
if [ -e /var/local/system/no_transitions ]; then
  bad "no_transitions present - USB mass storage is DISABLED (see fix-usb-disabled.sh)"
else
  ok "USB not blocked"
fi

echo >> "$OUT"
echo "===== $pass passed, $fail failed =====" >> "$OUT"
[ "$fail" = 0 ] && echo "Jailbreak looks healthy." >> "$OUT" \
                || echo "See the FAIL lines above." >> "$OUT"
chmod 666 "$OUT" 2>/dev/null
sync
[ -x /mnt/us/libkh/bin/fbink ] && /mnt/us/libkh/bin/fbink -y 2 "verify: $pass pass / $fail fail" 2>/dev/null
exit 0
