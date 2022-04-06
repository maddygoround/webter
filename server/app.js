const axios = require("axios");
const express = require("express");
const cors = require("cors");
const app = express();

app.use(cors());

app.get("/read10kbfile/:path", async (req, res) => {
    const resp = await axios("https://www.10kb.site/" + req.params.path + "-offer", {
        method: "GET",
    })

    res.send(resp.data);
})

app.listen("3001");