# Discogs to Spotify Playlist Converter

This project is a web application that converts a [Discogs](https://www.discogs.com/) URL into a Spotify playlist. It leverages Go for the backend logic and server-side rendering, enhanced with [HTMX](https://htmx.org/) for interactivity and styled using [Tailwind CSS](https://tailwindcss.com/).

![image](https://github.com/user-attachments/assets/f2eee859-ccad-4339-8596-a50e14307634)

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
   SPOTIFY_REDIRECT_URI=http://localhost:8080/auth/callback
   SPOTIFY_PROXY_URL=https://your-auth-proxy-url.com (optional, for local development)
   SESSION_KEY=random_secret_for_gorilla_session
   PORT=8080
   ENV=development
   ```

   ### 🔒 Auth Proxy Setup (Required for Local Development)

   Since Spotify deprecated localhost redirects, this app includes built-in auth proxy functionality for local development:

   **Setup Steps:**

   1. **Deploy your app** (if not already deployed):
      ```bash
      # Deploy to your preferred platform (Koyeb, Heroku, Railway, etc.)
      # Example for Koyeb using the included koyeb.yaml:
      koyeb app deploy
      ```

   2. **Configure environment variables**:
      ```env
      # In your local .env file:
      SPOTIFY_PROXY_URL=https://your-deployed-app-url.com
      ENV=development
      
      # In your production deployment:
      LOCAL_DEV_URL=http://localhost:8080  # Points to your local dev server
      ```

   3. **Update Spotify app settings**:
      - Redirect URI: `https://your-deployed-app-url.com/auth/proxy/callback/spotify`

   **How it works:**
   1. User clicks login on `localhost:8080`
   2. Spotify redirects to your **production app** at `/auth/proxy/callback/spotify`
   3. Your production app redirects the user back to `localhost:8080/auth/callback`
   4. Your local dev server handles the callback normally

   **Alternative: Use tunneling tools** (ngrok, cloudflared, etc.)

   1. Install and run a tunnel to your local server:
      ```bash
      # Using ngrok
      ngrok http 8080
      ```

   2. Use the provided HTTPS URL as your redirect URI in Spotify app settings.

   ⚠️ **Important**: The `SPOTIFY_REDIRECT_URI` should be set to your local callback URL (`http://localhost:8080/auth/callback`) **only for local development**. For production deployments, set `SPOTIFY_REDIRECT_URI` to your production callback URL (e.g., `https://your-deployed-app-url.com/auth/callback`). In both cases, the Spotify app settings should use your production app's proxy URL.

3. Install dependencies:

   ```bash
   go mod tidy
   ```

4. Run the application:

   ```bash
   go run main.go
   ```

5. Open your browser and navigate to `http://localhost:8080`.

## Production Deployment

For production deployment, you don't need the auth proxy. Simply:

1. Set your environment to production:
   ```env
   ENV=production
   SPOTIFY_PROXY_URL=  # Leave empty or remove
   ```

2. Configure your Spotify app redirect URI to your production callback URL:
   ```
   https://your-production-domain.com/auth/callback
   ```

## Project Structure

```
├── internal/
│   ├── adapters/          # External service adapters (Spotify, Discogs)
│   ├── core/              # Business logic and entities
│   ├── infrastructure/    # Configuration, server, session management
│   └── usecases/          # Application use cases
├── static/                # Static assets (CSS, images)
└── main.go                # Application entry point
```

## Contributing

Contributions and feedback are welcome! Feel free to submit a pull request or open an issue for discussion.
(🚧 CONTRIBUTE.md WIP 🚧)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Happy playlisting! 🎵
