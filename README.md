# tmplnator

[![Build Status](https://travis-ci.org/albertrdixon/tmplnator.svg?branch=master)](https://travis-ci.org/albertrdixon/tmplnator) [![GoDoc](https://godoc.org/github.com/albertrdixon/tmplnator?status.svg)](https://godoc.org/github.com/albertrdixon/tmplnator)

[![GoDoc](https://godoc.org/github.com/albertrdixon/tmplnator?status.svg)](https://godoc.org/github.com/albertrdixon/tmplnator)

~~A killing machine sent from the future~~ 
Yet another file generator using golang [text/template](http://golang.org/pkg/text/template/).

## Description

I wanted something were I could put all my application configuration templates in a single folder within my [docker](https://www.docker.com/) images, then run a single command to have them processed and placed without a ton of extra configuration.

Enter tmplnator. Templates describe where they should go, what their file mode should be and what user and group should own them. I'm sure this will not be for everyone, like I said, I made this for me.

## Install

In your Dockerfile do something like:

```
RUN curl -#kL https://github.com/albertrdixon/tmplnator/releases/download/<version>/tnator-linux-amd64.tar.gz |\
    tar xvz -C /usr/local/bin
```

## Usage

Help Menu:

```
usage: t2 [<flags>] <template-directory>

Flags:
  --help           Show help.
  -d, --debug      Enable debug mode
  -q, --quiet      Enable quiet mode
  -p, --print-tmp  Print out TmpDir
  -F, --force-tmp  Force all generated files to be written to TmpDir
  -T, --tmpdir="$TMPDIR/T2"
                   Set TmpDir
  --version        Show application version.

Args:
  <template-directory>
    Directory under which there are templates
```

Super simple. Use the following methods in the template to set up how the file should be generated from the template: `dir` `mode` `user` `group`

```
# example supervisor.conf template
{{ file "/etc/supervisor/supervisor.conf" }}
{{ mode 0644 }}
[supervisord]
nodaemon  = true
...
```

Add your templates to some directory: `ADD configs /templates`

Run tmplnator like so: `t2 /templates`

And that's it!

## Template Functions

Access environment variables in the template with `.Env` like so `.Env.VARIABLE`

Access etcd values with `.Get <key>` if key not found will look in ENV

`path "/path/to/destination/dir" <args...>`: Describe destination directory. Accepts printf style formatting in path string. **NOTE**: Templates without a described `dir` will use `default-dir` as their output directory.

`file "name" <args...>`: Describe name of generated file. Will set path elements to the path. Accepts printf style formatting of name string.

`mode <file_mode>`: Describe filemode for generated file

`dir_mode <file_mode>`: Describe mode for any generated directories

`to_json <input>`: Marshal JSON string

`from_json <string>`: Unmarshal JSON string

`from_json_array <string>`: Unmarshal JSON array

`first <slice>`: Return the first element of the slice

`last <slice>`: Return the last element of the slice

`file_exists <filename>`: True if file exists

`parseURL <string>`: Return [url.URL](https://golang.org/pkg/net/url/#URL) object of given url string

`has_key <map> <key>`: True if key exists in map

`default <value> <default_value>`: Output default_value if value is nil or empty string, otherwise output value

`fmt <format> <args...>`: fmt.Sprintf

`split <string> <separator>`: strings.Split

`join <slice> <separator>`: strings.Join

`has_suffix <string> <suffix>`: strings.HasSuffix

`contains <string> <substring>`: strings.Contains

`fields <string>`: strings.Fields

`downcase <string>`: strings.ToLower

`upcase <string>`: strings.ToUpper

`trim_suffix <string> <suffix>`: strings.TrimSuffix

`trim_space <string>`: strings.TrimSpace

## Shout Outs

My projects steals heavily from the wonderful projects below:

* [jwilder's dockerize](https://github.com/jwilder/dockerize)
* [kelseyhightower's confd](https://github.com/kelseyhightower/confd)
* [jwilder's docker-gen](https://github.com/jwilder/docker-gen)
