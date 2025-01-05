tailwind-build:
	tailwindcss -i ./static/css/styles.css -o ./static/css/output.css --minify

tailwind-watch:
	tailwindcss -i ./static/css/styles.css -o ./static/css/output.css --watch