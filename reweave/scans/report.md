# Vulnerability Report

```
Report date: 2025-09-15
Unique vulnerability count: 20
Images version: 2.10.0
```

## Scanner Details

```
Application:         grype
Version:             0.99.1
BuildDate:           2025-08-30T14:21:51Z
GitCommit:           6e1807dd7130c9eff23c1bad19a96cc6afef33aa
GitDescription:      v0.99.1
Platform:            linux/amd64
GoVersion:           go1.24.6
Compiler:            gc
Syft Version:        v1.32.0
Supported DB Schema: 6
```

## Vulnerabilities

### weave-kube: (19) 

```
NAME                              INSTALLED             FIXED IN                       TYPE       VULNERABILITY        SEVERITY  EPSS           RISK   
golang.org/x/crypto               v0.31.0               0.35.0                         go-module  GHSA-hcg3-q754-cr77  High      0.2% (40th)    0.1    
golang.org/x/oauth2               v0.10.0               0.27.0                         go-module  GHSA-6v2p-p543-phr9  High      < 0.1% (23rd)  < 0.1  
stdlib                            go1.23.1              *1.23.12, 1.24.6               go-module  CVE-2025-47907       High      < 0.1% (18th)  < 0.1  
stdlib                            go1.23.1              1.22.11, *1.23.5, 1.24.0-rc.2  go-module  CVE-2024-45336       Medium    < 0.1% (13th)  < 0.1  
stdlib                            go1.23.1              1.22.11, *1.23.5, 1.24.0-rc.2  go-module  CVE-2024-45341       Medium    < 0.1% (8th)   < 0.1  
stdlib                            go1.23.1              *1.23.8, 1.24.2                go-module  CVE-2025-22871       Critical  < 0.1% (1st)   < 0.1  
stdlib                            go1.23.1              *1.23.10, 1.24.4               go-module  CVE-2025-4673        Medium    < 0.1% (3rd)   < 0.1  
golang.org/x/net                  v0.33.0               0.38.0                         go-module  GHSA-vvgc-356p-c3xw  Medium    < 0.1% (4th)   < 0.1  
busybox                           1.37.0-r19                                           apk        CVE-2024-58251       Low       < 0.1% (7th)   < 0.1  
busybox-binsh                     1.37.0-r19                                           apk        CVE-2024-58251       Low       < 0.1% (7th)   < 0.1  
ssl_client                        1.37.0-r19                                           apk        CVE-2024-58251       Low       < 0.1% (7th)   < 0.1  
golang.org/x/net                  v0.33.0               0.36.0                         go-module  GHSA-qxp5-gwg8-xv66  Medium    < 0.1% (2nd)   < 0.1  
busybox                           1.37.0-r19                                           apk        CVE-2025-46394       Low       < 0.1% (4th)   < 0.1  
busybox-binsh                     1.37.0-r19                                           apk        CVE-2025-46394       Low       < 0.1% (4th)   < 0.1  
ssl_client                        1.37.0-r19                                           apk        CVE-2025-46394       Low       < 0.1% (4th)   < 0.1  
github.com/containerd/containerd  v1.7.11               1.7.27                         go-module  GHSA-265r-hfxg-fhmg  Medium    < 0.1% (1st)   < 0.1  
stdlib                            go1.23.1              *1.23.11, 1.24.5               go-module  CVE-2025-4674        High      < 0.1% (0th)   < 0.1  
stdlib                            go1.23.1              1.22.12, *1.23.6, 1.24.0-rc.3  go-module  CVE-2025-22866       Medium    < 0.1% (0th)   < 0.1  
github.com/docker/docker          v25.0.6+incompatible                                 go-module  GHSA-4vq8-7jfc-9cvp  Low       < 0.1% (1st)   < 0.1
```

### weave-npc: (17) 

