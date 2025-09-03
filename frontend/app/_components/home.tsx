"use client";

import { useState } from "react";
import { VendorContract } from "../lib/definitions";
import RenewalCalendar from "./renewal-calendar";
import UploadButton from "./upload-button";
import RenewalTable from "./renewal-table";

/**
 * HomePage component for the Renewal Calendar app.
 *
 * Displays the calendar, upload button, and renewal list.
 * Handles loading and error states, and updates renewal events when files are uploaded.
 *
 * Props:
 *   initialRenewalEvents: VendorContract[] - Initial array of renewal events to display.
 */
export default function HomePage({
  initialRenewalEvents,
}: {
  initialRenewalEvents: VendorContract[];
}) {
  const [renewalEvents, setRenewalEvents] =
    useState<VendorContract[]>(initialRenewalEvents);

  const [loading, setLoading] = useState(false);

  const [dataError, setDataError] = useState<string | null>(null);

  return (
    <>
      <h1 className="mt-8 mb-4 text-center text-3xl font-bold">
        Renewal Calendar
      </h1>
      <div className="flex flex-col items-center justify-center min-h-[10vh]">
        <UploadButton
          setRenewalEvents={setRenewalEvents}
          setLoading={setLoading}
          setDataError={setDataError}
        />
      </div>

      {dataError && (
        <div className="mx-auto my-4 max-w-lg p-4 bg-red-100 text-red-800 border border-red-300 rounded-lg shadow">
          {dataError}
        </div>
      )}

      {loading ? (
        <div>Loading...</div>
      ) : (
        <div>
          <RenewalCalendar renewalEvents={renewalEvents} />
          <RenewalTable renewalEvents={renewalEvents} />
        </div>
      )}
    </>
  );
}
