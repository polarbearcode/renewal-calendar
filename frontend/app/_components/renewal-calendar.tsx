"use client";

import dayGridPlugin from "@fullcalendar/daygrid";
import FullCalendar from "@fullcalendar/react";
import { VendorContract } from "../lib/definitions";

/**
 * RenewalCalendar component displays a calendar view of vendor contract renewal events using FullCalendar.
 *
 * Props:
 *   renewalEvents: Array<VendorContract> - Array of contract objects to display as calendar events.
 */
// Component for the renewal calendar display
export default function RenewalCalendar({
  renewalEvents,
}: {
  renewalEvents: Array<VendorContract>;
}) {
  return (
    <>
      <div className="calendar-wrapper w-[900px] mx-auto border rounded-lg shadow p-4 m-2 [&_.fc]:w-full [&_.fc-scrollgrid]:w-full">
        <FullCalendar
          plugins={[dayGridPlugin]}
          initialView="dayGridMonth"
          height="500px"
          events={renewalEvents.map((contract: VendorContract) => {
            return {
              title: contract.seller,
              date: new Date(contract.renewal_date).toISOString().split("T")[0],
            };
          })}
        />
      </div>
    </>
  );
}
