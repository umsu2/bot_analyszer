# Bot Analyser

1. records IP and other info of each req
2. record action of each IP (user)
3. return back 200 on all requests
4. analyse and store request routes, payload, timestamp


this is a trivial project.
I want to try the following
- grpc
- rabbitmq
- kubernetes
- gokit
- machinary


architecture

1. gateway
    - grpc to two different microservices (raw storage, processing)
2. processing
   - get the route, ip, important payload -> place on queue
3. storage
    - get the raw byte, place on queue
4. workers
    - consumes the queue and place payload either in sql or mongo
