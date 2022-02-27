import { BotClient } from '../entities/Client';
import { Command } from '../entities/Command';

class PingCommand extends Command {
  constructor(client: BotClient) {
    super(client, 'ping');
  }
}

export default PingCommand;
