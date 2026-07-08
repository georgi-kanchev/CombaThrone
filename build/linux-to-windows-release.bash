#!/bin/bash
set -e

# --- CONFIGURATION ---
GOOS=windows
GOARCH=amd64
CC=x86_64-w64-mingw32-gcc
CXX=x86_64-w64-mingw32-g++

# Target the parent directory
SCRIPT_DIR=$(pwd)
GO_PROGRAM_DIR=$(realpath ..) 
RELEASE_DIR="$SCRIPT_DIR/windows-release"
OUTPUT_EXE="$RELEASE_DIR/game.exe"

mkdir -p "$RELEASE_DIR"

# --- CREATE TEMPORARY WORKDIR ---
WORKDIR=$(mktemp -d)
echo "Using temporary workdir: $WORKDIR"

# --- CLONE AND BUILD RAYLIB ---
echo "Cloning Raylib..."
git clone --depth 1 https://github.com/raysan5/raylib.git "$WORKDIR/raylib"

echo "Building Raylib static library for Windows..."
mkdir -p "$WORKDIR/raylib/build-windows"
cd "$WORKDIR/raylib/build-windows"
cmake -G "Unix Makefiles" \
      -DCMAKE_SYSTEM_NAME=Windows \
      -DCMAKE_C_COMPILER=$CC \
      -DCMAKE_CXX_COMPILER=$CXX \
      -DBUILD_SHARED_LIBS=OFF \
      ..
make

# Note: check if path is libraylib.a or raylib/libraylib.a based on cmake version
RAYLIB_INCLUDE="$WORKDIR/raylib/src"
RAYLIB_LIB=$(find "$WORKDIR/raylib/build-windows" -name "libraylib.a" | head -n 1)

# --- TEMPORARILY INJECT CGO FLAGS ---
# We create the temp file in the parent dir so 'go build' sees it as part of the package
cd "$GO_PROGRAM_DIR"
TEMP_GO_FILE="./_tmp_build_inject.go"

echo "package main" > "$TEMP_GO_FILE"
echo "// #cgo CFLAGS: -I$RAYLIB_INCLUDE" >> "$TEMP_GO_FILE"
echo "// #cgo LDFLAGS: $RAYLIB_LIB -lopengl32 -lgdi32 -lwinmm -lshell32" >> "$TEMP_GO_FILE"
echo "import \"C\"" >> "$TEMP_GO_FILE"

# --- BUILD GO PROGRAM ---
echo "Building Windows Release..."
export CGO_ENABLED=1
export GOOS=$GOOS
export GOARCH=$GOARCH
export CC=$CC
export CXX=$CXX

# Instead of 'cat'ing everything, we just build the directory. 
# The injected CGO flags in the temp file will apply to the whole build.
go build -ldflags="-s -w -H=windowsgui" -o "$OUTPUT_EXE" .

# --- CLEAN UP ---
rm "$TEMP_GO_FILE"
rm -rf "$WORKDIR"

echo "-----------------------------------------------"
echo "Build complete: $OUTPUT_EXE"
echo "-----------------------------------------------"
