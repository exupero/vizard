build: vizard.go assets.go
	go build vizard.go assets.go

debug: vizard.go assets/index.html assets/styles.css assets/bundle.js
	go-bindata -debug -prefix "assets/" -o assets.go assets/
	go build vizard.go assets.go

assets.go: assets/index.html assets/styles.css assets/bundle.js
	go-bindata -prefix "assets/" -o assets.go assets/

assets/bundle.js: vizard.js
	browserify vizard.js -t babelify --outfile assets/bundle.js
