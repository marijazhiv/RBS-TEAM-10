import subprocess
import sys
import os
from datetime import datetime

def run_cmd(cmd):
    try:
        result = subprocess.run(cmd, shell=True, capture_output=True, text=True, timeout=30)
        return result.stdout, result.stderr, result.returncode
    except subprocess.TimeoutExpired:
        return None, "Command timed out", -1
    except Exception as e:
        return None, str(e), -1
    
def log_message(message, log_file):
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    formatted_message = f"[{timestamp}] {message}"
    # print(formatted_message)
    with open(log_file, "a") as f:
        f.write(formatted_message + "\n")

def check_os_info(log_file):
    log_message("[+] Checking OS information...", log_file)

    stdout, stderr, returncode = run_cmd("cat /etc/os-release")
    if stdout:
        log_message("[+] OS Release Info:", log_file)
        for line in stdout.split('\n'):
            if line.strip() and '=' in line:
                key, value = line.split('=', 1)
                if key in ['PRETTY_NAME', 'NAME', 'VERSION_ID']:
                    cleaned_value = value.strip().strip('"')
                    log_message(f"  {key}: {cleaned_value}", log_file)
    else:
        log_message("[-] Cannot read /etc/os-release", log_file)

    stdout, stderr, returncode = run_cmd("cat /etc/debian_version 2>/dev/null")
    if stdout and returncode == 0:
        log_message(f"[+] Debian Version: {stdout.strip()}", log_file)

    stdout, stderr, returncode = run_cmd("lsb_release -a 2>/dev/null")
    if stdout and returncode == 0:
        log_message("[+] LSB Release Info:", log_file)
        for line in stdout.split('\n'):
            if line.strip():
                log_message(f"  {line}", log_file)

    stdout, stderr, returncode = run_cmd("cat /etc/os-release | grep VERSION_ID=")
    if stdout:
        version = stdout.strip().split('=')[1].strip('"')
        if "7" in version or "8" in version:    
            log_message("[-] SECURITY ISSUE: Old OS version detected!", log_file)
            log_message("   Recommend: Upgrade to a supported version", log_file)

def check_kernel_info(log_file):
    log_message("[+] Checking kernel information...", log_file)

    stdout, stderr, returncode = run_cmd("uname -a")
    if stdout:
        kernel_info = stdout.strip()
        log_message(f"[+] Kernel: {kernel_info}", log_file)

        if "2.6." in kernel_info or "3." in kernel_info:
            log_message("[-] SECURITY ISSUE: Old kernel version detected!", log_file)
            log_message("   Recommendation: Upgrade kernel to supported version", log_file)

    stdout, stderr, returncode = run_cmd("uptime")
    if stdout:
        uptime_info = stdout.strip()
        log_message(f"[+] Uptime: {uptime_info}", log_file)

        if "days" in uptime_info:
            days = int(uptime_info.split(" days")[0].split()[-1])
            if days > 30:
                log_message("[-] SECURITY ISSUE: System has long uptime!", log_file)
                log_message("   Recommendation: Regular reboots after security updates", log_file)

    stdout, stderr, returncode = run_cmd("cat /proc/loadavg")
    if stdout:
        load_avg = stdout.strip()
        log_message(f"[+] Load average: {load_avg}", log_file)

        load_values = load_avg.split()[:3]
        if float(load_values[0]) > 5.0:
            log_message("[-] SECURITY ISSUE: High system load detected!", log_file)
            log_message("   Recommendation: Investigate high load causes", log_file)

def check_time_config(log_file):
    log_message("[+] Checking time configuration...", log_file)

    stdout, stderr, returncode = run_cmd("cat /etc/timezone 2>/dev/null || timedatectl show --property=Timezone 2>/dev/null | cut -d= -f2")
    if stdout:
        timezone = stdout.strip()
        log_message(f"[+] Timezone: {timezone}", log_file)

        if "UTC" not in timezone and "GMT" not in timezone:
            log_message("[-] SECURITY ISSUE: Non-UTC timezone detected!", log_file)
            log_message("   Recommendation: Use UTC timezone for servers", log_file)
    else:
        log_message("[-] Cannot determine timezone", log_file)

    stdout, stderr, returncode = run_cmd("ps -edf | grep ntp | grep -v grep")
    if stdout:
        log_message("[+] NTP service running:", log_file)
        for line in stdout.split('\n'):
            if line.strip():
                log_message(f"  {line.strip()}", log_file)
    else:
        log_message("[-] SECURITY ISSUE: NTP service is not running!", log_file)
        log_message("   Recommendation: Install and configure NTP service", log_file)

    stdout, stderr, returncode = run_cmd("ntpq -p 2>/dev/null || ntpstat 2>/dev/null")
    if stdout:
        log_message("[+] NTP status:", log_file)
        for line in stdout.split('\n'):
            if line.strip():
                log_message(f"  {line.strip()}", log_file)
    else:
        log_message("[-] SECURITY ISSUE: NTP not synchronized!", log_file)
        log_message("   Recommendation: Check NTP configuration", log_file)

    stdout, stderr, returncode = run_cmd("timedatectl status 2>/dev/null | grep synchronized")
    if stdout and "yes" in stdout:
        log_message("[+] System time is synchronized", log_file)
    elif stdout and "no" in stdout:
        log_message("[-] SECURITY ISSUE: System time not synchronized!", log_file)
        log_message("   Recommendation: Configure time synchronization", log_file)

