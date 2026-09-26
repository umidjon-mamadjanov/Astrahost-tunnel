from http.server import BaseHTTPRequestHandler, HTTPServer

class Handler(BaseHTTPRequestHandler):
    def do_POST(self):
        length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(length)

        response = (
            "method=" + self.command + "\n"
            "path=" + self.path + "\n"
            "content-type=" + self.headers.get("Content-Type", "") + "\n"
            "x-test-header=" + self.headers.get("X-Test-Header", "") + "\n"
            "body=" + body.decode() + "\n"
        ).encode()

        self.send_response(200)
        self.send_header("Content-Type", "text/plain")
        self.send_header("Content-Length", str(len(response)))
        self.end_headers()
        self.wfile.write(response)

print("POST test server listening on 127.0.0.1:8080")
HTTPServer(("127.0.0.1", 8080), Handler).serve_forever()
