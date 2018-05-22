const http = require('http')
const url = require('url')
const fs = require('fs')
const path = require('path')
const port = 3000

const contentMap = {
    '.ico': 'image/x-icon',
    '.html': 'text/html',
    '.js': 'text/javascript',
    '.json': 'application/json',
    '.css': 'text/css',
    '.png': 'image/png',
    '.jpg': 'image/jpeg',
    '.wav': 'audio/wav',
    '.mp3': 'audio/mpeg',
}

const cache = {}

function parseExtension(url) {

}

function serve(req, res) {
    console.log(req.url)
    let pathname = '.' + req.url
    const ext = path.parse(pathname).ext
    if (cache[pathname]) {
        res.setHeader('Content-type', contentMap[ext] || 'text/plain')
        res.end(cache[pathname])
    } else {
        fs.exists(pathname, function (exist) {
            if (!exist) {
                res.statusCode = 404
                res.end(`${pathname} not found!`)
                return
            }
            if (fs.statSync(pathname).isDirectory()) {
                pathname += '/index.html'
            }
            fs.readFile(pathname, function (err, data) {
                if (err) {
                    res.statusCode = 500
                    res.end(`internal server error`)
                    console.log(err)
                } else {
                    cache[pathname] = data
                    res.setHeader('Content-type', contentMap[ext] || 'text/plain')
                    res.end(data)
                }
            })
        })
    }
}

http.createServer(serve).listen(port)
console.log(`listening on port ${port}`)