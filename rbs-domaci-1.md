# Auditing in Software Systems

A **log file** represents a record of events in a computer system or network. Each entry contains information relevant to a specific event, valid or invalid login attempts, usage of resources such as creating, opening, or deleting files. Originally, log files were used for troubleshooting computer systems, but today log files have many more functions. They are used for optimizing systems and networks, monitoring user behavior, and generating data useful for investigating malicious program activities. Log files contain a large number of different pieces of information, so we can divide them into several categories depending on the installed services, and the most interesting for us is:

**Security software** (protection programs). It is most often used for detecting malicious program activities and protecting systems and data. It represents the main source of security log files. Types of security software are: Malware protection software, Intrusion Detection and Prevention Systems (IDPS), Remote access control software, Web proxies, Vulnerability management software, Authorization servers, Routers, Firewalls, Network quarantine servers.

**Audit mechanisms** represent monitoring, analyzing, and proving the flow of events in the system with the aim of detecting potential security risks, such as unauthorized access attempts, malicious behavior, authentication errors, and data manipulations.

In this work, I will focus on the analysis and research of how and in which ways it is possible to design a logging mechanism that fulfills the following requirements: efficient problem solving, actor identification or system non-repudiation, protection of sensitive data, reliability, time precision, and log clarity.  
---

### 1\. Information Required for Problem Solving

Efficient problem solving in information systems largely depends on the quality and content of log files. According to industry guidelines and best practices, logs should not represent only passive records of events, but an active tool for diagnostics and system monitoring.

According to the guide by Last9 company, each log entry should contain clear, structured, and contextual information that facilitates detecting the cause of the problem. The basic elements that a log entry must contain are:

* **Event time (timestamp)** \- The precise date and time when the event occurred enables temporal correlation of multiple events and discovery of cause-and-effect relationships. The standard format is ISO 8601 (e.g., 2025-07-09T15:23:44Z).  
* **Event type (log level)** \- Indicates the severity or purpose of the event, such as INFO, DEBUG, WARN, ERROR, or FATAL. This categorization allows filtering and focusing on the most critical entries, especially during the incident response process.  
* **System component (source/component)** \- The part of the application or module that generated the event. For example, AuthService or DatabaseHandler. Recognizing the component facilitates isolating the problem and assigning responsibility.  
* **Message with error context (message/exception)** \- A description of the event with additional information such as error code, exception, function call parameters, etc. This message should be clear and useful both for humans and for automated analysis tools.  
* **User or session identifier (user\_id, session\_id)** \- Enables linking logs to a specific user or session.

**Example of a Properly Formatted Log Entry:**

* `2025-07-09T15:23:44Z | ERROR | AuthService | Login failed for user_id=42 | reason=Invalid credential`

\*\*This log entry clearly shows that an error occurred in the authentication service at a precisely defined time (2025-07-09T15:23:44Z) and that it refers to the user with ID 42, along with an explanation of the failure.

Structured logs, especially in JSON format, further improve the possibilities of automated analysis and are well-suited for integrations with tools such as Logstash and Elasticsearch.

In this format, log data becomes easily searchable and analyzable within visualization tools like Kibana, which is of crucial importance for systematic detection of problems and security incidents.

---

### 2\. Non-repudiation and Actor Identification

Non-repudiation means that no actor (user, service, or system) can deny having performed a certain action.  
According to OWASP guidelines (2023), “systems must record enough information to determine which user performed a given action” where it is essential to log fields such as: unique user identifier (user\_id or username), session ID or token ID, IP address (if possible), and the type of action (e.g., data access, privilege change, or resource deletion).  
Log entries should record even failed authentication attempts, in order to detect potential attacks and allow for a complete reconstruction of user activity.

1. **Digital Signing and Chain of Verification**

One of the most reliable ways to ensure non-repudiation is digital signing of logs. According to HashiCorp Vault documentation, audit log devices can be configured so that every log event includes a digital signature, thereby ensuring that “no entry can be modified without being detected”.  
If a log entry has a digital signature, neither the user nor the administrator can later claim it did not exist \- because the signature can only be verified by an entity with the private key, while anyone can verify it with the public key.  
Vault also uses a hash chain, where each new entry includes the hash of the previous one, similar to blockchain \- which further guarantees the integrity and chronological order of logs. (A similar approach is used by Elastic’s Auditbeat).

2. **The Role of Authentication and Authorization**

Strong authentication is crucial for non-repudiation:

* token-based authentication (JWT, OAuth)

* two-factor authentication

* auditing of all failed and successful access attempts

If a user logged in using a valid token performs a DELETE operation, that token (or its identifier) is included in the log as evidence of the actor.

3. **Storing Logs in a System Users Cannot Modify**

Non-repudiation does not exist if a user can delete logs.  
AWS offers a so-called **Object Lock mechanism within its S3 service**, which allows log files to be stored in “Write Once Read Many” (WORM) mode.  
According to AWS documentation, with this approach users and administrators cannot delete or modify a record until the lock retention period has expired.

---

### 3\. Protection of Sensitive Data

One of the key aspects of secure logging is avoiding the recording of sensitive information in log files.  
Data that should never be logged:

* Passwords and PINs

* Access tokens (e.g., JWT, OAuth)

