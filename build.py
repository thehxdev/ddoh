#!/usr/bin/env python3

import os
import sys
import logging as log
import subprocess

log.basicConfig(format="[%(levelname)s] %(asctime)s %(name)s: %(message)s",
                    datefmt="%Y-%m-%d %H:%M:%S",
                    level=log.INFO)

if not sys.version_info >= (3, 8):
    log.error("this script requires python 3.8 or above")
    sys.exit(1)

BIN  = "ddoh"
GO   = "go"

ARCH = [ "amd64", "arm64" ]
OS   = [ "linux", "android", "windows", "darwin" ]

STDOUT = sys.stdout
STDERR = sys.stderr


def system(cmd: str):
    return subprocess.run(cmd, shell=True, stdout=STDOUT, stderr=STDERR)


def build(arch: str, osname: str) -> None:
    outName = f"{BIN}-{osname}-{arch}"
    envVars = f"GOOS={osname} GOARCH={arch}"
    log.info(f"building {outName}")

    if osname == "windows":
        envVars += " CGO_ENABLED=1"
        outName += ".exe"

    system(f"{envVars} {GO} build -ldflags='-s -buildid=' -o 'build/{outName}' .")


def main() -> None:
    os.makedirs("build")
    for osName in OS:
        for arch in ARCH:
            if osName == "android" and arch != "arm64":
                continue
            if osName == "darwin" and (arch == "386" or arch == "arm"):
                continue
            build(arch, osName)
    log.info(f"Done!")


if __name__ == "__main__":
    main()
