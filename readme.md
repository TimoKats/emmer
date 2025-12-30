# 🪣 emmer
[![Go Report Card](https://goreportcard.com/badge/github.com/TimoKats/emmer)](https://goreportcard.com/report/github.com/TimoKats/emmer)
[![Tests](https://github.com/TimoKats/emmer/actions/workflows/test.yaml/badge.svg)](https://github.com/TimoKats/emmer/actions/workflows/test.yaml)
[![License: EUPL](https://img.shields.io/badge/license-EUPL-blue.svg)](https://joinup.ec.europa.eu/collection/eupl/eupl-text-eupl-12)

Self-hosted API for CRUD-ing JSON data on different filesystems (local, S3, azure blob, ...). Built for data storage in small personal projects, or mocking an API for development. Advantages are simplicity, interoperabillity (JSON files, APIs, multiple filesystems) and performance (using the cache system).   

You can install emmer with the command below. After that, you can run `emmer` to start.

```console
foo@bar:~$ go install github.com/TimoKats/emmer@latest
foo@bar:~$ emmer
2025/09/10 16:41:20 set username to: admin
2025/09/10 16:41:20 set password to: ************
2025/09/10 16:41:20 selected local fs in: /home/user/.emmer
2025/09/10 16:41:20 server is running on http://localhost:8080
```

The API is based on your JSON structure. So the example below is for CRUD-ing `[key1][key2]` in `file.json`. The value (which can be anything) is then added to the body of the request. Moreover, there are helper functions for appending and incrementing values. 

```
DELETE/PUT/GET: /api/file/key1/key2/...
```

Note, please refer to the wiki for detailed documentation and examples.