* Encryption keys

* Personally identifiable information (e.g., national ID numbers, credit card numbers)

* Secret configuration data (API keys, session cookies)

* Health-related data and other categories as defined by GDPR and HIPAA standards

OWASP recommends that such data should never be logged in its original form, but instead masked or pseudonymized, or completely omitted.

Example:

*   
  `2025-07-09T16:02:55Z | WARN | AuthService | Failed login attempt for user_id=42 with password=[REDACTED]`

or

* `password=[MASKED]`

\*\*In the end, even if logs are generated without sensitive data, it is important to secure their physical and system-level access.  
(Set **read-only permissions** on log directories and files; log files should be rotated and archived automatically, e.g., with **logrotate**).

---

### 4\. Reliability, Availability, and Integrity of Logs

**4\. 1\. Log Reliability**  
The goal of a reliable logging mechanism is to ensure that **no critical records are lost**, regardless of application crash, system restart, or network outage.  
According to AWS CloudWatch best practices, log systems must use **asynchronous approaches** (so that logging takes place independently from the main application flow), because *blocking I/O* can significantly impact application performance and even cause log loss during delays.  
In the event of unavailability of the central logging system (e.g., Elasticsearch), AWS recommends **local caching of logs and retry mechanisms** when the service is restored \- the logs are sent to the central system as soon as the network recovers, while temporary storage (buffer) is used only during the outage (e.g., disk queue or memory).

**4\. 2\. Log Availability**  
To ensure the availability of logs, a combination of the following techniques is used:

* **Log rotation**: using tools such as *logrotate*, logs are automatically rotated every *n* hours or *n* MB, which prevents disk overflow; old logs are archived (moved, compressed, or stored elsewhere).

* **Backup and replication**: AWS CloudWatch allows automatic backup of logs to S3, where they can be additionally replicated across multiple locations (Cross-Region Replication).

* **Retention policies**: define how long logs are retained, e.g., 30 days active \+ 1 year archived.

* **Log centralization**: logs are sent to Elasticsearch/Logstash/Kibana (ELK), Grafana Loki, Graylog, or a SIEM solution \- where they are available 24/7.

**4\. 3\. Log Integrity**  
The goal is for logs to be protected from subsequent modification or deletion and to provide proof that they were not altered. This goal is achieved using cryptographic methods and controlled access.

**a) Hashing log entries**

* Each log entry is hashed (SHA-256, SHA-3).  
* Alternative: a chain of hash values (like blockchain), where each new entry includes the hash of the previous one.

**b) Digital signing**

* Log files are periodically signed using the system’s private key, while verification is performed with the public key (RSA, ECDSA).

**c) Writing to write-once storage (WORM)**

* AWS S3 allows *Object Lock* in Compliance mode – the file then cannot be modified or deleted during the defined period.

**d) Monitoring access to logs (Audit logging of logs)**

* Access to the log files themselves must be monitored. Linux *auditd* records every attempt to read, delete, or modify log files. These records can be sent to a central SIEM system and used for anomaly detection (e.g., sudden attempts to access archived logs by unauthorized users).

### ---

**5\. Accurate Event Timestamping**

Time must be expressed in the standardized ISO 8601 format, including the time zone or UTC indicator, for example:

* `2025-07-09T15:12:03Z`

In distributed systems (e.g., when multiple servers are sending logs), it is important that all systems are time-synchronized. Otherwise, events from different sources will have illogical timestamps.  
In such cases, it is recommended to use **NTP (Network Time Protocol)**.  
NTP is a network protocol that enables synchronization of the local server clock with reference time servers. (*w32time; ntpd/chronyd*).

---

### 6\. Log Neatness and Readability

1. **Consistent and structured format of log entries** (using **JSON format** \- compatible with tools such as the ELK stack; allows filtering by fields or **key-value pairs**). Easy machine analysis and parsing.

2. **Log level must be clearly indicated**: INFO, DEBUG, WARN, ERROR, AUDIT, FATAL.  
   (distinguishes significant events from "noise", enables filtering by priority in tools like Kibana)

3. **Log rotation and file size limiting**: (rotation by time/size)  
   The old log file is renamed, a new file is opened for further logging, and the old file is compressed and archived.  
   A **retention policy** (e.g., 30-day retention) defines:  
   when old logs are deleted, when they are sent to long-term archives (e.g., AWS S3).

---

## Logrotate: Log Rotation in Traditional Systems

**What is logrotate?**  
Logrotate is a standard Linux tool for rotating, archiving, and managing log files.  
The main goal is to prevent log files from growing to the point where they become a problem for storage and performance.

**How it works:**

1. **Rotation** – renames active files (e.g., `app.log` → `app.log.1`)

2. **Creating a new file** – uses the **`create`** directive to open a new log file with specified permissions

3. **Compression** – compresses old files

4. **History rotation** – keeps a defined number of archives (e.g., the last 7), and deletes everything beyond that

**Known Security Vulnerabilities**

**1\. Race Condition (TOCTOU)**  
In versions prior to 3.7.9, logrotate was vulnerable to TOCTOU (time-of-check-to-time-of-use) attacks. An attacker with system access could attempt to exploit the time gap between when logrotate checks a file and when it rotates, deletes, or replaces it.

