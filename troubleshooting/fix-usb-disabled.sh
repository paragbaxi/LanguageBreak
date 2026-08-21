#!/bin/sh
# fix-usb-disabled.sh — restore USB mass storage after demo mode.
#
# ADDRESSES: an undocumented dead end. The host enumerates "Amazon Kindle" on
# the USB bus but NO VOLUME EVER MOUNTS, which looks exactly like a bad cable
# and sends people hunting for hardware faults.
#
# CAUSE: /usr/bin/disableUSBInDemo.sh creates
# /var/local/system/no_transitions, which stops volumd mounting userstore.
# Amazon's own /usr/bin/enableUSBInDemo.sh only removes it IF DEMO_MODE still
# exists:
#
#     DEMO_MODE_FILE=/var/local/system/DEMO_MODE
#     NO_TRANSITIONS=/var/local/system/no_transitions
#     [ -e "$DEMO_MODE_FILE" ] && rm -f "$NO_TRANSITIONS"
#
# So once the demo flag is cleared, the blocker is STRANDED with nothing left
# that will ever remove it. This removes it directly.
#
# DIAGNOSIS on the host: `ioreg -p IOUSB` (macOS) or `lsusb` shows the device
# while `diskutil list external` / `lsblk` shows no volume. That is this bug,
# not a cable.
#
# HOW TO RUN IT — read this first, because it cannot bootstrap itself over
# USB: USB not mounting is the symptom. This is a ROOT-SHELL convenience, not
# a rescue tool. If your `;` commands still work you do not need it at all:
#
#   `;enter_demo` puts the device back in demo mode, which means DEMO_MODE
#   exists again, which re-arms the device's own enableUSBInDemo.sh; `;uzb` —
#   the same command the README uses at the hotfix step — then gets you a
#   mount. Leaving demo mode afterwards the sanctioned way (`;demo` -> Resell
#   Device) runs deleteDemoModeFlagFile.sh, which removes DEMO_MODE and
#   no_transitions together, so nothing is stranded a second time.
#
# Use this script when you have root and would rather not do that round trip
# — over ssh if dropbear is installed, from KUAL, or as /mnt/us/emergency.sh
# if you can still reach /mnt/us. It is one `rm` instead of two reboots.
#
# If the `;` channel is dead too (the managed-mode lockout), neither route is
# open. A ~40 second power hold with the cable UNPLUGGED restored demo mode
# and the `;` channel on my PW3; held while plugged in it only powers the
# device off.

LOG=/mnt/us/fix-usb.log
{
  echo "===== fix-usb-disabled $(date) ====="
  if [ -e /var/local/system/no_transitions ]; then
    rm -f /var/local/system/no_transitions
    if [ -e /var/local/system/no_transitions ]; then
      echo "  FAILED to remove no_transitions"
    else
      echo "  removed no_transitions - USB mass storage restored"
    fi
  else
    echo "  no_transitions not present; USB is not blocked by this"
  fi
  echo "  DEMO_MODE: $([ -e /var/local/system/DEMO_MODE ] && echo set || echo clear)"
} >> "$LOG" 2>&1
sync
[ -x /mnt/us/libkh/bin/fbink ] && /mnt/us/libkh/bin/fbink -y 2 "usb fix applied" 2>/dev/null
exit 0
