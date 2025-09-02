import { fetchCalendarData } from "./lib/api";
import { VendorContract } from "./lib/definitions";
import HomePage from "./_components/home";

export default async function Home() {
  let renewalEvents: VendorContract[] = await fetchCalendarData();

  return (
    <>
      <HomePage initialRenewalEvents={renewalEvents} />
    </>
  );
}
