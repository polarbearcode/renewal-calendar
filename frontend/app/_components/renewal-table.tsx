import { VendorContract } from "../lib/definitions";

/**
 * RenewalList component displays a table of vendor contract renewal events.
 *
 * Props:
 *   renewalEvents: Array<VendorContract> - Array of contract objects to display in the table.
 */
export default function RenewalTable({
  renewalEvents,
}: {
  renewalEvents: Array<VendorContract>;
}) {
  return (
    <div className="table-container">
      <table>
        <thead>
          <tr>
            <th>Seller</th>
            <th>Renewal Date</th>
          </tr>
        </thead>
        <tbody>
          {renewalEvents.map((event, i) => (
            <tr key={"renewal-list-" + i}>
              <td>{event.seller}</td>
              <td>{standardizeDate(event.renewal_date)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

/**
 * Standardizes the date format for display to MMMM, DD, YYYY
 * @param input The date string to standardize.
 * @returns The standardized date string.
 */
function standardizeDate(input: string): string {
  const date = new Date(input);

  if (isNaN(date.getTime())) {
    throw new Error(`Invalid date: ${input}`);
  }

  return new Intl.DateTimeFormat("en-US", {
    month: "long",
    day: "numeric",
    year: "numeric",
  }).format(date);
}
