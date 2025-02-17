# vcpkg-cache-http
[![build](https://github.com/lesomnus/vcpkg-cache-http/actions/workflows/build.yaml/badge.svg)](https://github.com/lesomnus/vcpkg-cache-http/actions/workflows/build.yaml)
[![test](https://github.com/lesomnus/vcpkg-cache-http/actions/workflows/test.yaml/badge.svg)](https://github.com/lesomnus/vcpkg-cache-http/actions/workflows/test.yaml)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/8535edf4b4af4238a2140994efe99d55)](https://app.codacy.com/gh/lesomnus/vcpkg-cache-http/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/8535edf4b4af4238a2140994efe99d55)](https://app.codacy.com/gh/lesomnus/vcpkg-cache-http/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

HTTP provider for [*vcpkg*](https://github.com/microsoft/vcpkg) binary caching.
```sh
$ vcpkg-cache-http
2023-07-12T17:37:14Z INF use default store store=file:vcpkg-cache
2023-07-12T17:37:14Z INF start server addr=0.0.0.0:15151
2023-07-12T17:39:04Z INF _=nAq80yfCqKVl method=GET url=/zlib/1.2.13/70a5ceda64f1b5c01c1f7afe7669a32bc11c11496d8aeb094d7389a43c946f4b
2023-07-12T17:39:04Z INF REQ GET _=nAq80yfCqKVl hash=70a5ceda64f1b5c01c1f7afe7669a32bc11c11496d8aeb094d7389a43c946f4b name=zlib version=1.2.13
2023-07-12T17:39:04Z WRN RES GET _=nAq80yfCqKVl dt=0.098659 status=404
2023-07-12T17:39:05Z INF _=6jFojxVafWyU method=PUT url=/zlib/1.2.13/70a5ceda64f1b5c01c1f7afe7669a32bc11c11496d8aeb094d7389a43c946f4b
2023-07-12T17:39:05Z INF REQ PUT _=6jFojxVafWyU hash=70a5ceda64f1b5c01c1f7afe7669a32bc11c11496d8aeb094d7389a43c946f4b name=zlib version=1.2.13
2023-07-12T17:39:05Z INF RES PUT _=6jFojxVafWyU dt=0.839125 status=200
2023-07-12T17:41:12Z INF _=JAjAKcKtGmfl method=GET url=/zlib/1.2.13/70a5ceda64f1b5c01c1f7afe7669a32bc11c11496d8aeb094d7389a43c946f4b
2023-07-12T17:41:12Z INF REQ GET _=JAjAKcKtGmfl hash=70a5ceda64f1b5c01c1f7afe7669a32bc11c11496d8aeb094d7389a43c946f4b name=zlib version=1.2.13
2023-07-12T17:41:12Z INF RES GET _=JAjAKcKtGmfl dt=0.467711 status=200
```

## Usage

Just start the server.
By default, it listens on port 15151 and creates a directory named `vcpkg-cache` to store the binary cache.

```sh
$ vcpkg-cache-http                                        
2023-07-12T17:37:14Z INF use default store store=file:vcpkg-cache
2023-07-12T17:37:14Z INF start server addr=0.0.0.0:15151
```

Set *vcpkg* binary source as `http,http://localhost:15151/{triplet}/{name}/{version}/{sha},readwrite`.
It can be set to environment variable `VCPKG_BINARY_SOURCES` or passed by *vcpkg* command with `--binarysource` flag.
Please see *vcpkg* official document about [Binary Caching](https://learn.microsoft.com/en-us/vcpkg/users/binarycaching) for details.

```sh
$ vcpkg install --binarysource="http,http://localhost:15151/{triplet}/{name}/{version}/{sha},readwrite" zlib
Computing installation plan...
The following packages will be built and installed:
    zlib:x64-linux -> 1.2.13
Detecting compiler hash for triplet x64-linux...
Restored 0 package(s) from /home/hypnos/.cache/vcpkg/archives in 7 us. Use --debug to see more details.
Restored 1 package(s) from HTTP servers in 12.9 ms. Use --debug to see more details.
Installing 1/1 zlib:x64-linux...
Elapsed time to handle zlib:x64-linux: 769 us
Total install time: 774 us
```

Note that `zlib` is cached on the server; message indicating that `zlib` has been restored from the HTTP server.

## Install

### From Source
```sh
go install github.com/lesomnus/vcpkg-cache-http@latest
```

### Docker
```sh
docker pull ghcr.io/lesomnus/vcpkg-cache-http:edge 
```

Currently, only `edge` tag built from the `main` branch is provided.

## Stores

`vcpkg-cache-http` accepts 1 positional argument that indicates the location to store the binary cache. The store is formed as `kind[:[path][,opt[=val]]]`.
For example, store that stores the binary cache to a directory where `./path/to/store` is as follow:

```
files:./path/to/store
```

Available stores are:

- `files:[vcpkg-cache]`
  
    Stores to a directory at the given path.

- `archives:[${HOME}/.cache/vcpkg/archives]`

    Use *vcpkg*'s `files` provider at the given path as a store.
