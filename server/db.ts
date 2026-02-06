import { drizzle } from "drizzle-orm/node-postgres";
import pg from "pg";
import * as schema from "@shared/schema";

const { Pool } = pg;

// Aiven Database Configuration
const poolConfig = {
  host: 'hurry-hurry.g.aivencloud.com',
  port: 22637,
  database: 'defaultdb',
  user: 'avnadmin',
  password: 'AVNS__v7YcrlbWVN6jtm03JL',
  ssl: {
    rejectUnauthorized: false // Required for Aiven
  }
};

export const pool = new Pool(poolConfig);
export const db = drizzle(pool, { schema });
