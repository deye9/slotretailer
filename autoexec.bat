@echo off
setlocal enabledelayedexpansion
Echo Launch dir: "%~dp0"
Echo Current dir: "%CD%"

(
    start cmd /C "wails serve"
    start cmd /C "cd frontend & npm run serve"
) | pause

echo Done