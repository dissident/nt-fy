import rabbit from 'amqplib';
import dotenv from 'dotenv';
import TelegramBot from 'node-telegram-bot-api';

import { Channel } from 'amqplib';

dotenv.config();

const amqpUrl = process.env.AMQP || 'amqp://localhost';
const token = process.env.TELEGRAM_TOKEN || 'token';

const bot = new TelegramBot(token);
const chatId = process.env.CHAT_ID || '666';

const q = 'notify';
const queueOptions = {
  durable: true,
  deadLetterExchange: "notify.dlx",
  deadLetterRoutingKey: "notify.dlx",
};
const prefetchCount = 5;

const open = rabbit.connect(amqpUrl);

const sendMessage = (ch: Channel, msg: string) => {
  ch.sendToQueue(q, Buffer.from(msg));
};

// Consumer
(async () => {
  const conn = await open;

  const ch = await conn.createChannel();
  await ch.prefetch(prefetchCount);
  await ch.assertQueue(q, queueOptions);
  await ch.consume(q, async (msg) => {
    if (msg !== null) {
      const content = msg.content.toString();
      try {
        bot.sendMessage(chatId, content);
      } catch(e) {
        console.log(e);
      }
      
      console.log(content);
      console.log(process.memoryUsage().rss);
      await ch.ack(msg);
    }
  });
})();
