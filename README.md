# libhosty

[![made-with-Go](https://img.shields.io/badge/made%20with-Go-1f425f.svg)](http://golang.org)
[![Go Report Card](https://goreportcard.com/badge/github.com/areYouLazy/libhosty)](https://goreportcard.com/report/github.com/areYouLazy/libhosty)
[![Build and Test](https://github.com/areYouLazy/libhosty/actions/workflows/build-and-test.yml/badge.svg?branch=main&event=push)](https://github.com/areYouLazy/libhosty/actions/workflows/build-and-test.yml)
![gopherbadger-tag-do-not-edit](coverage_badge.png)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/areYouLazy/libhosty)

## Description

libhosty is a pure golang library to manipulate hosts-like files. It is inspired by [txeh](https://github.com/txn2/txeh), with some enrichments.

## Table of Contents

- [Description](#description)
- [Table of Contents](#table-of-contents)
- [Main Features](#main-features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
  - [Templates](#templates)
- [Credits](#credits)
- [License](#license)

## Main Features

* Comment/Uncomment a line without removing it from the file
* Restore the default hosts file for your system
* Add/Remove Address lines
* Add/Remove Comment lines
* Add/Remove Empty lines
* Query by hostname
* Automatically handles duplicate entries

## Requirements

Requires go >= 1.21

## Installation

Ensure you have go on your system

```bash
> go version
go version go1.22.5 linux/amd64
```

and pull the library (note that we are on v2 now)

```bash
> go get github.com/areYouLazy/libhosty/v2
```

## Usage

To use the library, just import it and call the `Init()` method.

To load a custom hosts file use the `InitFromCustomPath(path string)` routine

To parse an inline hosts file use the `InitFromString(lines string)` routine

> Note: This code doesn't handle errors for readability purposes, but you SHOULD!

```go
package main

import "github.com/areYouLazy/libhosty"

func main() {
    //you can initialize libhosty with a custom path
    //
    // cnf, _ := libhosty.InitFromCustomPath("/path/to/my/custom/hosts/file")
    
    // or you can parse an inline hosts file
    //
    // hfile := `# Example of an inline hosts file
    // 127.0.0.1  localhost
    // ::1        localhost
    // 1.1.1.1    cloudflare.dns
    // 8.8.8.8    google.dns`
    // cnf, _ := libhosty.InitFromString(hfile)
    
    // load hosts file from default OS location
    hfl, _ := libhosty.Init()
    
    //add an empty line
    hfl.AddEmptyFileLine()
    
    //add a host with a comment
    hfl.AddHostFileLine("12.12.12.12", "my.host.name", "comment on my hostname!")
    
    //add a comment
    hfl.AddCommentFileLine("just a comment")
    
    //add an empty line
    hfl.AddEmptyFileLine()
    
    //add another host without comment
    hfl.AddHostsFileLine("13.13.13.13", "another.host.name", "")
    
    //add another fqdn to the previous ip
    hfl.AddHostsFileLine("12.12.12.12", "second.host.name", "")
    
    // comment for host lines can be done by hostname, row line
    // or IP (as net.IP or string)
    //
    // Comment the line with address 12.12.12.12
    //
    // By-Row-Number
    hfl.CommentHostsFileLineByRow(idx)
    //
    // By-Hostname
    hfl.CommentHostsFileLinesByHostname("second.host.name")
    //
    // By-Address-As-IP
    ip := net.ParseIP("12.12.12.12")
    hfl.CommentHostsFileLinesByIP(ip)
    //
    // By-Address-As-String
    hfl.CommentHostsFileLinesByAddress("12.12.12.12")
    
    // render the hosts file
    fmt.Println(hfl.RenderHostsFile())
    
    // write file to disk
    hfl.WriteHostsFile()
    
    // or to a custom location
    hfl.WriteHostsFileTo("/home/sonica/hosts-export.txt")
    
    // restore the original hosts file based on running OS
    hfl.RestoreTemplate()
    
    // render the hosts file
    fmt.Println(hfl.RenderHostsFile())
    
    // write to disk
    hfl.WriteHostsFile()
}
```

The 1st `fmt.Println()` should output something like this (in a linux host)

```console
# Do not remove the following line, or various programs
# that require network functionality will fail.
127.0.0.1               localhost.localdomain localhost
::1                     localhost6.localdomain6 localhost6

# 12.12.12.12           my.host.name second.host.name   #comment on my hostname!
# just a comment line

13.13.13.13             another.host.name
```

While the 2nd `fmt.Println()` should output the default template for your OS

```console
# Do not remove the following line, or various programs
# that require network functionality will fail.
127.0.0.1               localhost.localdomain localhost
::1                     localhost6.localdomain6 localhost6

```

If you handle errors properly, you'll notice that this example program will fail on the `SaveHostsFile()` call if started as a normal user, as editing the hosts file requires root privileges. This does not prevent libhosty from loading, managing, rendering and exporting the hosts file

## Contributing

Issues and PRs are more than welcome!

### Templates

If you find a hosts template (like the Docker one) that you think can be useful to have in this library feel free to open an Issue/Pull Request

## Credits

Project Contributors will be listed here

## License

Licenses under Apache License 2.0
