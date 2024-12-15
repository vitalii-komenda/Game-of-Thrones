Requirements
--
- docker

Setup
--
Spin up db and migrate tables
```
make bootstrap
make import-data
```
Running
--
Build docker container and run it on localhost:8080
```
make brun
```
To get elasticsearch working, run this and wait a minute for synchronizing from postgres

```
make brun-elastic
```

Swagger doc
--
http://localhost:8080/swagger/index.html

Coverage
--
[Coverage Report](./coverage.html)