def check_installed_packages(log_file):
    log_message("[+] Checking installed packages...", log_file)
    
    stdout, stderr, returncode = run_cmd("dpkg -l | wc -l")
    if stdout and returncode == 0:
        count = int(stdout.strip())
        log_message(f"[+] Number of installed packages (dpkg): {count}", log_file)

        if count > 1000:
            log_message("[-] SECURITY ISSUE: Too many packages installed!", log_file)
            log_message("   Recommendation: Remove unnecessary packages", log_file)
        
        server_inappropriate_packages = [
            "xserver", "xorg", "wayland", "gnome", "kde", "xfce", "lxde",
            "mate", "cinnamon", "screensaver", "wallpaper", "game", "steam",
            "libreoffice", "thunderbird", "firefox", "chromium", "vlc",
            "gimp", "inkscape", "blender", "audacity"
        ]
        
        found_packages = []
        for pkg in server_inappropriate_packages:
            stdout, stderr, returncode = run_cmd(f"dpkg -l | grep -i {pkg} | head -1")
            if stdout and stdout.strip():
                found_packages.append(stdout.strip())
        
        if found_packages:
            log_message("[-] SECURITY ISSUE: Inappropriate packages for a server!", log_file)
            log_message("   Recommendation: Remove GUI/games/multimedia packages", log_file)
            for pkg in found_packages:
                parts = pkg.split()
                if len(parts) >= 3:
                    log_message(f"   {parts[1]} - {parts[2]}", log_file)
        else:
            log_message("[+] No obviously inappropriate packages found", log_file)
    
    stdout, stderr, returncode = run_cmd("rpm -qa 2>/dev/null | wc -l")
    if stdout and returncode == 0:
        count = int(stdout.strip())
        log_message(f"[+] Number of installed packages (rpm): {count}", log_file)

        stdout, stderr, returncode = run_cmd("rpm -qa | grep -E -i '(x11|gnome|kde|games)' | head -5")
        if stdout and stdout.strip():
            log_message("[-] SECURITY ISSUE: Inappropriate RPM packages found!", log_file)
            log_message("   Recommendation: Remove GUI/games packages", log_file)
            for line in stdout.split('\n'):
                if line.strip():
                    log_message(f"  {line.strip()}", log_file)

def check_logging_config(log_file):
    log_message("[+] Checking logging configuration...", log_file)
    
    stdout, stderr, returncode = run_cmd("ps -edf | grep -E '(rsyslog|syslog-ng|systemd-journal)' | grep -v grep")
    if stdout:
        log_message("[+] Logging service is running:", log_file)
        for line in stdout.split('\n'):
            if line.strip():
                log_message(f"   {line.strip()}", log_file)
    else:
        log_message("[-] SECURITY ISSUE: No logging service detected!", log_file)
        log_message("   Recommendation: Install and configure rsyslog or syslog-ng", log_file)
    
    stdout, stderr, returncode = run_cmd("systemctl is-active systemd-journald 2>/dev/null")
    if stdout and "active" in stdout:
        log_message("[+] Systemd journal is active", log_file)
    else:
        log_message("[-] No systemd journal available", log_file)
    
    config_files = ["/etc/rsyslog.conf", "/etc/rsyslog.d/", "/etc/syslog-ng/"]
    config_found = False
    
    for config_file in config_files:
        if os.path.exists(config_file):
            config_found = True
            if os.path.isfile(config_file):
                stdout, stderr, returncode = run_cmd(f"grep -v '^#' {config_file} | grep -v '^$' | head -5")
                if stdout:
                    log_message(f"[+] {config_file} configuration:", log_file)
                    for line in stdout.split('\n'):
                        if line.strip():
                            log_message(f"   {line.strip()}", log_file)
            elif os.path.isdir(config_file):
                log_message(f"[+] Found logging config directory: {config_file}", log_file)
    
    if not config_found:
        log_message("[-] No logging configuration files found", log_file)
    
    stdout, stderr, returncode = run_cmd("grep -r '@' /etc/rsyslog.* /etc/syslog-ng/* 2>/dev/null | grep -v '#' | head -5")
    if stdout:
        log_message("[+] Remote logging configuration found:", log_file)
        for line in stdout.split('\n'):
            if line.strip():
                log_message(f"   {line.strip()}", log_file)
    else:
        log_message("[-] SECURITY ISSUE: No remote logging configuration found!", log_file)
        log_message("   Recommendation: Configure remote logging to central log server", log_file)
    
    stdout, stderr, returncode = run_cmd("find /var/log -name '*.log' -o -name 'messages' -o -name 'syslog' | head -5")
    if stdout:
        log_message("[+] Found log files:", log_file)
        for line in stdout.split('\n'):
            if line.strip():
                log_message(f"   {line.strip()}", log_file)
    else:
        log_message("[-] No log files found in /var/log", log_file)
    
    stdout, stderr, returncode = run_cmd("ls /etc/logrotate.d/ 2>/dev/null | head -5")
    if stdout:
        log_message("[+] Log rotation configurations found:", log_file)
        for line in stdout.split('\n'):
            if line.strip():
                log_message(f"   {line.strip()}", log_file)
    else:
        log_message("[-] SECURITY ISSUE: No log rotation configuration found!", log_file)
        log_message("   Recommendation: Configure log rotation", log_file)

