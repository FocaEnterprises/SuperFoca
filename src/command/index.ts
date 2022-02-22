import Command from "./Command"
import CommandContext from "./CommandContext"

import HelloCommand from "./commands/HelloCommand"
import PingCommand from "./commands/PingCommand"

const commands: Command[] = [ PingCommand, HelloCommand ]

export { commands, Command, CommandContext }