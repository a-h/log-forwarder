# log-forwarded

WARNING: This is a PoC, don't use this.

## Tasks

### create-log-stream

```
aws logs create-log-stream --log-group-name LogStashGroup --log-stream-name LogStashStream
```

### put-log-events

```
aws logs put-log-events --log-group-name LogStashGroup --log-stream-name LogStashStream --log-events file://events.json
```

### cdk-deploy

Dir: cdk

```
cdk deploy
```

### docker-build

Dir: function

```
docker build -t log-forwarder:latest .
```

### docker-run

Dir: function

```
docker run log-forwarder:latest
```
