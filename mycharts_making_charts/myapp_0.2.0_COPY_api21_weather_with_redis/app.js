var axios = require("axios").default;
var http = require('http');
var config = require('config')

var server = http.createServer(function (req, res) {
    res.writeHead(200, {
        'Content-Type': 'text/html'
    });
    var options = {
        method: 'GET',
        url: 'https://weatherapi-com.p.rapidapi.com/current.json',
        params: {
            q: req.url.trimStart("/")
        },
        headers: {
            'x-rapidapi-host': 'weatherapi-com.p.rapidapi.com',
            'x-rapidapi-key': process.env.APIKEY
        }
    };
    data = {}
    axios.request(options).then(function (response) {
        res.write(
            `<html>
                <style>
                    .content {
                    max-width: 500px;
                    margin: auto;
                    padding-top: 100px;
                    }
                    h1 {
                        font-size: 40pt;
                        margin: 0;
                    }
                    h2 {
                        font-size: 20pt;
                    }
                </style>
                <body>
                    <div class="content">
                        <h2>` + response.data.location.name + `, ` + response.data.location.country + `</h2>
                        <span style="font-size: 20pt;">` + response.data.current.condition.text + `</span><img style="vertical-align: middle" src='//cdn.weatherapi.com/weather/64x64/day/116.png' />
                        <h1>`
                            + response.data.current.temp_c + `&deg;C
                        </h1>
                    </div>
                </body>
            </html>`
        );
        res.end();
    }).catch(function (error) {
        res.write("Invalid request");
        res.end();
        console.error(error);
    });
});

server.listen(config.network.port);

console.log('Myapp is running at port ' + config.network.port)