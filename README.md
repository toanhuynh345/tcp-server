# How it works
1. Init server with address port :3000.
2. Start to listen connection tcp.
3. When there is a client connection, it is openned to read data from client.
4. Each of messages from client will be sent and then processed
5. Server will response to the client with sentences "thank you for your message!\n"
6. If errors, read errors and still accept connection, print errors to console but don't interrupt server.
