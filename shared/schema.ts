import { pgTable, text, serial, integer, boolean, timestamp } from "drizzle-orm/pg-core";
import { relations } from "drizzle-orm";
import { createInsertSchema } from "drizzle-zod";
import { z } from "zod";

// === TABLE DEFINITIONS ===
export const users = pgTable("users", {
  id: serial("id").primaryKey(),
  username: text("username").notNull(),
  discordId: text("discord_id").notNull().unique(),
  avatar: text("avatar"),
  isAdmin: boolean("is_admin").default(false),
  isAuthed: boolean("is_authed").default(false),
  createdAt: timestamp("created_at").defaultNow(),
});

export const accounts = pgTable("accounts", {
  id: serial("id").primaryKey(),
  userId: integer("user_id").notNull(), // Foreign key to users
  token: text("token").notNull(),
  prefix: text("prefix").default("!"),
  discordUsername: text("discord_username"),
  discordAvatar: text("discord_avatar"),
  guildsCount: integer("guilds_count").default(0),
  friendsCount: integer("friends_count").default(0),
  status: text("status").default("offline"), // online, offline, idle, dnd, invisible
  pid: integer("pid"),
  createdAt: timestamp("created_at").defaultNow(),
});

// === RELATIONS ===
export const usersRelations = relations(users, ({ many }) => ({
  accounts: many(accounts),
}));

export const accountsRelations = relations(accounts, ({ one }) => ({
  user: one(users, {
    fields: [accounts.userId],
    references: [users.id],
  }),
}));

// === BASE SCHEMAS ===
export const insertUserSchema = createInsertSchema(users).omit({ id: true, createdAt: true });
export const insertAccountSchema = createInsertSchema(accounts).omit({ id: true, createdAt: true, status: true, pid: true, discordUsername: true, discordAvatar: true, guildsCount: true, friendsCount: true });

// === EXPLICIT API CONTRACT TYPES ===
export type User = typeof users.$inferSelect;
export type InsertUser = z.infer<typeof insertUserSchema>;

export type Account = typeof accounts.$inferSelect;
export type InsertAccount = z.infer<typeof insertAccountSchema>;

export type CreateAccountRequest = InsertAccount;
export type UpdateAccountRequest = Partial<InsertAccount>;

// Admin view types
export type UserWithAccounts = User & { accounts: Account[] };
