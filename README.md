# Zero downtime with gRPC + Consul

Objective is to demonstrate how we can rely on gRPC client name-resolver + Consul to load balance requests accross available services for a given FQDN.
With this property demonstrated, our deployment scenario will become:

1. Deploy new version of the service (tagged as _blue_)
2. Await healthcheck for service tagged _blue_ to be ready
3. Gracefully stop the old version of the service (tagged _green_)

Service tagged _green_ will not accept any new request and will finalize ongoing requests first

## Setup Consul + gRPC servers
You can start Consul and our gRPC Greeter services via `docker compose`

```
docker compose up -d
docker compose logs -f greeter-server-blue greeter-server-green
```

If you go to http://0.0.0.0:8500/ui/dc1/services/greeter/instances you should see the 2 instances of the greeter service

We can check that the gRPC servers are working via

```
docker run \
    --volume ./helloworld:/helloworld \
    --network=zero_downtime_zero-downtime \
    fullstorydev/grpcurl \
        -import-path /helloworld \
        -proto helloworld.proto \
        -plaintext \
        -d '{"name": "foo"}' \
        greeting-server-blue:50000 helloworld.Greeter/SayHello
```

## gRPC client using Consul as name resolver + round robin load balancing

We can then use our custom client that rely on Consul name resolver by doing the following:

```
docker build . --file Dockerfile -t client --build-arg DIRECTORY=./client
docker run --network=zero_downtime_zero-downtime client -name=foo
```

You can check in our gRPC server logs that the requests that are reaching the 2 servers

If we stop one of the server, all the requests will go through the remaining one

```
docker compose stop greeter-server-blue
```

If the server is stopped while handling a request (while sleeping), it will await for the thread to terminate but will not accept any new request.


## FAQ

### Update Go files from .proto file

```
docker build . --file Dockerfile.buf -t buf
docker run --volume ./helloworld:/defs buf
```

