import "dotenv/config"

import client from "./SuperFoca"

client.init(process.env.TOKEN)