**Example:**  
A low-privilege user creates a symbolic link to a system directory, logrotate creates a file as root, the file is written with attacker-controlled content and signature, creating a privileged backdoor.  
The attacker's **symbolic link (symlink)** redirects log rotation to a sensitive file (e.g., /etc/passwd), and the log file gets overwritten, replaced, or moved.

\*The process\> logrotate runs as root (with the highest privileges), and when it rotates a log file, it renames or overwrites the existing log file and creates a new one where logging continues. That new file is created with root privileges.  
An attacker with low privileges creates a symbolic link (shortcut) pointing to a sensitive file or directory (ln \-s /etc/passwd /var/log/app.log), So, the file that looks like a log file actually points to /etc/passwd \- a critical system file.  
When *logrotate* rotates /var/log/app.log, it: checks the file, renames or deletes the existing file, creates a new file with the same name, starts writing the next logs to that file; but since /var/log/app.log is a symlink to /etc/passwd, these actions are actually performed on /etc/passwd.

*\*A symlink is a special type of file that acts as a shortcut or pointer to another file or directory.*

**Recommendations what to do:**

* Create new log files with the **`create`** directive and defined permissions; avoid using **`copytruncate`** (used when a process doesn’t support `SIGUSR1`; it creates a copy of the current log file and empties the original instead of renaming it \- race condition:  
  if the log file is modified while logrotate is making a copy and clearing the original, it can lead to inconsistent data)  
* Use the **`su`** directive in the configuration to perform log rotation as a less privileged user instead of root.   
  (New log files will have ownership and permissions that match the application. (If log files are owned by another user, e.g., www-data for a web server.))  
* Update logrotate to version **3.20.0 or higher** to avoid known vulnerabilities  
* Protect access to log directories by setting **proper permissions** and ownership:  
- Only the owner and possibly the admin group should have write (`w`) permission  
- Other users should have read-only (`r`) or no access at all

---

**2\. Log Injection (Log Forging)**  
If the application allows unfiltered user input to be written into a log file, an attacker can inject malicious lines.  
Unmodified and unfiltered user fields (e.g., username) may contain newlines, fabricated log entries, or even executable code, compromising log clarity, integrity, and analytical reliability.

* They can appear as system logs (fake ERROR messages, root login attempts, etc.)

* They can insert new lines or escape sequences that disrupt log structure

* They can deceive the person analyzing the logs later

**Recommendations what to do:**

**Never write user input directly into logs without sanitization.**

1. **Escape** special and problematic characters (\\n, \\r, ANSI codes \- which can break the log line and insert malicious content).  
   Convert them to safe representations or remove them entirely  
2. Use **structured log formats (e.g., JSON)** \- in JSON logs, injection is easier to detect because the format requires escaping special characters, and a parser will raise an error if the input is malformed.

---

### ELK Stack: Integration of Elasticsearch, Logstash, and Kibana for Log Ingestion, Processing, and Visualization with a Focus on Security Events

The **ELK stack** is a set of open-source tools commonly used for collecting, processing, and visualizing logs in real-time. It consists of:

* **Elasticsearch**: a distributed search and analytics database

* **Logstash**: a data processing and log collection system

* **Kibana**: a data visualization and analytics tool for data stored in Elasticsearch

#### How is a log ingestion pipeline configured?

**Logstash** acts as a central *pipeline* where logs arrive from various sources:

* Application logs (e.g., from Java, Python, Node.js applications)

* System logs (e.g., syslog from Linux servers)

* Web Application Firewall (WAF) logs and other security tools

A **Logstash pipeline** is defined in configuration files with three main blocks:

* input \- defines the log sources

* filter \- allows parsing, filtering, and data transformation

* output \- defines the destination of the data, typically **Elasticsearch**


#### Indexing and storing logs in Elasticsearch

Elasticsearch stores logs in the form of indices, which are logical collections of documents (log entries). According to the Elasticsearch documentation:

* Indices are time-oriented (e.g., daily indices logs-2025.07.10), which facilitates handling large volumes of data

* Each log is indexed as a **JSON document** with searchable and filterable fields

* You can define **mappings** for indices to specify field types (e.g., date, keyword, ip), which improves search and aggregation performance

* Elasticsearch automatically distributes data across cluster nodes to ensure **high availability and scalability**

* Configuration may also include a **retention policy** for automatically deleting or archiving old data

#### Visualization and filtering in Kibana

Kibana uses Elasticsearch as a data source and allows users to create **interactive dashboards**, **charts** (bar, pie, line, heatmap), **tables**, and **maps**. 

It supports:

* Visualizations for **security events**, such as: event timelines (geolocation maps of login attempts), geographic maps with IP addresses (number of failed logins per IP address (brute-force detection)), statistics on attack types or failed access attempts (top users triggering DELETE actions)

* Advanced filtering using **KQL (Kibana Query Language)** \- allows users to filter events precisely by user, event type, time, and other fields

* **Pivot tables** \- for summarizing, organizing, and analyzing large amounts of data in tabular form

* Creating **alerts and notifications** based on specific criteria \- enabling timely responses to security incidents

