[![Build Status](http://192.168.1.108:8080/api/badges/cpp/demo/status.svg)](http://192.168.1.108:8080/cpp/demo)

drone是一个基于容器的本地持续交付平台，和Jenkins是差不多的，然后配合轻量级的gogs来作为git管理，都是基于golang开发的很符合我的需求，我们来把它们结合作为一个完整的CI、CD平台。

首先我们要先安装docker，上次的篇幅我们已经说过了我就不赘述了。

需要的东西有：linux，docker，docker-compose，drone，gogs。

## 安装gogs和drone

配合[荣锋亮大哥](https://www.cnblogs.com/rongfengliang/p/9963311.html)的yml文件和docker-compose我们可以很容易安装他们：

```yaml
version: '3'
services:
  drone-server:
    image: drone/drone:latest
    ports:
      - "8080:80"
      - 8843:443
      - 9000
    volumes:
      - ./drone:/var/lib/drone/
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - DRONE_OPEN=true
      - DRONE_SERVER_HOST=drone-server
      - DRONE_DEBUG=true
      - DRONE_GIT_ALWAYS_AUTH=false
      - DRONE_GOGS=true
      - DRONE_GOGS_SKIP_VERIFY=false
      - DRONE_GOGS_SERVER=http://gogs:3000
      - DRONE_PROVIDER=gogs
      - DRONE_DATABASE_DATASOURCE=/var/lib/drone/drone.sqlite
      - DRONE_DATABASE_DRIVER=sqlite3
      - DRONE_SERVER_PROTO=http
      - DRONE_RPC_SECRET=ALQU2M0KdptXUdTPKcEw
      - DRONE_SECRET=ALQU2M0KdptXUdTPKcEw
  gogs:
    image: gogs/gogs:latest
    ports:
      - "10022:22"
      - "3000:3000"
    volumes:
      - ./data/gogs:/data
    depends_on:
      - mysql
  mysql:
    image: mysql:5.7.16
    volumes:
      - ./gogs/mysql:/var/lib/mysql
      - /var/run/docker.sock:/var/run/docker.sock
    ports:
      - 3308:3306
    command: --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
    environment:
      MYSQL_ROOT_PASSWORD: pass
      MYSQL_DATABASE: gogs
      MYSQL_USER: gogs
      MYSQL_PASSWORD: pass
      TZ: Asia/Shanghai
  drone-agent:
    image: drone/agent:latest
    depends_on:
      - drone-server
    environment:
      - DRONE_RPC_SERVER=http://drone-server
      - DRONE_RPC_SECRET=ALQU2M0KdptXUdTPKcEw
      - DRONE_DEBUG=true
      - DOCKER_HOST=tcp://docker-bind:2375
      - DRONE_SERVER=drone-server:9000
      - DRONE_SECRET=ALQU2M0KdptXUdTPKcEw
      - DRONE_MAX_PROCS=5
  docker-bind:
     image: docker:dind
     privileged: true
    #  command: --storage-driver=overlay
```

我们创建一个存放docker-compose.yml文件的目录比如就叫gogs，然后我们把这些yml保存成docker-compose.yml，然后执行docker-compose来安装：

```bash
$ docker-compose up -d
```

配合yml文件，我们就安装好了drone-server和drone-agent还有gogs，然后我们用浏览器打开`http://localhost:3000/`来进入gogs并初始化它。

![gogs初始化1](https://i.loli.net/2019/03/02/5c7a4622758c4.png)

域名和应用URL记得一样接着我们创建一个管理员用户，然后其他的都默认，点击立即安装完成。

![gogs初始化2](https://i.loli.net/2019/03/02/5c7a46de1ddff.png)

初始化成功之后我们可以在gogs里边创建一个仓库，然后登陆drone。

# drone

打开浏览器输入`http://localhost/`直接进入drone，密码是gogs的你的刚刚的账户和密码。

![drone](https://i.loli.net/2019/03/02/5c7a495a6c7a1.png)

我们会看到一个刚刚创建的仓库，激活它！

激活之后，我们回到gogs那边，仓库的设置里边的webhook应该已经配置好了

![web hook](https://i.loli.net/2019/03/06/5c7f4e9198fb6.png)

我们可以测试web hook，如果没有问题的话，应该会提示成功。

## 上传源码

测试没有问题之后，我们初始化我们的代码文件夹为git仓库，然后push到gogos上边

![test仓库](https://i.loli.net/2019/03/06/5c7f4f152f320.png)

然后为你的仓库加上`.drone.yml`配置文件，drone-server会自动读取这个文件进行CI、CD操作等。以下这个是我们的示例文件

```yaml
kind: pipeline
name: demo

steps:
  - name: build
    image: golang:1.11.4
    commands:
      - pwd
      - go version
      - go build .
      - go test demo/util

  #  - name: frontend
  #    image: node:6
  #    commands:
  #      - npm install
  #      - npm test

  - name: publish
    image: plugins/docker:latest
    settings:
      username:
        from_secret: docker_username
      password:
        from_secret: docker_password
      repo: example/demo
      tags: latest

  - name: deploy
    image: appleboy/drone-ssh
    pull: true
    settings:
      host: example.me
      user: root
      key:
        from_secret: deploy_key
      script:
        - cd /data
        - mkdir app/
        - cd /data/app
        - docker rmi -f example/demo
        - echo "login docker"
        - echo "login success, pulling..."
        - docker pull example/demo:latest
        - echo "image running"
        - docker run -p 8088:8088 -d example/demo
        - echo "run success"

```

我们首先进行简单的golang build和test操作然后根据Dockerfile文件把我们的程序构建成docker镜像，接着上传到docker hub中，然后通过drone-ssh插件部署这个镜像。



## 开始构建

有了配置文件之后，推送代码我们就可以去drone查看构建进度：

![drone](https://i.loli.net/2019/03/06/5c7f655ba2519.png)

# drone的设置

在进入drone的时候，选择一个项目我们可以进行一些必要的设置，比如配置secrets，定时任务和徽章等等。

比如配置文件需要的密钥，用户名和密码，一些环境变量都可以在secrets设置，构建状态徽章可以在你的项目README.md文件加上去。

![设置](https://i.loli.net/2019/03/06/5c7f66a50f9b1.png)

项目加上徽章：

![徽章](https://i.loli.net/2019/03/06/5c7f66ed1d946.png)

[示例代码](https://github.com/nagaame/drone-golang-example)，本文完。
