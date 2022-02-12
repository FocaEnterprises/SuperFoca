import { Interaction } from "discord.js"
import { REST } from "@discordjs/rest"
import { Routes } from "discord-api-types/v9"

const commands = [
  {
    name: "hello",
    description: "Respondo um World!"
  }
]

const rest = new REST({ version: '9' }).setToken(process.env.TOKEN)

async function refreshSlashCommands() {
  try {
    console.log("Trying to refresh slash commands")
    await rest.put(Routes.applicationGuildCommands(process.env.CLIENT_ID, process.env.GUILD_ID), { body: commands })
  } catch (error) {
    console.log("Failed to refresh slash commands")
    console.log(error)
  }
}

async function onInteraction(interaction: Interaction) {
  if(!interaction.isCommand()) return

  if(interaction.commandName === 'hello') {
    interaction.reply("World!")
  }
}

export { refreshSlashCommands, onInteraction }