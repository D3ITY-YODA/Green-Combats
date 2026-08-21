// schemas/updates.ts

import { z } from "zod";

export const UpdateFormSchema = z.object({
  place_id: z
    .string()
    .uuid("Please select a valid place."),
    
  topic_key: z
    .string()
    .min(1, "Topic is required."),
    
  type: z.enum(
    [
      "information",
      "important",
      "warning",
      "emergency",
      "community",
    ],
  ),

  title: z
    .string()
    .min(1, "Title is required.")
    .max(120, "Title cannot exceed 120 characters."),
    
  message: z
    .string()
    .min(1, "Message is required.")
    .max(1000, "Message cannot exceed 1000 characters."),
    
  valid_from: z
    .string()
    .datetime("Invalid start date and time."),
    
  valid_until: z
    .string()
    .datetime("Invalid end date and time.")
    .optional(),
    
  source_name: z
    .string()
    .max(200, "Source name cannot exceed 200 characters.")
    .optional(),
    
  display_on_today: z
    .boolean(),
});

// Infer the TypeScript type for use in React Hook Form
export type UpdateFormValues = z.infer<typeof UpdateFormSchema>;
