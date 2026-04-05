const server = Bun.serve({
    port: 3000,
    async fetch(request) {
        // Log headers of incoming request
        console.log("Received request with headers:",);
        console.table(request.headers);
        const path = new URL(request.url).pathname;

        if (request.method === "GET") {
            if (path === '/') {
                return new Response("Welcome to Bun!");
            } else if (path === '/health') {
                return new Response("OK", { status: 200 });
            }
        } else if (request.method === "POST") {
            console.log("Received POST with body:");
            const body = await request.text();
            console.log(body);
            console.log("Path: " + path);

            if (path === "/form") {
                return new Response("Success", { status: 200 });
            }
        } else {
            return new Response("Not Found", { status: 404 });
        }
    },
});

console.log(`Listening on ${server.url}`);