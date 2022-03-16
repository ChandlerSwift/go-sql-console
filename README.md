# go-sql-console

`go-sql-console` is a lightweight SQL debug console that is designed to be
easily added to web applications written in Go.

I wrote `go-sql-console` to fill a need for an RSVP form handler I wrote
recently. It was a quick, one-off program that was only intended to be deployed
once, and so it didn't have much of an admin panel (only a page listing who had
responded). However, the guest list was updated several times, and each time I
had to log into that server, fire up `sqlite`, and make the changes I wanted.
That's ridiculous! I already have a web app for that. Just let me submit SQL
queries, run them, and let me see the results! So `go-sql-console` was born.

### Using
Include the handler into your app:

```go
http.Handle("/sql-console", sql_console.Handler{
    DB: db,
})
```

For an example, check out the `example/` directory.

### Security
TODO: joke about how it's not SQL injection if there's no query to inject into

On its own, `go-sql-console` allows unauthenticated users to execute arbitrary
SQL with the full permissions of the application. Make sure to include some kind
of authentication middleware when you include this. Or don't run it on a
public-facing server. Or both!
