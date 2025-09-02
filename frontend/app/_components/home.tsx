"use client";
import { useState } from "react";
import { VendorContract } from "../lib/definitions";
import RenewalCalendar from "./renewal-calendar";
import UploadButton from "./upload-button";

export default function HomePage({
  initialRenewalEvents,
}: {
  initialRenewalEvents: VendorContract[];
}) {
  const [renewalEvents, setRenewalEvents] =
    useState<VendorContract[]>(initialRenewalEvents);

  return (
    <>
      <h1 className="mt-8 mb-4 text-center text-3xl font-bold">
        Renewal Calendar
      </h1>
      <div className="flex flex-col items-center justify-center min-h-[10vh]">
        <UploadButton setRenewalEvents={setRenewalEvents} />
      </div>

      <div>
        <RenewalCalendar renewalEvents={renewalEvents} />
      </div>
    </>
  );
}