* Kibana also supports plugins for integration with **SIEM** (Security Information and Event Management) tools, as well as machine learning for detecting anomalies in logs. Kibana, within the **Elastic Security module**, provides specialized dashboards, visualizations, and tools for analyzing security events.

# C. Multi-factor authentication

**Multi-factor authentication (MFA)** is a security method that requires users to provide two or more types of verification to access an account or a system. It ensures an extra layer of protection beyond just a password by combining something that the user knows (like a password) with something the user has (like a phone with a TOTP authenticator app or a physical security token) or with something they are (like a fingerprint or a face scan).

## Types of multi-factor authentication

There are five different types of factors for multi-factor authentication, and any combination of these can be used. However, in practice, only the first three are common when it comes to securing web applications. These five types are as follows:

1. **Something You Know**  
2. **Something You Have**  
3. **Something You Are**  
4. **Somewhere You Are**  
5. **Something You Do**

What we should keep in mind is that requiring multiple instances of the same authentication factor (like requiring both password and PIN) **does not constitute the MFA**.

### Something You Know

This is the most common type of authentication. It is based on something the user knows, typically a password. Due to their simplicity for both users and developers, they are extremely popular. These factors do not require any additional hardware or integration with any other service.

### Something You Have

Possession-based authentication is based on the user having a physical or digital item that is required to authenticate. This is the most common form of MFA, typically used in combination with a password. Most common types of possession-based authentication are hardware and software tokens, as well as digital certificates. If implemented properly, they make it significantly more complicated for an attacker to compromise. The most significant disadvantage of this approach is creating an additional burden for the user, as they must keep the authentication factor with them whenever they want to use it.

### Something You Are

This type of authentication is based on the physical attributes of the user. It is less common in web applications as it requires users to have specialized hardware and raises significant privacy concerns. This type of authentication is commonly used for operating system authentication, as well as in some mobile applications, because modern smartphones often include specialized hardware for biometric-based authentication (e.g, some banking apps can use fingerprints for authentication).

This type of authentication introduces the usage of biometrics, including fingerprint scans, facial recognition, iris scans, or voice recognition.

Well-implemented biometrics are hard to spoof and require a targeted attack, while also being fast and convenient for users. However, biometric authentication has several drawbacks, including the need for manual user enrollment, often requiring custom, and sometimes expensive, additional hardware, raising privacy concerns due to the storage of sensitive physical data, presenting challenges if compromised since biometric data is difficult to change, and exposing hardware to additional attack vectors.

### Somewhere You Are

This type of authentication is location-based, which means it uses the user’s physical location as a verification factor. While some argue that location merely influences whether the MFA is required, this essentially treats it as a factor of its right. Notable implementations of location-based authentication include Conditional Access Policies in Microsoft Azure or the Network Unlock feature in BitLocker.

There are several approaches to how location-based authentication can be implemented. Some of them are the following:

* **Source IP address:**   
  * This approach uses the IP address user is connected from as an authentication factor. This type of authentication is typically implemented in an allow-list-based approach, whether it is a static (such as corporate office ranges) or a dynamic list (such as previous IP addresses the user has authenticated from).   
  * Advantages of this approach are simplicity for users and requiring minimal configuration and management from the administration staff.   
  * However, there are also several drawbacks, including not providing any protection if the user’s system is compromised, not providing any protection against rogue insiders, and requirement for careful restriction of trusted IP addresses.  
* **Geolocation:**   
  * Rather than using the exact IP address of the user, the geographic location where the IP address is registered can also be used. This approach is less precise, but much more feasible to implement where IP addresses are not static. A common usage for this type of authentication is when a user attempts to authenticate outside of their normal country.  
  * The main advantage of this approach is its simplicity for the users, as it does not require any user action.  
  * However, there are several limitations, including not providing any protection if the user’s system is compromised and not providing any protection against rogue insiders. It can also be easily bypassed by using trusted IP addresses, and may be less accurate due to privacy tools like VPNs and Apple’s iCloud private relay.  
* **Geofencing:**  
  * This is a more precise version of geolocation, which allows users to define a specific area in which they are allowed to authenticate. This approach is often used in mobile applications, where the user’s location can be easily determined with a high degree of accuracy using geopositioning hardware like GPS.  
  * Advantages of geofencing include simplicity for users and offering a high level of protection against remote attacks.  
  * As well as the other approaches, geofencing has several limitations. It doesn’t provide any protection if the user’s system is compromised, nor does it offer any protection against rogue insiders. One more limitation is not offering protection against attackers who are physically close to a trusted location.

### Something You Do

This type of authentication is based on a user’s behaviour, such as the way they type, move their mouse, or use their mobile device. This is the least common form of MFA, and it’s combined with other forms of authentication to increase the level of assurance in the user identity. It is also the most difficult to implement, and it may also require specific hardware as well as a significant amount of data and processing power to analyze users’ behaviour.

This type of authentication can be implemented in several forms:

* **Behavioural profiling:**  
  * This is based on how user interacts with the application, such as the time of day they usually log in, the device that they use, and the way they navigate the application.  
* **Keystroke and mouse dynamics:**  
  * This approach is based on the way the user types and moves their mouse. However, this type of authentication is still mostly theoretical and not widely used in practice.  
