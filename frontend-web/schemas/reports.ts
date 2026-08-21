// schemas/reports.ts

import { z } from "zod";

export const SubmitReportSchema = z.object({
  type: z.enum(
    [
      "water_change",
      "flooding_visible",
      "unusually_dry",
      "vegetation_stress",
      "heat_impact",
      "infrastructure_change",
      "incorrect_information",
      "other",
    ],
    {
      required_error: "Please select what you are seeing.",
    }
  ),
  
  place_id: z.string().uuid("Invalid place selected."),
  
  location: z.object({
    latitude: z
      .number()
      .min(-90, "Invalid latitude.")
      .max(90, "Invalid latitude."),
    longitude: z
      .number()
      .min(-180, "Invalid longitude.")
      .max(180, "Invalid longitude."),
    accuracy_meters: z
      .number()
      .positive("Accuracy must be a positive number.")
      .optional(),
  }),
  
  observed_at: z.string().datetime("Invalid observation time."),
  
  description: z
    .string()
    .max(5000, "Description cannot exceed 5000 characters.")
    .optional(),
});

// Infer the TypeScript type from the schema for use in React Hook Form
export type SubmitReportFormValues = z.infer<typeof SubmitReportSchema>;
