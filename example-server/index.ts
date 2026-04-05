const server = Bun.serve({
    port: 3000,
    async fetch(request) {
        console.log("Received request with headers:");
        console.table(Object.fromEntries(request.headers));
        const url = new URL(request.url);
        const path = url.pathname;

        if (request.method === "GET") {
            if (path === "/") {
                const response = new Response("Welcome to Bun!", { status: 200 });
                // Append Set-Cookie so multiple cookies are preserved
                response.headers.append(
                    "Set-Cookie",
                    "sessionId=abc123; HttpOnly; Path=/; SameSite=Lax"
                );
                return response;
            }

            if (path === "/health") {
                return new Response("OK", { status: 200 });
            }

            if (path === "/cookies") {
                // Echo back cookies sent by the client
                const incoming = request.headers.get("cookie") || "";
                return new Response(JSON.stringify({ cookies: incoming }), {
                    status: 200,
                    headers: { "Content-Type": "application/json" },
                });
            }
        }

        if (request.method === "POST") {
            console.log("Received POST with body:");
            const body = await request.text();
            console.log(body);
            console.log("Path: " + path);

            if (path === "/form") {
                const cookies = request.headers.get("cookie") || "";
                console.log("Received cookies:", cookies);
                const res = new Response(JSON.stringify({ ok: true, receivedCookies: cookies }), {
                    status: 200,
                    headers: { "Content-Type": "application/json" },
                });
                // Set a cookie in response (append to allow multiple Set-Cookie)
                res.headers.append(
                    "Set-Cookie",
                    "formReceived=yes; HttpOnly; Path=/; SameSite=Lax"
                );
                return res;
            }
        }

        return new Response("Not Found", { status: 404 });
    },
});

console.log(`Listening on ${server.url}`);