* **Gait analysis:**  
  * Gait analysis is based on the way the user walks using cameras and sensors. This is often used in physical security systems, but is not widely used in web applications. Mobile devices can use the accelerometer to detect the user’s gait and use it as an additional factor, but this is still mostly theoretical.

## Password \+ TOTP

Passwords are the most commonly used form of authentication, due to the simplicity of implementing them. Most multifactor authentication systems use a password in combination with at least one additional factor.

Password-based authentication is simple, widely understood, and natively supported across virtually all authentication frameworks, which makes it very easy to implement. However, these advantages come with significant drawbacks. Users often tend to choose weak passwords, reuse them across multiple systems, and remain vulnerable to phishing attacks.

**One-Time Password (OTP) tokens** are a form of possession-based authentication. Here, the user is required to submit a constantly changing numeric code to authenticate. The most common variety of OTPs is TOTP, which stands for Time-based One-Time Password. TOTP can be both hardware and software-based.

### Hardware OTP Tokens

Hardware OTP Tokens generate constantly changing numeric codes, which have to be submitted when authenticating. The most common example is the RSA SecureID, which generates a six-digit numeric code that changes every 60 seconds.

Advantages of this approach include:

* Hardware OTP tokens offer strong security since they are separate physical devices, which makes them extremely difficult to compromise remotely.  
* Additionally, they offer a convenient alternative for users who prefer not to use a mobile phone or other personal devices for authentication.

Despite their security benefits, hardware OTP Tokens come with several challenges:

* Deploying them to users can be logistically challenging and introduce additional costs.  
* Replacing lost or broken tokens can cause delays due to procurement and shipping.  
* Some of these systems require connection with a backend server, which can introduce additional vulnerabilities and create a single point of failure.  
* Moreover, stolen tokens can be used without a PIN or unlock code, and despite the fact that their codes are very short-lived, they still can be susceptible to phishing attacks.

### Software OTP Tokens

Software OTP Tokens are a cheaper and easier alternative to using Hardware OTP Tokens. This approach means that the software is being used for generating Time-Based One-Time Password (TOTP) codes. This typically involves users installing a TOTP application of their choice (or recommended by a distinct platform) on their mobile phone and then scanning a QR code provided by the web application, which provides the initial seed. After that, the authenticator app starts generating s six-digit numeric code every 30-60 seconds in a similar way as a hardware OTP Token.

Nowadays, most web applications use standardized TOTP tokens, which allow users to install any application that supports TOTP. On the other hand, some applications use their variations of this, which require users to install a specific app to use the service. This practice should be avoided in favour of a standards-based approach.

Authentication based on software TOTPs offers several advantages:

* It eliminates the need for physical tokens, which can help reduce the costs and administrative burden.  
* If a user loses access to their TOTP app, a new one can easily be set up without the delays and costs of shipping a replacement physical token.  
* Since TOTP is widely used and accepted, most users already have a TOTP app installed on their phone.  
* Additionally, if a phone is stolen, the codes remain protected as long as the device is secured with a screen lock.

Despite its advantages, authentication based on software TOTPs comes with several disadvantages:

* Since the TOTP app is usually installed on a mobile device, it’s prone to compromises, especially when the same device is used for both generating the code and authentication.  
* Users may also store backup seeds insecurely, increasing the risk of unauthorized access.  
* Additionally, not all users have compatible mobile devices.  
* If a device is lost, stolen, or out of battery, the user may be locked out.  
* While TOTP codes are short-lived, they can still be prone to phishing attacks.

## Password \+ TOTP flow based on the example of Google Authenticator

The password and TOTP authentication flow consists of several steps, some of which occur only once, while others are repeated during each authentication attempt.

The password \+ TOTP authentication flow when Google Authenticator is used is split into two phases. These phases are the **initial setup** and **authentication flow.** In the next part of the text, the complete overview of the flow is given.

1. **Initial setup** occurs only once and consists of the following steps:  
   1. The user logs into the web application and then navigates to the account security settings.  
   2. The user selects the option to enable two-factor authentication using an authenticator app.  
   3. The web application displays a QR code containing the TOTP secret key. The user opens the Google Authenticator on their mobile device and scans the QR code.  
   4. The web application may provide backup codes in case the user loses access to the authenticator app. This step is optional, but recommended.  
   5. The user enters a six-digit code generated by Google Authenticator to confirm the setup was successful.  
2. **Authentication flow** is the phase of the password \+ TOTP flow that occurs whenever authentication is attempted and consists of the following steps:  
   1. The user starts the flow by entering their usual login credentials. After this step, the user is prompted to a page where the TOTP code should be entered for the authentication process to be completed.  
   2. The user opens the app on their mobile phone to retrieve the current six-digit code.  
   3. The user enters the TOTP code from the app on the MFA page in the web application.  
   4. If both the password and code are correct, access is granted. If not, the user has to start the authentication flow again.

## Password and TOTP implementation

Password-based authentication remains the most widely adopted form of identity verification across digital platforms, despite the emergence of more advanced methods.

Password-based authentication involves two stages:

1. **User registration**:  
   1. The user creates a unique identifier for themselves, such as an email or a username, and a password.  
   2. The password is then **hashed** and never stored in plain text.  
   3. Optionally, a **salt** is added to defend against rainbow table attacks.  
