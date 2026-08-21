// tests/e2e/onboarding.spec.ts

import { test, expect, Page } from '@playwright/test';

// --- Helper Functions ---

async function fillSignUpForm(page: Page) {
  await page.getByLabel(/full name/i).fill('Jane Doe');
  await page.getByLabel(/phone number/i).fill('1234567890');
  await page.getByLabel(/password/i).fill('SecurePass123!');
  await page.getByLabel(/confirm password/i).fill('SecurePass123!');
  
  // Accept terms checkbox
  await page.getByRole('checkbox', { name: /terms and privacy policy/i }).check();
  
  await page.getByRole('button', { name: /sign up/i }).click();
}

// --- Public Onboarding Flow ---

test.describe('Public Onboarding', () => {
  test('successfully completes the public user onboarding flow', async ({ page }) => {
    // 1. User opens landing page
    await page.goto('/');
    await expect(page.getByRole('heading', { name: /green compass/i })).toBeVisible();
    
    // 2. Selects language (via "Get started" or "Choose language")
    await page.getByRole('link', { name: /get started/i }).click();
    await expect(page).toHaveURL(/.*\/choose-language/);
    await page.getByRole('button', { name: /english/i }).click();
    await page.getByRole('button', { name: /continue/i }).click();

    // 3. Creates account
    await expect(page).toHaveURL(/.*\/sign-up/);
    await fillSignUpForm(page);

    // 4. Selects personal account
    await expect(page).toHaveURL(/.*\/account-choice/);
    await page.getByRole('button', { name: /personal account/i }).click();

    // 5. Selects place
    await expect(page).toHaveURL(/.*\/choose-place/);
    await page.getByPlaceholder(/search for a place/i).fill('Lower Valley');
    // Assuming a dropdown or list appears
    await page.getByRole('option', { name: /lower valley/i }).click(); 
    await page.getByRole('button', { name: /continue/i }).click();

    // 6. Skips optional interests
    await expect(page).toHaveURL(/.*\/interests/);
    await page.getByRole('button', { name: /skip for now/i }).click();

    // 7. Sets notification preferences
    await expect(page).toHaveURL(/.*\/notifications/);
    await page.getByRole('button', { name: /continue/i }).click();

    // 8. Completes setup & 9. Reaches Today
    await expect(page).toHaveURL(/.*\/today/);
    await expect(page.getByRole('heading', { name: /today/i })).toBeVisible();
    await expect(page.getByText(/lower valley/i)).toBeVisible();
  });
});

// --- Organization Onboarding Flow ---

test.describe('Organization Onboarding', () => {
  // Note: For E2E tests, it's best practice to mock the backend API responses 
  // or use a seeded test database to ensure the organization exists.
  
  test('successfully completes the organization onboarding flow', async ({ page }) => {
    // Start from the account choice page (simulating post-signup)
    await page.goto('/account-choice');
    
    // 1. User chooses organization path
    await page.getByRole('button', { name: /organization account/i }).click();

    // 2. Searches organization
    await expect(page).toHaveURL(/.*\/organization-search/);
    await page.getByPlaceholder(/search organizations/i).fill('Lower Valley Water Authority');
    
    // Mocking the API response for the search to ensure test stability
    await page.route('**/api/v1/organizations/search', route => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: [{ id: 'org-123', name: 'Lower Valley Water Authority', type: 'water_authority' }]
        }),
      });
    });
    
    await page.getByRole('option', { name: /lower valley water authority/i }).click();

    // 3. Requests access
    await page.getByRole('button', { name: /request access/i }).click();
    
    // 4. Receives invitation & 5. Accepts invitation
    // In a real E2E flow, this would involve checking a test email or using a seeded DB.
    // Here we simulate navigating to the invitation acceptance URL.
    await page.goto('/organization-invitation?token=mock-invite-token-123');
    await expect(page.getByText(/you have been invited/i)).toBeVisible();
    await page.getByRole('button', { name: /accept invitation/i }).click();

    // 6. Sees organization context & 7. Opens console
    await expect(page).toHaveURL(/.*\/console/);
    await expect(page.getByRole('heading', { name: /overview/i })).toBeVisible();
    await expect(page.getByText(/lower valley water authority/i)).toBeVisible();
  });
});
