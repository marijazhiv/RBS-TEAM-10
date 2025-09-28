#!/usr/bin/env python3
"""
Network Review Script
Performs comprehensive network configuration review for Debian systems
Generates system_audit.log file
"""

import os
import sys
import subprocess
import logging
from datetime import datetime

class NetworkReview:
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
            result = subprocess.run(cmd, shell=True, capture_output=True, text=True, timeout=30)
            
            if result.returncode == 0:
                self.logger.info(f"[SUCCESS] {description}")
            else:
                self.logger.warning(f"[WARNING] {description} - Return code: {result.returncode}")
            
            # Log the output
            if result.stdout.strip():
                self.logger.info(f"Output:\n{result.stdout}")
            if result.stderr.strip():
                self.logger.warning(f"Errors:\n{result.stderr}")
                
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
    
    def check_general_info(self):
        """Retrieve basic system configuration information"""
        self.logger.info("="*60)
        self.logger.info("GENERAL NETWORK INFORMATION")
        self.logger.info("="*60)
        
        commands = [
            ('ifconfig -a', 'Network interfaces (ifconfig)'),
            ('ip addr show', 'Network interfaces (ip addr)'),
            ('route -n', 'System routes'),
            ('netstat -rn', 'System routes (netstat)'),
            ('cat /etc/resolv.conf', 'DNS configuration'),
            ('cat /etc/hosts', 'Hosts file'),
            ('cat /etc/nsswitch.conf', 'Name service switch configuration'),
            ('cat /etc/network/interfaces', 'Network interfaces configuration')
        ]
        
        for cmd, desc in commands:
            self.run_command(cmd, desc)
    
    def check_firewall_rules(self):
        """Review firewall configuration"""
        self.logger.info("")
        self.logger.info("="*60)
        self.logger.info("FIREWALL CONFIGURATION")
        self.logger.info("="*60)
        
        # Check current iptables rules
        iptables_result = self.run_command('iptables -L -v -n', 'Current iptables rules')
        
        # Check if rules persist after reboot
        persistence_result = self.run_command(
            'cat /etc/network/if-pre-up.d/iptables 2>/dev/null || echo "File not found"', 
            'Firewall persistence script'
        )
        
        saved_rules_result = self.run_command(
            'cat /etc/iptables.up.rules 2>/dev/null || echo "File not found"', 
            'Saved firewall rules'
        )
        
        # Analyze firewall configuration
        self.analyze_firewall(iptables_result, persistence_result, saved_rules_result)
    
    def analyze_firewall(self, iptables_result, persistence_result, saved_rules_result):
        """Analyze firewall rules for security issues"""
        self.logger.info("")
        self.logger.info("FIREWALL SECURITY ANALYSIS")
        self.logger.info("-" * 40)
        
        if not iptables_result['success']:
            self.logger.error("[-] Cannot analyze firewall - iptables command failed")
            return
        
        iptables_output = iptables_result['stdout']
        
        # Check if SSH is open to everyone
        if 'tcp dpt:22' in iptables_output or 'tcp dpt:ssh' in iptables_output:
            if '0.0.0.0/0' in iptables_output or 'anywhere' in iptables_output:
                self.logger.error("[-] SECURITY ISSUE: SSH port 22 is open to everyone")
                self.logger.error("    RECOMMENDATION: Implement IP whitelisting for SSH access")
                self.results['security_issues'].append("SSH port 22 is open to everyone - implement IP whitelisting")
            else:
                self.logger.info("[+] SSH access is properly restricted")
        
        # Check HTTP port exposure
        if 'tcp dpt:80' in iptables_output or 'tcp dpt:www' in iptables_output:
            if '0.0.0.0/0' in iptables_output or 'anywhere' in iptables_output:
                self.logger.info("[+] HTTP port 80 is open (expected for web server)")
            else:
                self.logger.info("[+] HTTP port 80 is properly restricted")
        
        # Check HTTPS port exposure
        if 'tcp dpt:443' in iptables_output or 'tcp dpt:https' in iptables_output:
            self.logger.info("[+] HTTPS port 443 configuration detected")
        
        # Check outgoing traffic restrictions
        if 'Chain OUTPUT' in iptables_output:
            output_policy = self.extract_policy(iptables_output, 'OUTPUT')
            if output_policy == 'ACCEPT':
                self.logger.error("[-] SECURITY ISSUE: Outgoing traffic is not restricted")
                self.logger.error("    RECOMMENDATION: Implement outgoing traffic restrictions")
                self.results['security_issues'].append("Outgoing traffic is not restricted")
            else:
                self.logger.info("[+] Outgoing traffic is restricted (good practice)")
        
        # Check INPUT policy
        input_policy = self.extract_policy(iptables_output, 'INPUT')
        if input_policy == 'DROP' or input_policy == 'REJECT':
            self.logger.info(f"[+] INPUT chain policy is {input_policy} (secure)")
        else:
            self.logger.warning(f"[-] INPUT chain policy is {input_policy} (less secure)")
            self.results['warnings'].append(f"INPUT chain policy is {input_policy} - consider using DROP")
        
        # Check firewall persistence
        if 'File not found' in persistence_result.get('stdout', ''):
            self.logger.error("[-] SECURITY ISSUE: No firewall persistence script found")
            self.logger.error("    RECOMMENDATION: Create /etc/network/if-pre-up.d/iptables")
            self.results['security_issues'].append("No firewall persistence script - rules may not survive reboot")
        else:
            self.logger.info("[+] Firewall persistence script is configured")
        
        # Check if saved rules match current rules
        if saved_rules_result['success'] and 'File not found' not in saved_rules_result.get('stdout', ''):
            self.logger.info("[+] Saved firewall rules file exists")
        else:
            self.logger.warning("[-] No saved firewall rules file found")
            self.results['warnings'].append("No saved firewall rules file")
    
    def extract_policy(self, rules_text, chain_name):
        """Extract policy from iptables output"""
        for line in rules_text.split('\n'):
            if f'Chain {chain_name}' in line:
                parts = line.split()
                if len(parts) >= 4:
                    return parts[3].strip('()')
        return None
    
    def check_ipv6(self):
        """Check IPv6 configuration"""
        self.logger.info("")
        self.logger.info("="*60)
        self.logger.info("IPv6 CONFIGURATION")
        self.logger.info("="*60)
        
        ip6tables_result = self.run_command('ip6tables -L -v -n', 'IPv6 firewall rules')
        ipv6_status_result = self.run_command('sysctl net.ipv6.conf.all.disable_ipv6', 'IPv6 disable status')
        ipv6_addr_result = self.run_command('ip -6 addr show', 'IPv6 addresses')
        
        # Analyze IPv6 configuration
        if ip6tables_result['success']:
            if 'Chain INPUT (policy ACCEPT' in ip6tables_result['stdout']:
                self.logger.error("[-] SECURITY ISSUE: IPv6 firewall has ACCEPT policy")
                self.logger.error("    RECOMMENDATION: Apply same firewall rules to IPv6 or disable IPv6")
                self.results['security_issues'].append("IPv6 firewall has ACCEPT policy - apply restrictions")
            else:
                self.logger.info("[+] IPv6 firewall has restrictive policy")
        
        # Check if IPv6 is disabled
        if ipv6_status_result['success']:
            if 'net.ipv6.conf.all.disable_ipv6 = 1' in ipv6_status_result['stdout']:
                self.logger.info("[+] IPv6 is disabled system-wide")
            else:
                self.logger.warning("[-] IPv6 is enabled")
                self.results['warnings'].append("IPv6 is enabled - ensure proper firewall rules")
        
        # Check for IPv6 addresses
        if ipv6_addr_result['success'] and ipv6_addr_result['stdout'].strip():
            self.logger.info("[+] IPv6 addresses are configured")
            self.results['info_findings'].append("IPv6 addresses are configured on the system")
    
    def check_network_services(self):
        """Check running network services"""
        self.logger.info("")
        self.logger.info("="*60)
        self.logger.info("NETWORK SERVICES")
        self.logger.info("="*60)
        
        commands = [
            ('ss -tuln', 'Listening ports and services'),
            ('netstat -tuln', 'Listening ports (alternative)'),
            ('systemctl list-units --type=service --state=running', 'Running services'),
            ('ps aux | grep -E "(ssh|apache|nginx|mysql|postgres)" | grep -v grep', 'Key service processes')
        ]
        
        for cmd, desc in commands:
            result = self.run_command(cmd, desc)
            
            # Analyze listening services for security issues
            if 'ss -tuln' in cmd or 'netstat -tuln' in cmd:
                if result['success']:
                    self.analyze_listening_services(result['stdout'])
    
    def analyze_listening_services(self, services_output):
        """Analyze listening services for security issues"""
        self.logger.info("")
        self.logger.info("LISTENING SERVICES ANALYSIS")
        self.logger.info("-" * 40)
        
        for line in services_output.split('\n'):
            if 'LISTEN' in line:
                # Check for potentially risky services
                if ':22 ' in line and 'tcp' in line:
                    self.logger.info("[+] SSH service is listening (expected)")
                elif ':80 ' in line and 'tcp' in line:
                    self.logger.info("[+] HTTP service is listening (expected for web server)")
                elif ':443 ' in line and 'tcp' in line:
                    self.logger.info("[+] HTTPS service is listening (expected for web server)")
                elif ':21 ' in line and 'tcp' in line:
                    self.logger.warning("[-] FTP service is listening (consider using SFTP/SSH)")
                    self.results['warnings'].append("FTP service detected - consider using more secure alternatives")
                elif ':23 ' in line and 'tcp' in line:
                    self.logger.error("[-] SECURITY ISSUE: Telnet service is listening (insecure)")
                    self.results['security_issues'].append("Telnet service detected - replace with SSH")
                elif ':25 ' in line and 'tcp' in line:
                    self.logger.info("[+] SMTP service is listening (expected for mail server)")
    
    def generate_summary(self):
        """Generate summary report"""
        self.logger.info("")
        self.logger.info("="*60)
        self.logger.info("AUDIT SUMMARY")
        self.logger.info("="*60)
        
        # Display security issues
        if self.results['security_issues']:
            self.logger.error(f"SECURITY ISSUES FOUND: {len(self.results['security_issues'])}")
            for i, issue in enumerate(self.results['security_issues'], 1):
                self.logger.error(f"  {i}. {issue}")
        else:
            self.logger.info("[+] No critical security issues found!")
        
        # Display warnings
        if self.results['warnings']:
            self.logger.warning(f"WARNINGS: {len(self.results['warnings'])}")
            for i, warning in enumerate(self.results['warnings'], 1):
                self.logger.warning(f"  {i}. {warning}")
        else:
            self.logger.info("[+] No warnings identified")
        
        # Display informational findings
        if self.results['info_findings']:
            self.logger.info(f"INFORMATIONAL FINDINGS: {len(self.results['info_findings'])}")
            for i, finding in enumerate(self.results['info_findings'], 1):
                self.logger.info(f"  {i}. {finding}")
        
        # System information
        self.logger.info("")
        self.logger.info("SYSTEM INFORMATION")
        self.logger.info("-" * 40)
        self.run_command('uname -a', 'System information')
        self.run_command('lsb_release -a', 'Distribution information')
        
        self.logger.info("")
        self.logger.info("Audit completed. Full log saved to: system_audit.log")
        
        # Final summary
        self.logger.info("")
        self.logger.info("FINAL SCORE")
        self.logger.info("-" * 40)
        total_issues = len(self.results['security_issues']) + len(self.results['warnings'])
        if total_issues == 0:
            self.logger.info("[+] EXCELLENT: No security issues or warnings found")
        elif len(self.results['security_issues']) == 0:
            self.logger.info("[+] GOOD: No critical security issues, but review warnings")
        else:
            self.logger.error("[-] NEEDS ATTENTION: Critical security issues found")
    
    def run_full_review(self):
        """Execute complete network review"""
        self.logger.info("Starting Comprehensive Network Security Audit")
        self.logger.info(f"Audit Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        self.logger.info("")
        
        self.check_general_info()
        self.check_firewall_rules()
        self.check_ipv6()
        self.check_network_services()
        self.generate_summary()
        
        return self.results

def main():
    """Main function"""
    # Check if running as root
    if os.geteuid() != 0:
        print("WARNING: Some commands require root privileges")
        print("For complete results, run: sudo python3 network_review.py")
        print("Continuing with limited functionality...")
        print("")
    
    try:
        reviewer = NetworkReview()
        results = reviewer.run_full_review()
        
        # Print final message to console
        print("\n" + "="*60)
        print("AUDIT COMPLETE")
        print("="*60)
        print("Check system_audit.log for detailed results")
        
        if results['security_issues']:
            print(f"\nCRITICAL ISSUES: {len(results['security_issues'])}")
            for issue in results['security_issues']:
                print(f"  [-] {issue}")
        
        if results['warnings']:
            print(f"\nWARNINGS: {len(results['warnings'])}")
            for warning in results['warnings']:
                print(f"  [!] {warning}")
                
    except KeyboardInterrupt:
        print("\nAudit interrupted by user")
    except Exception as e:
        print(f"Error during audit: {e}")

if __name__ == "__main__":
    main()