2. **User login:**  
   1. User submits credentials.  
   2. The system retrieves the stored password hash and compares it against the hash of the password input.  
   3. If matched, access is granted.

As we already mentioned, passwords **should never be stored in plain text**. Instead of that, we are using **hashing algorithms.** Some examples of hashing algorithms are the following:

* **MD5** or **SHA1**: These algorithms were used by early systems because of their speed, but they are now obsolete.  
* **Bcrypt:** Bcrypt offers adaptive, built-in salt alongside moderate computational costs.  
* **Scrypt:** This hashing algorithm is memory-intensive to resist hardware attacks.  
* **Argon2:** This hashing algorithm was a winner of the Password Hashing Competition (PHC). It offers both memory-hardness and parallelism.

Most modern systems use bcrypt, scrypt, or Argon2.

Common vulnerabilities of password-based authentication are:

* **Brute-force attacks:** Manual or automated trials of password combinations.  
* **Dictionary attacks:** This type of attack is performed by using a list of common passwords.  
* **Credential stuffing:** This type of attack involves using leaked username/password combinations from past breaches.  
* **Keylogging or phishing:** For that kind of attack, malware or deception is used to steal credentials.  
* **Rainbow table attacks:** This type of attack is performed by using precomputed tables for reversing hashes.  
* **Insecure storage:** Storing passwords in an insecure way, without hashing or salting.

To strengthen password-based authentication and protect against common threats, several mitigation techniques can be employed. **Salting** is used to defend against precomputed hash attacks by ensuring that each user’s password hash is unique. **Rate limiting** and **CAPTCHA** mechanisms help block brute-force attacks by restricting the number of login attempts and verifying human input. Incorporating **multi-factor authentication (MFA)** adds a layer of security beyond the password itself. **Account lockouts** can throttle repeated failed login attempts, deterring automated attacks. Enforcing **strong password policies** with requirements for length, complexity, and periodic changes helps in reducing the likelihood of weak credentials. Finally, **credential leak detection** is implemented through integration with databases of known compromised passwords and can alert users and systems when a password may no longer be secure.

**One-Time Password (OTP) tokens** are a type of possession-based (*Something You Have*) authentication that requires submitting a numeric code that changes frequently. The most commonly used type is the **Time-Based One-Time Password (TOTP)**, which can be implemented through both hardware and software tokens.

### TOTP algorithm

**The Time-Based One-Time Password (TOTP) algorithm** is an extension of the HMAC-based One-Time Password (HOTP) algorithm. Unlike HOTP, which generates codes based on the counter, TOTP derives its codes from the current time, making them valid only for a short, predefined period of time.

At the core of the TOTP algorithm is the formula:

**TOTP \= HOTP(K, T)**

* **K** is a shared secret key between the client (user) and the server  
* **T** is a time-based counter calculated from the current time.

To compute **T**, the system uses the following formula:

**T \= (Current Unix Time \- T0) / X**

* **T0** is the starting point of the time counter, usually set to the Unix epoch, i.e., 00:00:00 UTC on January 1st, 1970  
* **X** is the time step in seconds, typically set to 30 seconds

For example, if the current time is 59 seconds and the time step is 30, the counter value is 1\. When the time reaches 60 seconds, it becomes 2\. This counter is calculated continuously, ensuring that a new code is generated at regular intervals.

The TOTP generation process is defined by the RFC 6238 standard and includes the following steps:

1. **Setting up a shared secret key:** During the initial setup (provisioning), the server and client agree on a shared secret key. This key is being stored securely on both sides and is often encoded in a QR code for easy setup.  
2. **Determining the current time:** This step includes retrieving the current Unix time in seconds (time since January 1st, 1970).  
3. **Generating a time-based counter:** This step includes dividing the current time by the configured time step (commonly 30 seconds) to compute the counter value.  
4. **Converting the counter to an 8-byte array:** Formatting the counter as an 8-byte array using big-endian byte order.  
5. **Applying HMAC-SHA1 hashing algorithm:** Hashing the counter using the HMAC-SHA1 hashing algorithm, with the shared secret key as an HMAC key. This produces a 20-byte hash.  
6. **Performing dynamic truncation:** By taking the last 4 bytes of the hash value (specifically, the 31st to 34th bytes) as an offset and performing a bitwise AND operation with a binary mask, we obtain a 4-byte dynamic truncation value.  
7. **Generating a numerical code:** By converting the truncated 4-byte chunk into an integer, then applying the modulo 10d (where **d** is the number of digits, usually 6), we are getting the final OTP.  
8. **Result formatting:** Padding the result value with leading zeros to achieve the desired number of digits.

### TOTP validation

Validating Time-Based One-Time Passwords (TOTP) requires careful handling of time synchronization and network delays to ensure both security and usability are maintained.

When a one-time password (OTP) arrives on the server, the system doesn’t know about the exact moment when the OTP was generated on the client. The server independently computes a code using the same shared secret and current Unix timestamp, applying HMAC alongside dynamic truncation. Because network latency may shift the reception time across time-step boundaries, the server must accommodate a small “time window” during validation.

