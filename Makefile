deploy:
	@ go mod tidy
	@ git add .
	@ git commit -m 'auto'
	@ git push origin master
build:
	@ go mod tidy
	@ git pull origin master
	@ cd assets && rice embed-go && cd ..
	@ ~/go/bin/go build -ldflags="-s -w"
	@ sudo systemctl restart embercat
