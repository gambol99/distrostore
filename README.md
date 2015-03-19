[![Build Status](https://travis-ci.org/gambol99/distrostore.svg?branch=master)](https://travis-ci.org/gambol99/distrostore)
[![GoDoc](http://godoc.org/github.com/gambol99/distrostore?status.png)](http://godoc.org/github.com/gambol99/distrostore)

### **Distro Store**
-----------------

The Distro Store is a wrapper for [Consul](https://github.com/hashicorp/consul); the use case being you want the functionality of the Consul (raft consensus, node membership and notification, distributed key/value store, but without having to run it as a separate / external service, i.e. you want it embed into your application.

#### **Usages**

Bootstrapping the cluster and or joining a cluster

	cfg := distrostore.DefaultContext()
	cfg.Bootstrap = true
	// set the bind interface
	cfg.BindAddress = "127.0.0.1"
	cfg.BindAdvertised = "127.0.0.1"
	// give some datadir 
	cfg.DataDir = "/var/lib/somedir
	store, err := distrostore.New(cfg)

	# bootstapping and connection to other members
	(same as above, but remove the Bootstrap flag)
	cfg.Members = []string{host1:port, host2:port}
	store, err := distrostore.New(cfg)

