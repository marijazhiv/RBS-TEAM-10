#!/usr/bin/env python3
"""
Services Security Audit Script
Performs comprehensive services and configuration review for Linux systems
Generates system_audit.log file with security issues display
"""

import os
import sys
import subprocess
import logging
import re
import configparser
from datetime import datetime

class ServicesAudit:
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
    
    def check_running_services(self):
        """Identify running services and network listeners"""
        self.logger.info("="*60)
        self.logger.info("RUNNING SERVICES REVIEW")
        self.logger.info("="*60)
        
        # Get running processes
        ps_result = self.run_command('ps -edf', 'Running processes')
        
        # Get network listeners
        tcp_listeners = self.run_command('lsof -i TCP -n -P | grep LISTEN', 'TCP listening services')
        udp_listeners = self.run_command('lsof -i UDP -n -P', 'UDP listening services')
        ss_result = self.run_command('ss -tuln', 'Network sockets (ss)')
        
        self.analyze_running_services(ps_result, tcp_listeners, udp_listeners, ss_result)
    
    def analyze_running_services(self, ps_result, tcp_listeners, udp_listeners, ss_result):
        """Analyze running services for security issues"""
        self.logger.info("")
        self.logger.info("SERVICES SECURITY ANALYSIS")
        self.logger.info("-" * 40)
        
        # Analyze TCP listeners
        if tcp_listeners['success']:
            services_found = {}
            for line in tcp_listeners['stdout'].split('\n'):
                if 'LISTEN' in line:
                    parts = line.split()
                    if len(parts) >= 9:
                        service = parts[0]
                        port = parts[8]
                        services_found[port] = service
            
            # Check for common services and their security implications
            for port, service in services_found.items():
                if ':22' in port and 'sshd' in service:
                    self.logger.info("[+] SSH service detected on port 22")
                elif ':80' in port and 'apache' in service:
                    self.logger.info("[+] Apache HTTP service detected on port 80")
                elif ':443' in port and 'apache' in service:
                    self.logger.info("[+] Apache HTTPS service detected on port 443")
                elif ':3306' in port and 'mysql' in service:
                    self.logger.warning("[-] MySQL listening on network interface")
                    self.results['warnings'].append("MySQL listening on network interface - consider binding to localhost")
                elif ':21' in port:
                    self.logger.error("[-] SECURITY ISSUE: FTP service detected")
                    self.results['security_issues'].append("FTP service running - use SFTP/SSH instead")
                elif ':23' in port:
                    self.logger.error("[-] CRITICAL: Telnet service detected")
                    self.results['security_issues'].append("Telnet service running - extremely insecure")
    
    def check_ssh_configuration(self):
        """Review SSH server configuration"""
        self.logger.info("")
        self.logger.info("="*60)
        self.logger.info("SSH CONFIGURATION REVIEW")
        self.logger.info("="*60)
        
        sshd_config = '/etc/ssh/sshd_config'
        
        if os.path.exists(sshd_config):
            config_result = self.run_command(f'cat {sshd_config}', 'SSH daemon configuration')
            self.analyze_ssh_config(config_result['stdout'])
        else:
            self.logger.warning("[-] SSH configuration file not found")
    
    def analyze_ssh_config(self, config_content):
        """Analyze SSH configuration for security issues"""
        self.logger.info("")
        self.logger.info("SSH CONFIGURATION ANALYSIS")
        self.logger.info("-" * 40)
        
        lines = config_content.split('\n')
        permit_root_login = None
        protocol_version = None
        port = '22'
        tcp_forwarding = None
        
        for line in lines:
            line = line.strip()
            if not line or line.startswith('#'):
                continue
                
            if 'PermitRootLogin' in line:
                permit_root_login = line.split()[1]
            elif 'Protocol' in line:
                protocol_version = line.split()[1]
            elif 'Port' in line:
                port = line.split()[1]
            elif 'AllowTcpForwarding' in line:
                tcp_forwarding = line.split()[1]
        
        # Check PermitRootLogin
        if permit_root_login and permit_root_login.lower() == 'yes':
            self.logger.error("[-] SECURITY ISSUE: Root login permitted via SSH")
            self.logger.error("    RECOMMENDATION: Set PermitRootLogin to no")
            self.results['security_issues'].append("SSH PermitRootLogin enabled - disable root SSH access")
        else:
            self.logger.info("[+] Root SSH login is disabled")
        
        # Check Protocol version
        if protocol_version and '1' in protocol_version:
            self.logger.error("[-] SECURITY ISSUE: SSH Protocol 1 enabled")
            self.logger.error("    RECOMMENDATION: Use Protocol 2 only")
            self.results['security_issues'].append("SSH Protocol 1 enabled - use Protocol 2 only")
        else:
            self.logger.info("[+] SSH Protocol 2 is configured")
        
        # Check default port
        if port == '22':
            self.logger.warning("[-] Using default SSH port 22")
            self.logger.warning("    RECOMMENDATION: Consider changing default port")
            self.results['warnings'].append("SSH using default port 22 - consider changing")
        
        # Check TCP forwarding
        if tcp_forwarding and tcp_forwarding.lower() == 'yes':
            self.logger.warning("[-] TCP forwarding enabled in SSH")
            self.logger.warning("    RECOMMENDATION: Set AllowTcpForwarding to no if not needed")
            self.results['warnings'].append("SSH TCP forwarding enabled")
    
    def check_mysql_configuration(self):
        """Review MySQL configuration and security"""
        self.logger.info("")
        self.logger.info("="*60)
        self.logger.info("MYSQL CONFIGURATION REVIEW")
        self.logger.info("="*60)
        
        # Check MySQL configuration files
        mysql_configs = [
            '/etc/mysql/my.cnf',
            '/etc/mysql/mysql.conf.d/mysqld.cnf',
            '/etc/my.cnf'
        ]
        
        for config_file in mysql_configs:
            if os.path.exists(config_file):
                config_result = self.run_command(f'cat {config_file}', f'MySQL config: {config_file}')
                self.analyze_mysql_config(config_result['stdout'], config_file)
        
        # Check if MySQL is accessible and analyze users
        self.analyze_mysql_users()
    
    def analyze_mysql_config(self, config_content, config_file):
        """Analyze MySQL configuration for security issues"""
        bind_address = None
        
        for line in config_content.split('\n'):
            if 'bind-address' in line and not line.strip().startswith('#'):
                bind_address = line.split('=')[1].strip()
                break
        
        if bind_address == '0.0.0.0' or not bind_address:
            self.logger.error("[-] SECURITY ISSUE: MySQL binding to all interfaces")
            self.logger.error("    RECOMMENDATION: Set bind-address = 127.0.0.1")
            self.results['security_issues'].append("MySQL binding to all interfaces - restrict to localhost")
        elif bind_address == '127.0.0.1':
            self.logger.info("[+] MySQL binding to localhost only")
    
    def analyze_mysql_users(self):
        """Analyze MySQL users and privileges"""
        # Try to connect to MySQL without password (common misconfiguration)
        mysql_connect = self.run_command('mysql -u root -e "SELECT version();" 2>/dev/null', 
                                       'MySQL connection test')
        
        if mysql_connect['success']:
            self.logger.error("[-] CRITICAL: MySQL root access without password")
            self.results['security_issues'].append("MySQL root access without password - CRITICAL")
            
            # Get MySQL users and privileges
            users_result = self.run_command('mysql -u root -e "SELECT user, host, authentication_string FROM mysql.user;" 2>/dev/null', 
                                          'MySQL users')
            file_priv_result = self.run_command('mysql -u root -e "SELECT user, host, file_priv FROM mysql.user WHERE file_priv=\'Y\';" 2>/dev/null', 
                                              'MySQL FILE privileges')
            
            if file_priv_result['success'] and 'wordpress' in file_priv_result['stdout']:
                self.logger.error("[-] SECURITY ISSUE: WordPress user has FILE privilege")
                self.results['security_issues'].append("WordPress MySQL user has FILE privilege - remove unnecessary privilege")
        
        # Check for debian-sys-maint password exposure
        debian_cnf = '/etc/mysql/debian.cnf'
        if os.path.exists(debian_cnf):
            debian_result = self.run_command(f'cat {debian_cnf}', 'MySQL Debian maintenance config')
            if 'password' in debian_result['stdout']:
                self.logger.warning("[-] MySQL debian-sys-maint password in config file")
                self.results['warnings'].append("MySQL debian-sys-maint password in config file - ensure proper permissions")
    
    def check_apache_configuration(self):
        """Review Apache web server configuration"""
        self.logger.info("")
        self.logger.info("="*60)
        self.logger.info("APACHE CONFIGURATION REVIEW")
        self.logger.info("="*60)
        
        # Check Apache running user
        user_result = self.run_command('ps aux | grep apache | grep -v grep', 'Apache processes')
        self.analyze_apache_processes(user_result['stdout'])
        
        # Check Apache configuration files
        apache_configs = [
            '/etc/apache2/apache2.conf',
            '/etc/apache2/envvars',
            '/etc/apache2/conf-enabled/security.conf',
            '/etc/apache2/sites-enabled/000-default.conf'
        ]
        
        for config_file in apache_configs:
            if os.path.exists(config_file):
                config_result = self.run_command(f'cat {config_file}', f'Apache config: {config_file}')
                self.analyze_apache_config(config_result['stdout'], config_file)
        
        # Check web directory permissions
        self.check_web_directory_permissions()
    
    def analyze_apache_processes(self, processes_output):
        """Analyze Apache processes for security issues"""
        if 'root' in processes_output and 'apache' in processes_output:
            self.logger.error("[-] SECURITY ISSUE: Apache running as root")
            self.results['security_issues'].append("Apache running as root - change to non-privileged user")
        else:
            self.logger.info("[+] Apache running as non-root user")
    
    def analyze_apache_config(self, config_content, config_file):
        """Analyze Apache configuration for security issues"""
        if 'envvars' in config_file:
            if 'APACHE_RUN_USER=www-data' in config_content:
                self.logger.info("[+] Apache configured to run as www-data")
        
        if 'security.conf' in config_file:
            if 'ServerTokens Prod' not in config_content:
                self.logger.warning("[-] Apache ServerTokens not set to Prod")
                self.results['warnings'].append("Apache ServerTokens not minimized - information disclosure risk")
            
            if 'ServerSignature Off' not in config_content:
                self.logger.warning("[-] Apache ServerSignature not disabled")
                self.results['warnings'].append("Apache ServerSignature enabled - information disclosure risk")
        
        if '000-default' in config_file:
            if 'Indexes' in config_content and 'Options' in config_content:
                if 'Indexes' in config_content and '-Indexes' not in config_content:
                    self.logger.error("[-] SECURITY ISSUE: Directory listing enabled")
                    self.results['security_issues'].append("Apache directory listing enabled - information disclosure risk")
    
    def check_web_directory_permissions(self):
        """Check web directory permissions"""
        web_dirs = ['/var/www', '/var/www/html', '/srv/www']
        
        for web_dir in web_dirs:
            if os.path.exists(web_dir):
                perm_result = self.run_command(f'ls -la {web_dir}', f'Web directory permissions: {web_dir}')
                
                # Check for world-writable files
                find_result = self.run_command(f'find {web_dir} -type f -perm -o+w 2>/dev/null', 
                                             f'World-writable files in {web_dir}')
                
                if find_result['success'] and find_result['stdout'].strip():
                    self.logger.error("[-] SECURITY ISSUE: World-writable files in web directory")
                    self.results['security_issues'].append("World-writable files in web directory - security risk")
    
    def check_php_configuration(self):
        """Review PHP configuration"""
        self.logger.info("")
        self.logger.info("="*60)
        self.logger.info("PHP CONFIGURATION REVIEW")
        self.logger.info("="*60)
        
        php_configs = [
            '/etc/php5/apache2/php.ini',
            '/etc/php/7.4/apache2/php.ini',
            '/etc/php/8.0/apache2/php.ini',
            '/etc/php/8.1/apache2/php.ini'
        ]
        
        for config_file in php_configs:
            if os.path.exists(config_file):
                config_result = self.run_command(f'cat {config_file}', f'PHP config: {config_file}')
                self.analyze_php_config(config_result['stdout'])
                break
        
        # Check for Suhosin
        self.check_suhosin_configuration()
    
    def analyze_php_config(self, config_content):
        """Analyze PHP configuration for security issues"""
        issues = []
        
        if 'expose_php = On' in config_content:
            issues.append("PHP expose_php enabled - information disclosure")
        
        if 'display_errors = On' in config_content:
            issues.append("PHP display_errors enabled - information disclosure")
        
        if 'allow_url_include = On' in config_content:
            issues.append("PHP allow_url_include enabled - security risk")
        
        if issues:
            for issue in issues:
                self.logger.error(f"[-] SECURITY ISSUE: {issue}")
                self.results['security_issues'].append(issue)
        else:
            self.logger.info("[+] PHP security settings are properly configured")
    
    def check_suhosin_configuration(self):
        """Check Suhosin PHP protection"""
        suhosin_config = '/etc/php5/conf.d/suhosin.ini'
        
        if os.path.exists(suhosin_config):
            config_result = self.run_command(f'egrep -v -e "^;|^$" {suhosin_config}', 
                                           'Suhosin configuration')
            if 'suhosin.executor.disable_eval' not in config_result['stdout']:
                self.logger.warning("[-] Suhosin eval protection not enabled")
                self.results['warnings'].append("Suhosin eval protection not configured")
        else:
            self.logger.info("[+] Suhosin not installed (consider installing for additional PHP protection)")
    
    def check_crontab_security(self):
        """Review cron jobs and their security"""
        self.logger.info("")
        self.logger.info("="*60)
        self.logger.info("CRONTAB SECURITY REVIEW")
        self.logger.info("="*60)
        
        # Get system crontabs
        crontab_result = self.run_command('cat /etc/crontab', 'System crontab')
        cron_daily_result = self.run_command('ls -la /etc/cron.daily/', 'Daily cron jobs')
        cron_hourly_result = self.run_command('ls -la /etc/cron.hourly/', 'Hourly cron jobs')
        
        # Check individual user crontabs
        users_result = self.run_command('ls /home/', 'System users')
        if users_result['success']:
            for user in users_result['stdout'].split():
                user_cron = self.run_command(f'crontab -u {user} -l 2>/dev/null', 
                                           f'Crontab for user {user}')
                if user_cron['success']:
                    self.analyze_crontab(user_cron['stdout'], user)
        
        # Check cron script permissions
        self.check_cron_script_permissions()
    
    def analyze_crontab(self, crontab_content, user):
        """Analyze crontab entries for security issues"""
        for line in crontab_content.split('\n'):
            if line.strip() and not line.startswith('#'):
                # Look for script executions
                if '.sh' in line or '.py' in line:
                    # Extract script path
                    parts = line.split()
                    for part in parts:
                        if '.sh' in part or '.py' in part:
                            script_path = part
                            if os.path.exists(script_path):
                                self.check_script_permissions(script_path, user)
                            break
    
    def check_script_permissions(self, script_path, user):
        """Check permissions on cron scripts"""
        perm_result = self.run_command(f'ls -la {script_path}', f'Permissions: {script_path}')
        
        if 'rwxrwxrwx' in perm_result['stdout'] or 'rwxr-xr-x' in perm_result['stdout']:
            if 'root' in user:
                self.logger.error("[-] CRITICAL: World-accessible root cron script")
                self.results['security_issues'].append(f"World-accessible root cron script: {script_path}")
            else:
                self.logger.warning(f"[-] World-accessible cron script for user {user}")
                self.results['warnings'].append(f"World-accessible cron script: {script_path}")
    
    def check_cron_script_permissions(self):
        """Check permissions on cron script directories"""
        cron_dirs = ['/etc/cron.d', '/etc/cron.daily', '/etc/cron.hourly', '/etc/cron.weekly', '/etc/cron.monthly']
        
        for cron_dir in cron_dirs:
            if os.path.exists(cron_dir):
                perm_result = self.run_command(f'ls -la {cron_dir}', f'Permissions: {cron_dir}')
    
    def generate_summary(self):
        """Generate summary report"""
        self.logger.info("")
        self.logger.info("="*60)
        self.logger.info("SERVICES AUDIT SUMMARY")
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
        
        critical_count = len([issue for issue in self.results['security_issues'] if 'CRITICAL' in issue])
        
        if critical_count > 0:
            self.logger.error("[-] CRITICAL: {} critical issues require immediate attention".format(critical_count))
        elif self.results['security_issues']:
            self.logger.error("[-] NEEDS ATTENTION: Security issues found that require remediation")
        else:
            self.logger.info("[+] GOOD: Services security configuration is acceptable")
    
    def run_full_audit(self):
        """Execute complete services audit"""
        self.logger.info("Starting Comprehensive Services Security Audit")
        self.logger.info("Audit Time: {}".format(datetime.now().strftime('%Y-%m-%d %H:%M:%S')))
        self.logger.info("")
        
        self.check_running_services()
        self.check_ssh_configuration()
        self.check_mysql_configuration()
        self.check_apache_configuration()
        self.check_php_configuration()
        self.check_crontab_security()
        self.generate_summary()
        
        return self.results

def main():
    """Main function"""
    # Check if running as root
    if os.geteuid() != 0:
        print("WARNING: This audit requires root privileges for complete results")
        print("Run with: sudo python3 services_audit.py")
        sys.exit(1)
    
    try:
        auditor = ServicesAudit()
        results = auditor.run_full_audit()
        
        # Print final summary to console
        print("\n" + "="*60)
        print("SERVICES AUDIT COMPLETE")
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
        print("Error during services audit: {}".format(e))

if __name__ == "__main__":
    main()