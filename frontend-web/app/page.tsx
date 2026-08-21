import { redirect } from "next/navigation";

export default function HomePage() {
  // Redirect to the marketing landing page or sign-in
  redirect("/sign-in");
}
