#Fxh.Go

A fast and simple blog engine with [GoInk](https://github.com/fuxiaohei/GoInk) framework in Golang.

[![Build Status](https://drone.io/github.com/fuxiaohei/GoBlog/status.png)](https://drone.io/github.com/fuxiaohei/GoBlog/latest)
[![GoWalker](http://b.repl.ca/v1/Go_Walker-API_Documentation-green.png)](http://gowalker.org/github.com/fuxiaohei/GoBlog)

Current version is **0.2.5** on 2014.02.28

Development board is in [Trello](https://trello.com/b/7AHrcQL8/fxh-go-with-goink).

**Notice: the project is planning to rebuild with sqlite database. I make sure the data are compatible and all functionalities are same - 2014.10.20**

### Overview

`Fxh.Go` is a dynamic blog engine written in Golang. It's fast and very simple configs. Fxh.Go persists data into pieces of json files and support compress them as backup zip for next upgrade or installation.

`Fxh.Go` supports markdown contents as articles or pages, ajax comments and dynamic administration.

`Fxh.Go` contains two kinds of content as article and page. They can be customized as you want.

### Installation

`Fxh.Go` requires **Go 1.2** or above.

##### Gobuild.io

[Gobuild.io](http://gobuild.io/) can build cross-platform executable file for pure go projects. You can download `Fxh.Go` binary from Gobuild.io.

[![Gobuild Download](http://gobuild.io/badge/github.com/fuxiaohei/GoBlog/download.png)](http://gobuild.io/github.com/fuxiaohei/GoBlog)

##### Manual

Use go get command:

    go get github.com/fuxiaohei/GoBlog

Then you can find binary file `GoBlog(.exe)` in `$GOPATH/bin`.

### Run

Make a new dir to run `Fxh.Go`:

    cd new_dir
    Goblog

Then it will unzip static files in `new_dir` , initialize raw data and start server at `localhost:9001`.

##### Admin

Visit `localhost:9001/login/` to enter administrator with username `admin` and password `admin`. You'd better change them after installed successfully.

##### Deployment

I prefer to use nginx as proxy. The server section in `nginx.conf`:

        server {
                listen       80;
                server_name  your_domain;
                charset utf-8;
                access_log  /var/log/nginx/your_domain.access.log;

                location / {
                    proxy_pass http://127.0.0.1:9001;
                }

                location /static {
                    root            /var/www/your_domain;  # binary file is in this directory
                    expires         1d;
                    add_header      Cache-Control public;
                    access_log      off;
                }
        }

### Questions

Create issues or pull requests here.

### Products

* [抛弃世俗之浮躁，留我钻研之刻苦](http://wuwen.org)
* [FuXiaoHei.Me](http://fuxiaohei.me)

### Thanks

* [@Unknwon](https://github.com/Unknwon) on testing and [zip library](https://github.com/Unknwon/cae) support.

### License

The MIT License

