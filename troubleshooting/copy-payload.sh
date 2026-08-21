#!/bin/sh
# copy-payload.sh — copy the LanguageBreak payload to the Kindle correctly.
#
# ADDRESSES: silent jailbreak failures caused by copying with a GUI file
# manager. The README says "copy the contents of the LanguageBreak folder",
# which is easy to do wrongly in two ways that produce no error:
#
#   1. .demo/ is a HIDDEN directory. macOS Finder does not show it and a
#      drag-copy leaves it behind entirely.
#   2. One payload file's NAME is the exploit itself:
#        a; export SLASH=$(awk 'BEGIN {print substr(ARGV[1], 0, 1)}' ${PWD}); sh ${SLASH}mnt${SLASH}us${SLASH}jb
#      containing ; $ ( ) { } quotes and spaces, and it must land byte-identical
#      on FAT32.
#
# Also sweeps macOS AppleDouble "._" sidecars, which otherwise litter the
# Kindle root.
#
# USAGE:  ./copy-payload.sh /path/to/extracted/LanguageBreak /Volumes/Kindle
#
# Run this at README step 7 AND again at step 11.

set -eu
SRC="${1:?usage: copy-payload.sh <LanguageBreak dir> <kindle mount>}"
DST="${2:?usage: copy-payload.sh <LanguageBreak dir> <kindle mount>}"

[ -d "$SRC" ] || { echo "no such directory: $SRC"; exit 1; }
[ -d "$DST" ] || { echo "Kindle not mounted at: $DST"; exit 1; }

echo "copying payload -> $DST"
# Trailing slash on SRC copies the CONTENTS, not the folder. This matters.
# COPYFILE_DISABLE stops macOS writing ._ sidecars for extended attributes,
# which FAT32 cannot hold.
COPYFILE_DISABLE=1 rsync -rl "$SRC/" "$DST/"
command -v dot_clean >/dev/null 2>&1 && dot_clean -m "$DST" 2>/dev/null || true
sync

echo
echo "verifying — all five must be present:"
rc=0
for f in .demo documents DONT_CHECK_BATTERY jb patchedUks.sqsh; do
  if [ -e "$DST/$f" ]; then echo "  OK      $f"; else echo "  MISSING $f"; rc=1; fi
done

echo
echo "verifying the exploit filename survived verbatim:"
# Compared in full, not grepped for a substring: a name that lost its trailing
# characters, gained a space or had a metacharacter rewritten still contains
# "export SLASH" while being useless as an exploit.
EXPLOIT="a; export SLASH=\$(awk 'BEGIN {print substr(ARGV[1], 0, 1)}' \${PWD}); sh \${SLASH}mnt\${SLASH}us\${SLASH}jb"
if [ -e "$DST/documents/dictionaries/$EXPLOIT" ]; then
  echo "  OK      shell-injection dictionary file present, name byte-identical"
else
  echo "  MISSING exploit filename did not survive the copy"
  echo "  expected: $EXPLOIT"
  echo "  found in documents/dictionaries/:"
  ls -1 "$DST/documents/dictionaries/" 2>/dev/null | sed 's/^/    /' || echo "    (directory missing)"
  rc=1
fi

echo
echo "stray AppleDouble sidecars in root:"
ls -1a "$DST" 2>/dev/null | grep '^\._' || echo "  (none - clean)"

if [ "$rc" = 0 ]; then
  printf '\nPayload copied correctly. Eject before continuing.\n'
else
  printf '\n*** COPY INCOMPLETE - do not continue, the jailbreak will fail silently. ***\n'
fi
exit $rc
