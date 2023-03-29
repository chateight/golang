//
// local webserver
//
'use strict';

main()

var html = require('fs').readFileSync('views/index.html');

function main() {
const http = require('http');

const hostname = '127.0.0.1';
const port = 3000;

const server = http.createServer((req, res) => {
  res.statusCode = 200;
  res.setHeader('Content-Type',{'Content-Type':'text/html'});
  res.end(html);
});

server.listen(port, hostname, () => {
  console.log(`Server running at http://${hostname}:${port}/`);
});

}