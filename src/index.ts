import "dotenv/config"

import { Client, Interaction } from "discord.js"
import { REST } from "@discordjs/rest"
import { Routes, GatewayIntentBits } from "discord-api-types/v9"

const commands = [
  {
    name: "hello",
    description: "Respondo um World!"
  }
]

const rest = new REST({ version: '9' }).setToken(process.env.TOKEN)

const postSlashCommands = async () => {
  console.log("Posting slash commands")

  try {
    await rest.put(Routes.applicationGuildCommands(process.env.CLIENT_ID, process.env.GUILD_ID), { body: commands })
  } catch (error) {
    console.log("Failed to post slash commands")
    console.log(error)
  }
}

const client = new Client({ intents: [GatewayIntentBits.Guilds] })

const start = async () => {
  await postSlashCommands()
  await client.login(process.env.TOKEN)
}

client.on('ready', () => console.log(`Logged-in @${client.user.tag}!`))

client.on('interactionCreate', (interacion: Interaction) => {
  if(!interacion.isCommand()) return

  if(interacion.commandName === 'hello') {
    interacion.reply("World!")
  }
})

start()