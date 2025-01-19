# Discogs to Spotify Playlist Converter

This project is a web application that converts a [Discogs](https://www.discogs.com/) URL into a Spotify playlist. It leverages Go for the backend logic and server-side rendering, enhanced with [HTMX](https://htmx.org/) for interactivity and styled using [Tailwind CSS](https://tailwindcss.com/).

![image](https://github.com/user-attachments/assets/b0f843b2-f2f7-4745-9a5e-80f40b737e42)


## Usage

1. Login to Spotify. Your playlist will be created automatically in your account.
2. Introduce a Collection, Wantlist or List URL from Discogs.

    Examples of valid URLs:
   - Collection: `https://www.discogs.com/es/user/username/collection`
   - Wantlist: `https://www.discogs.com/es/wantlist?user=username`
   - List: `https://www.discogs.com/es/lists/SomeList/1545836`
3. Enjoy the music.

## Tech Stack

- **Backend**: [Go](https://go.dev/), [gin](https://github.com/gin-gonic/gin), [gorilla/sessions](https://github.com/gorilla/sessions)
- **Frontend**: [HTMX](https://htmx.org/), [Tailwind CSS](https://tailwindcss.com/)
- **APIs**: [Discogs API](https://www.discogs.com/developers/), [Spotify API](https://developer.spotify.com/documentation/web-api)

## Getting Started

### Prerequisites

- Go 1.23
- Spotify Developer Account with a [registered application](https://developer.spotify.com/documentation/web-api/concepts/apps).

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/discogs-to-spotify.git
   cd discogs-to-spotify
   ```

2. Set up environment variables:

   Create a `.env` file with the following keys:

   ```env
   SPOTIFY_CLIENT_ID=your_spotify_client_id
   SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
   SPOTIFY_REDIRECT_URI=matches_the_settings_on_spotify_app (e.g. http://localhost:8080/auth/callback)
   SESSION_KEY=random_secret_for_gorilla_session
   PORT=optionally_set_the_port (8080 by default)
   ```
   2.1. ⚠️ Don't forget to set the _Redirect URI_ on your Spotify app settings.

3. Install dependencies:

   ```bash
   go mod tidy
   ```

4. Run the application:

   ```bash
   go run main.go
   ```

5. Open your browser and navigate to `http://localhost:8080`.


## Contributing

Contributions and feedback are welcome! Feel free to submit a pull request or open an issue for discussion.
(🚧 CONTRIBUTE.md WIP 🚧)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Happy playlisting! 🎵
