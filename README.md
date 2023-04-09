In my previous post, [Understanding Unix Domain Sockets in Golang](https://www.kungfudev.com/posts/understanding-unix-domain-sockets-in-golang/), I mentioned that one potential use case for Unix domain sockets is to communicate between containers in Kubernetes. I received requests for an example of how to do this, so in this post, I'll provide a simple example using two Go applications that you can find in this [repository](https://github.com/douglasmakey/go-sockets-uds-network-pprof).

Using Unix domain sockets in Kubernetes can be an effective way to communicate containers within the same pod.

Some advantages of using Unix domain sockets for communication between containers within a pod in K8s are:

-   Faster communication than using network sockets. This can be useful when containers need to communicate frequently or transfer large amounts of data.
-   No need for a network interface
-   Secure transmission, the communication is restricted to the local host.
-   Simplicity, no need for IP addresses or port numbers. Ok, this maybe isn't a great advantage but YOLO!

In K8s, you can achieve this by sharing a volume between the containers and using the socket file within the volume as the communication channel.

### Sharing a Volume in Kubernetes:

To share the Unix domain socket file between two containers in the same pod, you must create an `emptyDir` volume and mount it in both containers. In Kubernetes, you can do this using a `volumeMount` in the container specification in the yaml.

> An emptyDir volume is a temporary volume created when a Pod is assigned to a node and exists as long as that Pod runs on that node. An emptyDir volume is initialized with an empty directory and can be used to store data shared between the containers in the Pod.

Here is an example of how you might create a volume and mount it in two containers in the same pod:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kungfudev-deploy
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: nethttp
        image: douglasmakey/simple-http:latest
        volumeMounts:
        - mountPath: /tmp/
          name: socket-volume
      - name: unixhttp
        image: douglasmakey/simple-http-uds:latest
        volumeMounts:
        - mountPath: /tmp/
          name: socket-volume
      volumes:
      - name: socket-volume
        emptyDir: {}
```

In this example, both `nethttp` and `unixhttp` have a volume mounted in the container’s filesystem in `/tmp`, which allows them to access the same files within the volume.

### Running

Ready to try this in your own Kubernetes cluster? Check out this [repository](https://github.com/douglasmakey/go-sockets-uds-network-pprof), which has everything you need to get started! The included `Earthfile` in the repo will help you build and create container images for use in your local K8s cluster, or you can simply use the images published on Docker Hub. Follow the example YAML above to set up your own Unix domain socket communication between containers.

> Earthly is a tool for building, testing, and deploying applications using containers. It provides a simple and easy-to-use command-line interface for defining and automating the build, test, and deployment steps of a project using a script called an Earthfile.

Once the deployment is up and running on your K8s cluster, it's time to put it to the test! Head over to the `nethttp` container and make a request to `/test` using curl to see how it performs.

```bash
$ k exec -it kungfudev-deploy -c nethttp bash
---
$ curl localhost:8000/test
Hello kung fu developer from a server running on UDS!
```

### Show me the code!

The `unixhttp` app code:

```go
package main

import (...)

const socketPath = "/tmp/httpecho.sock"

func main() {
	// Create a Unix domain socket and listen for incoming connections.
	socket, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal(err)
	}

	// Cleanup the sockfile.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove(socketPath)
		os.Exit(1)
	}()

	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi kung fu developer from a server running on UDS! \n"))
	})

	server := http.Server{Handler: m}
	if err := server.Serve(socket); err != nil {
		log.Fatal(err)
	}
}

```

The `nethttp` app code:

```go
package main
import (...)

var (
	socketPath = "/tmp/httpecho.sock"
	// Creating a new HTTP client that is configured to make HTTP requests over a Unix domain socket.
	httpClient = http.Client{
		Transport: &http.Transport{
			// Set the DialContext field to a function that creates
			// a new network connection to a Unix domain socket
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}
)

func test(w http.ResponseWriter, req *http.Request) {
	resp, err := httpClient.Get("http://unix/")
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.Write(b)
}

func main() {
	http.HandleFunc("/test", test)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

```

That's it! With just a few simple steps, you can leverage the power of Unix domain sockets to communicate between containers in a Kubernetes pod. This may be a simple example, but the technique can be applied to all sorts of real-world scenarios. Thanks for following along, and I hope you found this example helpful and easy to understand. Happy coding!