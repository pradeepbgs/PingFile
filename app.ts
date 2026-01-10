


Bun?.serve({
    port:3000,
    fetch: (req:Request) => {
        const {pathname} = new URL(req.url)
        if(pathname === "/") {
        return new Response("Hello", {status:500})
        }
        else if (pathname === "/post") {
            return new Response("post")
        }
    }
})
