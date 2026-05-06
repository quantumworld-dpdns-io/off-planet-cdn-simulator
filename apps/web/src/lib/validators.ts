import { z } from "zod";

export const CreateSiteSchema = z.object({
  name: z.string().min(1).max(100),
  location: z.string().optional(),
  description: z.string().optional(),
});

export const CreateCacheObjectSchema = z.object({
  name: z.string().min(1),
  source_url: z.string().url(),
  priority_class_id: z.string().uuid(),
  site_id: z.string().uuid(),
  size_bytes: z.number().positive(),
  tags: z.array(z.string()).default([]),
});

export type CreateSite = z.infer<typeof CreateSiteSchema>;
export type CreateCacheObject = z.infer<typeof CreateCacheObjectSchema>;
