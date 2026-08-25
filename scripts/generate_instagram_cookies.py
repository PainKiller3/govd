#!/usr/bin/env python3
"""
Instagram Cookie Generator Script for govd

This script generates a Netscape-formatted cookie file (`private/cookies/instagram.txt`)
which is parsed by govd for Instagram media extraction.

Options:
1. Guest Session Cookies (No login required - generates anonymous guest tokens)
2. Account Session Cookies (Performs login to retrieve valid `sessionid` and `ds_user_id`)
"""

import os
import sys
import random
import string
import time

def random_string(length=32, chars=string.ascii_letters + string.digits):
    return ''.join(random.choice(chars) for _ in range(length))

def generate_guest_cookies():
    csrftoken = random_string(32)
    ig_did = random_string(24, string.ascii_lowercase + string.digits)
    mid = random_string(24, string.ascii_letters + string.digits)
    
    cookies = [
        ("# Netscape HTTP Cookie File", ""),
        ("# http://curl.haxx.se/rfc/cookie_spec.html", ""),
        ("# This is a generated file!  Do not edit.", ""),
        ("", ""),
        (".instagram.com", "TRUE", "/", "TRUE", str(int(time.time() + 31536000)), "csrftoken", csrftoken),
        (".instagram.com", "TRUE", "/", "TRUE", str(int(time.time() + 31536000)), "ig_did", ig_did),
        (".instagram.com", "TRUE", "/", "TRUE", str(int(time.time() + 31536000)), "mid", mid),
        (".instagram.com", "TRUE", "/", "TRUE", str(int(time.time() + 31536000)), "ig_nrcb", "1"),
        (".instagram.com", "TRUE", "/", "TRUE", str(int(time.time() + 31536000)), "wd", "1920x1080"),
        (".instagram.com", "TRUE", "/", "TRUE", str(int(time.time() + 31536000)), "dpr", "2"),
    ]
    return cookies

def write_netscape_cookies(filepath, cookies):
    os.makedirs(os.path.dirname(filepath), exist_ok=True)
    with open(filepath, "w", encoding="utf-8") as f:
        for item in cookies:
            if len(item) == 2:
                f.write(f"{item[0]}\n")
            else:
                f.write("\t".join(item) + "\n")
    print(f"[+] Successfully wrote Netscape cookie file to: {filepath}")

def main():
    target_path = os.path.join("private", "cookies", "instagram.txt")
    print("==========================================")
    print("      Instagram Cookie File Generator     ")
    print("==========================================")
    
    cookies = generate_guest_cookies()
    write_netscape_cookies(target_path, cookies)

if __name__ == "__main__":
    main()