To address this issue, servers commonly accept OTPs not only for the current time step, but also for adjacent ones (usually one step before and one step after the current one). For example, when a 30-second step is used, the server may accept codes generated within the window that lasts from 30 seconds before to 30 seconds after the current time. RFC 6238 recommends at most **one adjacent time step** to balance tolerance for delays against security.

### TOTP resynchronization

In Time-Based One-Time Password (TOTP) systems, both the user and server should use closely synchronized clocks for generating and validating one-time codes. However, slight differences in these device clocks can occur over time, causing a valid OTP code to be rejected. This phenomenon is called **clock drift**.

To handle this, the server implements a **resynchronization mechanism** that allows a small time discrepancy. Instead of validating OTP using only the current time step, the server is also able to check a limited number of previous or future time steps. For example, if the time step is set to 30 seconds, and the server allows up to two steps backwards, it will accept OTPs generated up to 89 seconds earlier. 

If an OTP is successfully validated from one of these adjacent steps, the server can record the time offset in the number of steps between itself and the client. Future OTPs can then be verified using this offset, effectively keeping the client and the server synchronized even if their clocks are slightly different.

In the case of a user not authenticating for a long period, the accumulated clock drift may become too large to be corrected automatically. In these cases, automatic resynchronization may fail, and additional methods, such as email or SMS verification, may be required for manually resynchronizing the user’s device with the server.

## Common mistakes in implementing MFA

Although using MFA has a lot of advantages, some common mistakes can happen during MFA implementation that can lead to serious security risks. Some of those common mistakes are:

* SMS is prone to SIM swapping, interception, etc., so it’s not considered a valuable option nowadays. Experts recommend using SMS only as a fallback option (for one-time use).  
* Accepting the same code multiple times weakens our security, so the good practice would be to track the last used time or invalidate the code after one use.  
* If recovery codes are short or poorly protected, they can be prone to brute-force attacks. Additionally, if recovery options are not tied to identity verification, that can lead to their exploitation.  
* MFA Fatigue and push bombing is an increasingly common approach to defeating MFA, where the user is bombarded with many requests to accept a login until the user eventually accepts one. When MFA applications are configured to send push notifications to the end user, attackers can send a flood of login attempts in the hope that the user will click on one of them. In September 2022, Uber's security was breached using a multi-factor authentication fatigue attack.  
* Excessive MFA prompts can frustrate users. Besides that, some users don’t want to use MFA applications and refuse to use them.  
* Partial enforcement (only for some users or apps) leaves security gaps.  
* Some older systems don’t support the MFA, which leads to them becoming a weak link. Risk increases significantly if MFA is skipped for critical legacy apps.  
* A common practice of storing OTPs in the same password manager as users’ login credentials defeats the principle of MFA independence. If the password manager is compromised, that will lead to the fall of both factors at once.  
* Whitelisting trusting IPs or skipping MFA by context (e.g., user’s geolocation) can be spoofed or misused.

### Potential challenges with using TOTP

Despite their widespread use and effectiveness, TOTP apps come with several potential challenges that can affect both usability and security:

* As we mentioned earlier, TOTP relies on precise time synchronization between the user’s device and the server. Even slight discrepancies can cause mismatched code, resulting in failed login attempts. This happens often when devices drift out due to the timezone changes or poor internet time updates.  
* Since TOTP apps are installed on mobile devices, the loss or theft of these devices can lead to user being locked out of their accounts. Access recovery often includes contacting customer service, which can be time-consuming, especially in cases when multiple services are affected.  
* Not all TOTP apps are compatible with all services. Some services require proprietary solutions, leading to users needing to install multiple authentication apps. This practice can confuse users and, therefore, lead to an increasing number of user mistakes.  
* Users tend not to store backup codes securely. This malpractice can lead to difficulties in the process of recovering access after device replacement.

Another serious TOTP vulnerability is backup codes theft. Backup codes are an essential part of the MFA system. They allow user to regain access to their accounts if primary methods aren’t available, for example, due to a stolen mobile device with the TOTP app installed. During MFA setup, the user is given a set of backup codes that should be securely stored and used in cases like this. Backup code is an alphanumeric code that is meant to be stored securely, preferably in a location separated from the user’s primary device. Improper storage practices can defeat the whole point of MFA. For example, if backup codes are stored in the same vault or password manager as the user’s credentials and a data breach occurs, both factors are going to be exposed to the attackers, and they could gain access to the user’s accounts despite MFA settings. Best practices to avoid backup code theft are:

* Generating long, one-time use backup codes. Recommended backup code length is 8+ alphanumeric characters.  
* Implementing strict rate limiting to prevent guessing attacks.  
* Invalidating each code after only one use and regenerating sets after every use.  
* Users should be properly educated not to keep backup codes in plaintext in the Notes app. Codes should be stored securely or encrypted in a database.  
* Keeping backup codes offline, printed, or in secure vaults, not on the device alongside the user’s credentials.

## Integrating MFA in an ELK environment

The ELK stack (Elasticsearch, Logstash, Kibana) is one of the most popular platforms for managing and visualizing log data. Like with any system that works with potentially sensitive information, ensuring secure access to ELK components, specifically Kibana, is essential. Adding MFA offers an effective security measure.

