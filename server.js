const http = require('http');
const jsonFile = require('./schemas/content-tree.schema.json')

const server = http.createServer((req, res) => {
  if (req.url === '/content-tree' && req.method === 'GET') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify(jsonFile, null, 2));
  } else {
    res.writeHead(404);
    res.end();
  }
});

server.listen(3000, () => {
  console.log('Server running at http://localhost:5000');
});
