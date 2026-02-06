import { z } from 'zod';
import { insertAccountSchema, accounts, users } from './schema';

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
  accounts: {
    list: {
      method: 'GET' as const,
      path: '/api/accounts',
      responses: {
        200: z.array(z.custom<typeof accounts.$inferSelect>()),
        401: errorSchemas.unauthorized,
      },
    },
    create: {
      method: 'POST' as const,
      path: '/api/accounts',
      input: insertAccountSchema,
      responses: {
        201: z.custom<typeof accounts.$inferSelect>(),
        400: errorSchemas.validation,
        401: errorSchemas.unauthorized,
      },
    },
    get: {
      method: 'GET' as const,
      path: '/api/accounts/:id',
      responses: {
        200: z.custom<typeof accounts.$inferSelect>(),
        404: errorSchemas.notFound,
        401: errorSchemas.unauthorized,
      },
    },
    delete: {
      method: 'DELETE' as const,
      path: '/api/accounts/:id',
      responses: {
        204: z.void(),
        404: errorSchemas.notFound,
        401: errorSchemas.unauthorized,
      },
    },
    start: {
      method: 'POST' as const,
      path: '/api/accounts/:id/start',
      responses: {
        200: z.custom<typeof accounts.$inferSelect>(),
        404: errorSchemas.notFound,
      },
    },
    stop: {
      method: 'POST' as const,
      path: '/api/accounts/:id/stop',
      responses: {
        200: z.custom<typeof accounts.$inferSelect>(),
        404: errorSchemas.notFound,
      },
    },
    validate: {
      method: 'POST' as const,
      path: '/api/accounts/validate',
      responses: {
        200: z.object({ message: z.string() }),
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
    accounts: {
      method: 'GET' as const,
      path: '/api/admin/accounts',
      responses: {
        200: z.array(z.custom<typeof accounts.$inferSelect>()),
        403: errorSchemas.unauthorized,
      },
    },
    authUser: {
      method: 'POST' as const,
      path: '/api/admin/users/:id/auth',
      responses: {
        200: z.custom<typeof users.$inferSelect>(),
        403: errorSchemas.unauthorized,
      },
    },
    validateInstances: {
      method: 'POST' as const,
      path: '/api/admin/accounts/validate',
      responses: {
        200: z.object({ message: z.string() }),
      },
    }
  },
  stats: {
    method: 'GET' as const,
    path: '/api/stats',
    responses: {
      200: z.object({
        activeBots: z.number(),
        totalUsers: z.number(),
        uptime: z.string()
      })
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

export type AccountResponse = z.infer<typeof api.accounts.create.responses[201]>;
export type UserResponse = z.infer<typeof api.auth.me.responses[200]>;
