# Ampersand-Rethink-Go

Simple REST API with AmpersandJS, Go, and RethinkDB.

### Setup

Instal Ampersand.js and use the Ampersand.js CLI to generate a new app:


```bash
$ npm i -g ampersand
$ ampersand
```

It will generate a new app. Choose Hapi or Express as the server, it doesn't matter since this is the server you'll actually use.

Once you've generated the app, you need to tell the frontend to use `http://localhost:8000` as the URL for the backend. To do that, go to `config/default.json` and change the value of `apiUrl` to the address above. Then, you can use you use this server.

To start this server you'll need Go and RethinkDB installed. Start RethinkDB, then in a different tab, run:

```bash
$ go run main.go
```

If that works, you're all set! Open your browser to `http://localhost:3000/` to see it working.

If you need to check the contents of your DB, go to `http://localhost:8080/` to view the RethinkDB admin panel.
