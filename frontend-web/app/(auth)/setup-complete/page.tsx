import Link from "next/link";
export default function SetupCompletePage() {
  return (
    <div className="text-center space-y-6">
      <div className="mx-auto h-16 w-16 rounded-full bg-status-normal/10 flex items-center justify-center">
        <svg className="h-8 w-8 text-status-normal" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" /></svg>
      </div>
      <h1 className="text-2xl font-bold text-forest-deep">You're all set!</h1>
      <p className="text-text-muted">Your profile is ready. You can now access your personalized environmental updates.</p>
      <Link href="/today" className="block w-full rounded-xl bg-forest py-3 font-medium text-white hover:bg-forest-deep">Go to Dashboard</Link>
    </div>
  );
}
