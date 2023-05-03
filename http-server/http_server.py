import http.server
import socketserver
import json

# PORT = 8000


    

class MyHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/ping':
            self.send_response(200)
            self.send_header('Content-type', 'text/plain')
            self.end_headers()
            self.wfile.write(b'pong')
        else:
            self.send_response(400)
            self.send_header('Content-type', 'text/plain')
            self.end_headers()
            self.wfile.write(b'BAD-REQUEST')



if __name__ == "__main__":
    with open('config.json') as f:
        config = json.load(f)

    servers = config['servers']
    for server in servers:
        HOST = server['host']
        PORT = server['port']
        with socketserver.TCPServer((HOST, PORT), MyHandler) as httpd:
            print(f"Serving on port {PORT}")
            httpd.serve_forever()