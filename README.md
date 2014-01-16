#Fxh.Go

A fast and simple blog engine with GoInk framework in Golang.

Current version is **0.1.1-beta**

### Overview

Fxh.Go is a dynamic blog engine written in Golang. It's fast and very simple configs. Fxh.Go persists data into pieces of json files and support compress them as backup zip for next upgrade or installation.

Fxh.Go supports markdown contents as articles or pages, ajax comments and dynamic administration.

Fxh.Go contains two kinds of content as article and page. They can be customized as you want.

Documentation is writing.

### Installation

Fxh.Go is written in Golang with support for Windows, Linux and Mac OSX.

The stable release can use `go get` to install:

    go get github.com/fuxiaohei/GoBlog

remember set `$GOPATH/bin` to global environment variables.

The newest sources is in branch develop, you can download and build in manual (not-recommended).

### Setup

If installed, `GoBlog` binary file is built in `$GOPATH/bin`.

make a new dir to install Fxh.Go:

    cd empty_dir
    Goblog

then it will unzip static files in `empty_dir` , initialize original data and start server at `localhost:9000`

#### Administration

visit `localhost:9000/login/` to enter administrator with username `admin` and password `admin`. You'd better change them after installed successfully.

### Suggestion and Contribution

create issues or pull requests here.

### Thanks

gladly thank for [@Unknwon](https://github.com/Unknwon) on testing and [zip library](https://github.com/Unknwon/cae) support.

### License

The MIT License

