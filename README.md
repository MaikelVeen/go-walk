# Go Walk!

This project takes GPX track files, processes them into GeoJSON format, and displays them as paths on an interactive Mapbox map in a web browser.

## Prerequisites

*   [Go](https://go.dev/doc/install) (version 1.18 or later recommended)
*   [Node.js and npm](https://nodejs.org/en/download/) (for frontend build)
*   A [Mapbox Account](https://account.mapbox.com/auth/signup/) and a public access token.

## Setup & Usage

### Obtain GPX Data:

Download or acquire your GPX track files, for example from Garmin Connect.Place them into a directory. For example, you could create a directory named `gpx_data` in the project root: `mkdir gpx_data` and put your files there.

### Configure Mapbox Token:

Create a file named `.env` inside the `app/` directory: `touch app/.env`
Go to your Mapbox account and create a **new** public access token and add your new token to the `app/.env` file:

```dotenv
MAPBOX_ACCESS_TOKEN=YOUR_NEW_MAPBOX_ACCESS_TOKEN_HERE
```

### Process GPX Data:

Run the `transform` command. Specify the directory containing your GPX files (`-d`) and the desired output directory (`-O`), which should be `app/data` for the frontend to find it.

```sh
go build -o go-walk && ./go-walk transform -d path/to/your/gpx_data -O app/data
```

This will generate `.geojson` files and a `manifest.json` inside the `app/data/` directory.

### Install Frontend:

Navigate to the `app` directory and install the dependencies:
```sh
cd app && npm install && npm run build
```

This executes `build.js`, which reads `app/data/manifest.json`, embeds the GeoJSON data and the Mapbox token, and creates the final bundle in `app/dist/app.js`.

### View the Map:

Open the `app/index.html` file directly in your default web browser: `open app/index.html`