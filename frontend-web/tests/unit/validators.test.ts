// tests/unit/validators.test.ts

import { describe, it, expect } from 'vitest';
import { SignInSchema, SignUpSchema } from '@/schemas/auth';
import { SubmitReportSchema } from '@/schemas/reports';
import { UpdateFormSchema } from '@/schemas/updates';

describe('Auth Validators', () => {
  describe('SignInSchema', () => {
    it('should validate a correct sign-in payload', () => {
      const validData = {
        phone: '254712345678',
        password: 'SecurePass123!',
      };
      const result = SignInSchema.safeParse(validData);
      expect(result.success).toBe(true);
    });

    it('should reject a phone number with letters', () => {
      const invalidData = {
        phone: '254-ABC-1234',
        password: 'SecurePass123!',
      };
      const result = SignInSchema.safeParse(invalidData);
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toContain('digits');
      }
    });

    it('should reject a password shorter than 8 characters', () => {
      const invalidData = {
        phone: '254712345678',
        password: 'short',
      };
      const result = SignInSchema.safeParse(invalidData);
      expect(result.success).toBe(false);
    });
  });

  describe('SignUpSchema', () => {
    it('should validate a correct sign-up payload', () => {
      const validData = {
        full_name: 'Jane Doe',
        phone: '254712345678',
        password: 'SecurePass123!',
        confirm_password: 'SecurePass123!',
        accept_terms: true,
      };
      const result = SignUpSchema.safeParse(validData);
      expect(result.success).toBe(true);
    });

    it('should reject if passwords do not match', () => {
      const invalidData = {
        full_name: 'Jane Doe',
        phone: '254712345678',
        password: 'SecurePass123!',
        confirm_password: 'DifferentPass123!',
        accept_terms: true,
      };
      const result = SignUpSchema.safeParse(invalidData);
      expect(result.success).toBe(false);
      if (!result.success) {
        // Check that the error is attached to the confirm_password field
        const passwordError = result.error.issues.find(i => i.path.includes('confirm_password'));
        expect(passwordError?.message).toBe('Passwords do not match');
      }
    });

    it('should reject if terms are not accepted', () => {
      const invalidData = {
        full_name: 'Jane Doe',
        phone: '254712345678',
        password: 'SecurePass123!',
        confirm_password: 'SecurePass123!',
        accept_terms: false,
      };
      const result = SignUpSchema.safeParse(invalidData);
      expect(result.success).toBe(false);
    });
  });
});

describe('Report Validators', () => {
  describe('SubmitReportSchema', () => {
    it('should validate a correct community report', () => {
      const validData = {
        type: 'water_change',
        place_id: '123e4567-e89b-12d3-a456-426614174000', // Valid UUID
        location: {
          latitude: -1.2921,
          longitude: 36.8219,
          accuracy_meters: 15,
        },
        observed_at: new Date().toISOString(),
        description: 'The river level has dropped significantly.',
      };
      const result = SubmitReportSchema.safeParse(validData);
      expect(result.success).toBe(true);
    });

    it('should reject invalid latitude coordinates', () => {
      const invalidData = {
        type: 'water_change',
        place_id: '123e4567-e89b-12d3-a456-426614174000',
        location: {
          latitude: 999, // Invalid
          longitude: 36.8219,
        },
        observed_at: new Date().toISOString(),
      };
      const result = SubmitReportSchema.safeParse(invalidData);
      expect(result.success).toBe(false);
    });

    it('should reject descriptions exceeding 5000 characters', () => {
      const invalidData = {
        type: 'water_change',
        place_id: '123e4567-e89b-12d3-a456-426614174000',
        location: { latitude: 0, longitude: 0 },
        observed_at: new Date().toISOString(),
        description: 'a'.repeat(5001),
      };
      const result = SubmitReportSchema.safeParse(invalidData);
      expect(result.success).toBe(false);
    });
  });
});

describe('Update Validators', () => {
  describe('UpdateFormSchema', () => {
    it('should validate a correct institutional update', () => {
      const validData = {
        place_id: '123e4567-e89b-12d3-a456-426614174000',
        topic_key: 'water_outlook',
        type: 'warning',
        title: 'Flood Warning',
        message: 'Heavy rainfall expected.',
        valid_from: new Date().toISOString(),
        valid_until: new Date(Date.now() + 86400000).toISOString(),
        display_on_today: true,
      };
      const result = UpdateFormSchema.safeParse(validData);
      expect(result.success).toBe(true);
    });

    it('should reject titles exceeding 120 characters', () => {
      const invalidData = {
        place_id: '123e4567-e89b-12d3-a456-426614174000',
        topic_key: 'water_outlook',
        type: 'warning',
        title: 'a'.repeat(121),
        message: 'Heavy rainfall expected.',
        valid_from: new Date().toISOString(),
        display_on_today: true,
      };
      const result = UpdateFormSchema.safeParse(invalidData);
      expect(result.success).toBe(false);
    });

    it('should reject invalid update types', () => {
      const invalidData = {
        place_id: '123e4567-e89b-12d3-a456-426614174000',
        topic_key: 'water_outlook',
        type: 'critical_alert', // Not in enum
        title: 'Flood Warning',
        message: 'Heavy rainfall expected.',
        valid_from: new Date().toISOString(),
        display_on_today: true,
      };
      const result = UpdateFormSchema.safeParse(invalidData);
      expect(result.success).toBe(false);
    });
  });
});
