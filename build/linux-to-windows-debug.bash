#!/bin/bash
set -e

# --- CONFIGURATION ---
GOOS=windows
GOARCH=amd64
CC=x86_64-w64-mingw32-gcc
CXX=x86_64-w64-mingw32-g++
OUTPUT_DIR_NAME="windows-debug"
OUTPUT_EXE="game.exe"

# 1. Capture the absolute path of the directory containing THIS script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# 2. Define the output path relative to the script directory
FINAL_OUTPUT_DIR="$SCRIPT_DIR/$OUTPUT_DIR_NAME"

# 3. Move to the project root for building
cd "$SCRIPT_DIR/.."
GO_PROGRAM_DIR=$(pwd)

# Create the folder inside the script directory
mkdir -p "$FINAL_OUTPUT_DIR"

# --- CREATE TEMPORARY WORKDIR ---
WORKDIR=$(mktemp -d)
echo "Using temporary workdir: $WORKDIR"

# --- CLONE AND BUILD RAYLIB ---
echo "Cloning Raylib..."
git clone --depth 1 https://github.com/raysan5/raylib.git "$WORKDIR/raylib"

echo "Building Raylib static library for Windows (Debug)..."
mkdir -p "$WORKDIR/raylib/build-windows"
cd "$WORKDIR/raylib/build-windows"
cmake -G "Unix Makefiles" \
      -DCMAKE_SYSTEM_NAME=Windows \
      -DCMAKE_C_COMPILER=$CC \
      -DCMAKE_CXX_COMPILER=$CXX \
      -DBUILD_SHARED_LIBS=OFF \
      -DCMAKE_BUILD_TYPE=Debug \
      ..
make

RAYLIB_INCLUDE="$WORKDIR/raylib/src"
RAYLIB_LIB=$(find "$WORKDIR/raylib/build-windows" -name "libraylib.a" | head -n 1)

# --- INJECT CGO FLAGS ---
cd "$GO_PROGRAM_DIR"
TEMP_GO_FILE="./_debug_build_inject.go"

echo "package main" > "$TEMP_GO_FILE"
echo "// #cgo CFLAGS: -I$RAYLIB_INCLUDE" >> "$TEMP_GO_FILE"
echo "// #cgo LDFLAGS: $RAYLIB_LIB -lopengl32 -lgdi32 -lwinmm -lshell32" >> "$TEMP_GO_FILE"
echo "import \"C\"" >> "$TEMP_GO_FILE"

# --- BUILD GO PROGRAM ---
echo "Building Windows Debug..."
export CGO_ENABLED=1
export GOOS=$GOOS
export GOARCH=$GOARCH
export CC=$CC
export CXX=$CXX

# Output directly to the folder inside the script directory
go build -gcflags "all=-N -l" -o "$FINAL_OUTPUT_DIR/$OUTPUT_EXE" .

# --- CLEAN UP ---
rm "$TEMP_GO_FILE"
rm -rf "$WORKDIR"

echo "-----------------------------------------------"
echo "Build complete: $FINAL_OUTPUT_DIR/$OUTPUT_EXE"
echo "-----------------------------------------------"
