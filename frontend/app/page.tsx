import { fetchCalendarData } from "./lib/api";
import { VendorContract } from "./lib/definitions";
import HomePage from "./_components/home";

// This is the app's homepage. Fetches data from the database and displays UI container.
export default async function Home() {
  const renewalEvents: VendorContract[] = await fetchCalendarData();

  return (
    <>
      <HomePage initialRenewalEvents={renewalEvents} />
    </>
  );
}
