var axios = require("axios").default;
var http = require('http');
var config = require('config')
const redis = require("redis");
const redisPort = config.redis.port
const client = redis.createClient(redisPort,config.redis.server);

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
    let callresponse = (res, payload) => {
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
                        <h2>` + payload.location.name + `, ` + payload.location.country + `</h2>
                        <span style="font-size: 20pt;">` + payload.current.condition.text + `</span><img style="vertical-align: middle" src='//cdn.weatherapi.com/weather/64x64/day/116.png' />
                        <h1>`
            + payload.current.temp_c + `&deg;C
                        </h1>
                    </div>
                </body>
            </html>`
        );
    }
    const searchTerm = options.params.q;
    client.get(searchTerm, async (err, data) => {
        if (err) throw err;
        if (data) {
            payload = JSON.parse(data)
            callresponse(res, payload)
            res.end();
        } else {
            axios.request(options).then(function (response) {
                payload = response.data
                callresponse(res, payload)
                res.end();
                client.setex(searchTerm, 600, JSON.stringify(payload));
            }).catch(function (error) {
                res.write("Invalid request");
                res.end();
                console.error(error);
            });
        }
    });
});

server.listen(config.network.port);

console.log('Myapp is running at port ' + config.network.port)