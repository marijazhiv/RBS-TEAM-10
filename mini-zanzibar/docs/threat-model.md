# Threat Model - Mini-Zanzibar Authorization System

## System Overview

Mini-Zanzibar is a global authorization system that manages access control through:
- LevelDB for storing ACL tuples
- ConsulDB for namespace configuration with versioning
- REST API for ACL and namespace management

## Assets

### Primary Assets:
1. **Authorization Data (ACL Tuples)** - Critical
   - Format: object#relation@user
   - Stored in LevelDB
   - Controls access to protected resources

2. **Namespace Configurations** - High
   - Defines authorization rules and relations
   - Stored in ConsulDB with versioning
   - Controls authorization logic

3. **API Service** - High
   - Provides authorization decisions
   - Manages ACL and namespace operations

### Supporting Assets:
1. **Database Systems** (LevelDB, ConsulDB)
2. **Configuration Data**
3. **Log Files**
4. **System Infrastructure**

## Threat Actors

### External Attackers
- **Skill Level**: Medium to High
- **Motivation**: Data theft, system disruption, privilege escalation
- **Access**: Network access to API endpoints

### Malicious Insiders
- **Skill Level**: High
- **Motivation**: Data theft, sabotage, unauthorized access
- **Access**: Potential system-level access

### Compromised Applications
- **Skill Level**: Variable
- **Motivation**: Lateral movement, data access
- **Access**: API access through compromised client applications

## Attack Vectors and Threats

### 1. API Security Threats

#### T1.1: Unauthorized API Access
- **Description**: Attackers bypass authentication to access API endpoints
- **Impact**: High - Unauthorized ACL modifications, data exposure
- **Likelihood**: Medium
- **Mitigations**: 
  - TODO: Implement JWT authentication
  - TODO: API rate limiting
  - TODO: Request validation

#### T1.2: Injection Attacks
- **Description**: SQL/NoSQL injection through API parameters
- **Impact**: High - Database compromise, data manipulation
- **Likelihood**: Medium
- **Mitigations**:
  - TODO: Input validation and sanitization
  - TODO: Parameterized queries
  - TODO: Database access controls

#### T1.3: Authorization Bypass
- **Description**: Flaws in authorization logic allow privilege escalation
- **Impact**: Critical - Complete access control bypass
- **Likelihood**: High
- **Mitigations**:
  - TODO: Comprehensive authorization testing
  - TODO: Secure coding practices
  - TODO: Regular security reviews

### 2. Data Storage Threats

#### T2.1: Database Compromise
- **Description**: Direct access to LevelDB or ConsulDB
- **Impact**: Critical - Complete ACL data exposure/manipulation
- **Likelihood**: Medium
- **Mitigations**:
  - TODO: Database encryption at rest
  - TODO: Access controls and monitoring
  - TODO: Network segmentation

#### T2.2: Configuration Tampering
- **Description**: Unauthorized modification of namespace configurations
- **Impact**: High - Authorization logic manipulation
- **Likelihood**: Medium
- **Mitigations**:
  - TODO: Configuration versioning and integrity checks
  - TODO: Administrative access controls
  - TODO: Change auditing

### 3. Network Security Threats

#### T3.1: Man-in-the-Middle Attacks
- **Description**: Interception of API communications
- **Impact**: High - Credential theft, data exposure
- **Likelihood**: Medium
- **Mitigations**:
  - TODO: TLS encryption for all communications
  - TODO: Certificate validation
  - TODO: Network monitoring

#### T3.2: Denial of Service
- **Description**: Resource exhaustion attacks on API endpoints
- **Impact**: Medium - Service unavailability
- **Likelihood**: High
- **Mitigations**:
  - TODO: Rate limiting implementation
  - TODO: Resource monitoring
  - TODO: Load balancing

### 4. Application Logic Threats

#### T4.1: Race Conditions
- **Description**: Concurrent ACL operations causing inconsistent state
- **Impact**: Medium - Data inconsistency, temporary privilege escalation
- **Likelihood**: Medium
- **Mitigations**:
  - TODO: Proper concurrency controls
  - TODO: Database transactions
  - TODO: State validation

#### T4.2: Business Logic Flaws
- **Description**: Flaws in ACL evaluation logic
- **Impact**: High - Incorrect authorization decisions
- **Likelihood**: Medium
- **Mitigations**:
  - TODO: Comprehensive testing
  - TODO: Code reviews
  - TODO: Formal verification methods

## Process Models

### ACL Creation Process
1. **Input**: API request with object, relation, user
2. **Validation**: Format validation, authentication check
3. **Authorization**: Check permissions for ACL management
4. **Storage**: Store tuple in LevelDB
5. **Logging**: Audit log entry
6. **Response**: Confirmation or error

**Threats**: Unauthorized creation, injection attacks, data corruption

### Authorization Check Process
1. **Input**: API request with object, relation, user
2. **Validation**: Format validation
3. **Namespace Lookup**: Retrieve relation configuration from Consul
4. **Evaluation**: Apply authorization rules (direct, computed usersets)
5. **Caching**: Store result for performance (TODO)
6. **Response**: Authorization decision

**Threats**: Logic bypass, cache poisoning, performance attacks

### Namespace Management Process
1. **Input**: API request with namespace configuration
2. **Validation**: Configuration format and logic validation
3. **Authorization**: Administrative access check
4. **Versioning**: Create new version in Consul
5. **Activation**: Update latest version pointer
6. **Logging**: Configuration change audit

**Threats**: Configuration tampering, privilege escalation, logic flaws

## Risk Assessment Matrix

| Threat | Impact | Likelihood | Risk Level | Priority |
|--------|---------|------------|------------|----------|
| T1.3 - Authorization Bypass | Critical | High | Critical | 1 |
| T2.1 - Database Compromise | Critical | Medium | High | 2 |
| T1.1 - Unauthorized API Access | High | Medium | High | 3 |
| T1.2 - Injection Attacks | High | Medium | High | 4 |
| T2.2 - Configuration Tampering | High | Medium | High | 5 |
| T3.1 - MITM Attacks | High | Medium | High | 6 |
| T4.2 - Business Logic Flaws | High | Medium | High | 7 |
| T4.1 - Race Conditions | Medium | Medium | Medium | 8 |
| T3.2 - Denial of Service | Medium | High | Medium | 9 |

## Recommendations

### High Priority:
1. Implement comprehensive authorization logic with proper testing
2. Add authentication and authorization middleware
3. Implement input validation and sanitization
4. Add database encryption and access controls

### Medium Priority:
1. Implement TLS for all communications
2. Add comprehensive audit logging
3. Implement rate limiting and DDoS protection
4. Add configuration integrity checks

### Low Priority:
1. Implement caching with security considerations
2. Add monitoring and alerting
3. Implement backup and recovery procedures
4. Add performance testing and optimization