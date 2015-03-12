# tmplnator

[![Build Status](https://travis-ci.org/albertrdixon/tmplnator.svg?branch=master)](https://travis-ci.org/albertrdixon/tmplnator)

~~A killing machine sent from the future~~ 
Yet another file generator using golang [text/template](http://golang.org/pkg/text/template/).

## Description

I wanted something were I could put all my application configuration templates in a single folder within my [docker](https://www.docker.com/) images, then run a single command to have them processed and placed without a ton of extra configuration.

Enter tmplnator. Templates describe where they should go, what their file mode should be and what user and group should own them. I'm sure this will not be for everyone, like I said, I made this for me.

## Install

In your Dockerfile do something like:

```
RUN curl -#kL https://github.com/albertrdixon/tmplnator/releases/download/<version>/tnator-linux-amd64-<version>.tar.gz |\
    tar xvz -C /usr/local/bin
```

## Usage

Help Menu:

```
Usage of ./t2:
  -bpool-size="": Size of write buffer pool
  -default-dir="": Default output directory
  -delete="": Remove templates after processing
  -etcd-peers=: etcd peers in host:port (can be provided multiple times)
  -namespace="": etcd key namespace
  -template-dir="": Template directory
  -threads="": Number of processing threads
  -v="": Verbosity (0:quiet output, 1:default, 2:debug output)
  -version="": show version
```

Super simple. Use the following methods in the template to set up how the file should be generated from the template: `dir` `mode` `user` `group`

```
# example supervisor.conf template
{{ dir "/etc/supervisor" }}
{{ mode 0644 }}
{{ user 0 }}
{{ group 0 }}
[supervisord]
nodaemon  = true
...
```

Add your templates to some directory: `ADD configs /templates`

Run tmplnator like so: `t2 -template-dir /templates`

And that's it!

## Template Functions

Access environment variables in the template with `.Env` like so `.Env.VARIABLE`

Access etcd values with `.Var <key>` if key not found will look in ENV

`dir "/path/to/destination/dir"`: Describe destination directory

`mode <file_mode>`: Describe filemode for generated file

`user <uid>`: Describe uid for generated file

`group <gid>`: Describe gid for generated file

`to_json <input>`: Marshal JSON string

`from_json <string>`: Unmarshal JSON string

`from_json_array <string>`: Unmarshal JSON array

`first <slice>`: Return the first element of the slice

`last <slice>`: Return the last element of the slice

`file_exists <filename.`: True if file exists

`parseURL <string>`: Return url.URL object of given url string

`has_key <map> <key>`: True if key exists in map

`default <value> <default_value>`: Output default_value if value is nil or empty string, otherwise output value

`split <string>`: strings.Split

`join <slice>`: strings.Join

`has_suffix <string> <suffix>`: strings.HasSuffix

`contains <string> <substring>`: strings.Contains

`fields <string>`: strings.Fields

`downcase <string>`: strings.ToLower

`upcase <string>`: strings.ToUpper

`trim_suffix <string> <suffix>`: strings.TrimSuffix

`trim_space <string>`: strings.TrimSpace
