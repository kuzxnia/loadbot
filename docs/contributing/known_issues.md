


* srv not working with some DNS servers - golng 1.13+ issue see [this](https://github.com/golang/go/issues/37362) and [this](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#hdr-Potential_DNS_Issues)

    > Old versions of kube-dns and the native DNS resolver (systemd-resolver) on Ubuntu 18.04 are known to be non-compliant in this manner. 
