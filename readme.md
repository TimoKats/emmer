# 🪣 emmer
[![Go Report Card](https://goreportcard.com/badge/github.com/TimoKats/emmer)](https://goreportcard.com/report/github.com/TimoKats/emmer)
[![Tests](https://github.com/TimoKats/emmer/actions/workflows/test.yaml/badge.svg)](https://github.com/TimoKats/emmer/actions/workflows/test.yaml)
[![License: EUPL](https://img.shields.io/badge/license-EUPL-blue.svg)](https://joinup.ec.europa.eu/collection/eupl/eupl-text-eupl-12)
<img width="88" height="31" alt="button(2)" style="float:right" src="https://github.com/user-attachments/assets/197197c9-213f-4246-888f-4e7e34f800aa" />


Self-hosted API for CRUD-ing JSON files on different filesystems (local, S3, azure blob, ...). Built for data storage in small personal projects, or mocking an API for development. Advantages are simplicity, interoperabillity and performance (using the cache system).   

```console
foo@bar:~$ go install github.com/TimoKats/emmer@latest
foo@bar:~$ emmer
2025/09/10 16:41:20 set username to: admin
2025/09/10 16:41:20 set password to: ************
2025/09/10 16:41:20 selected local fs in: /home/user/.emmer
2025/09/10 16:41:20 server is running on http://localhost:8080
```

The API is based on your JSON structure. So the example below is for CRUD-ing `[geralt][city]` in `file.json`. The value (which can be anything) is then added to the body of the request. Moreover, there are helper functions for appending and incrementing values. Refer to the wiki for detailed documentation and examples. 

```
DELETE/PUT/GET: /api/users/geralt/city
```

<img width="100%" alt="emmer-improved" src="https://github.com/user-attachments/assets/24ded58d-33ba-48dc-906f-d285630001eb" />

