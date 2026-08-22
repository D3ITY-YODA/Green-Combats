import { z } from "zod";
export const UpdateFormSchema = z.object({
  place_id: z.string().uuid(),
  topic_key: z.string(),
  type: z.enum(["information", "important", "warning", "emergency", "community"]),
  title: z.string().min(1).max(120),
  message: z.string().min(1).max(1000),
  valid_from: z.string(),
  valid_until: z.string().optional(),
  display_on_today: z.boolean(),
});
