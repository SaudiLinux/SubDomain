@echo off
REM Example script for using Sub tool on Windows
REM Created by SayerLinux (SaudiSayer@gmail.com)

echo [+] Basic subdomain discovery for example.com
echo sub.exe -t example.com -o discovered_subdomains.txt
echo.

echo [+] Advanced subdomain discovery with custom wordlist and increased threads
echo sub.exe -t example.com -w wordlists\large_wordlist.txt -c 200 -o discovered_subdomains.txt -v
echo.

echo [+] Scanning discovered subdomains for services and extracting files
echo sub.exe scan -t discovered_subdomains.txt -o output_directory
echo.

echo [+] Targeted scan of a specific subdomain
echo sub.exe scan -t admin.example.com -o admin_output
echo.

echo [+] Full workflow example
echo # Step 1: Discover subdomains
echo sub.exe -t example.com -w wordlists\default.txt -c 100 -o discovered.txt -v
echo.
echo # Step 2: Filter live subdomains (using PowerShell)
echo powershell -Command "Get-Content discovered.txt | Where-Object {$_ -notmatch 'Not Found'} | Set-Content live_subdomains.txt"
echo.
echo # Step 3: Scan live subdomains for services
echo sub.exe scan -t live_subdomains.txt -p -e=false -o services_output
echo.
echo # Step 4: Extract files from interesting subdomains (using PowerShell)
echo powershell -Command "Get-Content live_subdomains.txt | Where-Object {$_ -match 'admin|portal|api'} | Set-Content interesting_subdomains.txt"
echo sub.exe scan -t interesting_subdomains.txt -p=false -e -o files_output
echo.

echo [+] Note: These are example commands. Replace 'example.com' with your target domain.
echo [+] For more information, run 'sub.exe --help' or 'sub.exe scan --help'

pause