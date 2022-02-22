import CommandContext from "./CommandContext"

interface Command {
  readonly name: String
  readonly description: String

  execute: (context: CommandContext) => Promise<void>
}

export default Command