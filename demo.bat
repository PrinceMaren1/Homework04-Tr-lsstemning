start go run Client/Client.go -port 1000
start go run Client/Client.go -port 1001 -peers "1000"
start go run Client/Client.go -port 1002 -peers "1000,1001"
start go run Client/Client.go -port 1003 -peers "1000,1001,1002"