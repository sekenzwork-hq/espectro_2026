# Prerequisites

-Git (should be installed on you machine)
-IDE (Use any IDE as you wish)


# Install Go SDK from official site

You can download the SDK from golang official website : https://go.dev/doc/install. After installing SDK, check 
version of the Go by using this command in cmd or terminal

```bash
 
    go version

```

Expected output : go version go1.27.0 windows/amd6

# Create a Folder

Create a folder where you want to clone this project and clone this repository by using this command in cmd or terminal

```bash

   git clone https://github.com/sekenzwork-hq/espectro_2026

```

Navigate to the folder where you cloned it and open it in an IDE

# Required Setup

After opening the IDE, open the terminal and use these commands to run the server

```bash
    go mod tidy
    
    go mod download

    go install github.com/air-verse/air@latest

    go tool air

```