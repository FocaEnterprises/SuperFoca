import { CommandInteraction } from 'discord.js';

import { BotClient } from './Client';

export class Command {
  readonly name: string;
  readonly description?: string;

  private client: BotClient;

  constructor(client: BotClient, name: string, description?: string) {
    this.name = name;
    this.description = description;
    this.client = client;
  }

  execute(interaction: CommandInteraction): void {
    throw new Error('Method not implemented.');
  }

  getHelp(): object {
    throw new Error('Method not implemented.');
  }
}
