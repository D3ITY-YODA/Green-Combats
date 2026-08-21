// app/report/page.tsx
"use client";

import { useState } from "react";
import { submitReport } from "@/lib/api/reports-api";
import { Send, CheckCircle2, Droplets, AlertTriangle, Sun, Sprout, MoreHorizontal } from "lucide-react";

const OBSERVATION_OPTIONS = [
  { id: "water", label: "Water has changed", icon: Droplets },
  { id: "flood", label: "Flooding is visible", icon: AlertTriangle },
  { id: "dry", label: "Conditions are unusually dry", icon: Sun },
  { id: "crops", label: "Plants or crops are under stress", icon: Sprout },
  { id: "other", label: "Something else", icon: MoreHorizontal },
];

export default function ReportPage() {
  const [selectedOption, setSelectedOption] = useState<string | null>(null);
  const [details, setDetails] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);

  const handleSubmit = async () => {
    if (!selectedOption) return;
    
    setIsSubmitting(true);
    try {
      await submitReport({
        type: selectedOption,
        description: selectedOption === "other" ? details : undefined,
        place_id: "", // TODO: source from the user's selected place store
        observed_at: new Date().toISOString(),
        location: {
          latitude: -1.29,
          longitude: 36.82,
          accuracy_meters: 10,
        },
      });
      setIsSuccess(true);
    } catch (error) {
      console.error("Failed to submit report", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  // Success State
  if (isSuccess) {
    return (
      <main className="flex min-h-screen flex-col items-center justify-center p-6 text-center">
        <CheckCircle2 className="h-16 w-16 text-status-normal mb-4" />
        <h1 className="text-page font-bold text-forest-deep mb-2">Observation Recorded</h1>
        <p className="text-body text-text-muted max-w-md">
          Thank you for sharing what you are seeing. Your report helps your community stay informed.
        </p>
        <button 
          onClick={() => { setIsSuccess(false); setSelectedOption(null); setDetails(""); }}
          className="mt-8 px-6 py-3 rounded-xl bg-forest text-white font-medium hover:bg-forest-deep transition-colors"
        >
          Submit Another Report
        </button>
      </main>
    );
  }

  return (
    <main className="flex min-h-screen flex-col p-6 md:p-10 max-w-2xl mx-auto">
      {/* Header */}
      <header className="mb-8">
        <h1 className="text-page font-bold text-forest-deep mb-1">Report</h1>
        <p className="text-metadata text-text-muted">
          Share what you are seeing in your area
        </p>
      </header>

      {/* Observation Options (Large tap targets) */}
      <div className="space-y-3 mb-6">
        {OBSERVATION_OPTIONS.map((option) => {
          const isSelected = selectedOption === option.id;
          const Icon = option.icon;
          
          return (
            <button
              key={option.id}
              onClick={() => setSelectedOption(option.id)}
              className={`w-full text-left p-4 rounded-xl border flex items-center gap-4 transition-all ${
                isSelected 
                  ? "border-forest bg-forest/5 text-forest-deep" 
                  : "border-background-stone bg-background hover:bg-background-mist text-text-charcoal"
              }`}
            >
              <Icon className={`h-5 w-5 flex-shrink-0 ${isSelected ? "text-forest" : "text-text-muted"}`} />
              <span className={`text-body font-medium ${isSelected ? "text-forest-deep" : "text-text-charcoal"}`}>
                {option.label}
              </span>
            </button>
          );
        })}
      </div>

      {/* Conditional Text Area for "Something else" */}
      {selectedOption === "other" && (
        <div className="mb-6">
          <label className="block text-card font-medium text-text-charcoal mb-2">
            Please describe what you are seeing:
          </label>
          <textarea
            value={details}
            onChange={(e) => setDetails(e.target.value)}
            rows={4}
            className="w-full rounded-xl border border-background-stone bg-background p-4 text-body text-text-charcoal focus:outline-none focus:ring-2 focus:ring-forest/50 focus:border-forest transition-all"
            placeholder="Type your observation here..."
          />
        </div>
      )}

      {/* Submit Button */}
      <button 
        onClick={handleSubmit}
        disabled={!selectedOption || isSubmitting}
        className="w-full py-4 rounded-xl bg-forest text-white font-medium disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 hover:bg-forest-deep transition-colors"
      >
        {isSubmitting ? (
          "Submitting..."
        ) : (
          <>
            <Send className="h-4 w-4" /> Submit Observation
          </>
        )}
      </button>
    </main>
  );
}
