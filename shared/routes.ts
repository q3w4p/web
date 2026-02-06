import { z } from 'zod';
import { insertBotSchema, bots, users } from './schema';

// ============================================
// SHARED ERROR SCHEMAS
// ============================================
export const errorSchemas = {
  validation: z.object({
    message: z.string(),
    field: z.string().optional(),
  }),
  notFound: z.object({
    message: z.string(),
  }),
  internal: z.object({
    message: z.string(),
  }),
  unauthorized: z.object({
    message: z.string(),
  }),
};

// ============================================
// API CONTRACT
// ============================================
export const api = {
  auth: {
    me: {
      method: 'GET' as const,
      path: '/api/user',
      responses: {
        200: z.custom<typeof users.$inferSelect>(),
        401: errorSchemas.unauthorized,
      },
    },
    logout: {
      method: 'POST' as const,
      path: '/api/auth/logout',
      responses: {
        200: z.object({ message: z.string() }),
      },
    },
  },
  bots: {
    list: {
      method: 'GET' as const,
      path: '/api/bots',
      responses: {
        200: z.array(z.custom<typeof bots.$inferSelect>()),
        401: errorSchemas.unauthorized,
      },
    },
    create: {
      method: 'POST' as const,
      path: '/api/bots',
      input: insertBotSchema,
      responses: {
        201: z.custom<typeof bots.$inferSelect>(),
        400: errorSchemas.validation,
        401: errorSchemas.unauthorized,
      },
    },
    get: {
      method: 'GET' as const,
      path: '/api/bots/:id',
      responses: {
        200: z.custom<typeof bots.$inferSelect>(),
        404: errorSchemas.notFound,
        401: errorSchemas.unauthorized,
      },
    },
    delete: {
      method: 'DELETE' as const,
      path: '/api/bots/:id',
      responses: {
        204: z.void(),
        404: errorSchemas.notFound,
        401: errorSchemas.unauthorized,
      },
    },
    start: {
      method: 'POST' as const,
      path: '/api/bots/:id/start',
      responses: {
        200: z.custom<typeof bots.$inferSelect>(),
        404: errorSchemas.notFound,
      },
    },
    stop: {
      method: 'POST' as const,
      path: '/api/bots/:id/stop',
      responses: {
        200: z.custom<typeof bots.$inferSelect>(),
        404: errorSchemas.notFound,
      },
    }
  },
  admin: {
    users: {
      method: 'GET' as const,
      path: '/api/admin/users',
      responses: {
        200: z.array(z.custom<typeof users.$inferSelect>()),
        403: errorSchemas.unauthorized,
      },
    },
    bots: {
      method: 'GET' as const,
      path: '/api/admin/bots',
      responses: {
        200: z.array(z.custom<typeof bots.$inferSelect>()),
        403: errorSchemas.unauthorized,
      },
    }
  }
};

// ============================================
// HELPER
// ============================================
export function buildUrl(path: string, params?: Record<string, string | number>): string {
  let url = path;
  if (params) {
    Object.entries(params).forEach(([key, value]) => {
      if (url.includes(`:${key}`)) {
        url = url.replace(`:${key}`, String(value));
      }
    });
  }
  return url;
}

export type BotResponse = z.infer<typeof api.bots.create.responses[201]>;
export type UserResponse = z.infer<typeof api.auth.me.responses[200]>;
