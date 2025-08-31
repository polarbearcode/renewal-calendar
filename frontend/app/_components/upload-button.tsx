"use client";
import { Dispatch, useRef } from "react";
import { uploadToS3 } from "../lib/api";

// Component for uploading PDF files button
export default function UploadButton() {
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const handleButtonClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files: FileList | null = event.target.files;
    if (files && files.length > 0) {
      uploadToS3(files);
      parseFiles(files);
    }
  };

  return (
    <div>
      <button
        className="bg-blue-500 text-white p-2 rounded text-xl"
        onClick={handleButtonClick}
      >
        Upload PDF
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
