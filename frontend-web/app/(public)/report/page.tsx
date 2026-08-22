"use client";

import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { SubmitReportSchema, type ReportFormValues } from "@/schemas/reports";
import { submitObservation } from "@/lib/api/observations";

export default function ReportPage() {
  const [isSuccess, setIsSuccess] = useState(false);
  const [submitError, setSubmitError] = useState("");
  const [location, setLocation] = useState<{ lat: number; lon: number } | null>(null);

  const { register, handleSubmit, formState: { errors, isSubmitting }, setValue } = useForm<ReportFormValues>({
    resolver: zodResolver(SubmitReportSchema),
    defaultValues: {
      type: "water_change",
      description: "",
      place_id: "1",
      location: { latitude: 0, longitude: 0 },
      observed_at: new Date().toISOString(),
    },
  });

  // Try to get geolocation on mount
  useEffect(() => {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (pos) => {
          const loc = { lat: pos.coords.latitude, lon: pos.coords.longitude };
          setLocation(loc);
          setValue("location", { latitude: loc.lat, longitude: loc.lon });
        },
        () => { /* Geolocation denied — user can still submit */ },
        { enableHighAccuracy: true, timeout: 5000 },
      );
    }
  }, [setValue]);

  const onSubmit = async (data: ReportFormValues) => {
    setSubmitError("");
    try {
      await submitObservation({
        category: mapReportTypeToCategory(data.type),
        description: data.description || "",
        place_id: data.place_id || undefined,
        lat: data.location.latitude || undefined,
        lon: data.location.longitude || undefined,
      });
      setIsSuccess(true);
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : "Submission failed. Please try again.");
    }
  };

  if (isSuccess) {
    return (
      <div className="mx-auto max-w-2xl px-4 py-16 text-center">
        <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-full bg-status-normal/10">
          <svg className="h-8 w-8 text-status-normal" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" /></svg>
        </div>
        <h1 className="text-2xl font-bold text-text-charcoal">Thank you!</h1>
        <p className="mt-2 text-text-muted">Your update has been received. It will be reviewed by the relevant local team.</p>
        <button onClick={() => setIsSuccess(false)} className="mt-8 rounded-xl bg-forest px-5 py-3 font-medium text-white hover:bg-forest-deep">Submit another update</button>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-2xl px-4 py-8">
      <header className="mb-8">
        <p className="text-sm text-text-muted">Report an update</p>
        <h1 className="mt-1 text-3xl font-semibold text-text-charcoal">What are you seeing?</h1>
        {location && (
          <p className="mt-2 text-xs text-status-normal">📍 Location detected ({location.lat.toFixed(4)}, {location.lon.toFixed(4)})</p>
        )}
      </header>
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        <fieldset>
          <legend className="text-base font-semibold text-text-charcoal mb-3">Select the type of update</legend>
          <div className="grid gap-3">
            {[
              { value: "water_change", label: "Water has changed" },
              { value: "flooding_visible", label: "Flooding is visible" },
              { value: "unusually_dry", label: "Conditions are unusually dry" },
              { value: "vegetation_stress", label: "Plants or crops are under stress" },
              { value: "other", label: "Something else" }
            ].map((opt) => (
              <label key={opt.value} className="flex cursor-pointer items-center gap-3 rounded-xl border border-background-stone bg-background p-4 hover:bg-background-mist has-[:checked]:border-forest has-[:checked]:bg-forest/5">
                <input type="radio" value={opt.value} {...register("type")} className="h-4 w-4 text-forest focus:ring-forest" />
                <span className="text-base font-medium text-text-charcoal">{opt.label}</span>
              </label>
            ))}
          </div>
          {errors.type && <p className="mt-2 text-sm text-status-emergency">{errors.type.message}</p>}
        </fieldset>

        <label className="block">
          <span className="text-sm font-medium text-text-charcoal">Tell us more (optional)</span>
          <textarea {...register("description")} className="mt-2 min-h-32 w-full rounded-xl border border-background-stone bg-background p-3 text-text-charcoal focus:border-forest focus:ring-1 focus:ring-forest" placeholder="Describe what you observed..." />
          {errors.description && <p className="mt-2 text-sm text-status-emergency">{errors.description.message}</p>}
        </label>

        {submitError && (
          <p className="text-sm text-status-emergency">{submitError}</p>
        )}

        <button type="submit" disabled={isSubmitting} className="w-full rounded-xl bg-forest px-5 py-3 font-medium text-white hover:bg-forest-deep disabled:opacity-50">
          {isSubmitting ? "Sending..." : "Send update"}
        </button>
      </form>
    </div>
  );
}

/** Map the report form type to backend observation category. */
function mapReportTypeToCategory(type: string): "flood" | "drought" | "water_quality" | "crop_damage" | "air_quality" | "other" {
  const map: Record<string, "flood" | "drought" | "water_quality" | "crop_damage" | "air_quality" | "other"> = {
    water_change: "water_quality",
    flooding_visible: "flood",
    unusually_dry: "drought",
    vegetation_stress: "crop_damage",
    other: "other",
  };
  return map[type] || "other";
}
