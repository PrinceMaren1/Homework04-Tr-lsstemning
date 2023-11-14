To start the program:

- Run the first Client with a given port. "go run Client.go -port 1000"
- Run subsequent Clients with a given port and the ports of existing Client seperated by a comma. 
- "go run Client.go -port 1001 -peers"1000"
- "go run Client.go -port 1002 -peers"1000,1001"


If you want to run a quick demo and you are on windows, you can run the demo.bat file.