deploy:
	@ go mod tidy
	@ git add .
	@ git commit -m 'auto'
	@ git push origin main
build:
	@ go mod tidy
	@ git pull origin main
	@ ~/go/bin/go build -ldflags="-s -w"
	@ sudo systemctl restart bloodpressure
