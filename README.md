# net-cat

"net-cat" is a custom version of Linux's net-cat, written in Go. It is a networking tool that allows users to create TCP servers and connect to remote TCP servers.

## Facility

To install "net-cat", you must have Go installed on your system. You can then clone this repository and build the executable with the following command:
```
bash go build -o net-cat
```


## Use

### Create a server

To create a server, run "net-cat" with the port number you want the server to listen on. By default, the server listens on port 8989 if no port is specified. The program works with localhost.
```
./net-cat <port>
```

The client connects using the following command:
```
./net-cat localhost <port>
```

## Contribution

Contributions are welcome. If you would like to contribute to this project, please submit a pull request with your changes.

## Licence

"net-cat" is distributed under the MIT license. See the LICENSE file for more information.