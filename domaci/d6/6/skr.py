#!/usr/bin/env python3
"""
Filesystem Security Audit Script
Performs comprehensive filesystem security review for Linux systems
Generates system_audit.log file with security issues display
"""

import os
import sys
import subprocess
import logging
import stat
from datetime import datetime

class FilesystemAudit:
    def __init__(self):
        self.setup_logging()
        self.results = {
            'timestamp': datetime.now().isoformat(),
            'security_issues': [],
            'warnings': [],
            'info_findings': []
        }
        
    def setup_logging(self):
        """Setup logging to system_audit.log"""
        logging.basicConfig(
            level=logging.INFO,
            format='%(asctime)s - %(levelname)s - %(message)s',
            handlers=[
                logging.FileHandler('system_audit.log', mode='w'),
                logging.StreamHandler(sys.stdout)
            ]
        )
        self.logger = logging.getLogger()
    
    def run_command(self, cmd, description):
        """Execute a shell command and log the result"""
        self.logger.info(f"Executing: {cmd}")
        try:
            result = subprocess.run(cmd, shell=True, capture_output=True, text=True, timeout=60)
            
            if result.returncode == 0:
                self.logger.info(f"[SUCCESS] {description}")
            else:
                self.logger.warning(f"[WARNING] {description} - Return code: {result.returncode}")
            
            return {
                'success': result.returncode == 0,
                'stdout': result.stdout,
                'stderr': result.stderr,
                'returncode': result.returncode
            }
        except subprocess.TimeoutExpired:
            self.logger.error(f"[ERROR] {description} - TIMEOUT")
            return {'success': False, 'error': 'Command timed out'}
        except Exception as e:
            self.logger.error(f"[ERROR] {description} - ERROR: {str(e)}")
            return {'success': False, 'error': str(e)}
    
    def check_mounted_partitions(self):
        """Review mounted partitions and fstab configuration"""
        self.logger.info("="*60)
        self.logger.info("MOUNTED PARTITIONS REVIEW")
        self.logger.info("="*60)
        
        # Check current mounts
        mounts_result = self.run_command('mount', 'Current mounted filesystems')
        fstab_result = self.run_command('cat /etc/fstab', 'Filesystem table configuration')
        
        self.analyze_mount_options(mounts_result, fstab_result)
    
    def analyze_mount_options(self, mounts_result, fstab_result):
        """Analyze mount options for security issues"""
        self.logger.info("")
        self.logger.info("MOUNT OPTIONS SECURITY ANALYSIS")
        self.logger.info("-" * 40)
        
        # Analyze current mounts
        if mounts_result['success']:
            for line in mounts_result['stdout'].split('\n'):
                if ' on ' in line and ' type ' in line:
                    parts = line.split()
                    if len(parts) >= 6:
                        device = parts[0]
                        mount_point = parts[2]
                        options = parts[5].strip('()')
                        
                        # Check for noatime (security issue)
                        if 'noatime' in options:
                            self.logger.error("[-] SECURITY ISSUE: noatime option found on {}".format(mount_point))
                            self.logger.error("    RECOMMENDATION: Remove noatime to preserve access times for forensics")
                            self.results['security_issues'].append("noatime option on {} - removes access time tracking".format(mount_point))
                        
                        # Check for missing security options
                        if mount_point in ['/tmp', '/var/tmp']:
                            if 'noexec' not in options:
                                self.logger.error("[-] SECURITY ISSUE: Missing noexec on {}".format(mount_point))
                                self.logger.error("    RECOMMENDATION: Add noexec to prevent binary execution")
                                self.results['security_issues'].append("Missing noexec on {}".format(mount_point))
                            if 'nosuid' not in options:
                                self.logger.error("[-] SECURITY ISSUE: Missing nosuid on {}".format(mount_point))
                                self.logger.error("    RECOMMENDATION: Add nosuid to prevent SUID execution")
                                self.results['security_issues'].append("Missing nosuid on {}".format(mount_point))
                        
                        if mount_point == '/home':
                            if 'nosuid' not in options:
                                self.logger.warning("[-] WARNING: Missing nosuid on /home")
                                self.results['warnings'].append("Missing nosuid on /home")
        
        # Analyze fstab for persistent configuration
        if fstab_result['success']:
            for line in fstab_result['stdout'].split('\n'):
                line = line.strip()
                if line and not line.startswith('#'):
                    parts = line.split()
                    if len(parts) >= 4:
                        mount_point = parts[1]
                        options = parts[3]
                        
                        if 'noatime' in options:
                            self.logger.error("[-] SECURITY ISSUE: noatime in fstab for {}".format(mount_point))
                            self.results['security_issues'].append("noatime in fstab for {}".format(mount_point))
    
    def check_sensitive_files(self):
        """Check permissions on sensitive files"""
        self.logger.info("")
        self.logger.info("="*60)
        self.logger.info("SENSITIVE FILES PERMISSIONS")
        self.logger.info("="*60)
        
        sensitive_files = [
            '/etc/shadow',
            '/etc/gshadow',
            '/etc/mysql/my.cnf',
            '/etc/ssl/private',
            '/root/.ssh',
            '/etc/passwd',
            '/etc/group'
        ]
        
        # Check common backup files
        backup_patterns = [
            '/etc/shadow.backup',
            '/etc/shadow.bak',
            '/etc/passwd.backup',
            '/etc/passwd.bak',
            '*.backup',
            '*.bak'
        ]
        
        for file_path in sensitive_files:
            self.check_file_permissions(file_path)
        
        # Check for backup files
        self.logger.info("")
        self.logger.info("BACKUP FILES CHECK")
        self.logger.info("-" * 40)
        for pattern in backup_patterns:
            self.find_backup_files(pattern)
        
        # Check web directory permissions
        self.check_web_directories()
    
    def check_file_permissions(self, file_path):
        """Check permissions for a specific file"""
        if os.path.exists(file_path):
            try:
                file_stat = os.stat(file_path)
                mode = file_stat.st_mode
                
                # Check if world-readable
                if mode & stat.S_IROTH:
                    self.logger.error("[-] SECURITY ISSUE: {} is world-readable".format(file_path))
                    self.logger.error("    RECOMMENDATION: Restrict to root access only")
                    self.results['security_issues'].append("{} is world-readable".format(file_path))
                
                # Check if world-writable
                if mode & stat.S_IWOTH:
                    self.logger.error("[-] CRITICAL ISSUE: {} is world-writable".format(file_path))
                    self.logger.error("    RECOMMENDATION: Immediate action required")
                    self.results['security_issues'].append("{} is world-writable - CRITICAL".format(file_path))
                
                # Check group permissions
                if mode & stat.S_IWGRP and file_path in ['/etc/shadow', '/etc/gshadow']:
                    self.logger.error("[-] SECURITY ISSUE: {} is group-writable".format(file_path))
                    self.results['security_issues'].append("{} is group-writable".format(file_path))
                    
            except Exception as e:
                self.logger.warning("[-] Could not check permissions for {}: {}".format(file_path, e))
        else:
            self.logger.info("[+] {} does not exist".format(file_path))
    
    def find_backup_files(self, pattern):
        """Find backup files with insecure permissions"""
        result = self.run_command('find / -name "{}" 2>/dev/null | grep -v "/proc"'.format(pattern), 
                                'Finding backup files: {}'.format(pattern))
        
        if result['success'] and result['stdout'].strip():
            for file_path in result['stdout'].split('\n'):
                file_path = file_path.strip()
                if file_path and os.path.exists(file_path):
                    self.check_backup_file_permissions(file_path)
    
    def check_backup_file_permissions(self, file_path):
        """Check permissions on backup files"""
        try:
            file_stat = os.stat(file_path)
            mode = file_stat.st_mode
            
            if mode & stat.S_IROTH:
                self.logger.error("[-] SECURITY ISSUE: Backup file {} is world-readable".format(file_path))
                self.logger.error("    RECOMMENDATION: Secure or remove backup files")
                self.results['security_issues'].append("Backup file {} is world-readable".format(file_path))
                
        except Exception as e:
            self.logger.warning("[-] Could not check backup file {}: {}".format(file_path, e))
    
    def check_web_directories(self):
        """Check web directory permissions"""
        web_dirs = ['/var/www', '/var/www/html', '/srv/www', '/usr/share/nginx/html']
        
        for web_dir in web_dirs:
            if os.path.exists(web_dir):
                self.check_directory_permissions(web_dir, "Web directory")
    
    def check_directory_permissions(self, dir_path, description):
        """Check directory permissions"""
        try:
            dir_stat = os.stat(dir_path)
            mode = dir_stat.st_mode
            
            if mode & stat.S_IWOTH:
                self.logger.error("[-] SECURITY ISSUE: {} {} is world-writable".format(description, dir_path))
                self.results['security_issues'].append("{} {} is world-writable".format(description, dir_path))
                
            if mode & stat.S_IROTH and dir_path in ['/backup', '/var/backups']:
                self.logger.error("[-] SECURITY ISSUE: Backup directory {} is world-readable".format(dir_path))
                self.results['security_issues'].append("Backup directory {} is world-readable".format(dir_path))
                
        except Exception as e:
            self.logger.warning("[-] Could not check directory {}: {}".format(dir_path, e))
    
    def check_setuid_files(self):
        """Find and analyze setuid files"""
        self.logger.info("")
        self.logger.info("="*60)
        self.logger.info("SETUID FILES REVIEW")
        self.logger.info("="*60)
        
        # Find all setuid files
        result = self.run_command('find / -perm -4000 -ls 2>/dev/null | grep -v "/proc"', 
                                'Finding setuid files')
        
        if result['success']:
            setuid_files = []
            for line in result['stdout'].split('\n'):
                if line.strip():
                    setuid_files.append(line)
                    self.logger.info("[FOUND] {}".format(line))
            
            self.analyze_setuid_files(setuid_files)
        else:
            self.logger.error("[-] Failed to find setuid files")
    
    def analyze_setuid_files(self, setuid_files):
        """Analyze setuid files for security risks"""
        self.logger.info("")
        self.logger.info("SETUID FILES SECURITY ANALYSIS")
        self.logger.info("-" * 40)
        
        legitimate_setuid = [
            '/bin/su',
            '/bin/ping',
            '/bin/ping6',
            '/bin/umount',
            '/bin/mount',
            '/usr/bin/passwd',
            '/usr/bin/sudo',
            '/usr/sbin/uuidd'
        ]
        
        suspicious_locations = [
            '/tmp/',
            '/var/tmp/',
            '/dev/shm/',
            '/home/',
            '/var/www/'
        ]
        
        for file_line in setuid_files:
            for legit in legitimate_setuid:
                if legit in file_line:
                    self.logger.info("[+] Legitimate setuid: {}".format(legit))
                    break
            else:
                # Check if file is in suspicious location
                for location in suspicious_locations:
                    if location in file_line:
                        self.logger.error("[-] SUSPICIOUS: Setuid file in risky location: {}".format(file_line))
                        self.results['security_issues'].append("Suspicious setuid file: {}".format(file_line))
                        break
                else:
                    self.logger.warning("[-] UNKNOWN setuid file: {}".format(file_line))
                    self.results['warnings'].append("Unknown setuid file: {}".format(file_line))
    
    def check_world_writable_files(self):
        """Find world-writable files"""
        self.logger.info("")
        self.logger.info("="*60)
        self.logger.info("WORLD-WRITABLE FILES REVIEW")
        self.logger.info("="*60)
        
        # Find world-writable files excluding /proc
        result = self.run_command('find / -type f -perm -002 ! -path "/proc/*" 2>/dev/null', 
                                'Finding world-writable files')
        
        if result['success']:
            world_writable_files = []
            for line in result['stdout'].split('\n'):
                if line.strip() and os.path.exists(line.strip()):
                    file_path = line.strip()
                    world_writable_files.append(file_path)
                    
                    # Check if it's a critical file
                    if any(critical in file_path for critical in ['/etc/', '/bin/', '/sbin/', '/usr/bin/', '/usr/sbin/']):
                        self.logger.error("[-] CRITICAL: World-writable system file: {}".format(file_path))
                        self.results['security_issues'].append("World-writable system file: {}".format(file_path))
                    elif any(web_dir in file_path for web_dir in ['/var/www/', '/srv/www/']):
                        self.logger.warning("[-] World-writable web file: {}".format(file_path))
                        self.results['warnings'].append("World-writable web file: {}".format(file_path))
            
            self.logger.info("[+] Found {} world-writable files".format(len(world_writable_files)))
    
    def check_backup_directories(self):
        """Check backup directories and permissions"""
        self.logger.info("")
        self.logger.info("="*60)
        self.logger.info("BACKUP DIRECTORIES REVIEW")
        self.logger.info("="*60)
        
        backup_dirs = ['/backup', '/var/backups', '/root/backup', '/home/backup', '/tmp/backup']
        
        for backup_dir in backup_dirs:
            if os.path.exists(backup_dir):
                self.logger.info("[FOUND] Backup directory: {}".format(backup_dir))
                self.check_directory_permissions(backup_dir, "Backup directory")
                
                # Check contents of backup directory
                self.check_backup_directory_contents(backup_dir)
    
    def check_backup_directory_contents(self, backup_dir):
        """Check contents of backup directory"""
        try:
            for item in os.listdir(backup_dir):
                item_path = os.path.join(backup_dir, item)
                if os.path.isfile(item_path):
                    # Check for sensitive backup files
                    if any(sensitive in item for sensitive in ['shadow', 'passwd', 'config', '.cnf', '.key', '.pem']):
                        self.check_backup_file_permissions(item_path)
        except Exception as e:
            self.logger.warning("[-] Could not list backup directory {}: {}".format(backup_dir, e))
    
    def generate_summary(self):
        """Generate summary report"""
        self.logger.info("")
        self.logger.info("="*60)
        self.logger.info("FILESYSTEM AUDIT SUMMARY")
        self.logger.info("="*60)
        
        # Display security issues
        if self.results['security_issues']:
            self.logger.error("SECURITY ISSUES FOUND: {}".format(len(self.results['security_issues'])))
            for i, issue in enumerate(self.results['security_issues'], 1):
                self.logger.error("  {}. {}".format(i, issue))
        else:
            self.logger.info("[+] No critical security issues found!")
        
        # Display warnings
        if self.results['warnings']:
            self.logger.warning("WARNINGS: {}".format(len(self.results['warnings'])))
            for i, warning in enumerate(self.results['warnings'], 1):
                self.logger.warning("  {}. {}".format(i, warning))
        else:
            self.logger.info("[+] No warnings identified")
        
        # Final assessment
        self.logger.info("")
        self.logger.info("FINAL ASSESSMENT")
        self.logger.info("-" * 40)
        
        critical_issues = len([issue for issue in self.results['security_issues'] if 'CRITICAL' in issue])
        
        if critical_issues > 0:
            self.logger.error("[-] CRITICAL: {} critical issues require immediate attention".format(critical_issues))
        elif self.results['security_issues']:
            self.logger.error("[-] NEEDS ATTENTION: Security issues found that require remediation")
        else:
            self.logger.info("[+] GOOD: Filesystem security configuration is acceptable")
    
    def run_full_audit(self):
        """Execute complete filesystem audit"""
        self.logger.info("Starting Comprehensive Filesystem Security Audit")
        self.logger.info("Audit Time: {}".format(datetime.now().strftime('%Y-%m-%d %H:%M:%S')))
        self.logger.info("")
        
        self.check_mounted_partitions()
        self.check_sensitive_files()
        self.check_setuid_files()
        self.check_world_writable_files()
        self.check_backup_directories()
        self.generate_summary()
        
        return self.results

def main():
    """Main function"""
    # Check if running as root
    if os.geteuid() != 0:
        print("WARNING: This audit requires root privileges for complete results")
        print("Run with: sudo python3 filesystem_audit.py")
        sys.exit(1)
    
    try:
        auditor = FilesystemAudit()
        results = auditor.run_full_audit()
        
        # Print final summary to console
        print("\n" + "="*60)
        print("FILESYSTEM AUDIT COMPLETE")
        print("="*60)
        print("Detailed results saved to: system_audit.log")
        
        if results['security_issues']:
            print("\nCRITICAL FINDINGS:")
            for issue in results['security_issues'][:5]:  # Show first 5 issues
                print("  [-] {}".format(issue))
            if len(results['security_issues']) > 5:
                print("  ... and {} more issues".format(len(results['security_issues']) - 5))
                
    except KeyboardInterrupt:
        print("\nAudit interrupted by user")
    except Exception as e:
        print("Error during filesystem audit: {}".format(e))

if __name__ == "__main__":
    main()