"use client";
import { Dispatch, SetStateAction, useRef } from "react";
import { uploadToS3, parseFiles, fetchCalendarData } from "../lib/api";
import { VendorContract } from "../lib/definitions";

/**
 * UploadButton component for uploading PDF files and updating renewal events.
 *
 * Props:
 *   setRenewalEvents: Dispatch to update the renewal events array after upload and parsing.
 *   setLoading: Dispatch to toggle loading state during upload and parsing. Page displays loading
 *   when loading is true.
 *   setDataError: Dispatch to set error messages if upload or parsing fails.
 */
export default function UploadButton({
  setRenewalEvents,
  setLoading,
  setDataError,
}: {
  setRenewalEvents: Dispatch<SetStateAction<VendorContract[]>>;
  setLoading: Dispatch<SetStateAction<boolean>>;
  setDataError: Dispatch<SetStateAction<string | null>>;
}) {
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  // For the button to connect to the hidden input element.
  const handleButtonClick = () => {
    fileInputRef.current?.click();
  };

  // User picks files and they are uploaded, parsed, and then displayed.
  const handleFileChange = async (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    const files: FileList | null = event.target.files;
    if (files && files.length > 0) {
      try {
        setLoading(true);
        await uploadToS3(files);
        await parseFiles(files);
        const data = await fetchCalendarData();
        setRenewalEvents(data);
      } catch (error) {
        console.error("Error uploading files:", error);

        setDataError(
          "File(s): " +
            Array.from(files)
              .map((file) => file.name)
              .join(", ") +
            " failed to parse.",
        );

        // auto-clear after 5 seconds
        setTimeout(() => setDataError(null), 5000);
      } finally {
        setLoading(false);
      }
    }
  };

  return (
    <div>
      <button
        className="bg-blue-500 text-white p-2 rounded text-xl"
        onClick={handleButtonClick}
      >
        Upload PDFs
      </button>
      <input
        type="file"
        multiple
        ref={fileInputRef}
        onChange={handleFileChange}
        accept="application/pdf"
        className="hidden"
      />
    </div>
  );
}
