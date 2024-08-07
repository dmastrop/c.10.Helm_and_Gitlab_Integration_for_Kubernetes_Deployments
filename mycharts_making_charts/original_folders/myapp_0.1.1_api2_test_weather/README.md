This is a simple web server written for "Helm: The Kubernetes package manager hands-on course"
## Running it natively
1. Install NodeJS
2. Run `npm install`
3. Modify the config.json to your preferences
4. Make sure you have a valid API key from https://rapidapi.com/weatherapi/api/weatherapi-com/
4. Run `APIKEY=your_key node app.js`
5. Assuming that you selected port `8080` and you want to check the current weather conditions in Paris, navigate to `localhost:8080/paris`

## Running it through Docker
You can build your own image using `docker build -t your_image_name:tag`
Or you can just use the already-pushed image at `afakharany/hellonodejs:2.0.0`. For example:
```bash
docker run -p 8080:80 -d -e APIKEY=your_key -v config.json:/app/config/default.json afakharany/hellonodejs:2.0.0
```