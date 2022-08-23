module main

go 1.18

replace github.com/alexconrey/go-hs-api => ../go-hs-api

require (
	github.com/alexconrey/go-hs-api v0.0.0-00010101000000-000000000000
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
)

require github.com/felixge/httpsnoop v1.0.1 // indirect
