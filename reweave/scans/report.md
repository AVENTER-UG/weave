# Vulnerability Report

```
Report date: 2024-04-04
Unique vulnerability count: 17
Images version: 2.8.5-beta1
```

## Scanner Details

```
Application:         grype
Version:             0.74.7
BuildDate:           2024-02-26T18:24:14Z
GitCommit:           987238519b8d6e302130ab715f20daed6634da68
GitDescription:      v0.74.7
Platform:            linux/amd64
GoVersion:           go1.21.7
Compiler:            gc
Syft Version:        v0.105.1
Supported DB Schema: 5
```

## Vulnerabilities

### weave-kube: (17) 

```
NAME           INSTALLED   FIXED-IN  TYPE  VULNERABILITY   SEVERITY 
busybox        1.36.1-r15            apk   CVE-2023-42366  Medium    
busybox        1.36.1-r15            apk   CVE-2023-42365  Medium    
busybox        1.36.1-r15            apk   CVE-2023-42364  Medium    
busybox        1.36.1-r15            apk   CVE-2023-42363  Medium    
busybox-binsh  1.36.1-r15            apk   CVE-2023-42366  Medium    
busybox-binsh  1.36.1-r15            apk   CVE-2023-42365  Medium    
busybox-binsh  1.36.1-r15            apk   CVE-2023-42364  Medium    
busybox-binsh  1.36.1-r15            apk   CVE-2023-42363  Medium    
curl           8.5.0-r0              apk   CVE-2024-0853   Medium    
curl           8.5.0-r0              apk   CVE-2024-2466   Unknown   
curl           8.5.0-r0              apk   CVE-2024-2398   Unknown   
curl           8.5.0-r0              apk   CVE-2024-2004   Unknown   
libuv          1.47.0-r0             apk   CVE-2024-24806  High      
ssl_client     1.36.1-r15            apk   CVE-2023-42366  Medium    
ssl_client     1.36.1-r15            apk   CVE-2023-42365  Medium    
ssl_client     1.36.1-r15            apk   CVE-2023-42364  Medium    
ssl_client     1.36.1-r15            apk   CVE-2023-42363  Medium
```

### weave-npc: (12) 

```
NAME           INSTALLED   FIXED-IN  TYPE  VULNERABILITY   SEVERITY 
busybox        1.36.1-r15            apk   CVE-2023-42366  Medium    
busybox        1.36.1-r15            apk   CVE-2023-42365  Medium    
busybox        1.36.1-r15            apk   CVE-2023-42364  Medium    
busybox        1.36.1-r15            apk   CVE-2023-42363  Medium    
busybox-binsh  1.36.1-r15            apk   CVE-2023-42366  Medium    
busybox-binsh  1.36.1-r15            apk   CVE-2023-42365  Medium    
busybox-binsh  1.36.1-r15            apk   CVE-2023-42364  Medium    
busybox-binsh  1.36.1-r15            apk   CVE-2023-42363  Medium    
ssl_client     1.36.1-r15            apk   CVE-2023-42366  Medium    
ssl_client     1.36.1-r15            apk   CVE-2023-42365  Medium    
ssl_client     1.36.1-r15            apk   CVE-2023-42364  Medium    
ssl_client     1.36.1-r15            apk   CVE-2023-42363  Medium
```

### weave: (17) 

```
NAME           INSTALLED   FIXED-IN  TYPE  VULNERABILITY   SEVERITY 
busybox        1.36.1-r15            apk   CVE-2023-42366  Medium    
busybox        1.36.1-r15            apk   CVE-2023-42365  Medium    
busybox        1.36.1-r15            apk   CVE-2023-42364  Medium    
busybox        1.36.1-r15            apk   CVE-2023-42363  Medium    
busybox-binsh  1.36.1-r15            apk   CVE-2023-42366  Medium    
busybox-binsh  1.36.1-r15            apk   CVE-2023-42365  Medium    
busybox-binsh  1.36.1-r15            apk   CVE-2023-42364  Medium    
busybox-binsh  1.36.1-r15            apk   CVE-2023-42363  Medium    
curl           8.5.0-r0              apk   CVE-2024-0853   Medium    
curl           8.5.0-r0              apk   CVE-2024-2466   Unknown   
curl           8.5.0-r0              apk   CVE-2024-2398   Unknown   
curl           8.5.0-r0              apk   CVE-2024-2004   Unknown   
libuv          1.47.0-r0             apk   CVE-2024-24806  High      
ssl_client     1.36.1-r15            apk   CVE-2023-42366  Medium    
ssl_client     1.36.1-r15            apk   CVE-2023-42365  Medium    
ssl_client     1.36.1-r15            apk   CVE-2023-42364  Medium    
ssl_client     1.36.1-r15            apk   CVE-2023-42363  Medium
```

### weaveexec: (17) 

```
NAME           INSTALLED   FIXED-IN  TYPE  VULNERABILITY   SEVERITY 
busybox        1.36.1-r15            apk   CVE-2023-42366  Medium    
busybox        1.36.1-r15            apk   CVE-2023-42365  Medium    
busybox        1.36.1-r15            apk   CVE-2023-42364  Medium    
busybox        1.36.1-r15            apk   CVE-2023-42363  Medium    
busybox-binsh  1.36.1-r15            apk   CVE-2023-42366  Medium    
busybox-binsh  1.36.1-r15            apk   CVE-2023-42365  Medium    
busybox-binsh  1.36.1-r15            apk   CVE-2023-42364  Medium    
busybox-binsh  1.36.1-r15            apk   CVE-2023-42363  Medium    
curl           8.5.0-r0              apk   CVE-2024-0853   Medium    
curl           8.5.0-r0              apk   CVE-2024-2466   Unknown   
curl           8.5.0-r0              apk   CVE-2024-2398   Unknown   
curl           8.5.0-r0              apk   CVE-2024-2004   Unknown   
libuv          1.47.0-r0             apk   CVE-2024-24806  High      
ssl_client     1.36.1-r15            apk   CVE-2023-42366  Medium    
ssl_client     1.36.1-r15            apk   CVE-2023-42365  Medium    
ssl_client     1.36.1-r15            apk   CVE-2023-42364  Medium    
ssl_client     1.36.1-r15            apk   CVE-2023-42363  Medium
```

