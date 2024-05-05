build:
		if [ -f taxi-fare ]; then rm -rf taxi-fare; fi;
		if [ -f taxifare.log ]; then rm -rf taxifare.log; fi;
		GOOS=linux GOARCH=arm64 go build -o taxi-fare;
