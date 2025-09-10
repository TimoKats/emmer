# 🪣 emmer
[![Go Report Card](https://goreportcard.com/badge/github.com/TimoKats/emmer)](https://goreportcard.com/report/github.com/TimoKats/emmer)
[![Tests](https://github.com/TimoKats/emmer/actions/workflows/test.yaml/badge.svg)](https://github.com/TimoKats/emmer/actions/workflows/test.yaml)
[![License: EUPL](https://img.shields.io/badge/license-EUPL-blue.svg)](https://joinup.ec.europa.eu/collection/eupl/eupl-text-eupl-12)

An API for creating and querying JSON data on different filesystems (local, S3, azure blob, ...). You can install emmer with the command below. After that, you can run `emmer` to start adding/querying data using the API.

```console
foo@bar:~$ go install github.com/TimoKats/emmer@latest
foo@bar:~$ emmer
2025/09/10 16:41:20 set username to: admin
2025/09/10 16:41:20 set password to: ************
2025/09/10 16:41:20 selected local fs in: /home/user/.emmer
2025/09/10 16:41:20 server is running on http://localhost:8080
```

**The complete documentation can be found on [my website](https://timokats.xyz/pages/emmer.php).**
