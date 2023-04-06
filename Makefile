db:
	docker run -d \
		--name cockroach \
		-p 26257:26257 \
			cockroachdb/cockroach:v22.2.7 \
				start-single-node \
				--insecure

tables:
	cockroach sql --insecure < example/create.sql

debug:
	go run crdbgen.go \
		--url "postgres://root@localhost:26257/curious_cupcake?ssl_mode=disable" \
		--debug