import { drizzle } from "drizzle-orm/node-postgres";
import pg from "pg";
import * as schema from "@shared/schema";

const { Pool } = pg;

// Aiven Database Configuration
const poolConfig = {
  connectionString: process.env.DATABASE_URL || "postgres://avnadmin:AVNS__v7YcrlbWVN6jtm03JL@hurry-hurry.g.aivencloud.com:22637/defaultdb?sslmode=require",
  ssl: {
    rejectUnauthorized: false // Required for Aiven
  }
};

export const pool = new Pool(poolConfig);
export const db = drizzle(pool, { schema });
