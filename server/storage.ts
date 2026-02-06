import { users, bots, type User, type InsertUser, type Bot, type InsertBot } from "@shared/schema";
import { db } from "./db";
import { eq } from "drizzle-orm";

export interface IStorage {
  getUser(id: number): Promise<User | undefined>;
  getUserByDiscordId(discordId: string): Promise<User | undefined>;
  createUser(user: InsertUser): Promise<User>;
  
  getBots(userId: number): Promise<Bot[]>;
  getBot(id: number): Promise<Bot | undefined>;
  createBot(bot: InsertBot): Promise<Bot>;
  deleteBot(id: number): Promise<void>;
  updateBotStatus(id: number, status: string): Promise<Bot>;
  
  // Admin
  getAllUsers(): Promise<User[]>;
  getAllBots(): Promise<Bot[]>;
}

export class DatabaseStorage implements IStorage {
  async getUser(id: number): Promise<User | undefined> {
    const [user] = await db.select().from(users).where(eq(users.id, id));
    return user;
  }

  async getUserByDiscordId(discordId: string): Promise<User | undefined> {
    const [user] = await db.select().from(users).where(eq(users.discordId, discordId));
    return user;
  }

  async createUser(insertUser: InsertUser): Promise<User> {
    const [user] = await db.insert(users).values(insertUser).returning();
    return user;
  }

  async getBots(userId: number): Promise<Bot[]> {
    return await db.select().from(bots).where(eq(bots.userId, userId));
  }

  async getBot(id: number): Promise<Bot | undefined> {
    const [bot] = await db.select().from(bots).where(eq(bots.id, id));
    return bot;
  }

  async createBot(insertBot: InsertBot): Promise<Bot> {
    const [bot] = await db.insert(bots).values(insertBot).returning();
    return bot;
  }

  async deleteBot(id: number): Promise<void> {
    await db.delete(bots).where(eq(bots.id, id));
  }

  async updateBotStatus(id: number, status: string): Promise<Bot> {
    const [bot] = await db.update(bots)
      .set({ status })
      .where(eq(bots.id, id))
      .returning();
    return bot;
  }

  async getAllUsers(): Promise<User[]> {
    return await db.select().from(users);
  }

  async getAllBots(): Promise<Bot[]> {
    return await db.select().from(bots);
  }
}

export const storage = new DatabaseStorage();
