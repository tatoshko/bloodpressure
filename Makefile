deploy:
	@ go mod tidy
	@ git add .
	@ git commit -m 'auto'
	@ git push origin main
build:
	@ go mod tidy
	@ git pull origin main
	@ go build -ldflags="-s -w" -o bloodpressure
	@ sudo systemctl restart bloodpressure
