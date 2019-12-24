.PHONY: stations
stations:
	env GOOS=linux GOARCH=amd64 go build -o cmd/lambda/getstations/getstations cmd/lambda/getstations/main.go
	zip -j /tmp/stations.zip cmd/lambda/getstations/getstations
	aws lambda update-function-code --function-name getstations --zip-file fileb:///tmp/stations.zip
	aws lambda invoke --function-name getstations /dev/stdout
	rm cmd/lambda/getstations/getstations

.PHONY: times
times:
	env GOOS=linux GOARCH=amd64 go build -o cmd/lambda/times/times cmd/lambda/times/main.go
	zip -j /tmp/times.zip cmd/lambda/times/times
	aws lambda update-function-code --function-name times --zip-file fileb:///tmp/times.zip
	aws lambda invoke --function-name times /dev/stdout
	rm cmd/lambda/times/times

.PHONY: lambda
lambda: stations times
	echo "Updated lambda"
