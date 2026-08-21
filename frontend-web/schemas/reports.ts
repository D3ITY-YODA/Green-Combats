import { z } from "zod";

export const SubmitReportSchema = z.object({
  type: z.enum([
    "water_change", "flooding_visible", "unusually_dry", 
    "vegetation_stress", "heat_impact", "infrastructure_change", 
    "incorrect_information", "other",
  ]),
  place_id: z.string(),
  location: z.object({
    latitude: z.number().min(-90).max(90),
    longitude: z.number().min(-180).max(180),
  }),
  observed_at: z.string(),
  description: z.string().max(5000).optional(),
});

export type ReportFormValues = z.infer<typeof SubmitReportSchema>;
