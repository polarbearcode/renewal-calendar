import UploadButton from "@/app/_components/upload-button";
import RenewalCalendar from "@/app/_components/renewal-calendar";

export default function Home() {
  const renewalEvents = [
    {
      vendor: "Vendor A",
      effectiveDate: "2023-01-01",
      renewalDate: "2025-08-29",
      autoRenew: true,
      amount: 1000,
    },
    {
      vendor: "Vendor B",
      effectiveDate: "2023-02-01",
      renewalDate: "2024-02-01",
      autoRenew: false,
      amount: 2000,
    },
  ];

  return (
    <>
      <h1 className="mt-8 mb-4 text-center text-3xl font-bold">
        Renewal Calendar
      </h1>
      <div className="flex flex-col items-center justify-center min-h-[10vh]">
        <UploadButton />
      </div>

      <div>
        <RenewalCalendar renewalEvents={renewalEvents} />
      </div>
    </>
  );
}
