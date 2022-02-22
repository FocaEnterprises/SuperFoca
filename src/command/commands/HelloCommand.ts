import Command from "../Command"
import CommandContext from "../CommandContext"

const HelloCommand:Command = {
  name: "hello",
  description: "Respondo um olá, mundo!",

  async execute(context: CommandContext) {
    context.reply("Hello, World!")
  }
}

export default HelloCommand