Although Kibana doesn’t natively support MFA, it can be implemented in a very simple way through identity federation using protocols such as SAML (Security Assertion Markup Language) in combination with a third-party identity provider like Okta, Azure Active Directory or Keycloak. These providers often support MFA, which can be enforced during the authentication process.

If an organization wants to secure administrator access to Kibana using Okta. The first step of this process involves enabling ELK’s native security features by modifying configuration files for both Elasticsearch and Kibana. These configurations allow Kibana to defer authentication responsibility to the identity provider. The administrator is then registered in Okta and an application is used for Kibana using the SAML protocol. The application specifies endpoints such as Assertion Consumer Service (ACS) URL and identity ID, which correspond to Kibana’s internal authentication mechanism. 

Once this setup is in place, Okta can enforce MFA policies on the Kibana application. For example, when an administrator attempts to log in to the Kibana application, they are redirected to Okta login page. After entering basic credentials (username and password), they are prompted for an additional verification step. An additional verification step can be approving push notifications via Okta Verify, entering a time-based one-time password (TOTP) or confirming an SMS code. Only after the MFA challenge is completed, the user is redirected back to Kibana with an active, authenticated session.

This approach offers several advantages:

* Centralizes identity and access management, making it easier to enforce security policies consistently across services.  
* Enhances protection against credential theft, phishing and brute-force attacks.  
* Provides seamless user experience when integrated with Single Sign-On (SSO) systems.

Alternatively, if the organization doesn’t utilize a SAML-compatible identity provider, a reverse-proxy-based solution can be implemented. For example, a combination of Nginx and Authelia can be deployed in front of Kibana. In this approach, Nginx handles traffic routing, while Authelia performs authentication and MFA enforcement. Only authenticated users who complete the MFA step are allowed to access the Kibana interface. This method provides flexibility and doesn’t rely on Elastic’s commercial features and external identity providers.

## References

* **ZT10 \- Forenzička analiza LOG datoteka.pdf:** https://singipedia.singidunum.ac.rs/izdanje/40135-forenzicka-analiza-log-datoteka  
* **Last9**, *Application Logging Best Practices*, available at: [https://last9.io/blog/application-logs](https://last9.io/blog/application-logs)

* **OWASP Logging Cheat Sheet**: [https://cheatsheetseries.owasp.org/cheatsheets/Logging\_Cheat\_Sheet.html](https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html)

* **MITRE CWE-117**: *Improper Output Neutralization for Logs* – [https://cwe.mitre.org/data/definitions/117.html](https://cwe.mitre.org/data/definitions/117.html)

* **OWASP Logging and Alerting Controls**

* [https://developer.hashicorp.com/vault/docs/audit](https://developer.hashicorp.com/vault/docs/audit)

* [https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-lock.html\#object-lock-overview](https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-lock.html#object-lock-overview)

* [https://owasp.org/www-project-top-ten/2017/A3\_2017-Sensitive\_Data\_Exposure](https://owasp.org/www-project-top-ten/2017/A3_2017-Sensitive_Data_Exposure)

* **AWS Docs** – *CloudWatch Logs Best Practices*

* [https://help.ubuntu.com/community/LogRotation](https://help.ubuntu.com/community/LogRotation)

* [https://www.elastic.co/docs/api/doc/elasticsearch/group/endpoint-indices](https://www.elastic.co/docs/api/doc/elasticsearch/group/endpoint-indices)  
* [https://www.elastic.co/guide/en/kibana/index.html](https://www.elastic.co/guide/en/kibana/index.html)  
* [\[OWASP\] Multi-factor authentication cheat sheet](https://cheatsheetseries.owasp.org/cheatsheets/Multifactor_Authentication_Cheat_Sheet.html)  
* [Research insight: Password-based authentication \- Foundations, Risks and Modern practices](https://medium.com/@nikhilmane372/research-insight-password-based-authentication-foundations-risks-modern-practices-adcfbbce8cb8)  
* [TOTP: Time-based One-Time Password Algorithm (RFC 6238\)](https://datatracker.ietf.org/doc/html/rfc6238)  
* [Two-factor authentication using TOTP](https://medium.com/@puran.joshi307/2-factor-authentication-using-totp-a9f1ff1e0b1a)  
* [The crucial interplay of time and IT: Multi-factor authentication, TOTP, and potential challenges](%20https://blog.waterloointuition.com/the-crucial-interplay-of-time-and-it-multi-factor-authentication-totp-and-potential-challenges/)  
* [Common pitfalls of MFA and how to avoid them](https://workos.com/blog/common-pitfalls-of-mfa-and-how-to-avoid-them)  
* [Security settings in Elasticsearch](https://www.elastic.co/guide/en/elasticsearch/reference/current/security-settings.html)  
* [SAML guide in Elasticsearch](https://www.elastic.co/guide/en/elasticsearch/reference/current/saml-guide.html)  
* [How to secure Kibana with SAML single sign-on](https://www.elastic.co/blog/how-to-secure-kibana-with-saml-single-sign-on)  
* [SAML authentication guide](https://developer.okta.com/docs/guides/saml-authentication/)  
* [Securing NGINX Plus using multi-factor authentication](https://www.nginx.com/blog/securing-nginx-plus-using-multi-factor-authentication/)