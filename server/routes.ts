import type { Express } from "express";
import { createServer, type Server } from "http";
import { storage } from "./storage";
import { api } from "@shared/routes";
import { z } from "zod";
import passport from "passport";
import { Strategy as DiscordStrategy } from "passport-discord";

export async function registerRoutes(
  httpServer: Server,
  app: Express
): Promise<Server> {
  
  // Auth Setup
  const DISCORD_CLIENT_ID = process.env.DISCORD_CLIENT_ID;
  const DISCORD_CLIENT_SECRET = process.env.DISCORD_CLIENT_SECRET;
  const DISCORD_CALLBACK_URL = process.env.DISCORD_CALLBACK_URL || "https://" + process.env.REPL_SLUG + "." + process.env.REPL_OWNER + ".repl.co/api/auth/discord/callback";

  if (DISCORD_CLIENT_ID && DISCORD_CLIENT_SECRET) {
    passport.use(new DiscordStrategy({
      clientID: DISCORD_CLIENT_ID,
      clientSecret: DISCORD_CLIENT_SECRET,
      callbackURL: DISCORD_CALLBACK_URL,
      scope: ['identify', 'email']
    }, async (accessToken, refreshToken, profile, done) => {
      try {
        let user = await storage.getUserByDiscordId(profile.id);
        if (!user) {
          user = await storage.createUser({
            username: profile.username,
            discordId: profile.id,
            avatar: `https://cdn.discordapp.com/avatars/${profile.id}/${profile.avatar}.png`,
          });
        }
        return done(null, user);
      } catch (err) {
        return done(err as Error, undefined);
      }
    }));

    passport.serializeUser((user: any, done) => {
      done(null, user.id);
    });

    passport.deserializeUser(async (id: number, done) => {
      try {
        const user = await storage.getUser(id);
        done(null, user);
      } catch (err) {
        done(err, null);
      }
    });

    app.use(passport.initialize());
    app.use(passport.session());

    // Auth Routes
    app.get('/api/auth/discord', passport.authenticate('discord'));

    app.get('/api/auth/discord/callback', 
      passport.authenticate('discord', { failureRedirect: '/' }),
      (req, res) => {
        res.redirect('/dashboard');
      }
    );

    app.post('/api/auth/logout', (req, res, next) => {
      req.logout((err) => {
        if (err) return next(err);
        res.json({ message: "Logged out" });
      });
    });

    app.get('/api/user', (req, res) => {
      if (!req.isAuthenticated()) return res.status(401).json({ message: "Unauthorized" });
      res.json(req.user);
    });

  } else {
    console.warn("Discord Auth not configured. Missing DISCORD_CLIENT_ID or DISCORD_CLIENT_SECRET.");
  }

  // API Routes
  const requireAuth = (req: any, res: any, next: any) => {
    if (req.isAuthenticated()) return next();
    res.status(401).json({ message: "Unauthorized" });
  };

  const requireAdmin = (req: any, res: any, next: any) => {
    if (req.isAuthenticated() && req.user.isAdmin) return next();
    res.status(403).json({ message: "Forbidden" });
  };

  // Bots
  app.get(api.bots.list.path, requireAuth, async (req, res) => {
    const bots = await storage.getBots((req.user as any).id);
    res.json(bots);
  });

  app.post(api.bots.create.path, requireAuth, async (req, res) => {
    try {
      const input = api.bots.create.input.parse(req.body);
      const bot = await storage.createBot({
        ...input,
        userId: (req.user as any).id,
      });
      res.status(201).json(bot);
    } catch (err) {
      if (err instanceof z.ZodError) {
        return res.status(400).json({
          message: err.errors[0].message,
          field: err.errors[0].path.join('.'),
        });
      }
      throw err;
    }
  });

  app.get(api.bots.get.path, requireAuth, async (req, res) => {
    const bot = await storage.getBot(Number(req.params.id));
    if (!bot) return res.status(404).json({ message: "Not found" });
    if (bot.userId !== (req.user as any).id && !(req.user as any).isAdmin) {
      return res.status(403).json({ message: "Forbidden" });
    }
    res.json(bot);
  });

  app.delete(api.bots.delete.path, requireAuth, async (req, res) => {
    const bot = await storage.getBot(Number(req.params.id));
    if (!bot) return res.status(404).json({ message: "Not found" });
    if (bot.userId !== (req.user as any).id && !(req.user as any).isAdmin) {
      return res.status(403).json({ message: "Forbidden" });
    }
    await storage.deleteBot(Number(req.params.id));
    res.status(204).send();
  });

  // Start/Stop (Simulation)
  app.post(api.bots.start.path, requireAuth, async (req, res) => {
    const bot = await storage.getBot(Number(req.params.id));
    if (!bot) return res.status(404).json({ message: "Not found" });
    if (bot.userId !== (req.user as any).id && !(req.user as any).isAdmin) {
      return res.status(403).json({ message: "Forbidden" });
    }
    const updated = await storage.updateBotStatus(bot.id, "online");
    res.json(updated);
  });

  app.post(api.bots.stop.path, requireAuth, async (req, res) => {
    const bot = await storage.getBot(Number(req.params.id));
    if (!bot) return res.status(404).json({ message: "Not found" });
    if (bot.userId !== (req.user as any).id && !(req.user as any).isAdmin) {
      return res.status(403).json({ message: "Forbidden" });
    }
    const updated = await storage.updateBotStatus(bot.id, "offline");
    res.json(updated);
  });

  // Admin Routes
  app.get(api.admin.users.path, requireAdmin, async (req, res) => {
    const users = await storage.getAllUsers();
    res.json(users);
  });

  app.get(api.admin.bots.path, requireAdmin, async (req, res) => {
    const bots = await storage.getAllBots();
    res.json(bots);
  });

  return httpServer;
}
