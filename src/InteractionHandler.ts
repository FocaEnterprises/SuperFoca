import { Interaction } from "discord.js"
import { REST } from "@discordjs/rest"
import { Routes } from "discord-api-types/v9"

import { commands, CommandContext, Command } from "./command"

const rest = new REST({ version: '9' }).setToken(process.env.TOKEN)

async function refreshSlashCommands() {
  try {
    console.log("Trying to refresh slash commands")
    await rest.put(Routes.applicationGuildCommands(process.env.CLIENT_ID, process.env.GUILD_ID), { body: commands as Command[] })
  } catch (error) {
    console.log("Failed to refresh slash commands")
    console.log(error)
  }
}

async function onInteraction(interaction: Interaction) {
  if(!interaction.isCommand()) return

  const name = interaction.commandName

  for(const command of commands) {
    if(command.name === name) {
      const context = new CommandContext(interaction)
      return command.execute(context)
    }
  }
}

export { refreshSlashCommands, onInteraction }