# lbtool — host-side helpers in Go

Optional. The shell scripts in `../` are the primary path and need no toolchain.

This covers the two jobs that run on **your computer** while the Kindle is
plugged in over USB:

```
go build -o lbtool .

./lbtool verify                 # did the jailbreak actually take?
./lbtool payload ./LanguageBreak   # copy the payload, then prove the copy
```

Both take an optional trailing volume path; without one the Kindle is found by
looking for a mounted volume containing `documents/` and `system/` (the label is
user-changeable and localised, so it is not used for detection).

## Why this exists alongside the shell scripts

One file in the payload has a **name** that is the shell injection carrying the
exploit — semicolons, `$()`, `${}`, and a trailing space. Every hop through a
shell, a GUI file manager, or an archive tool is a chance to alter it, and the
failure is silent: the jailbreak simply does nothing and you are left debugging
the device.

Go moves filenames as bytes. There is no shell in the copy path. `main_test.go`
round-trips that exact name, trailing space included, so a regression fails the
build rather than a stranger's Kindle.

It also drops `._` AppleDouble sidecars (FAT32 cannot hold them) and fsyncs
every file rather than relying on a trailing `sync`, because people unplug the
moment a copy looks finished.

## What is deliberately NOT here

`fix-managed-mode.sh` and `fix-usb-disabled.sh` have no Go equivalent and should
not get one. They delete files under `/var/local` as root and reboot, so they
run **on the device** via `/mnt/us/emergency.sh`. Shell is the right tool there,
and a compiled binary would be the wrong one.
