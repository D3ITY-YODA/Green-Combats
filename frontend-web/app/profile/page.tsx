// app/profile/page.tsx
"use client";

import { useState } from "react";
import { MapPin, Bell, Globe, Shield, HelpCircle, ChevronRight, Plus, Accessibility } from "lucide-react";

export default function ProfilePage() {
  const [language, setLanguage] = useState("English");
  const [notifications, setNotifications] = useState(true);

  return (
    <main className="flex min-h-screen flex-col p-6 md:p-10 max-w-2xl mx-auto">
      {/* Header */}
      <header className="mb-8">
        <h1 className="text-page font-bold text-forest-deep mb-1">Profile</h1>
        <p className="text-metadata text-text-muted">Manage your preferences and support options</p>
      </header>

      <div className="space-y-8">
        {/* 1. Selected Locations */}
        <section>
          <h2 className="text-section font-semibold text-text-charcoal mb-3">Selected Locations</h2>
          <div className="space-y-2">
            <div className="flex items-center justify-between rounded-xl border border-background-stone bg-background p-4">
              <div className="flex items-center gap-3">
                <MapPin className="h-5 w-5 text-forest" />
                <span className="text-body text-text-charcoal">Lower Valley</span>
              </div>
              <span className="text-metadata text-status-normal font-medium">Active</span>
            </div>
            <button className="w-full flex items-center justify-center gap-2 rounded-xl border border-dashed border-background-stone bg-background-mist p-4 text-body text-text-muted hover:bg-background-stone transition-colors">
              <Plus className="h-4 w-4" /> Add another location
            </button>
          </div>
        </section>

        {/* 2. Preferences (Language & Notifications) */}
        <section>
          <h2 className="text-section font-semibold text-text-charcoal mb-3">Preferences</h2>
          <div className="space-y-2">
            {/* Language Selector */}
            <div className="flex items-center justify-between rounded-xl border border-background-stone bg-background p-4">
              <div className="flex items-center gap-3">
                <Globe className="h-5 w-5 text-forest" />
                <span className="text-body text-text-charcoal">Language</span>
              </div>
              <select 
                value={language} 
                onChange={(e) => setLanguage(e.target.value)}
                className="text-body text-text-muted bg-transparent focus:outline-none cursor-pointer text-right"
              >
                <option>English</option>
                <option>French</option>
                <option>Swahili</option>
              </select>
            </div>

            {/* Notifications Toggle */}
            <div className="flex items-center justify-between rounded-xl border border-background-stone bg-background p-4">
              <div className="flex items-center gap-3">
                <Bell className="h-5 w-5 text-forest" />
                <span className="text-body text-text-charcoal">Notifications</span>
              </div>
              <button 
                onClick={() => setNotifications(!notifications)}
                className={`w-12 h-6 rounded-full p-1 transition-colors ${notifications ? 'bg-forest' : 'bg-background-stone'}`}
                aria-label="Toggle notifications"
              >
                <div className={`w-4 h-4 rounded-full bg-white transition-transform ${notifications ? 'translate-x-6' : 'translate-x-0'}`} />
              </button>
            </div>
          </div>
        </section>

        {/* 3. Support (PDF: "Accessibility, Privacy, Help") */}
        <section>
          <h2 className="text-section font-semibold text-text-charcoal mb-3">Support</h2>
          <div className="space-y-2">
            <button className="w-full flex items-center justify-between rounded-xl border border-background-stone bg-background p-4 hover:bg-background-mist transition-colors text-left">
              <div className="flex items-center gap-3">
                <Shield className="h-5 w-5 text-forest" />
                <span className="text-body text-text-charcoal">Privacy & Security</span>
              </div>
              <ChevronRight className="h-4 w-4 text-text-muted" />
            </button>
            
            <button className="w-full flex items-center justify-between rounded-xl border border-background-stone bg-background p-4 hover:bg-background-mist transition-colors text-left">
              <div className="flex items-center gap-3">
                <HelpCircle className="h-5 w-5 text-forest" />
                <span className="text-body text-text-charcoal">Help & About</span>
              </div>
              <ChevronRight className="h-4 w-4 text-text-muted" />
            </button>

            <button className="w-full flex items-center justify-between rounded-xl border border-background-stone bg-background p-4 hover:bg-background-mist transition-colors text-left">
              <div className="flex items-center gap-3">
                <Accessibility className="h-5 w-5 text-forest" />
                <span className="text-body text-text-charcoal">Accessibility</span>
              </div>
              <ChevronRight className="h-4 w-4 text-text-muted" />
            </button>
          </div>
        </section>
      </div>
    </main>
  );
}
