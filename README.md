# WWorkShop Sample Application  (Model Viewer) 

This repository contains a simple web application built with Go that serves a World of Warcraft (WoW) model viewer. The application allows users to view WoW character models and item models using embedded templates and static resources.

## Features

- Character model viewer with customizable race, gender, and hairstyle options.
- Item model viewer that displays WoW items based on their display ID.
- Tailwind CSS and jQuery for UI styling and interactivity.
- Wowhead Model Viewer integration to display WoW models in 3D.

## Setup

### Requirements

- [Go](https://golang.org/doc/install) version 1.20 or higher
- [Node.js](https://nodejs.org/) for serving static resources

### Environment Variables

The following environment variables can be set to customize the app's behavior:

- `PORT`: The port on which the app will run (default: 8080).

### Installation

1. Clone this repository:

    ```bash
    git clone https://github.com/your-username/wow-model-viewer.git
    cd wow-model-viewer
    ```

2. Install dependencies:

    ```bash
    go mod download
    ```

3. Serve the application:

    ```bash
    go run main.go
    ```

4. Open your browser and navigate to `http://localhost:8080`.

### Running with Docker

To build and run the application locally using Docker, follow these steps:

1. Build the Docker image:

    ```bash
    docker build -t wow-model-viewer .
    ```

2. Run the container:

    ```bash
    docker run -p 8080:8080 --env PORT=8080 --env FLY_REGION=local wow-model-viewer
    ```

3. Open your browser and navigate to `http://localhost:8080`.


## Endpoints

### `/`
Serves the main page where users can view and customize character models and item models.

### `/item-lookup`
- **Method**: GET
- **Query Param**: `item` (required)
- Looks up an item by its ID and returns the display ID in JSON format.

### `/get-races`
- **Method**: GET
- Returns a JSON object of available WoW races with their corresponding race IDs.

### `/broken`
- **Method**: GET
- Simulates an internal server error.

## Static Resources

Static resources such as images or JavaScript files are served from the `/static` directory.

## Templates

HTML templates are embedded in the Go binary using the `embed` package and located under the `templates/` directory.

## Example Usage

1. Visit `http://localhost:8080`.
2. Customize your character by selecting a race, gender, and hairstyle.
3. View a 3D model of an item by selecting one from the dropdown in the item viewer.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.


