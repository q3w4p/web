# BotHost - Discord Bot Hosting Platform

## Overview

BotHost is a full-stack web application for hosting and managing Discord bots. Users authenticate via Discord OAuth2, then can create, deploy, start/stop, and delete Discord bot instances from a dashboard. The app includes an admin panel for managing all users and bots across the platform.

The project follows a monorepo structure with three main directories:
- `client/` — React SPA frontend
- `server/` — Express API backend
- `shared/` — Shared types, schemas, and API contract definitions used by both client and server

## User Preferences

Preferred communication style: Simple, everyday language.

## System Architecture

### Frontend (client/)
- **Framework**: React 18 with TypeScript
- **Routing**: Wouter (lightweight client-side router)
- **State Management**: TanStack React Query for server state (caching, mutations, refetching)
- **UI Components**: shadcn/ui (new-york style) built on Radix UI primitives
- **Styling**: Tailwind CSS with CSS variables for theming, dark-mode-first design with a purple/violet accent palette
- **Animations**: Framer Motion for page transitions and interactive elements
- **Forms**: React Hook Form with Zod validation via @hookform/resolvers
- **Build**: Vite with React plugin, outputs to `dist/public`
- **Path aliases**: `@/` maps to `client/src/`, `@shared/` maps to `shared/`

Key pages:
- `/` — Landing page (Home) with Discord login CTA
- `/dashboard` — Authenticated user's bot management
- `/admin` — Admin-only panel for viewing all users and bots

### Backend (server/)
- **Framework**: Express 5 on Node.js with TypeScript (runs via tsx in dev)
- **Session Management**: express-session with MemoryStore (dev), connect-pg-simple available for production
- **Authentication**: Passport.js with Discord OAuth2 strategy (passport-discord)
- **API Pattern**: RESTful JSON API under `/api/` prefix
- **Shared Contract**: API routes and validation schemas defined in `shared/routes.ts`, consumed by both client and server for type safety
- **Build**: esbuild bundles server to `dist/index.cjs` for production

### Database
- **ORM**: Drizzle ORM with PostgreSQL dialect
- **Schema**: Defined in `shared/schema.ts` using Drizzle's `pgTable` definitions
- **Migrations**: Managed via `drizzle-kit push` (schema push approach, not migration files)
- **Connection**: `DATABASE_URL` environment variable required, uses `pg.Pool`
- **Tables**:
  - `users` — id, username, discordId (unique), avatar, isAdmin, createdAt
  - `bots` — id, userId (FK to users), name, token, status (online/offline/error), createdAt
- **Relations**: One-to-many from users to bots
- **Validation**: Zod schemas auto-generated from Drizzle tables via `drizzle-zod`

### Storage Layer
- `server/storage.ts` defines an `IStorage` interface and `DatabaseStorage` implementation
- Methods: getUser, getUserByDiscordId, createUser, getBots, getBot, createBot, deleteBot, updateBotStatus, getAllUsers, getAllBots

### API Routes
Defined as a typed contract in `shared/routes.ts`:
- `GET /api/user` — Get current authenticated user
- `POST /api/auth/logout` — Logout
- `GET /api/auth/discord` — Initiate Discord OAuth
- `GET /api/auth/discord/callback` — Discord OAuth callback
- `GET /api/bots` — List user's bots
- `GET /api/bots/:id` — Get specific bot
- `POST /api/bots` — Create a bot
- `DELETE /api/bots/:id` — Delete a bot
- `POST /api/bots/:id/start` — Start a bot
- `POST /api/bots/:id/stop` — Stop a bot
- `GET /api/admin/users` — Admin: list all users
- `GET /api/admin/bots` — Admin: list all bots

### Dev vs Production
- **Dev**: Vite dev server with HMR served through Express middleware (`server/vite.ts`)
- **Production**: Vite builds static assets to `dist/public`, Express serves them (`server/static.ts`); server bundled with esbuild to `dist/index.cjs`

## External Dependencies

### Required Services
- **PostgreSQL Database**: Connected via `DATABASE_URL` environment variable. Must be provisioned before the app can start.
- **Discord OAuth2 Application**: Requires `DISCORD_CLIENT_ID`, `DISCORD_CLIENT_SECRET`, and optionally `DISCORD_CALLBACK_URL` environment variables for authentication.

### Environment Variables
| Variable | Required | Description |
|---|---|---|
| `DATABASE_URL` | Yes | PostgreSQL connection string |
| `DISCORD_CLIENT_ID` | Yes (for auth) | Discord application client ID |
| `DISCORD_CLIENT_SECRET` | Yes (for auth) | Discord application client secret |
| `DISCORD_CALLBACK_URL` | No | OAuth callback URL (auto-generated from Replit env if not set) |
| `SESSION_SECRET` | No | Session encryption secret (defaults to "dev_secret") |

### Key npm Dependencies
- **Server**: express, passport, passport-discord, drizzle-orm, pg, express-session, memorystore, zod
- **Client**: react, wouter, @tanstack/react-query, framer-motion, react-hook-form, shadcn/ui (Radix primitives), tailwindcss, lucide-react, react-icons
- **Shared**: drizzle-zod, zod
- **Build**: vite, esbuild, tsx, drizzle-kit