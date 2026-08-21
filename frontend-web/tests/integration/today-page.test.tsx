// tests/integration/today-page.test.tsx

import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { TodayStatus } from '@/components/today/today-status';
import { UpdateCard } from '@/components/updates/update-card';
import type { PublicUpdate, DataStatus } from '@/types/common'; // Adjust paths to your types

// --- Mock Data ---

const mockNormalUpdate: PublicUpdate = {
  id: 'update-1',
  title: 'Good morning',
  message: "Here's what's happening in Lower Valley.",
  type: 'information',
  priority: 'normal',
  place_name: 'Lower Valley',
  updated_at: '2023-10-25T08:00:00Z',
  status: 'published',
  display_on_today: true,
  display_in_feed: true,
  locale: 'en',
  topic_key: 'local_outlook',
  valid_from: '2023-10-25T08:00:00Z',
};

const mockImportantUpdate: PublicUpdate = {
  ...mockNormalUpdate,
  id: 'update-2',
  title: 'Water level rising',
  message: 'Water level at River Njoro is rising slowly. Be observant of risk.',
  type: 'warning',
  priority: 'important',
};

const mockTodayData = {
  place: { id: 'place-1', name: 'Lower Valley', type: 'region', country_code: 'KE' },
  status: {
    title: 'Partly cloudy',
    message: '24°C. Rain likely later.',
    data_status: 'current' as DataStatus,
    updated_at: '2023-10-25T08:00:00Z',
  },
  updates: [mockNormalUpdate, mockImportantUpdate],
  sections: [],
  updated_at: '2023-10-25T08:00:00Z',
};

// --- Helper to render the "Page" ---
// In a real app, this might be a client-side wrapper or the components composed together.
function renderTodayPage(data: typeof mockTodayData) {
  return render(
    <div>
      <TodayStatus
          title={data.status.title}
          message={data.status.message}
          dataStatus={data.status.data_status}
          updatedAt={data.status.updated_at}
        />
      <section>
        {data.updates.map((update) => (
          <UpdateCard key={update.id} update={update} />
        ))}
      </section>
    </div>
  );
}

// --- Test Cases (Section 22) ---

describe('Today Page Integration', () => {
  it('renders the normal state correctly', () => {
    renderTodayPage(mockTodayData);

    // Check TodayStatus
    expect(screen.getByText('Partly cloudy')).toBeInTheDocument();
    expect(screen.getByText('24°C. Rain likely later.')).toBeInTheDocument();
    
    // Check Place Name context (if rendered in header, here we check update context)
    expect(screen.getByText('Lower Valley')).toBeInTheDocument();
  });

  it('renders an important update with visual indicators', () => {
    renderTodayPage(mockTodayData);

    // The important update should be visible
    expect(screen.getByText('Water level rising')).toBeInTheDocument();
    expect(screen.getByText(/Water level at River Njoro/)).toBeInTheDocument();
    
    // Check for the "Read update" button
    const readButtons = screen.getAllByText('Read update');
    expect(readButtons.length).toBeGreaterThan(0);
  });

  it('renders multiple updates in a list', () => {
    renderTodayPage(mockTodayData);

    // Both updates should be present
    expect(screen.getByText('Good morning')).toBeInTheDocument();
    expect(screen.getByText('Water level rising')).toBeInTheDocument();
  });

  it('renders the delayed state when data_status is delayed', () => {
    const delayedData = {
      ...mockTodayData,
      status: {
        ...mockTodayData.status,
        data_status: 'delayed' as DataStatus,
        title: 'Information delayed',
        message: 'The latest update for this place is not available yet.',
      },
      updates: [], // Usually no updates if delayed
    };

    renderTodayPage(delayedData);

    expect(screen.getByText('Information delayed')).toBeInTheDocument();
    expect(screen.getByText(/not available yet/)).toBeInTheDocument();
    expect(screen.getByText(/Last reliable update/)).toBeInTheDocument();
  });

  it('does not render irrelevant topics or updates if none exist', () => {
    const emptyData = {
      ...mockTodayData,
      updates: [],
    };

    renderTodayPage(emptyData);

    // Should show status but no updates
    expect(screen.getByText('Partly cloudy')).toBeInTheDocument();
    expect(screen.queryByText('Good morning')).not.toBeInTheDocument();
    expect(screen.queryByText('Water level rising')).not.toBeInTheDocument();
  });
});
