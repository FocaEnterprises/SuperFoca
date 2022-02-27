import path from 'path';
import 'dotenv/config';

import { BotClient } from './entities/Client';

const TOKEN = process.env.TOKEN;

if (!TOKEN) throw new Error('No token provided');

const client = new BotClient({
  intents: ['GUILDS'],
  commandsDir: path.join(__dirname, 'commands'),
});

client.init(TOKEN);
