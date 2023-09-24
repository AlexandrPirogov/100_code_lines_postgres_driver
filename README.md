# Postgres driver in 100 lines of code

This repository was created for educational purposes to learn how to create a Postgres driver. The protocol used turned out to be quite simple.

## How it works

For communication with Postgres, Postgres uses protocal `Frontend-Backend`. `Frontend` sends `Message` to `Backend` and reads response.
Generally speacking `Messages` are consist of :
1) Tag
2) Message Len
3) Payload

The exception is `Startup message` as it doesn't contain `Tag`.

To start communication with Postgres backed you must first send `Startup` Message. 
Backend sends response with `Tag 'R'` and payload. Depends on Postgres auth configuration payload may have different value.
In this case as auth method set to `TRUST` payload will be equal `0` that means we are authenticated successfully.
Then Postgres backend sends couple of messages. We are interested in `Tag 'Z'` that means Postgres is ready to handle our queries.

After that `Front` is able to send query messages to back and read response.

The whole process of communication between `Frontend` and Postgres `Backend` more complex, i described it in few words and simplifed.

## Try it yourself!

Before you try it yourself, please note the following:
1) First of all take attention at `pg_hba.conf` - it uses `trust` auth method. `pg_hba.conf` exists only for running `postgres` container with no auth to make writing Postgres driver easy.
2) *Not all queries and responses will work correct*. Right now driver is able to execute simple queries such as `SELECT`, `INSERT`, `UPDATE`, `ALTER`, `CREATE/DROP TABLE`.
Errors and some reponses are not handled!


If you would to try it all you need is docker and pre-installed Go:
1) clone repo
2) go to directory with cloned repo
3) run the docker-compose. If ports are busy just change it.
```
docker-compose up
```
4)
```
go run main.go pretty.go
```

You'll see inviting char '->'. Feel free to write queries.

## References
1) [Postgres manual](https://www.postgresql.org/docs/current/protocol.html)
2) [PDF Postgres on wire](https://beta.pgcon.org/2014/schedule/attachments/330_postgres-for-the-wire.pdf)
