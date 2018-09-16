build:
	GOOS=linux GOARCH=mipsle go build -ldflags="-s -w" -o weather-station main.go 
	
copy: 
	rsync -P -a weather-station root@omega-5D69.local:/root/go