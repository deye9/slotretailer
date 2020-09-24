@echo off
setlocal enabledelayedexpansion

(
    start cmd /C "cd C:\Users\gbenga\Desktop\Retailer\retailmanager\frontend & npm run serve"
    start cmd /C "cd C:\Users\gbenga\Desktop\Retailer\retailmanager && wails serve"
) | pause

echo Done