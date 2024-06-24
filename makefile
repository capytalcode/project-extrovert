PORT?=8080

build: templ
	go build -o bin/www

build/static: templ
	go run ./cmd/build

dev/templ:
	go run github.com/a-h/templ/cmd/templ@v0.2.707 generate --watch \
		--proxy=http://localhost:$(PORT) \
		--open-browser=false

dev/server:
	go run github.com/air-verse/air@v1.52.2 \
		--build.cmd "go build -o tmp/bin/main" \
		--build.bin "tmp/bin/main" \
		--build.exclude_dir "node_modules" \
		--build.include_ext "go" \
		--build.stop_on_error "false" \
		--misc.clean_on_exit true \
		-- -p $(PORT) -d

dev/sync_assets:
	go run github.com/air-verse/air@v1.52.2 \
		--build.cmd "go run github.com/a-h/templ/cmd/templ@v0.2.707 generate --notify-proxy" \
		--build.bin "true" \
		--build.delay "100" \
		--build.exclude_dir "" \
		--build.include_dir "static" \
		--build.include_ext "js,css"

dev:
	make -j3 dev/templ dev/server dev/sync_assets

run: build
	./bin/www

all: build

clean:
	if [[ -d "dist" ]]; then rm -r ./dist; fi
	if [[ -d "tmp" ]]; then rm -r ./tmp; fi
	if [[ -d "bin" ]]; then rm -r ./bin; fi
	rm ./static/uno.css
	rm $(TEMPL_FILES)

# For some reason "templ generate" does not detect the files in CI,
# so this is a workaround.
TEMPL_FILES=$(patsubst %.templ, %_templ.go, $(wildcard **/*.templ))
templ: $(TEMPL_FILES)
	@echo Generating templ files
%_templ.go: %.templ
	go run github.com/a-h/templ/cmd/templ@v0.2.707 generate -f $^ > /dev/null
