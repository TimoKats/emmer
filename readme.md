# 🪣 emmer
[![Go Report Card](https://goreportcard.com/badge/github.com/TimoKats/emmer)](https://goreportcard.com/report/github.com/TimoKats/emmer)
[![Tests](https://github.com/TimoKats/emmer/actions/workflows/test.yaml/badge.svg)](https://github.com/TimoKats/emmer/actions/workflows/test.yaml)
[![License: EUPL](https://img.shields.io/badge/license-EUPL-blue.svg)](https://joinup.ec.europa.eu/collection/eupl/eupl-text-eupl-12)

An API for creating and querying JSON data. You can install emmer with the command below. After that, you can run emmer and the server will run on 8080 (unless other is specified). By default, the server will generate a password (see logs) and use the local filesystem at `~/.emmer`.

&nbsp;

```console
foo@bar:~$ go install github.com/TimoKats/emmer@latest
foo@bar:~$ emmer -p 1234
```

&nbsp;

**An overview of all settings and API endpoints can be found in the documentation on my website.**
