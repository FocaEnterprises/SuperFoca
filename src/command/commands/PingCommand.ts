import Command from "../Command"
import CommandContext from "../CommandContext"

const PingCommand: Command = {
  name: "ping",
  description: "Respondo com Pong!",

  async execute(context: CommandContext) {
    context.reply("Pong!")
  }
}

export default PingCommand