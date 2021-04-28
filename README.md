# Go R3-D3

![Gopher image](https://golang.org/doc/gopher/fiveyears.jpg)
*Gopher image by [Renee French][rf], licensed under [Creative Commons 3.0 Attributions license][cc3-by].*

# Getting started

Clone the repository

    git clone https://github.com/fittipaldi/go-r3-d3.git

Switch to the repo folder

    cd go-r3-d3 - this is the Project Root

# Executing

You can run the go code as developing
    
    go run main.go
    
    
# Testing

You have to go to path `app/apiv1/controllersv1` and inside this folder you can run the command below

    go test

# Loading balance using Nginx

The idea for this project is that you can create multi instance of the same service that you can use a proxy to manager the balance.
this is how you can add into the Nginx to manager the process balance.

    upstream backend {
        server localhost:8001;
        server localhost:8002;
        server localhost:8003;
    }
        