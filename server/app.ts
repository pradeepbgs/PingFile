import { type ContextType, Diesel } from 'diesel-core'


const app = new Diesel()

const port = process.env.PORT ?? 3000

app.get('/', (ctx: ContextType) => {
    return ctx.send("hello world");
})

app.post('/err', async (ctx: ContextType) => {
    const body = await ctx.body

    const statuses = [200, 400, 500]
    const status = statuses[Math.floor(Math.random() * statuses.length)]

    if (status === 200) {
        return ctx.send({ success: true, body }, status)
    }

    if (status === 400) {
        return ctx.send({ error: 'Bad Request' }, status)
    }

    return ctx.send({ error: 'Internal Server Error' }, status)
})


app.listen(port, () => {
    console.log("example server running for pinggfile to ping on me")
})