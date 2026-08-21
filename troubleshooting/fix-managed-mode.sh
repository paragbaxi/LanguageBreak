#!/bin/sh
# fix-managed-mode.sh — escape "managed mode" after a LanguageBreak install.
#
# ADDRESSES: LanguageBreak issue #67 ("stuck in managed mode"), and the same
# symptom reported in #58 and #52. The only answer in #67 today is a bare
# "solved with a factory reset", which is unhelpful when Settings is greyed out
# and the factory-reset option is exactly what you cannot reach.
#
# SYMPTOM
#   * Settings greyed out, "contact your administrator"
#   * USB mass storage unavailable
#   * ;log, ;uzb, ;enter_demo all fall through to a library search
#   * so the documented exit (;demo -> Resell Device) is unreachable, because
#     the ; command channel is itself what is broken
#
# HOW IT WORKS
#   /opt/var/local is a READ-ONLY SQUASHFS holding the pristine factory copy of
#   /var/local (visible as a /dev/loop mount in `mount`). The system repopulates
#   /var/local from it. Wiping /var/local therefore performs the same reset the
#   Settings menu would — and root can do it with Settings locked.
#
# THE JAILBREAK SURVIVES, two independent ways:
#   1. this script skips mkk, rp and linkfonts, and
#   2. /etc/upstart/bridge.conf restores /var/local/mkk from /mnt/us/mkk when it
#      is missing — the jailbreak authors built that self-heal in already.
#
# HOW TO RUN IT WITHOUT A SHELL
#   You almost certainly cannot ssh in at this point. Use the bridge's root
#   hook: copy this file to /mnt/us/emergency.sh and reboot. It runs as root at
#   boot. It renames itself when done, so it runs exactly once.
#
#   *** WARNING: this wipes device state — registration, settings, caches.
#   *** It does NOT touch /mnt/us, so your documents and the jailbreak payload
#   *** are untouched. Verified on PW3 / 5.16.2.1.1.

LOG=/mnt/us/fix-managed-mode.log
KEEP="mkk rp linkfonts"

# ---- 1. Disarm FIRST, while /mnt/us is definitely writable. ----
# If this script stops the framework or dies partway, /mnt/us (a fuse.fsp mount)
# can disappear, leaving the file in place to re-run on every boot. That failure
# mode produces a permanent boot loop and is much worse than the bug.
mv /mnt/us/emergency.sh /mnt/us/emergency.sh.done 2>/dev/null
sync

{
  echo "===== fix-managed-mode $(date) ====="
  id
  echo
  echo "-- factory template present? (must be a read-only squashfs) --"
  mount | grep -i "opt/var/local" || echo "  WARNING: /opt/var/local not mounted; reset may not repopulate"
  echo
  echo "-- /var/local before --"
  ls -la /var/local 2>/dev/null
  echo
  echo "-- jailbreak backup on /mnt/us (used by bridge.conf to self-heal) --"
  ls -la /mnt/us/mkk/ 2>/dev/null | head
  echo
  echo "-- wiping /var/local, keeping: $KEEP --"
} > "$LOG" 2>&1

for entry in /var/local/* /var/local/.[!.]* ; do
  [ -e "$entry" ] || continue
  base=$(basename "$entry")
  skip=0
  for k in $KEEP; do [ "$base" = "$k" ] && skip=1; done
  [ "$skip" = 1 ] && { echo "  keep   $base" >> "$LOG"; continue; }
  chattr -i "$entry" 2>/dev/null
  rm -rf "$entry" 2>/dev/null
  echo "  remove $base" >> "$LOG"
done
sync

echo "" >> "$LOG"
echo "-- rebooting; /var/local repopulates from /opt/var/local --" >> "$LOG"
sync
sleep 2

# On-screen confirmation if fbink is available (LanguageBreak ships it).
[ -x /mnt/us/libkh/bin/fbink ] && /mnt/us/libkh/bin/fbink -y 2 "managed-mode fix applied, rebooting" 2>/dev/null

reboot
