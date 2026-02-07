import { users, accounts, type User, type InsertUser, type Account, type InsertAccount } from "@shared/schema";
import { db } from "./db";
import { eq } from "drizzle-orm";

export interface IStorage {
  getUser(id: number): Promise<User | undefined>;
  getUserByDiscordId(discordId: string): Promise<User | undefined>;
  createUser(user: InsertUser): Promise<User>;
  updateUser(id: number, user: Partial<User>): Promise<User>;
  updateUserAuth(id: number, isAuthed: boolean): Promise<User>;
  
  getAccounts(userId: number): Promise<Account[]>;
  getAccount(id: number): Promise<Account | undefined>;
  createAccount(account: InsertAccount): Promise<Account>;
  deleteAccount(id: number): Promise<void>;
  updateAccountStatus(id: number, status: string, pid?: number | null): Promise<Account>;
  updateAccountDetails(id: number, details: Partial<Account>): Promise<Account>;
  
  // Admin
  getAllUsers(): Promise<User[]>;
  getAllAccounts(): Promise<Account[]>;
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

  async updateUser(id: number, userUpdate: Partial<User>): Promise<User> {
    const [user] = await db.update(users).set(userUpdate).where(eq(users.id, id)).returning();
    return user;
  }

  async updateUserAuth(id: number, isAuthed: boolean): Promise<User> {
    const [user] = await db.update(users).set({ isAuthed }).where(eq(users.id, id)).returning();
    return user;
  }

  async getAccounts(userId: number): Promise<Account[]> {
    return await db.select().from(accounts).where(eq(accounts.userId, userId));
  }

  async getAccount(id: number): Promise<Account | undefined> {
    const [account] = await db.select().from(accounts).where(eq(accounts.id, id));
    return account;
  }

  async createAccount(insertAccount: InsertAccount): Promise<Account> {
    const [account] = await db.insert(accounts).values(insertAccount).returning();
    return account;
  }

  async deleteAccount(id: number): Promise<void> {
    await db.delete(accounts).where(eq(accounts.id, id));
  }

  async updateAccountStatus(id: number, status: string, pid?: number | null): Promise<Account> {
    const [account] = await db.update(accounts)
      .set({ status, pid: pid !== undefined ? pid : undefined })
      .where(eq(accounts.id, id))
      .returning();
    return account;
  }

  async updateAccountDetails(id: number, details: Partial<Account>): Promise<Account> {
    const [account] = await db.update(accounts)
      .set(details)
      .where(eq(accounts.id, id))
      .returning();
    return account;
  }

  async getAllUsers(): Promise<User[]> {
    return await db.select().from(users);
  }

  async getAllAccounts(): Promise<Account[]> {
    return await db.select().from(accounts);
  }
}

export const storage = new DatabaseStorage();