```
NAME                 INSTALLED   FIXED IN                       TYPE       VULNERABILITY        SEVERITY  EPSS           RISK   
golang.org/x/crypto  v0.31.0     0.35.0                         go-module  GHSA-hcg3-q754-cr77  High      0.2% (40th)    0.1    
golang.org/x/oauth2  v0.10.0     0.27.0                         go-module  GHSA-6v2p-p543-phr9  High      < 0.1% (23rd)  < 0.1  
stdlib               go1.23.1    *1.23.12, 1.24.6               go-module  CVE-2025-47907       High      < 0.1% (18th)  < 0.1  
stdlib               go1.23.1    1.22.11, *1.23.5, 1.24.0-rc.2  go-module  CVE-2024-45336       Medium    < 0.1% (13th)  < 0.1  
stdlib               go1.23.1    1.22.11, *1.23.5, 1.24.0-rc.2  go-module  CVE-2024-45341       Medium    < 0.1% (8th)   < 0.1  
stdlib               go1.23.1    *1.23.8, 1.24.2                go-module  CVE-2025-22871       Critical  < 0.1% (1st)   < 0.1  
stdlib               go1.23.1    *1.23.10, 1.24.4               go-module  CVE-2025-4673        Medium    < 0.1% (3rd)   < 0.1  
golang.org/x/net     v0.33.0     0.38.0                         go-module  GHSA-vvgc-356p-c3xw  Medium    < 0.1% (4th)   < 0.1  
busybox              1.37.0-r19                                 apk        CVE-2024-58251       Low       < 0.1% (7th)   < 0.1  
busybox-binsh        1.37.0-r19                                 apk        CVE-2024-58251       Low       < 0.1% (7th)   < 0.1  
ssl_client           1.37.0-r19                                 apk        CVE-2024-58251       Low       < 0.1% (7th)   < 0.1  
golang.org/x/net     v0.33.0     0.36.0                         go-module  GHSA-qxp5-gwg8-xv66  Medium    < 0.1% (2nd)   < 0.1  
busybox              1.37.0-r19                                 apk        CVE-2025-46394       Low       < 0.1% (4th)   < 0.1  
busybox-binsh        1.37.0-r19                                 apk        CVE-2025-46394       Low       < 0.1% (4th)   < 0.1  
ssl_client           1.37.0-r19                                 apk        CVE-2025-46394       Low       < 0.1% (4th)   < 0.1  
stdlib               go1.23.1    *1.23.11, 1.24.5               go-module  CVE-2025-4674        High      < 0.1% (0th)   < 0.1  
stdlib               go1.23.1    1.22.12, *1.23.6, 1.24.0-rc.3  go-module  CVE-2025-22866       Medium    < 0.1% (0th)   < 0.1
```

### weave: (18) 

```
NAME                              INSTALLED             FIXED IN                       TYPE       VULNERABILITY        SEVERITY  EPSS           RISK   
golang.org/x/crypto               v0.31.0               0.35.0                         go-module  GHSA-hcg3-q754-cr77  High      0.2% (40th)    0.1    
stdlib                            go1.23.1              *1.23.12, 1.24.6               go-module  CVE-2025-47907       High      < 0.1% (18th)  < 0.1  
stdlib                            go1.23.1              1.22.11, *1.23.5, 1.24.0-rc.2  go-module  CVE-2024-45336       Medium    < 0.1% (13th)  < 0.1  
stdlib                            go1.23.1              1.22.11, *1.23.5, 1.24.0-rc.2  go-module  CVE-2024-45341       Medium    < 0.1% (8th)   < 0.1  
stdlib                            go1.23.1              *1.23.8, 1.24.2                go-module  CVE-2025-22871       Critical  < 0.1% (1st)   < 0.1  
stdlib                            go1.23.1              *1.23.10, 1.24.4               go-module  CVE-2025-4673        Medium    < 0.1% (3rd)   < 0.1  
golang.org/x/net                  v0.33.0               0.38.0                         go-module  GHSA-vvgc-356p-c3xw  Medium    < 0.1% (4th)   < 0.1  
busybox                           1.37.0-r19                                           apk        CVE-2024-58251       Low       < 0.1% (7th)   < 0.1  
busybox-binsh                     1.37.0-r19                                           apk        CVE-2024-58251       Low       < 0.1% (7th)   < 0.1  
ssl_client                        1.37.0-r19                                           apk        CVE-2024-58251       Low       < 0.1% (7th)   < 0.1  
golang.org/x/net                  v0.33.0               0.36.0                         go-module  GHSA-qxp5-gwg8-xv66  Medium    < 0.1% (2nd)   < 0.1  
busybox                           1.37.0-r19                                           apk        CVE-2025-46394       Low       < 0.1% (4th)   < 0.1  
busybox-binsh                     1.37.0-r19                                           apk        CVE-2025-46394       Low       < 0.1% (4th)   < 0.1  
ssl_client                        1.37.0-r19                                           apk        CVE-2025-46394       Low       < 0.1% (4th)   < 0.1  
github.com/containerd/containerd  v1.7.11               1.7.27                         go-module  GHSA-265r-hfxg-fhmg  Medium    < 0.1% (1st)   < 0.1  
stdlib                            go1.23.1              *1.23.11, 1.24.5               go-module  CVE-2025-4674        High      < 0.1% (0th)   < 0.1  
stdlib                            go1.23.1              1.22.12, *1.23.6, 1.24.0-rc.3  go-module  CVE-2025-22866       Medium    < 0.1% (0th)   < 0.1  
github.com/docker/docker          v25.0.6+incompatible                                 go-module  GHSA-4vq8-7jfc-9cvp  Low       < 0.1% (1st)   < 0.1
```

