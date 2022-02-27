import fs from 'fs';
import path from 'path';
import { Client, ClientOptions, Collection, CommandInteraction, Interaction } from 'discord.js';

import { Command } from './Command';

type BotClientOptions = ClientOptions & {
  commandsDir: string;
};

export class BotClient extends Client {
  public commands: Collection<string, any> = new Collection();

  private commandsDir: string;
  private isCommandsReady: boolean = false;
  private clientToken?: string;

  constructor(options: BotClientOptions) {
    super(options);

    this.on('ready', this.onReady);
    this.on('interactionCreate', this.onInteraction);

    this.commandsDir = options.commandsDir;

    this.loadCommands();
  }

  public init(token: string): void {
    this.clientToken = token;
  }

  private onReady(): void {
    console.log(`Logged in as ${this.user?.tag}`);
  }

  private onInteraction(interaction: Interaction): void {
    switch (interaction.type) {
      case 'APPLICATION_COMMAND':
        this.commands.get((interaction as CommandInteraction).commandName)?.execute(interaction);
        break;

      default:
        break;
    }
  }

  private onClientFeatureReady(): void {
    if (this.isCommandsReady && this.clientToken) this.login(this.clientToken);
  }

  private updateCommands(): void {
    /** @todo implement slash commands updater */
  }

  private async loadCommands(): Promise<void> {
    const files = await this.getFiles(this.commandsDir);

    for (const filePath of files) {
      const { default: commandClass } = await import(filePath);

      if (typeof commandClass !== 'function') {
        throw new Error('Command class must be a class');
      }

      const command = new commandClass(this);

      if (!(command instanceof Command)) {
        throw new Error('Command class must be a subclass of Command');
      }

      this.commands.set(command.name, command);
    }

    this.isCommandsReady = true;
    this.onClientFeatureReady();
    console.log(`Successfully loaded ${files.length} commands`);
  }

  private async getFiles(dir: string): Promise<string[]> {
    const files: string[] = [];

    const dirPaths = await fs.promises.readdir(dir);

    for (const dirPath of dirPaths) {
      const itemPath = path.join(dir, dirPath);
      const itemInfo = await fs.promises.lstat(itemPath);

      itemInfo.isDirectory()
        ? files.push(...(await this.getFiles(itemPath)))
        : itemInfo.isFile() && files.push(itemPath);
    }

    return files;
  }
}
