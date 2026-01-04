# 🪣 emmer 

[![Go Report Card](https://goreportcard.com/badge/github.com/TimoKats/emmer)](https://goreportcard.com/report/github.com/TimoKats/emmer)
[![Tests](https://github.com/TimoKats/emmer/actions/workflows/test.yaml/badge.svg)](https://github.com/TimoKats/emmer/actions/workflows/test.yaml)
[![Docker image](https://badgen.net/badge/icon/docker?icon=docker&label)](https://hub.docker.com/r/tiemster/emmer)



Self-hosted API for CRUD-ing JSON files on different filesystems (local, S3, MinIO, ...). Built for data storage in small personal projects, or mocking an API for development. Advantages are simplicity, interoperabillity and performance (using the cache system).   

```console
foo@bar:~$ go install github.com/TimoKats/emmer@latest
foo@bar:~$ emmer
2025/09/10 16:41:20 set username to: admin
2025/09/10 16:41:20 set password to: ************
2025/09/10 16:41:20 selected local fs in: /home/user/.emmer
2025/09/10 16:41:20 server is running on http://localhost:8080
```

In summary, he API is based on your JSON structure. So the example below is for CRUD-ing `[geralt][city]` in `file.json`. The value (which can be anything) is then added to the body of the request. Moreover, there are helper functions for appending and incrementing values. Refer to the wiki for detailed documentation and examples. 

```
DELETE/PUT/GET: /api/users/geralt/city
```
&nbsp;  

<div align="center" >
  <!-- <img width="100%"  alt="emmer drawio" src="https://github.com/user-attachments/assets/9bd584e1-fa81-432a-9455-11b3bf6b9e2a" /> -->
<img width="100%" alt="big-emmer-1 drawio" src="https://github.com/user-attachments/assets/3162486d-d77d-4026-8853-fca6cecbd87c" />  
</div>

<!-- <img width="88" height="31" alt="button(1)" src="https://github.com/user-attachments/assets/cc0ebf4f-3033-487c-bed7-7184a1f2ec9a" /> -->
