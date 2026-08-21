# **LanguageBreak** 

Jailbreak for any kindle running FW 5.16.2.1.1 or **LOWER**

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/E1E1QLG4D)

**The exploit works best around version 5.16.2, so if you are on lower firmware you should consider updating**

Big thanks to Bluebotlabs for all the help along the way and GeorgeYellow and bulltricks for bringing the vulnerability to light

The most up to date release can be found [here](https://github.com/notmarek/LanguageBreak/releases/latest).

# Installation

## Warnings
- ⚠️ This method **will** __DELETE__ all content on your device. Backup important files.
- ⚠️ This method works up to firmware version `5.16.2.1.1`. It is strongly suggested you upgrade to this firmware before starting.
- ⚠️ Do not update your firmware beyond `5.16.2.1.1` after jailbreaking with LanguageBreak, as future versions of the firmware patch it.
- ⚠️ Your device must have no password lock. Remove it before starting.

## Before you start, make sure:
- Airplane Mode is **ON**
- The device (esp. the root folder) should have no `.bin` or `update.bin.tmp.partial` files before you start. This may be a pending OTA update. Delete them.

## Jailbreaking your Kindle

1. Type `;enter_demo` into the Kindle search bar, press enter, then reboot the device.
2. Once the device boots, dismiss the wifi selection dialog, type whatever you'd like into the text fields, then continue.
3. Select `Skip`, then `Standard`, then `Done`
4. The device will take a few minutes to go into demo mode. When it's done, use [this gesture](https://www.youtube.com/watch?v=JzuIGbGPpig) to access the main screen.
5. Type `;demo` into the Kindle search bar, then press enter to access the Demo Mode Configuration screen.
6. Select the `Sideload Content` option.
7. Connect your Kindle to a PC, then copy the contents of the `LanguageBreak` folder from the [latest release's *.tar.gz file](https://github.com/notmarek/LanguageBreak/releases) to the Kindle's root directory (if prompted, overrwrite the existing files).
    - ⚠️ **Do not copy with a graphical file manager.** `.demo` is a *hidden* folder that Finder/Explorer will silently leave behind, and one payload file's *filename* is the exploit itself (it contains `;`, `$`, `(`, `)`, quotes and spaces) and must land byte-identical. Both failures are silent — the jailbreak simply does nothing. Use `rsync`/`cp` from a terminal, or [`troubleshooting/copy-payload.sh`](troubleshooting/copy-payload.sh), which copies and then verifies.
8. Eject, then *unplug your Kindle*, then return to the Demo Configuration Screen (using the method in step 5, if required).
9. Select `Resell Device`, then confirm.
10. As soon as the "Press the Power Button" screen appears, plug your Kindle back into your computer. (This is time sensitive, be quick.)
11. Copy the contents of the `LanguageBreak` folder to the Kindle's root directory, again. (if prompted, overrwrite the existing files).
12. After all files have been written, eject your Kindle, then press and hold the power button until the unit reboots.
13. At this point, a language selection screen appears. Select `简体中文` (Chinese), which should appear above an entry called `Pseudot` and below Japanese and click the Chinese text button that appears in the middle of the screen.
14. Your Kindle should reboot, and some log messages should appear in the top right-hand corner.
    
## Applying the Hotfix

1. After the device has rebooted, type `;uzb` into the Kindle search bar to enable USB access within demo mode, then press enter.
2. Connect the device to a PC and copy whichever `Update_hotfix_languagebreak-{languge/locale}.bin` file from the release matches your language to the Kindle's root directory.
3. Eject your Kindle, then `;dsts` into the Kindle search bar to access the settings page. Locate the `Update your Kindle` option and press it, then confirm.
 - Using your phone to translate the menus is helpful if you don't happen to be able to read Chinese.
 - If you see an infoscreen that states "Update Error" or similar then try to apply the hotfix again (last three steps)
 - **On PW3 the hotfix commonly has to be applied twice.** The first pass can leave the device looking fine (`;log` works, KUAL opens) while actually being broken: no `mkk` folder, permission errors on the rootfs, and KUAL apps that launch and silently die. If `mkk` is missing, apply the hotfix again — you do **not** need to re-enter demo mode, just copy the `.bin` and update. See [#48](https://github.com/notmarek/LanguageBreak/issues/48).

This will reboot the device out of Demo mode. Your device may go into Managed mode after completing these steps. Managed devices have some settings greyed out, and ask the user to contact their system administrator. See below for steps on how to restore functionality (and the correct language) to your device.
    
## Restoring the Correct Language and Exiting Managed Mode

> **If the `;` commands below do nothing**, the command channel itself is broken
> and neither procedure can work — `;demo`, `;uzb` and `;enter_demo` all fall
> through to an ordinary library search. This is the state reported in
> [#67](https://github.com/notmarek/LanguageBreak/issues/67),
> [#58](https://github.com/notmarek/LanguageBreak/issues/58) and
> [#52](https://github.com/notmarek/LanguageBreak/issues/52), usually alongside
> greyed-out Settings and USB that never mounts.
>
> Two things that recover it:
>
> - **A 40-second power-hold with the cable UNPLUGGED** often restores demo mode
>   and the `;` channel. Holding it while connected only powers the device off.
> - **[`troubleshooting/fix-managed-mode.sh`](troubleshooting/fix-managed-mode.sh)**
>   resets `/var/local` from root, which is what a factory reset does internally
>   — `/opt/var/local` is a read-only squashfs holding the pristine copy. This
>   works with Settings locked. **The jailbreak survives**, because
>   `bridge.conf` restores `/var/local/mkk` from `/mnt/us/mkk`.
    
### If your kindle **is not** registered to an Amazon account:
1. Type `;demo` into the Kindle search bar.
2. You will get a prompt with two buttons. Press the right-most button.
3. The device will reboot. If all is well, your Kindle should have a folder named `mkk` in the root directory.

### If your kindle **is** registered to an Amazon account:
1. Enter `;enter_demo` into the Kindle search bar, then reboot your device.
2. The device will be back in full "Demo Mode". Use [the same gesture](https://www.youtube.com/watch?v=JzuIGbGPpig) to access the main screen.
3. Enter `;demo` into the Kindle search bar.
4. Select `Resell device`, then confirm.
5. The device will reboot. If all is well, your Kindle should have a folder named `mkk` in the root directory.


# FAQ

Q: It didn't work!!!!!! 

A: Not a question.. but checkout issue [#48](https://github.com/notmarek/LanguageBreak/issues/48)

Q: How do I verify my installation?

A: (before applying hotfix): Install hotfix, if you can do that then it worked.
**Note that `;log` does NOT work yet at this stage** — before the hotfix it falls through to an ordinary library search, and that is *expected*. It is not a sign of failure, and several issues in this tracker are people concluding otherwise.

A: (after applying hotfix): Type `;log` into the Kindle search bar, some text should appear at the top right side of the screen.

A: (either stage, no shell needed): run [`troubleshooting/verify-jailbreak.sh`](troubleshooting/verify-jailbreak.sh). It checks the artefacts the exploit itself writes — `LanguageBreakRan`, `languagebreak_log` ending in `Finished installing jailbreak!`, and whether `mkk` exists — and writes the result to `/mnt/us/verify-jailbreak.txt`, readable over USB.

Q: Where are the hotfix files?

A: The structure of the tarball is as follows
```
LanguageBreak.tar.gz
|-- LanguageBreak
|	|-- documents
|	|	|-- dictionaries
|	|	|	|-- a; export SLASH=$(awk 'BEGIN {print substr(ARGV[1], 0, 1)}' ${PWD}); sh ${SLASH}mnt${SLASH}us${SLASH}jb
|	|	|	|-- amisane
|	|-- DONT_CHECK_BATTERY
|	|-- jb
|	|-- patchedUks
|	|-- .demo
|	|	|-- boot.flag
|-- Update_hotfix_languagebreak-*.bin
```

# Recovering a device you cannot reach

If the `;` commands are dead, Settings is locked and USB will not mount, it can
look like there is no way back in. There is, and it is built into the jailbreak:

```sh
BRIDGE_EMERGENCY="/mnt/us/emergency.sh"
if [ -f "${BRIDGE_EMERGENCY}" ] ; then
    [ -x "${BRIDGE_EMERGENCY}" ] || chmod +x "${BRIDGE_EMERGENCY}"
    /bin/sh "${BRIDGE_EMERGENCY}"
fi
```

**Any script placed at `/mnt/us/emergency.sh` runs as root at boot** (see
`mkk/bridge.conf`). No KUAL, no MRPI, no USBNetwork. Copy a script there over
USB, reboot, and it runs.

That is the intended route for the scripts here that need a boot-time root
shell — but it is not a universal one, and in particular it cannot work for
`fix-usb-disabled.sh`, whose symptom is that USB never mounts. See the
`HOW TO RUN IT` header in that script for the routes that do apply to it.

Two rules if you write your own:

1. **Rename the script before doing anything destructive.** If it dies partway,
   `/mnt/us` (a `fuse.fsp` mount) can go with it, leaving the file in place to
   re-run on every boot — a permanent boot loop that is worse than the original
   problem.
2. **Call `/var/local/mkk/bridge.sh` at the end** if you are not rebooting, as
   the emergency path makes the bridge skip its own setup.

### USB never mounts after demo mode

The host sees the device on the USB bus but no volume ever appears, which looks
exactly like a bad cable. `disableUSBInDemo.sh` creates
`/var/local/system/no_transitions`, and Amazon's own `enableUSBInDemo.sh` only
removes it *if `DEMO_MODE` still exists* — so once the demo flag is cleared the
blocker is stranded permanently.
[`troubleshooting/fix-usb-disabled.sh`](troubleshooting/fix-usb-disabled.sh)
removes it.

# Troubleshooting
Having general issues?

This method works best around firmware version `5.16.2`. Consider updating to this version to avoid compatability issues.

Download the standard firmware for your device using the URL templates below, replacing the `X`s with the version numbers as appropriate. 

```
PW5: https://s3.amazonaws.com/firmwaredownloads/update_kindle_all_new_paperwhite_11th_5.XX.X.bin
PW4: https://s3.amazonaws.com/firmwaredownloads/update_kindle_all_new_paperwhite_v2_5.XX.X.bin
PW3: https://s3.amazonaws.com/firmwaredownloads/update_kindle_all_new_paperwhite_5.XX.X.bin
Kindle 11th Gen: https://s3.amazonaws.com/firmwaredownloads/update_kindle_11th_5.XX.X.bin
Kindle 10th Gen: https://s3.amazonaws.com/firmwaredownloads/update_kindle_10th_5.XX.X.bin
Kindle 8th Gen: https://s3.amazonaws.com/firmwaredownloads/update_kindle_8th_5.XX.X.bin
Scribe: https://s3.amazonaws.com/firmwaredownloads/update_kindle_scribe_5.XX.X.bin
Oasis 10th Gen: https://s3.amazonaws.com/firmwaredownloads/update_kindle_all_new_oasis_v2_5.XX.X.bin
Oasis 9th Gen: https://s3.amazonaws.com/firmwaredownloads/update_kindle_all_new_oasis_5.XX.X.bin
Oasis 8th Gen: https://s3.amazonaws.com/firmwaredownloads/update_kindle_oasis_5.XX.X.bin
```

So version 5.16.2.1.1 for PW4 would be [https://s3.amazonaws.com/firmwaredownloads/update_kindle_all_new_paperwhite_v2_5.16.2.1.1.bin](]https://s3.amazonaws.com/firmwaredownloads/update_kindle_all_new_paperwhite_v2_5.16.2.1.1.bin)

**To install the file, place it into the root directory of the Kindle, then select `Update your Kindle` in settings. It should also apply the update on reboot if the menu is inaccessible for some reason.**


[Consider buying me a coffee :)](https://ko-fi.com/notmarek)
