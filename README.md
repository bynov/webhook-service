# Webhook Service

## URLs

* POST http://127.0.0.1:8000/webhooks - Slave
* POST http://127.0.0.1:9000/webhooks - Master

## Start locally

To start this project locally just run

```bash
docker-compose up
```

It will run 2 applications in master & slave mode and 1 postgres database for each service.
