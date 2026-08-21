// lib/api/reports-api.ts

export interface ReportPayload {
  observationType: string;
  details?: string;
  location: string;
}

export interface ReportResponse {
  success: boolean;
  message: string;
}

const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export async function submitReport(payload: ReportPayload): Promise<ReportResponse> {
  await delay(500);
  
  console.log("Mock API: Report submitted to backend:", payload);
  
  return {
    success: true,
    message: "Observation recorded successfully.",
  };
}
