# xac
Find exactly the right assets with Xac
*xac currently under construction*

# kv

`kv` is a simple command line key-value utility. 
Ever find yourself wishing bash had python-style dicts? `kv` fills that niche.

## Installation

Download the [appropriate release](https://github.com/xkortex/xac/releases) for your os/arch. Extract, rename to `kv` and place in your path. That's it!
Run `kv --help` to view usage instructions.
## but how does it work?

Setting a key creates a file on the filesystem where the key is the filename and the value is the file contents. 
It's that simple! 
The kv store is located in your [AppDir](https://godoc.org/github.com/Wessie/appdirs). This is os-dependent. See AppDirs link for details.

Since the store is just files, you can share state by using bind mounts. I recommend using read only. 
This is for very simple use cases! For more critical purposes, look at `redis`, `consul`, `etcd` and the like.


Namespaces `-n` allow you to use the same key in different contexts.
