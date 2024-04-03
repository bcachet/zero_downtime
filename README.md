# Zero downtime deployment of Podman container


```
docker compose up -d
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

We can then use our custom client that rely on Consul resolved by doing the following:

```
docker build . --file Dockerfile -t client --build-arg DIRECTORY=./client
docker run --network=zero_downtime_zero-downtime client -name=foo
```

## Update Go files from .proto file

```
docker build . --file Dockerfile.buf -t buf
docker run --volume ./helloworld:/defs buf
```
