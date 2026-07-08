@echo off
:: Enable strict error checking (closest equivalent to set -e)
setlocal enabledelayedexpansion

:: Create the release directory if it doesn't exist
if not exist windows-release (
    mkdir windows-release
)

echo Building Windows Release...

:: Set CGO to 1 to ensure zero-allocation mesh updates
set CGO_ENABLED=1

:: -s -w strips debug info (smaller file)
:: -H=windowsgui hides the ugly black terminal box when launching the game
go build -ldflags="-s -w -H=windowsgui" -o windows-release/game.exe ..\

:: Check if the build succeeded (exit code 0)
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo [ERROR] Build failed! Check your compiler/PATH settings.
    goto end
)

echo.
echo Build complete: .\windows-release\game.exe

:end
echo Press Enter to exit...
pause > nul