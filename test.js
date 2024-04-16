const http = require('http');

// Create HTTP server
const server = http.createServer((req, res) => {
    // Check if the request method is POST
    if (req.method === 'POST') {
        let body = '';
        // Listen for data event to receive data chunks
        req.on('data', chunk => {
            body += chunk.toString(); // Convert Buffer to string
        });
        // Listen for the end event to ensure all data has been received
        req.on('end', () => {
            console.log('Received POST request with body:', body);
            console.log('Received POST request with req:', req);
            res.writeHead(200, { 'Content-Type': 'application/json' });
            res.end('POST request received');
        });
    } else {
        // Respond to non-POST requests
        res.writeHead(405, { 'Content-Type': 'application/json' });
        res.end('Only POST method is supported');
    }
});

const PORT = 3000; // Port number for the HTTP server

// Start the server
server.listen(PORT, () => {
    console.log(`Server is listening on port ${PORT}`);
});
