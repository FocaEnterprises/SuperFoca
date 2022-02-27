import { BotClient } from '../../entities/Client';
import { Command } from '../../entities/Command';

class HelloCommand extends Command {
  constructor(client: BotClient) {
    super(client, 'hello');
  }
}

export default HelloCommand;
