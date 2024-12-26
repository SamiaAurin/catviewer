# Cat Viewer - Beego Web Application

## Table of Contents
1. [Introduction](#introduction)
2. [Features](#features)
3. [Prerequisites](#prerequisites) 
4. [Installation](#installation)
5. [Project Structure](#project-structure)
6. [Running the Application](#running-the-application)
7. [Unit Tests](#unit-tests) 


## Introduction
The **Cat Viewer** is a web application built with **Beego** for the backend and Vanilla JavaScript for the frontend. The project replicates the functionality seen on [The Cat API](https://thecatapi.com), where users can view random cat images, save their favorites, and vote on images.


## Features

- **Like and Dislike Cat Pictures**:  
  Users can like and dislike cat pictures, and the voted pics are displayed in the "Voted Pics" section.
  
- **Add and View Favorite Cat Pictures**:  
  Users can add cat pictures to their favorites, and these are shown in the "Favs" section. The "Favs" button in the header section allows users to view all their favorite cat pics.

- **Delete Favorite Pictures**:  
  Users have the option to delete cat pictures from their favorites section.

- **Search Cat Breeds**:  
  The app allows users to search for different cat breeds in the "Breeds" section in header, providing information on various types of cats.


## Prerequisites
Before running the project, ensure you have the following installed on your machine:

1. **Go (Golang)**:
   - Install Go from the official site: [Download Go](https://golang.org/dl/)
   - For detailed installation instructions, check out this [YouTube tutorial on installing Go](https://youtu.be/9IbfeyFlfeU?si=D6S6gQbI7rUxQfv8) for beginners.
   - Verify installation:
      ```bash
       go version
      ```
2. **Beego Framework**:
   - To install Beego, run the following command:
     ```bash
     go install github.com/beego/bee/v2@latest
     ```

3. **Other Dependencies**:
   - Make sure you have all Go dependencies set up in your Go project directory. You can use `go mod` to manage the dependencies after clone the repository:
     ```bash
     go mod tidy
     ```


## Installation

1. Clone this repository to your local machine: 
   To get started, clone the repository into your Go workspace `(go/src)` to ensure the project is placed in the correct directory for your Go workspace. If your Go workspace does not contain a src folder, please create one.

   ```bash
   git clone https://github.com/yourusername/catviewer.git
   cd catviewer
   ```
   
2. Install Go dependencies:

   ```bash
   go mod tidy
   ``` 
3. Set up your The Cat API key:

   - You will need an API key from [The Cat API](https://thecatapi.com/).
   - Once you have the key, place it in your Beego config file (`app.conf`). 
   - Open the configuration file located at `conf/app.conf` and update the following:

      ```bash
      appname = catviewer
      httpport = 8080
      runmode = dev
      catapi_key = live_xDKhFremvCoTHT3aHP6IfaxZA5XOTjhdvunrnecnwHyXfy9mrU3b8Yeu4NTNfQ0i
      catapi_url = https://api.thecatapi.com/v1
      ```

## Project Structure
The project has the following structure:
```bash
catviewer
├── config
│   └── app.conf
├── controllers
│   └── cat_controller.go
├── routers
│   └── router.go
├── views
│   ├── catviewer.tpl
│   └── index.tpl
├── static
│   ├── js
│   │   └── catviewer.js
│   └── css
│       └── styles.css
├── tests
│   └── cat_controller_test.go
├── main.go
├── go.mod
└── go.sum
```

- controllers: Contains the Go controllers, including logic for interacting with The Cat API.
- views: Contains the HTML templates for rendering the frontend.
- static: Stores static files such as JavaScript and CSS.
- config: Stores the Beego configuration files.
- tests: Contains unit tests for the project.


## Running the Application

### Step 1: Start the Beego Server
Run the following command to start the Beego application:
```bash
bee run
```

### Step 2: Access the Application
Open your browser and navigate to: http://localhost:8080/cat/vote


## Unit Tests
This project includes unit tests to ensure code reliability.

### Run all tests:
```bash
go test -v ./tests/cat_controller_test.go
```

### Generate coverage report:
```bash
go test -coverprofile coverage.out ./...
go tool cover -html coverage.out
```



