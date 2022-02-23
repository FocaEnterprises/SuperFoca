import { Client } from "discord.js"
import { GatewayIntentBits } from "discord-api-types/v9"
import { onInteraction, onMessage, refreshSlashCommands } from "./InteractionHandler"

const client = new Client({ intents: [GatewayIntentBits.Guilds] })

client.on('ready', () => console.log(`Logged-in @${client.user.tag}!`))
client.on('interactionCreate', onInteraction)
client.on('message', console.log)

async function init(token: string) {
  await client.login(token)
  await refreshSlashCommands()
}

export default { init }