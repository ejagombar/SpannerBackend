# Spanner: A Spotify Analysis Backend

https://github.com/ejagombar/SpannerBackend/assets/77460324/0cd6d9d1-2a8b-443a-ad08-6397535586c6

Welcome to the repository for the Spanner backend\! Spanner is a web application designed to provide Spotify users with insightful analysis of their listening habits, including their top tracks, favorite artists, and detailed playlist metrics.

This backend is written in Go and works in tandem with the [Spanner Frontend](https://github.com/ejagombar/SpannerFrontend), which is built with Typescript and React.

**Note: This project is no longer live, due to changes in the API and since I am moving around a bit at the moment, I have not set my server up again.**

## About the Project

This project was undertaken as a way to learn more about web development, REST APIs, and the Go programming language. It has been a rewarding experience that has involved building a full-stack application from the ground up.

## Key Features

  * **User Authentication**: Securely authenticates users with their Spotify account using the OAuth 2.0 flow.
  * **Profile Information**: Fetches and displays user profile information, including display name and follower count.
  * **Top Tracks & Artists**: Retrieves a user's top tracks and artists across different time ranges (short, medium, and long term).
  * **Playlist Analysis**: Provides in-depth analysis of a user's playlists, including audio features like danceability, energy, and valence.
  * **Playlist Browse**: Allows users to view all of their public and private playlists.

### ~~Planned Features~~
> [!CAUTION]
Unfortunately, Spotify have removed the Audio Analysis endpoints from their API making the following features much more difficult to develop.

The following features are planned for future development:
  * Lyric Sentiment Analysis
  * Metrics and Server Error Alerts
  * Playlist Comparison

## Technical Details

The backend is built using the [Fiber](https://gofiber.io/) web framework for Go. It exposes a REST API that the frontend consumes to fetch data from Spotify.

1.  **Configuration**: The application uses Viper for configuration management, loading necessary API keys and settings from an `app.env` file or environment variables.
2.  **Authentication**: When a user logs in, the backend redirects them to Spotify's authorization page. After authorization, Spotify calls back to the backend with an authorization code, which is exchanged for an access token and refresh token. These tokens are stored securely in a session cookie.
3.  **API Communication**: The backend uses the `zmb3/spotify` library to communicate with the Spotify API. All requests to Spotify from the backend include the user's access token to fetch their personal data.
4.  **API Endpoints**: The backend provides a set of API endpoints for fetching user data, top tracks, artists, and playlist information.

## Setup

To run the backend server locally, follow these steps:

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/ejagombar/SpannerBackend.git
    cd SpannerBackend
    ```

2.  **Create a configuration file:**
    Create a file named `app.env` in the root of the project and add the following, filling in your own Spotify Developer credentials:

    ```env
    CLIENT_ID=YOUR_SPOTIFY_CLIENT_ID
    CLIENT_SECRET=YOUR_SPOTIFY_CLIENT_SECRET
    PORT=8080
    ```

3.  **Install dependencies and run the server:**

    ```bash
    go mod tidy
    go run main.go
    ```

    The server will start on the port specified in your `app.env` file.
