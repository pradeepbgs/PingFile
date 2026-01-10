


Bun?.serve({
    port:3000,
    fetch: () => {
        return new Response("Hello", {status:500})
    }
})
