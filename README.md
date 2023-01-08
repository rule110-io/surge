
![](https://github.com/rule110-io/surge/blob/development/documentation/img/surge_color.png?raw=true) 
# Surge - P2P on steroids (...and NKN)


## About
This is the code of [Surge](https://getsurge.io). Surge is a 100% decentralized and secure P2P file sharing client. It utilizes the [NKN blockchain network](https://nkn.org) to send, publish and receive files.  

## Anonymous P2P? Are you kidding me?!
No, we are not :) 
By completely bypassing current internet technologies and using the NKN network each client is identified by an ID and nothing more. Here are some skribbles to show you how it all works: 

![When sending a file through the NKN network the sender just needs to contact one NKN node and tell him to what client the file should be delivered to. To communicate the IP is still needed but instantly dropped by the entry node.](https://github.com/rule110-io/surge/blob/development/documentation/img/surge1.png?raw=true)

![NKN now routes the file through its network using the optimal an fastest route to the client. If there is a client who is connected to the recipient it delivers the file to it.](https://github.com/rule110-io/surge/blob/development/documentation/img/surge2.png?raw=true)

## Where to download?
Find the current version for your operation system [here](https://github.com/rule110-io/surge/releases)


## Wanna build Surge on your own?

Besides using our pre-built executables you can build surge by your own.

Prerequisites:
- A running Golang environment
- A running [WailsV2](https://wails.io/docs/gettingstarted/installation/) environment

1. clone this repository in your go projects
2. run ``wails build -p``
3. check the ``build`` directory

Other helpful commands

``` bash
# start surge backend
$ wails serve

# start surge frontend
$ cd frontend
$ npm run serve
```

For detailed explanation on how things work, checkout [Wails docs](https://wails.app/gettingstarted/).

## Contribute

Surge is an open source project so everyone is invited and welcome to help. If you want to get in contact with us just jump into the [NKN Discord](https://discord.gg/hAxzRUV7DN)
