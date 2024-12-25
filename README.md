# Cat Viewer - Beego Web Application

## Table of Contents
1. [Introduction](#introduction)
2. [Prerequisites](#prerequisites)
3. [Installation](#installation)
4. [Project Structure](#project-structure)
5. [Features](#features)
6. [Running the Application](#running-the-application)
7. [Unit Tests](#unit-tests) 


## Introduction
The **Cat Viewer** is a web application built with **Beego** for the backend and Vanilla JavaScript for the frontend. The project replicates the functionality seen on [The Cat API](https://thecatapi.com), where users can view random cat images, save their favorites, and vote on images.

## Prerequisites
Before running the project, ensure you have the following installed on your machine:

1. **Go (Golang)**:
   - Install Go from the official site: [Download Go](https://golang.org/dl/)
   - For detailed installation instructions, check out this [YouTube tutorial on installing Go](https://youtu.be/9IbfeyFlfeU?si=D6S6gQbI7rUxQfv8) for beginners.

2. **Beego Framework**:
   - To install Beego, run the following command:
     ```bash
     go get github.com/beego/beego/v2
     ```

3. **Other Dependencies**:
   - Make sure you have all Go dependencies set up in your Go project directory. You can use `go mod` to manage the dependencies:
     ```bash
     go mod tidy
     ```

## Installation

1. Clone this repository to your local machine:
   ```bash
   git clone https://github.com/yourusername/catviewer.git
   cd catviewer
   ```
2. Install Go dependencies:
   ```bash
   go mod tidy
   ```
3. Install Beego if not done already:
   ```bash
   go get github.com/beego/beego/v2
   ```   
4. Set up your The Cat API key:
   - You will need an API key from [The Cat API](https://thecatapi.com/).
   - Once you have the key, place it in your Beego config file (`app.conf`). 


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

## Features

- **Like and Dislike Cat Pics**:  
  Users can like and dislike cat pictures, and the voted pics are displayed in the "Voted Pics" section.
  
- **Add and View Favorite Cat Pics**:  
  Users can add cat pictures to their favorites, and these are shown in the "Favs" section. The "Favs" button in the header section allows users to view all their favorite cat pics.

- **Delete Favorite Pics**:  
  Users have the option to delete cat pictures from their favorites section.

- **Search Cat Breeds**:  
  The app allows users to search for different cat breeds in the "Breeds" section in header, providing information on various types of cats.





