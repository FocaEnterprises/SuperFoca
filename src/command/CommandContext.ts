import { CommandInteraction, GuildMember } from "discord.js"

class CommandContext {
  readonly interaction: CommandInteraction

  constructor(interaction: CommandInteraction) {
    this.interaction = interaction
  }

  async reply(reply: string): Promise<void> {
    return this.interaction.reply(reply)
  }

  async getMember(): Promise<GuildMember> {
    const members = await this.interaction.guild.members.fetch()
    return members.get(this.interaction.user.id)
  }
}

export default CommandContext