### weaveexec: (18) 

```
NAME                              INSTALLED             FIXED IN                       TYPE       VULNERABILITY        SEVERITY  EPSS           RISK   
golang.org/x/crypto               v0.31.0               0.35.0                         go-module  GHSA-hcg3-q754-cr77  High      0.2% (40th)    0.1    
stdlib                            go1.23.1              *1.23.12, 1.24.6               go-module  CVE-2025-47907       High      < 0.1% (18th)  < 0.1  
stdlib                            go1.23.1              1.22.11, *1.23.5, 1.24.0-rc.2  go-module  CVE-2024-45336       Medium    < 0.1% (13th)  < 0.1  
stdlib                            go1.23.1              1.22.11, *1.23.5, 1.24.0-rc.2  go-module  CVE-2024-45341       Medium    < 0.1% (8th)   < 0.1  
stdlib                            go1.23.1              *1.23.8, 1.24.2                go-module  CVE-2025-22871       Critical  < 0.1% (1st)   < 0.1  
stdlib                            go1.23.1              *1.23.10, 1.24.4               go-module  CVE-2025-4673        Medium    < 0.1% (3rd)   < 0.1  
golang.org/x/net                  v0.33.0               0.38.0                         go-module  GHSA-vvgc-356p-c3xw  Medium    < 0.1% (4th)   < 0.1  
busybox                           1.37.0-r19                                           apk        CVE-2024-58251       Low       < 0.1% (7th)   < 0.1  
busybox-binsh                     1.37.0-r19                                           apk        CVE-2024-58251       Low       < 0.1% (7th)   < 0.1  
ssl_client                        1.37.0-r19                                           apk        CVE-2024-58251       Low       < 0.1% (7th)   < 0.1  
golang.org/x/net                  v0.33.0               0.36.0                         go-module  GHSA-qxp5-gwg8-xv66  Medium    < 0.1% (2nd)   < 0.1  
busybox                           1.37.0-r19                                           apk        CVE-2025-46394       Low       < 0.1% (4th)   < 0.1  
busybox-binsh                     1.37.0-r19                                           apk        CVE-2025-46394       Low       < 0.1% (4th)   < 0.1  
ssl_client                        1.37.0-r19                                           apk        CVE-2025-46394       Low       < 0.1% (4th)   < 0.1  
github.com/containerd/containerd  v1.7.11               1.7.27                         go-module  GHSA-265r-hfxg-fhmg  Medium    < 0.1% (1st)   < 0.1  
stdlib                            go1.23.1              *1.23.11, 1.24.5               go-module  CVE-2025-4674        High      < 0.1% (0th)   < 0.1  
stdlib                            go1.23.1              1.22.12, *1.23.6, 1.24.0-rc.3  go-module  CVE-2025-22866       Medium    < 0.1% (0th)   < 0.1  
github.com/docker/docker          v25.0.6+incompatible                                 go-module  GHSA-4vq8-7jfc-9cvp  Low       < 0.1% (1st)   < 0.1
```

