package main

import (
	"atlas/webapp/agents"
	"atlas/webapp/db"
	"atlas/webapp/evals"
	"atlas/webapp/handlers"
	"atlas/webapp/messaging"
	"atlas/webapp/search"
	"atlas/webapp/server"
	"atlas/webapp/traces"
	"embed"
	"fmt"
	"mime"
	"os"
)

// THE FONTS ARE NAMED HERE OR THEY DO NOT SHIP. This directive lists each
// directory explicitly, so a new one is invisible to the binary until it is
// added -- static/fonts/ would have 404'd silently and the console would have
// gone on falling back to Segoe UI with nothing to say it had.
//
//go:embed static/index.html static/css/* static/js/* static/fonts/*
var staticFiles embed.FS

func main() {
	// Go's builtin MIME table knows .css, .js, .svg and .wasm -- NOT .woff2.
	// Without this the file server sniffs the bytes and answers
	// application/octet-stream. Browsers are lenient about font types in
	// @font-face and would probably still render it, and "probably" is not a
	// thing to ship a typeface on.
	mime.AddExtensionType(".woff2", "font/woff2")

	port := "8091"
	if p := os.Getenv("ATLAS_WEB_PORT"); p != "" {
		port = p
	}

	database, err := db.Open("data/webapp.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "db open: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	store := traces.NewStore(database)
	agentReg := agents.NewRegistry(database)
	evalEngine := evals.NewEngine(database, store)
	searchEngine := search.NewEngine(store, agentReg)
	msgBus := messaging.NewBus()
	msgBus.SetDB(database)

	h := handlers.New(store, agentReg, evalEngine, searchEngine, msgBus)

	// THE SERVICE WIRE TO THE DOOR (2026-09-24, his word: "build out the
	// auth"). This line read `ConfigureAuth(true, "", ...)` with the comment
	// "the service wire to the door stays empty -- the door is unchanged",
	// true when it was written and false since `--auth` existed: every
	// send-site guards on `h.service != ""`, so the glass sent NO
	// Authorization header at all. Arming the door would have answered every
	// one of this page's calls with 401 -- records, rack, worlds, Version
	// control -- and the PIN would not have helped, because that is a
	// different gate.
	//
	// THE SAME VARIABLE THE DOOR READS (`cmd/atlas-mcp/main.go`: --auth-service
	// or ATLAS_SERVICE), so the two agree by reading one place rather than by
	// someone remembering to set two.
	//
	// AN ENVIRONMENT VARIABLE AND NOT A FLAG, on purpose: RULE 7 is that keys
	// are never passed on a command line where `ps` can read them. Unset is
	// the ordinary case and is exactly the behaviour this glass had before --
	// no header, and a door without `--auth` does not ask for one.
	service := os.Getenv("ATLAS_SERVICE")

	// THE LOCK IS ON (2026-09-21, his words: "Simple login system for now,
	// user/pin to start", "make the thing at least semi-secure", "This PC
	// only"). One user, one PIN; the first person at this computer sets it up.
	h.ConfigureAuth(true, service, "data/sessions.json")
	h.ConfigureLock("data/user.json")
	srv := server.New(h, port, staticFiles)

	// PRESENCE, NEVER THE VALUE (RULE 7 / LAW 9; the same shape `team_status`
	// already uses). A boot line that said nothing would leave the one failure
	// this piece exists to prevent -- an armed door and a glass with no key --
	// looking exactly like an ordinary start.
	wire := "no service wire (set ATLAS_SERVICE to arm the door)"
	if service != "" {
		wire = "service wire held"
	}
	fmt.Printf("atlas-webapp listening on 127.0.0.1:%s -- this PC only, "+
		"opened with a PIN, %s\n", port, wire)
	if err := srv.ListenAndServe(); err != nil {
		fmt.Fprintf(os.Stderr, "server: %v\n", err)
		os.Exit(1)
	}
}
