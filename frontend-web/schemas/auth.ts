// schemas/auth.ts

import { z } from "zod";

// --- Sign In Schema ---

export const SignInSchema = z.object({
  phone: z
    .string()
    .min(1, "Phone number is required")
    .regex(/^[0-9]+$/, "Phone number must contain only digits"),
  password: z
    .string()
    .min(8, "Password must be at least 8 characters"),
});

export type SignInFormValues = z.infer<typeof SignInSchema>;

// --- Sign Up Schema ---

export const SignUpSchema = z
  .object({
    full_name: z
      .string()
      .min(2, "Full name must be at least 2 characters")
      .max(100, "Full name is too long"),
    phone: z
      .string()
      .min(1, "Phone number is required")
      .regex(/^[0-9]+$/, "Phone number must contain only digits"),
    password: z
      .string()
      .min(8, "Password must be at least 8 characters")
      .regex(/[A-Z]/, "Password must contain at least one uppercase letter")
      .regex(/[0-9]/, "Password must contain at least one number"),
    confirm_password: z.string(),
    accept_terms: z.boolean().refine((val) => val === true, {
      message: "You must accept the Terms and Privacy Policy",
    }),
  })
  .refine(
    (data) => data.password === data.confirm_password,
    {
      message: "Passwords do not match",
      path: ["confirm_password"], // Attaches the error to the confirm_password field
    }
  );

export type SignUpFormValues = z.infer<typeof SignUpSchema>;

// --- Forgot Password Schema (Optional but recommended for Screen 6) ---

export const ForgotPasswordSchema = z.object({
  phone: z
    .string()
    .min(1, "Phone number is required")
    .regex(/^[0-9]+$/, "Phone number must contain only digits"),
});

export type ForgotPasswordFormValues = z.infer<typeof ForgotPasswordSchema>;