def check_sudo_access(log_file):
    log_message("[+] Checking sudo access and system privileges...", log_file)

    stdout, stderr, returncode = run_cmd("cat /etc/sudoers 2>/dev/null | grep -v '^#' | grep -v '^$'")
    if stdout:
        log_message("[+] Sudoers configuration:", log_file)
        for line in stdout.split('\n'):
            if line.strip():
                log_message(f"   {line.strip()}", log_file)

    log_message("[+] Checking sudo privileges for all users...", log_file)

    stdout, stderr, returncode = run_cmd("getent passwd | grep -v '/false$' | grep -v '/nologin$' | cut -d: -f1")
    if stdout:
        users = [user.strip() for user in stdout.split('\n') if user.strip()]

        for user in users:
            stdout, stderr, returncode = run_cmd(f"sudo -l -U {user} 2>/dev/null")
            if stdout and "may run the following commands" in stdout:
                log_message(f"[!] User {user} has sudo privileges:", log_file)

                lines = stdout.split('\n')
                commands_section = False
                for line in lines:
                    if "may run the following commands" in line:
                        commands_section = True
                        continue
                    if commands_section and line.strip():
                        if "NOPASSWD" in line:
                            log_message(f"  SECURITY ISSUE: {user} can run without password: {line.strip()}", log_file)
                        else:
                            log_message(f"  {line.strip()}", log_file)

    stdout, stderr, returncode = run_cmd("grep -r 'NOPASSWD' /etc/sudoers* 2>/dev/null")
    if stdout:
        log_message("[-] SECURITY ISSUE: NOPASSWD sudo rules found!", log_file)
        log_message("   Recommendation: Require passwords for all sudo commands", log_file)
        for line in stdout.split('\n'):
            if line.strip():
                log_message(f"  {line.strip()}", log_file)

    stdout, stderr, returncode = run_cmd("getent passwd | cut -d: -f1,3 | grep ':0' | grep -v '^root:'")
    if stdout:
        log_message("[-] SECURITY ISSUE: Other users with UID 0 found!", log_file)
        log_message("   Recommendation: Remove or change UID for these users", log_file)
        for line in stdout.split('\n'):
            if line.strip():
                log_message(f"  {line.strip()}", log_file)

    stdout, stderr, returncode = run_cmd("getent shadow | cut -d: -f1,2 | grep '::' | grep -v '^root:'")
    if stdout:
        log_message("[-] SECURITY ISSUE: Users without passwords found!", log_file)
        log_message("   Recommendation: Set passwords for these accounts", log_file)
        for line in stdout.split('\n'):
            if line.strip():
                log_message(f"   {line.strip()}", log_file)

    stdout, stderr, returncode = run_cmd("grep -r 'PermitRootLogin\\|PasswordAuthentication' /etc/ssh/ 2>/dev/null | grep -v '#'")
    if stdout:
        log_message("[+] SSH configuration:", log_file)
        for line in stdout.split('\n'):
            if line.strip():
                if "PermitRootLogin yes" in line or "PasswordAuthentication yes" in line:
                    log_message(f"[-] SECURITY ISSUE: {line.strip()}", log_file)
                    log_message("  Recommendation: Disable root login and password authentication", log_file)
                else:
                    log_message(f"  {line.strip()}", log_file)

def module1_os_kernel_packages(log_file):
    log_message("=" * 60, log_file)
    log_message("MODULE 1: OS, KERNEL, PACKAGES & LOGGING", log_file)
    log_message("="*60, log_file)

    check_os_info(log_file)
    check_kernel_info(log_file)
    check_time_config(log_file)
    check_installed_packages(log_file)
    check_logging_config(log_file)
    check_sudo_access(log_file)

    log_message("[+] Module 1 completed", log_file)
    log_message("="*60, log_file)

if __name__ == '__main__':
    if os.geteuid() != 0:
        print("[-] This script requires root privileges. Run with sudo.")
        sys.exit(1)

    log_file = "system_audit.log"
    module1_os_kernel_packages(log_file)