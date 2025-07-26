#!/bin/bash

# Example script for using Sub tool
# Created by SayerLinux (SaudiSayer@gmail.com)

# Basic subdomain discovery
echo "[+] Basic subdomain discovery for example.com"
echo "./sub -t example.com -o discovered_subdomains.txt"
echo ""

# Advanced subdomain discovery with custom wordlist and increased threads
echo "[+] Advanced subdomain discovery with custom wordlist"
echo "./sub -t example.com -w wordlists/large_wordlist.txt -c 200 -o discovered_subdomains.txt -v"
echo ""

# Scanning discovered subdomains for services and extracting files
echo "[+] Scanning discovered subdomains for services and extracting files"
echo "./sub scan -t discovered_subdomains.txt -o output_directory"
echo ""

# Targeted scan of a specific subdomain
echo "[+] Targeted scan of a specific subdomain"
echo "./sub scan -t admin.example.com -o admin_output"
echo ""

# Full workflow example
echo "[+] Full workflow example"
echo "# Step 1: Discover subdomains"
echo "./sub -t example.com -w wordlists/default.txt -c 100 -o discovered.txt -v"
echo ""
echo "# Step 2: Filter live subdomains"
echo "cat discovered.txt | grep -v 'Not Found' > live_subdomains.txt"
echo ""
echo "# Step 3: Scan live subdomains for services"
echo "./sub scan -t live_subdomains.txt -p -e=false -o services_output"
echo ""
echo "# Step 4: Extract files from interesting subdomains"
echo "grep -i 'admin\|portal\|api' live_subdomains.txt > interesting_subdomains.txt"
echo "./sub scan -t interesting_subdomains.txt -p=false -e -o files_output"
echo ""

# Example of using with other tools
echo "[+] Integration with other tools"
echo "# Combining with other subdomain discovery tools"
echo "./sub -t example.com -o sub_results.txt"
echo "subfinder -d example.com -o subfinder_results.txt"
echo "amass enum -d example.com -o amass_results.txt"
echo "cat sub_results.txt subfinder_results.txt amass_results.txt | sort -u > all_subdomains.txt"
echo "./sub scan -t all_subdomains.txt -o combined_output"
echo ""

echo "[+] Note: These are example commands. Replace 'example.com' with your target domain."
echo "[+] For more information, run './sub --help' or './sub scan --help'"